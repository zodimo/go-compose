package state

import "sync"

// Subscription represents an active change observation that can be cancelled.
// When no longer needed, call Unsubscribe() to stop receiving notifications.
type Subscription interface {
	Unsubscribe()
}

// StateChangeNotifier is implemented by state objects that can notify subscribers
// when their value changes. This enables push-based invalidation for derived states.
type StateChangeNotifier interface {
	// Subscribe registers a callback to be invoked when the state changes.
	// Returns a Subscription that can be used to stop receiving notifications.
	// The callback may be invoked from any goroutine - implementations must be thread-safe.
	Subscribe(callback func()) Subscription
}

// InvalidationNotifier is an optional interface for derived states that support
// a separate subscription for invalidation events. This allows derived states
// to subscribe to invalidation events (for chain propagation) rather than
// value-change events (for user callbacks).
type InvalidationNotifier interface {
	// SubscribeForInvalidation registers a callback to be invoked when this
	// state is invalidated (not when its value changes).
	SubscribeForInvalidation(callback func()) Subscription
}

// subscription is the default implementation of Subscription
type subscription struct {
	unsubscribe func()
	once        sync.Once
}

func (s *subscription) Unsubscribe() {
	s.once.Do(func() {
		if s.unsubscribe != nil {
			s.unsubscribe()
		}
	})
}

// NewSubscription creates a Subscription with the given unsubscribe function.
func NewSubscription(unsubscribe func()) Subscription {
	return &subscription{unsubscribe: unsubscribe}
}

func NewNoOpSubscription() Subscription {
	return &subscription{unsubscribe: func() {}}
}
