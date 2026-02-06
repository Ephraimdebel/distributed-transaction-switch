package messaging

import (
	"time"

	"github.com/Ephraimdebel/transaction-switch/internal/transaction/domain/event"
)
type DeadLetterEvent struct {
	Event event.DomainEvent
	Reason string
	FailedAt time.Time
	RetryCount int
}

type DeadLetterQueue struct{
	events []DeadLetterEvent
}

func NewDeadLetterQueue() *DeadLetterQueue {
	return &DeadLetterQueue {
		events: make([]DeadLetterEvent, 0),
	}
}

func (d *DeadLetterQueue) Add(e DeadLetterEvent) {
	d.events = append(d.events, e)
}

func (d *DeadLetterQueue) List() []DeadLetterEvent {
	return d.events
}
