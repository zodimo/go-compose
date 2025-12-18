package flow

import (
	"context"
	"sync"
)

// StateFlow defines the read-only behavior
type StateFlow[T any] interface {
	Value() T
	Flow[T]
}

// MutableStateFlow is the hot, stateful producer
type MutableStateFlow[T any] struct {
	mu          sync.RWMutex
	value       T
	subscribers []chan T
}

func NewMutableStateFlow[T any](initial T) *MutableStateFlow[T] {
	return &MutableStateFlow[T]{
		value:       initial,
		subscribers: make([]chan T, 0),
	}
}

// Value returns the current state (Thread-safe)
func (s *MutableStateFlow[T]) Value() T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.value
}

// Emit updates the value and notifies all collectors
func (s *MutableStateFlow[T]) Emit(value T) {
	s.mu.Lock()
	s.value = value
	subs := make([]chan T, len(s.subscribers))
	copy(subs, s.subscribers)
	s.mu.Unlock()

	for _, ch := range subs {
		// Non-blocking send (Conflation)
		// If the subscriber is slow, we drain the old value and send the new one
		select {
		case ch <- value:
		default:
			// Buffer is full (slow collector); drain the old and replace with new
			select {
			case <-ch:
			default:
			}
			ch <- value
		}
	}
}

// Collect follows the Kotlin pattern: it blocks until the context is cancelled
func (s *MutableStateFlow[T]) Collect(ctx context.Context, collector func(T)) error {
	// 1. Immediately emit the current value to the collector
	current := s.Value()
	collector(current)

	// 2. Register a new subscriber channel with a buffer of 1 for conflation
	ch := make(chan T, 1)
	s.mu.Lock()
	s.subscribers = append(s.subscribers, ch)
	s.mu.Unlock()

	// 3. Cleanup when collection ends
	defer func() {
		s.mu.Lock()
		for i, sub := range s.subscribers {
			if sub == ch {
				s.subscribers = append(s.subscribers[:i], s.subscribers[i+1:]...)
				break
			}
		}
		s.mu.Unlock()
		close(ch)
	}()

	// 4. Listen for updates
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case val := <-ch:
			collector(val)
		}
	}
}

// Update provides an atomic way to modify state (like Kotlin's .update { ... })
func (s *MutableStateFlow[T]) Update(f func(current T) T) {
	s.mu.Lock()
	next := f(s.value)
	s.mu.Unlock()
	s.Emit(next)
}
