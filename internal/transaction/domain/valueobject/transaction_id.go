package valueobject

type TransactionID string

func NewTransactionId(id string) TransactionID{
	return TransactionID(id)
}

func (id TransactionID) String() string {
	return string(id)
}
