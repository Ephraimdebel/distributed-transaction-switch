package entity

import (
	"time"

	domainErr "github.com/Ephraimdebel/transaction-switch/internal/transaction/domain/error"
	"github.com/Ephraimdebel/transaction-switch/internal/transaction/domain/event"
	"github.com/Ephraimdebel/transaction-switch/internal/transaction/domain/state"
	"github.com/Ephraimdebel/transaction-switch/internal/transaction/domain/valueobject"
)

type Transaction struct {
	id        valueobject.TransactionID
	amount    valueobject.Amount
	state     state.State
	createdAt time.Time
	updatedAt time.Time

	events []event.DomainEvent
}

func NewTransaction(id valueobject.TransactionID, amount valueobject.Amount) *Transaction {
	now := time.Now()

	return &Transaction{
		id:        id,
		amount:    amount,
		state:     state.StateInitiated,
		createdAt: now,
		updatedAt: now,
	}
}

// reserve
func (t *Transaction) Reserve() error {
	if t.state != state.StateInitiated {
		return domainErr.ErrInvalidStateTransition
	}
	t.state = state.StateReserved
	t.updatedAt = time.Now()

	// handle events
	t.events = append(t.events, event.NewTransactionReserved(t.id,t.amount))

	return nil
}

// offline completion
func (t *Transaction) CompleteOffline() error {
	if t.state != state.StateReserved {
		return domainErr.ErrInvalidStateTransition
	}

	t.state = state.StateOfflineCompleted
	t.updatedAt = time.Now()

	return nil
}

// mark pending sync
func (t *Transaction) MarkPendingSync() error {
	if t.state != state.StateOfflineCompleted {
		return domainErr.ErrInvalidStateTransition
	}
	t.state = state.StatePendingSync
	t.updatedAt = time.Now()

	return nil
}

// validate
func (t *Transaction) Validate() error {
	if t.state != state.StateReserved && t.state != state.StatePendingSync {
		return domainErr.ErrInvalidStateTransition
	}

	t.state = state.StateValidating
	t.updatedAt = time.Now()

	return nil
}

// commit
func (t *Transaction) Commit() error {
	if t.state != state.StateValidating {
		return domainErr.ErrInvalidStateTransition
	}

	t.state = state.StateCommitted
	t.updatedAt = time.Now()

	return nil
}

// complete
func (t *Transaction) Complete() error {
	if t.state != state.StateCommitted {
		return domainErr.ErrInvalidStateTransition
	}

	t.state = state.StateCompleted
	t.updatedAt = time.Now()

	return nil
}

// fail
func (t *Transaction) Fail() {
	t.state = state.StateFailed
	t.updatedAt = time.Now()
}

// reject
func (t *Transaction) Reject() {
	t.state = state.StateRejected
	t.updatedAt = time.Now()
}

// expire
func (t *Transaction) Expire() {
	t.state = state.StateExpired
	t.updatedAt = time.Now()

}

// events

// func (t *Transaction) pullEvents() []event.DomainEvent {
// 	events := t.events
// 	t.events = nil
// 	return events
// }

func (t *Transaction) DomainEvents() []event.DomainEvent {
	return t.events
}

func (t *Transaction) ClearDomainEvents() {
	t.events = nil
}


func (t *Transaction) ID() valueobject.TransactionID {
	return t.id
}

func (t *Transaction) Amount() valueobject.Amount {
	return t.amount
}

func (t *Transaction) State() state.State {
	return t.state
}

func (t *Transaction) CreatedAt() time.Time {
	return t.createdAt
}

func (t *Transaction) UpdatedAt() time.Time {
	return t.updatedAt
}

func (t *Transaction) Clone() *Transaction {
	eventsCopy := make([]event.DomainEvent, len(t.events))
	copy(eventsCopy, t.events)

	return &Transaction{
		id:        t.id,
		amount:    t.amount,
		state:     t.state,
		createdAt: t.createdAt,
		updatedAt: t.updatedAt,
		events:    eventsCopy,
	}
}
