package events

import "sync"

// Bus defines a minimal pub/sub interface for domain events
type Bus interface {
	Publish(evt SessionEvent)
	Subscribe(handler func(SessionEvent)) (unsubscribe func())
}

// InMemoryBus is a simple, threadsafe in-memory event bus
type InMemoryBus struct {
	mu       sync.RWMutex
	handlers map[int]func(SessionEvent)
	nextID   int
}

// NewInMemoryBus creates a new in-memory event bus
func NewInMemoryBus() *InMemoryBus {
	return &InMemoryBus{handlers: make(map[int]func(SessionEvent))}
}

func (b *InMemoryBus) Publish(evt SessionEvent) {
	// Metrics
	observePublished(evt.Type)
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, h := range b.handlers {
		// Call handlers in separate goroutines to avoid blocking
		go h(evt)
	}
}

func (b *InMemoryBus) Subscribe(handler func(SessionEvent)) (unsubscribe func()) {
	b.mu.Lock()
	id := b.nextID
	b.nextID++
	b.handlers[id] = handler
	b.mu.Unlock()
	return func() {
		b.mu.Lock()
		delete(b.handlers, id)
		b.mu.Unlock()
	}
}
