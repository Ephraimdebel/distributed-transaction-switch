package application

import (
	"context"

	"github.com/Ephraimdebel/transaction-switch/internal/transaction/domain/entity"
	"github.com/Ephraimdebel/transaction-switch/internal/transaction/domain/event"
)

type TransactionRepository interface {
	Save(ctx context.Context, tx *entity.Transaction) error
	GetByID(ctx context.Context, id string) (*entity.Transaction, error)
}

type EventPublisher interface {
	Publish(ctx context.Context, events []event.DomainEvent) error
}
