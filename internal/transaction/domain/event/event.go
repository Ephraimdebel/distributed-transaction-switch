package event

import "time"

type DomainEvent interface {
	EventID() string
	EventName() string
	OccurredAt() time.Time
}
