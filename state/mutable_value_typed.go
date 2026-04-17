package state

import (
	"fmt"
	"reflect"
	"sync"
)

var _ MutableValueTyped[any] = &MutableValueTypedWrapper[any]{}
var _ MutableValue = &MutableValueTypedWrapper[any]{}
var _ StateChangeNotifier = &MutableValueTypedWrapper[any]{}

var _ MutableValueTyped[any] = &mutableValueTyped[any]{}
var _ MutableValue = &mutableValueTyped[any]{}
var _ StateChangeNotifier = &mutableValueTyped[any]{}

// MutableValueTypedOption configures a mutableValueTyped
type MutableValueTypedOption[T any] func(*mutableValueTypedConfig[T])

// mutableValueTypedConfig holds configuration for mutableValueTyped
type mutableValueTypedConfig[T any] struct {
	changeNotifier func(T)
	policy         MutationPolicy[T]
}

// WithChangeNotifier sets a callback to be invoked when the value changes.
func WithChangeNotifier[T any](notifier func(T)) MutableValueTypedOption[T] {
	return func(c *mutableValueTypedConfig[T]) {
		c.changeNotifier = notifier
	}
}

// WithPolicy sets a custom MutationPolicy for change detection.
// If not set, StructuralEqualityPolicy is used by default.
func WithPolicy[T any](policy MutationPolicy[T]) MutableValueTypedOption[T] {
	return func(c *mutableValueTypedConfig[T]) {
		c.policy = policy
	}
}

// MutableValue is a state container that notifies subscribers when its value changes.
type mutableValueTyped[T any] struct {
	cell           T
	changeNotifier func(T)
	mu             sync.RWMutex // RWMutex for thread-safe access (following go-frp Behavior pattern)
	policy         MutationPolicy[T]

	// Subscription support for push-based invalidation
	subscribers *SubscriptionManager
}

// NewMutableState creates a new typed mutable state with optional configuration.
// This is the Kotlin-aligned API using MutationPolicy.
func newMutableState[T any](initial T, opts ...MutableValueTypedOption[T]) MutableValueTyped[T] {
	config := &mutableValueTypedConfig[T]{
		changeNotifier: func(T) {},
		policy:         StructuralEqualityPolicy[T](),
	}
	for _, opt := range opts {
		opt(config)
	}

	return &mutableValueTyped[T]{
		cell:           initial,
		changeNotifier: config.changeNotifier,
		policy:         config.policy,
		subscribers:    NewSubscriptionManager(),
	}
}

func (mv *mutableValueTyped[T]) Get() T {
	NotifyRead(mv)
	mv.mu.RLock()
	value := mv.cell
	mv.mu.RUnlock()
	return value
}

func (mv *mutableValueTyped[T]) Set(value T) {
	mv.mu.Lock()
	changed := !mv.policy.Equivalent(mv.cell, value)
	if changed {
		mv.cell = value
	}
	changeNotifier := mv.changeNotifier
	mv.mu.Unlock()

	if changed {
		// Notify legacy change notifier
		if changeNotifier != nil {
			changeNotifier(value)
		}

		// Notify all subscribers (push invalidation to derived states)
		mv.subscribers.NotifyAll()
	}
}

func (mv *mutableValueTyped[T]) CompareAndSet(expect, update T) bool {
	mv.mu.Lock()
	current := mv.cell

	// Check if current matches expected
	if !mv.policy.Equivalent(current, expect) {
		mv.mu.Unlock()
		return false
	}

	// Check if update is same as current (no change needed)
	if mv.policy.Equivalent(current, update) {
		mv.mu.Unlock()
		return true // CAS succeeded but no notification needed
	}

	// Perform the update
	mv.cell = update
	changeNotifier := mv.changeNotifier
	mv.mu.Unlock()

	// Notify legacy change notifier
	if changeNotifier != nil {
		changeNotifier(update)
	}
	// Notify all subscribers
	mv.subscribers.NotifyAll()
	return true
}

func (mv *mutableValueTyped[T]) Update(f func(T) T) {
	for {
		current := mv.Get()
		newValue := f(current)
		if mv.CompareAndSet(current, newValue) {
			return
		}
	}
}

// update then get, return new
func (mv *mutableValueTyped[T]) UpdateAndGet(f func(T) T) T {
	for {
		current := mv.Get()
		newValue := f(current)
		if mv.CompareAndSet(current, newValue) {
			return newValue
		}
	}
}

// get then update, return old
func (mv *mutableValueTyped[T]) GetAndUpdate(f func(T) T) T {
	for {
		current := mv.Get()
		newValue := f(current)
		if mv.CompareAndSet(current, newValue) {
			return current
		}
	}
}

func (mv *mutableValueTyped[T]) Unwrap() MutableValue {
	return MutableValueTypedToUntyped(mv)
}

// Subscribe registers a callback to be invoked when the value changes.
// Returns a Subscription that can be used to stop receiving notifications.
func (mv *mutableValueTyped[T]) Subscribe(callback func()) Subscription {
	return mv.subscribers.Subscribe(callback)
}

// Wrapper

type MutableValueTypedWrapper[T any] struct {
	mv *mutableValue
}

func MutableValueToTyped[T any](mv MutableValue) (MutableValueTyped[T], error) {
	mvTyped, ok := mv.(*mutableValue)
	if !ok {
		return nil, fmt.Errorf("cell is not of type %T, got %T", mvTyped, mv)
	}

	if mvTyped.cell != nil {
		_, ok = mvTyped.cell.(T)
		if !ok {
			var zero T
			return nil, fmt.Errorf("cell is not of type %T, got %T", zero, mvTyped.cell)
		}
	} else {
		var zero T
		t := reflect.TypeOf(&zero).Elem()
		if t != nil {
			switch t.Kind() {
			case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Func, reflect.Chan, reflect.Interface:
				// nillable, ok
			default:
				return nil, fmt.Errorf("cell is nil but type %T is not nillable", zero)
			}
		}
	}

	return &MutableValueTypedWrapper[T]{
		mv: mvTyped,
	}, nil
}

func (w *MutableValueTypedWrapper[T]) Get() T {
	val := w.mv.Get()
	if val == nil {
		var zero T
		return zero
	}
	return val.(T)
}

func (w *MutableValueTypedWrapper[T]) Set(value T) {
	w.mv.Set(value)
}

func (w *MutableValueTypedWrapper[T]) CompareAndSet(expect, update T) bool {
	return w.mv.CompareAndSet(expect, update)
}

func (w *MutableValueTypedWrapper[T]) Update(f func(T) T) {
	w.mv.Update(func(current any) any {
		var tCurrent T
		if current != nil {
			tCurrent = current.(T)
		}
		return f(tCurrent)
	})
}

func (w *MutableValueTypedWrapper[T]) UpdateAndGet(f func(T) T) T {
	res := w.mv.UpdateAndGet(func(current any) any {
		var tCurrent T
		if current != nil {
			tCurrent = current.(T)
		}
		return f(tCurrent)
	})
	if res == nil {
		var zero T
		return zero
	}
	return res.(T)
}

func (w *MutableValueTypedWrapper[T]) GetAndUpdate(f func(T) T) T {
	res := w.mv.GetAndUpdate(func(current any) any {
		var tCurrent T
		if current != nil {
			tCurrent = current.(T)
		}
		return f(tCurrent)
	})
	if res == nil {
		var zero T
		return zero
	}
	return res.(T)
}

func (w *MutableValueTypedWrapper[T]) Subscribe(callback func()) Subscription {
	return w.mv.Subscribe(callback)
}

func (w *MutableValueTypedWrapper[T]) Unwrap() MutableValue {
	return w.mv
}
