package formengine

import (
	"sync"

	"github.com/zodimo/go-compose/state"
)

// ValueStore is the read/write seam between the form engine and typed value
// storage. Keeping the value behind this interface keeps the engine free of any
// direct dependency on the composer's remember lifecycle (design D2).
type ValueStore[T any] interface {
	Get() T
	Set(T)
}

// PlainValueStore is a mutex-guarded in-memory value holder, used for pure-Go
// (non-UI) use and tests.
type PlainValueStore[T any] struct {
	mu    sync.RWMutex
	value T
}

// NewPlainValueStore creates a PlainValueStore holding initial.
func NewPlainValueStore[T any](initial T) *PlainValueStore[T] {
	return &PlainValueStore[T]{value: initial}
}

// Get returns the current value.
func (s *PlainValueStore[T]) Get() T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.value
}

// Set replaces the current value.
func (s *PlainValueStore[T]) Set(v T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.value = v
}

// mutableValueValueStore adapts a state.MutableValueTyped[T] to ValueStore[T].
// Delegating Set to the typed mutable value drives the framework's whole-tree
// recomposition (design D2).
type mutableValueValueStore[T any] struct {
	mv state.MutableValueTyped[T]
}

// NewMutableValueValueStore adapts mv to a ValueStore[T] for UI use.
func NewMutableValueValueStore[T any](mv state.MutableValueTyped[T]) ValueStore[T] {
	return &mutableValueValueStore[T]{mv: mv}
}

func (s *mutableValueValueStore[T]) Get() T  { return s.mv.Get() }
func (s *mutableValueValueStore[T]) Set(v T) { s.mv.Set(v) }
