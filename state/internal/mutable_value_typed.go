package internal

import (
	"fmt"
	"sync"

	"github.com/zodimo/go-compose/state/core"
)

var _ core.MutableValueTyped[any] = &MutableValueTypedWrapper[any]{}
var _ core.MutableValue = &MutableValueTypedWrapper[any]{}
var _ core.StateChangeNotifier = &MutableValueTypedWrapper[any]{}

var _ core.MutableValueTyped[any] = &MutableValueTyped[any]{}
var _ core.MutableValue = &MutableValueTyped[any]{}
var _ core.StateChangeNotifier = &MutableValueTyped[any]{}

// MutableValue is a state container that notifies subscribers when its value changes.
type MutableValueTyped[T any] struct {
	cell           T
	changeNotifier func(T)
	mu             sync.RWMutex // RWMutex for thread-safe access (following go-frp Behavior pattern)
	policy         core.MutationPolicy[T]

	// Subscription support for push-based invalidation
	subscribers *core.SubscriptionManager
}

// NewMutableState creates a new typed mutable state with optional configuration.
// This is the Kotlin-aligned API using MutationPolicy.
func NewMutableState[T any](initial T, opts ...core.MutableValueTypedOption[T]) core.MutableValueTyped[T] {
	config := &core.MutableValueTypedConfig[T]{
		ChangeNotifier: func(T) {},
		Policy:         core.StructuralEqualityPolicy[T](),
	}
	for _, opt := range opts {
		opt(config)
	}

	return &MutableValueTyped[T]{
		cell:           initial,
		changeNotifier: config.ChangeNotifier,
		policy:         config.Policy,
		subscribers:    core.NewSubscriptionManager(),
	}
}

func (mv *MutableValueTyped[T]) Get() T {
	NotifyRead(mv)
	mv.mu.RLock()
	value := mv.cell
	mv.mu.RUnlock()
	return value
}

func (mv *MutableValueTyped[T]) Set(value T) {
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

func (mv *MutableValueTyped[T]) CompareAndSet(expect, update T) bool {
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

func (mv *MutableValueTyped[T]) Update(f func(T) T) {
	for {
		current := mv.Get()
		newValue := f(current)
		if mv.CompareAndSet(current, newValue) {
			return
		}
	}
}

// update then get, return new
func (mv *MutableValueTyped[T]) UpdateAndGet(f func(T) T) T {
	for {
		current := mv.Get()
		newValue := f(current)
		if mv.CompareAndSet(current, newValue) {
			return newValue
		}
	}
}

// get then update, return old
func (mv *MutableValueTyped[T]) GetAndUpdate(f func(T) T) T {
	for {
		current := mv.Get()
		newValue := f(current)
		if mv.CompareAndSet(current, newValue) {
			return current
		}
	}
}

func (mv *MutableValueTyped[T]) Unwrap() core.MutableValue {
	return MutableValueTypedToUntyped(mv)
}

// Subscribe registers a callback to be invoked when the value changes.
// Returns a Subscription that can be used to stop receiving notifications.
func (mv *MutableValueTyped[T]) Subscribe(callback func()) core.Subscription {
	return mv.subscribers.Subscribe(callback)
}

// Wrapper

type MutableValueTypedWrapper[T any] struct {
	mv       *MutableValue
	nillable bool
}

func MutableValueToTyped[T any](mv core.MutableValue) (core.MutableValueTyped[T], error) {
	mvTyped, ok := mv.(*MutableValue)
	if !ok {
		return nil, fmt.Errorf("cell is not of type %T, got %T", mvTyped, mv)
	}

	nillable := isNillableType[T]()
	_, err := typeAssert[T](mvTyped.cell, nillable)
	if err != nil {
		return nil, fmt.Errorf("could not convert cell to type %T: %w", mvTyped.cell, err)
	}

	return &MutableValueTypedWrapper[T]{
		mv:       mvTyped,
		nillable: nillable,
	}, nil
}

func (w *MutableValueTypedWrapper[T]) Get() T {
	val := w.mv.Get()
	return typeAssertUnsafe[T](val, w.nillable)
}

func (w *MutableValueTypedWrapper[T]) Set(value T) {
	w.mv.Set(value)
}

func (w *MutableValueTypedWrapper[T]) CompareAndSet(expect, update T) bool {
	return w.mv.CompareAndSet(expect, update)
}

func (w *MutableValueTypedWrapper[T]) Update(f func(T) T) {
	w.mv.Update(func(current any) any {
		return f(typeAssertUnsafe[T](current, w.nillable))
	})
}

func (w *MutableValueTypedWrapper[T]) UpdateAndGet(f func(T) T) T {
	res := w.mv.UpdateAndGet(func(current any) any {
		return f(typeAssertUnsafe[T](current, w.nillable))
	})
	if res == nil {
		var zero T
		return zero
	}
	return res.(T)
}

func (w *MutableValueTypedWrapper[T]) GetAndUpdate(f func(T) T) T {
	res := w.mv.GetAndUpdate(func(current any) any {
		return f(typeAssertUnsafe[T](current, w.nillable))
	})
	if res == nil {
		var zero T
		return zero
	}
	return res.(T)
}

func (w *MutableValueTypedWrapper[T]) Subscribe(callback func()) core.Subscription {
	return w.mv.Subscribe(callback)
}

func (w *MutableValueTypedWrapper[T]) Unwrap() core.MutableValue {
	return w.mv
}
