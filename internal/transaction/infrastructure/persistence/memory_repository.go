package persistence

import (
	"context"
	"errors"
	"sync"

	"github.com/Ephraimdebel/transaction-switch/internal/transaction/domain/entity"
)

// MemoryTransactionRepository is an in-memory, thread-safe repository
type MemoryTransactionRepository struct {
	mu           sync.RWMutex
	transactions map[string]*entity.Transaction
}

// NewMemoryTransactionRepository initializes the repository
func NewMemoryTransactionRepository() *MemoryTransactionRepository {
	return &MemoryTransactionRepository{
		transactions: make(map[string]*entity.Transaction),
	}
}

// Save stores or updates a transaction
func (r *MemoryTransactionRepository) Save(
	ctx context.Context,
	tx *entity.Transaction,
) error {

	// context check (important for distributed systems)
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	id := tx.ID().String()

	// store a copy to prevent external mutation
	r.transactions[id] = tx.Clone()

	return nil
}

// GetByID retrieves a transaction by ID
func (r *MemoryTransactionRepository) GetByID(
	ctx context.Context,
	id string,
) (*entity.Transaction, error) {

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	tx, exists := r.transactions[id]
	if !exists {
		return nil, errors.New("transaction not found")
	}

	// return a copy
	return tx.Clone(), nil
}
