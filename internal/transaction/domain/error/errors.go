package domain

import "errors"

var (
	ErrInvalidStateTransition = errors.New("invalid state transition")
	ErrInvalidAmount          = errors.New("Invalid transaction amount")
	ErrAlreadyFinalized       = errors.New("transaction already finalized")
)
