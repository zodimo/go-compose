package state

// SupportState defines the interface for state management within a composition.
// It allows composables to remember values and manage state that persists across recompositions.
type SupportState interface {
	// Remember stores a value computed by calc.
	// The value is calculated only once (or when the key changes, if applicable in future implementations)
	// and returned on subsequent recompositions.
	// This corresponds to transient state that is attached to the current position in the composition.
	Remember(key string, calc func() any) any

	// State returns a MutableValue that persists across recompositions.
	// When the value inside the MutableValue changes, it triggers a recomposition.
	//
	// Parameters:
	//   - key: A unique string to identify this state.
	//   - initial: A function that provides the initial value if the state does not exist.
	State(key string, initial func() any) MutableValue
}
