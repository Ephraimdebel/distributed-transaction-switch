package messaging

import (
	"time"

	"github.com/Ephraimdebel/transaction-switch/internal/transaction/domain/event"
)

type RetryEvent struct {
	Event       event.DomainEvent
	Attempts    int
	MaxAttempts int
	NextRetry   time.Time
}
