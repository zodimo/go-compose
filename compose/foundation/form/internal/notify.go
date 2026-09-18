package formengine

import (
	"sync"

	"github.com/zodimo/go-compose/state"
)

// typedStream is a per-node, per-kind change stream with a typed payload. It is
// backed by state.SubscriptionManager, which only stores func() callbacks: the
// stream captures the latest payload and hands it to every subscriber on notify.
// Each node owns three streams (value, status, touch).
type typedStream[T any] struct {
	sm      *state.SubscriptionManager
	mu      sync.RWMutex
	current T
}

func newTypedStream[T any](initial T) *typedStream[T] {
	return &typedStream[T]{sm: state.NewSubscriptionManager(), current: initial}
}

// subscribe registers fn, which is invoked with the current payload on each
// notify. It returns a state.Subscription for idempotent unsubscription.
func (s *typedStream[T]) subscribe(fn func(T)) state.Subscription {
	return s.sm.Subscribe(func() {
		s.mu.RLock()
		v := s.current
		s.mu.RUnlock()
		fn(v)
	})
}

// notify stores the payload and fans out to all subscribers.
func (s *typedStream[T]) notify(v T) {
	s.mu.Lock()
	s.current = v
	s.mu.Unlock()
	s.sm.NotifyAll()
}
