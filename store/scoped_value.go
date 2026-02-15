package store

import (
	"context"
)

type ScopedValue struct {
	cancelFunc context.CancelFunc // Context cancel functions for ViewModelScope

	onForgotten func() // remember observer
	onCleared   func() // view model observer
}

func NewScopedValue(
	cancelFunc context.CancelFunc,
	onForgotten func(),
	onCleared func(),
) ScopedValue {
	return ScopedValue{
		cancelFunc:  cancelFunc,
		onForgotten: onForgotten,
		onCleared:   onCleared,
	}
}

func (sv ScopedValue) Cleanup() {
	sv.cancelFunc()
	sv.onForgotten()
	sv.onCleared()
}
