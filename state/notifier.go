package state

import (
	"github.com/zodimo/go-compose/state/core"
)

// Subscription represents an active change observation that can be cancelled.
// When no longer needed, call Unsubscribe() to stop receiving notifications.
type Subscription = core.Subscription

// StateChangeNotifier is implemented by state objects that can notify subscribers
// when their value changes. This enables push-based invalidation for derived states.
type StateChangeNotifier = core.StateChangeNotifier

// InvalidationNotifier is an optional interface for derived states that support
// a separate subscription for invalidation events. This allows derived states
// to subscribe to invalidation events (for chain propagation) rather than
// value-change events (for user callbacks).
type InvalidationNotifier = core.InvalidationNotifier

// NewSubscription creates a Subscription with the given unsubscribe function.
var NewSubscription = core.NewNoOpSubscription

var NewNoOpSubscription = core.NewNoOpSubscription
