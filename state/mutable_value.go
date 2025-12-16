package state

// MutableValue is a wrapper around a value that can be read and written.
// Changes to the value are propagated to the composition system to trigger updates.
type MutableValue interface {
	// Get retrieves the current value.
	Get() any

	// Set updates the value and notifies listeners (e.g., the composition system) of the change.
	Set(value any)
}
