package messaging

import (
	"time"

	"github.com/Ephraimdebel/transaction-switch/internal/transaction/domain/event"
)

type EventPublisher struct {
	queue    chan RetryEvent
	handlers map[string][]func(event.DomainEvent) error
}

func NewEventPublisher(workCount int, queueSize int) *EventPublisher {
	p := &EventPublisher{
		queue:    make(chan RetryEvent, queueSize),
		handlers: make(map[string][]func(event.DomainEvent) error),
	}

	// start workers

	for i := 0; i < workCount; i++ {
		go p.worker(i)
	}

	return p
}

func (p *EventPublisher) Register(
	eventName string,
	handler func(event.DomainEvent) error,
) {
	p.handlers[eventName] = append(p.handlers[eventName], handler)

}

func (p *EventPublisher) Publish(events ...event.DomainEvent) {
	for _, e := range events {
		p.queue <- RetryEvent{
			Event:       e,
			Attempts:    0,
			MaxAttempts: 3,
		    NextRetry: time.Now(),
		}
	}
}

func (p *EventPublisher) worker(id int) {
	for retryEvent := range p.queue {
		if time.Now().Before(retryEvent.NextRetry) {
			time.Sleep(time.Until(retryEvent.NextRetry))
		}
		handlers := p.handlers[retryEvent.Event.EventName()]

		for _, handler := range handlers {
			err := handler(retryEvent.Event)

			if err != nil {
				p.handleRetry(retryEvent)
				break
			}
		}
	}

}

func (p *EventPublisher) handleRetry(re RetryEvent) {
	re.Attempts++

	if re.Attempts >= re.MaxAttempts {
		// later -> dead letter queue
		return
	}
	re.NextRetry = time.Now().Add(time.Duration(re.Attempts) * time.Second)

	p.queue <- re
}
