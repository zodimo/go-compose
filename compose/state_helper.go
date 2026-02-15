package compose

import "github.com/zodimo/go-compose/state"

// local alias for state.State
func Remember[T any](c state.SupportState, key string, initial func() T, options ...state.StateTypedOption[T]) (state.MutableValueTyped[T], error) {
	return state.Remember[T](c, key, initial, options...)
}

// local alias for state.MustState
func MustRemember[T any](c state.SupportState, key string, initial func() T, options ...state.StateTypedOption[T]) state.MutableValueTyped[T] {
	return state.MustRemember[T](c, key, initial, options...)
}

// local alias for state.State
// Deprecated: use Remember instead
func State[T any](c state.SupportState, key string, initial func() T, options ...state.StateTypedOption[T]) (state.MutableValueTyped[T], error) {
	return state.State[T](c, key, initial, options...)
}

// local alias for state.MustState
// Deprecated: use MustRemember instead
func MustState[T any](c state.SupportState, key string, initial func() T, options ...state.StateTypedOption[T]) state.MutableValueTyped[T] {
	return state.MustState[T](c, key, initial, options...)
}

func DerivedStateOf[T any](calculation func() T) *state.DerivedState[T] {
	return state.DerivedStateOf(calculation)
}

func DerivedStateWithPolicy[T any](calculation func() T, policy state.MutationPolicy[T]) *state.DerivedState[T] {
	return state.DerivedStateWithPolicy(calculation, policy)
}
