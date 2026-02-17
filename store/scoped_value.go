package store

import (
	"context"
)

type ScopedValue struct {
	cancelFuncs []context.CancelFunc // Context cancel functions for ViewModelScope

	onForgotten func() // remember observer
	onCleared   func() // view model observer
}

func NewScopedValue(
	cancelFuncs []context.CancelFunc,
	onForgotten func(),
	onCleared func(),
) ScopedValue {
	return ScopedValue{
		cancelFuncs: cancelFuncs,
		onForgotten: onForgotten,
		onCleared:   onCleared,
	}
}

func (sv ScopedValue) Cleanup() {
	for _, cancelFunc := range sv.cancelFuncs {
		cancelFunc()
	}
	sv.onForgotten()
	sv.onCleared()
}
