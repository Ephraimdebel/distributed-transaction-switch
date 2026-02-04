package event

import (
	"time"
	"github.com/Ephraimdebel/transaction-switch/internal/transaction/domain/valueobject"
)

type TransactionReserved struct {
	TransactionID valueobject.TransactionID
	Amount        valueobject.Amount
	occurredAt    time.Time
}

func NewTransactionReserved(txID valueobject.TransactionID, amount valueobject.Amount) TransactionReserved {
	return TransactionReserved{
		TransactionID: txID,
		Amount:        amount,
		occurredAt:    time.Now(),
	}
}

func (e TransactionReserved) EventName() string {
	return "TransactionReserved"
}

func (e TransactionReserved) OccurredAt() time.Time {
	return e.occurredAt
}
