package viewmodel

import (
	"context"

	"github.com/zodimo/go-compose/compose/effect"
	"github.com/zodimo/go-compose/lifecycle"
)

var _ lifecycle.ViewModel = (*ViewModel)(nil)
var _ lifecycle.HasCoroutineScope = (*ViewModel)(nil)
var _ lifecycle.HasViewModelScope = (*ViewModel)(nil)

type ViewModel struct {
	effect.CoroutineScope
}

func (v *ViewModel) OnCleared() {
}

func (v *ViewModel) SetViewModelScope(ctx context.Context) {
	v.CoroutineScope.SetCoroutineScope(ctx)
}
func (v *ViewModel) MustContext() context.Context {
	if ctx, ok := v.CoroutineScope.Context(); ok {
		return ctx
	}
	panic("CoroutineScope not initialized, use SetViewModelScope")
}
