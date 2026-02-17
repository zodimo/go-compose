package lifecycle

import "context"

// HasCoroutineScope is an optional interface for ViewModels that need a managed context/scope.
// When a Remembered value implements this interface, it will receive a context that is
// cancelled when the leaving the composition.
type HasCoroutineScope interface {
	// SetCoroutineScope provides a context that is cancelled when the remembered value is leaving the composition.
	// This is analogous to viewModelScope in Kotlin, which is a CoroutineScope.
	SetCoroutineScope(ctx context.Context)
}
