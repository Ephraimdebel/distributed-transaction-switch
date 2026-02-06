package messaging

import (
	"sync"
)

type ProcessedEventStore struct {
	mu     sync.Mutex
	events map[string]struct{}
}

func NewProcessedEventStore() *ProcessedEventStore {
	return &ProcessedEventStore{
		events: make(map[string]struct{}),
	}
}

func (s *ProcessedEventStore) Has(eventID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_,ok := s.events[eventID]
	return ok
}

func (s *ProcessedEventStore) Mark(eventID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events[eventID] = struct{}{}
} 