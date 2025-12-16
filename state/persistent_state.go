package state

// PersistentState is the interface for the underlying store that holds state values.
// It manages the lifecycle of MutableValues and allows the runtime to react to state changes.
type PersistentState interface {
	// GetState retrieves or creates a MutableValue for the given key.
	// If the state for the key does not exist, it is initialized using the initial function.
	GetState(key string, initial func() any) MutableValue

	// SetOnStateChange registers a callback that is invoked whenever any state managed by this store changes.
	// This is typically used by the runtime to trigger a new frame or recomposition.
	SetOnStateChange(callback func())
}
