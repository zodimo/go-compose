package flow

import (
	"context"
	"sync"
	"sync/atomic"
)

// StateFlow defines the read-only behavior
type StateFlow[T any] interface {
	Value() T
	Flow[T]
}

// MutableStateFlow is the hot, stateful producer
type MutableStateFlow[T any] struct {
	mu          sync.RWMutex
	value       atomic.Value
	subscribers []chan T
}

func NewMutableStateFlow[T any](initial T) *MutableStateFlow[T] {
	flow := &MutableStateFlow[T]{
		subscribers: make([]chan T, 0),
	}
	flow.value.Store(initial)
	return flow
}

// Value returns the current state (Thread-safe)
func (s *MutableStateFlow[T]) Value() T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.value.Load().(T)
}

// Emit updates the value and notifies all collectors
func (s *MutableStateFlow[T]) Emit(value T) {
	s.mu.Lock()
	s.value.Store(value)
	subs := make([]chan T, len(s.subscribers))
	copy(subs, s.subscribers)
	s.mu.Unlock()
	s.notifySubscribers(value, subs)

}

func (s *MutableStateFlow[T]) notifySubscribers(value T, subscribers []chan T) {
	for _, ch := range subscribers {
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
	ch := make(chan T, 1)

	s.mu.Lock()
	// Capture the value AND register the channel in one atomic step
	current := s.value.Load().(T)
	s.subscribers = append(s.subscribers, ch)
	s.mu.Unlock()

	collector(current)

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
	value := f(s.value.Load().(T))
	s.value.Store(value)
	subs := make([]chan T, len(s.subscribers))
	copy(subs, s.subscribers)
	s.mu.Unlock()
	s.notifySubscribers(value, subs)
}
