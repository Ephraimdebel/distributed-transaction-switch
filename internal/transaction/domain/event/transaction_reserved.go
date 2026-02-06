package event

import (
	"time"

	"github.com/Ephraimdebel/transaction-switch/internal/transaction/domain/valueobject"
	"github.com/google/uuid"
)

type TransactionReserved struct {
	eventID       string
	TransactionID valueobject.TransactionID
	Amount        valueobject.Amount
	occurredAt    time.Time
}

func NewTransactionReserved(txID valueobject.TransactionID, amount valueobject.Amount) *TransactionReserved {
	return &TransactionReserved{
		eventID:       uuid.NewString(),
		TransactionID: txID,
		Amount:        amount,
		occurredAt:    time.Now(),
	}
}

func (e TransactionReserved) EventID() string{
	return e.eventID
}

func (e TransactionReserved) EventName() string {
	return "TransactionReserved"
}

func (e TransactionReserved) OccurredAt() time.Time {
	return e.occurredAt
}
