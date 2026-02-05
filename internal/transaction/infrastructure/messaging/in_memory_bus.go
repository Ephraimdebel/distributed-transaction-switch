package messaging

import (
	"context"

	"github.com/Ephraimdebel/transaction-switch/internal/transaction/domain/event"
)

type InMemoryEventPublisher struct {
	handlers map[string][]func(event.DomainEvent)
}

func NewInMemoryEventPublisher() *InMemoryEventPublisher {
	return &InMemoryEventPublisher{
		handlers: make(map[string][]func(event.DomainEvent)),
	}
}

func (p *InMemoryEventPublisher) Register(
	eventName string,
	handler func(event.DomainEvent),
) {
	p.handlers[eventName] = append(p.handlers[eventName], handler)
}

func (p *InMemoryEventPublisher) Publish(ctx context.Context, events []event.DomainEvent) error{

	select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

	for _, e := range events {
		handlers := p.handlers[e.EventName()]

		for _, handler := range handlers {
			// async, non-blocking
			go handler(e)
			
		}
	}
	return nil
}
