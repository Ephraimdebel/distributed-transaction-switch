package valueobject

import (
	domainErr "github.com/Ephraimdebel/transaction-switch/internal/transaction/domain/error"
)
type Amount int64

func NewAmount(value int64) (Amount, error) {
	if value <= 0 {
		return 0, domainErr.ErrInvalidAmount
	}
	return Amount(value),nil
}