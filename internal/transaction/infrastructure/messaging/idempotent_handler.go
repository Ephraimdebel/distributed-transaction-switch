package messaging

import (
	"github.com/Ephraimdebel/transaction-switch/internal/transaction/domain/event"
)

func IdempotentHandler(
	store *ProcessedEventStore,
	handler func(event.DomainEvent) error,
) func (event.DomainEvent) error {
	return func(e event.DomainEvent) error {
		if store.Has(e.EventID()) {
			return nil
		}

		if err := handler(e); err != nil {
			return err
		}

		store.Mark(e.EventID())

		return nil
	}
}