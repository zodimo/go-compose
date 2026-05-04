package state

import (
	"sync"

	"github.com/zodimo/go-compose/state/core"
)

// MutableValueTypedOption configures a mutableValueTyped
type MutableValueTypedOption[T any] = core.MutableValueTypedOption[T]

// WithChangeNotifier sets a callback to be invoked when the value changes.
func WithChangeNotifier[T any](notifier func(T)) MutableValueTypedOption[T] {
	return core.WithChangeNotifier(notifier)
}

// WithPolicy sets a custom MutationPolicy for change detection.
// If not set, StructuralEqualityPolicy is used by default.
func WithPolicy[T any](policy MutationPolicy[T]) MutableValueTypedOption[T] {
	return core.WithPolicy(policy)
}

// MutableValue is a state container that notifies subscribers when its value changes.
type mutableValueTyped[T any] struct {
	cell           T
	changeNotifier func(T)
	mu             sync.RWMutex // RWMutex for thread-safe access (following go-frp Behavior pattern)
	policy         MutationPolicy[T]

	// Subscription support for push-based invalidation
	subscribers *core.SubscriptionManager
}
