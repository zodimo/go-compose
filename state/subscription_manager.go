package state

import "sync"

// SubscriptionManager manages a list of subscriber callbacks with thread-safe
// subscription and notification.
type SubscriptionManager struct {
	mu          sync.RWMutex
	subscribers map[*subscription]func()
	nextID      uint64
}

// NewSubscriptionManager creates a new SubscriptionManager.
func NewSubscriptionManager() *SubscriptionManager {
	return &SubscriptionManager{
		subscribers: make(map[*subscription]func()),
	}
}

// Subscribe adds a callback and returns a Subscription to remove it.
func (sm *SubscriptionManager) Subscribe(callback func()) Subscription {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sub := &subscription{}
	sm.subscribers[sub] = callback

	sub.unsubscribe = func() {
		sm.mu.Lock()
		defer sm.mu.Unlock()
		delete(sm.subscribers, sub)
	}

	return sub
}

// NotifyAll calls all registered subscriber callbacks.
// Callbacks are invoked with the read lock held to prevent concurrent modification.
func (sm *SubscriptionManager) NotifyAll() {
	sm.mu.RLock()

	subs := []func(){}
	for _, callback := range sm.subscribers {
		subs = append(subs, callback)
	}
	sm.mu.RUnlock()

	for _, callback := range subs {
		callback()
	}
}

// Clear removes all subscribers.
func (sm *SubscriptionManager) Clear() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.subscribers = make(map[*subscription]func())
}

// Count returns the number of active subscribers.
func (sm *SubscriptionManager) Count() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return len(sm.subscribers)
}
