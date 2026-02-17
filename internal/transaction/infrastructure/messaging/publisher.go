package messaging

import (
	"context"
	"fmt"
	"time"

	"github.com/Ephraimdebel/transaction-switch/internal/transaction/domain/event"
)

type EventPublisher struct {
	queue    chan RetryEvent
	handlers map[string][]func(event.DomainEvent) error
	dlq      *DeadLetterQueue
}

func NewEventPublisher(workCount int, queueSize int, dlq *DeadLetterQueue) *EventPublisher {
	p := &EventPublisher{
		queue:    make(chan RetryEvent, queueSize),
		handlers: make(map[string][]func(event.DomainEvent) error),
		dlq:      dlq,
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

func (p *EventPublisher) Publish(
	ctx context.Context,
	events []event.DomainEvent,
) error {
	for _, e := range events {
		select {
		case p.queue <- RetryEvent{
			Event:       e,
			Attempts:    0,
			MaxAttempts: 3,
			NextRetry:   time.Now(),
		}:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
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
				p.handleRetry(retryEvent, err)
				break
			}
		}
	}
	fmt.Print("worker ", id)

}

func (p *EventPublisher) handleRetry(re RetryEvent, err error) {
	re.Attempts++

	if re.Attempts >= re.MaxAttempts {
		// dead letter queue
		p.dlq.Add(DeadLetterEvent{
			Event:      re.Event,
			Reason:     err.Error(),
			FailedAt:   time.Now(),
			RetryCount: re.Attempts,
		})
		return
	}
	re.NextRetry = time.Now().Add(time.Duration(re.Attempts) * time.Second)

	p.queue <- re
}
