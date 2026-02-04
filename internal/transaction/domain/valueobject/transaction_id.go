package valueobject

type TransactionID string

func NewTransactionId(id string) TransactionID{
	return TransactionID(id)
}