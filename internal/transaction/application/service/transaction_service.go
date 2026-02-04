package service

import (
	"context"

	"github.com/Ephraimdebel/transaction-switch/internal/transaction/application"
	"github.com/Ephraimdebel/transaction-switch/internal/transaction/application/command"
	"github.com/Ephraimdebel/transaction-switch/internal/transaction/domain/entity"
	"github.com/Ephraimdebel/transaction-switch/internal/transaction/domain/valueobject"
)

type TransactionService struct {
	repo      application.TransactionRepository
	publisher application.EventPublisher
}

func NewTransactionService(
	repo application.TransactionRepository,
	publisher application.EventPublisher,
) *TransactionService {
	return &TransactionService{repo: repo, publisher: publisher}
}

func (s *TransactionService) Reserve(ctx context.Context, cmd command.ReserveTransaction) error {
	id := valueobject.NewTransactionId(cmd.TransactionID)
	amount, err := valueobject.NewAmount(cmd.Amount)
	if err != nil {
		return err
	}

	tx := entity.NewTransaction(id, amount)

	if err := tx.Reserve(); err != nil {
		return err
	}

	if err := s.repo.Save(ctx, tx); err != nil {
		return err
	}

	events := tx.DomainEvents()
	if err := s.publisher.Publish(ctx, events); err!= nil {
		return err
	}
	tx.ClearDomainEvents()
	
	return nil
}
