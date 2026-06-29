package orchestration

import (
	"sync"

	"github.com/michaelxu2288/agentctl/internal/domain"
)

type EventBus struct {
	mu      sync.RWMutex
	history []domain.Event
}

func NewEventBus() *EventBus {
	return &EventBus{history: make([]domain.Event, 0, 256)}
}

func (e *EventBus) Publish(ev domain.Event) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.history = append(e.history, ev)
}

func (e *EventBus) History() []domain.Event {
	e.mu.RLock()
	defer e.mu.RUnlock()
	out := make([]domain.Event, len(e.history))
	copy(out, e.history)
	return out
}
