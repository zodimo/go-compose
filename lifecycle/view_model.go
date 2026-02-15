package lifecycle

import "context"

// ViewModel is an interface that ViewModels can implement to be notified when they are cleared.
// In Android/Kotlin, ViewModel.onCleared() is called when the ViewModel is no longer used.
type ViewModel interface {
	// OnCleared is called when the ViewModel is no longer used and is being destroyed.
	// Use this to clean up resources such as cancelling background tasks or closing connections.
	OnCleared()
}

// HasViewModelScope is an optional interface for ViewModels that need a managed context/scope.
// When a ViewModel implements this interface, it will receive a context that is
// cancelled when the ViewModel is cleared, enabling graceful shutdown of goroutines.
type HasViewModelScope interface {
	// SetViewModelScope provides a context that is cancelled when the ViewModel is cleared.
	// This is analogous to viewModelScope in Kotlin, which is a CoroutineScope.
	SetViewModelScope(ctx context.Context)
}
