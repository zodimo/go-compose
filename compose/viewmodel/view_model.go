package viewmodel

import (
	"context"
	"fmt"

	"github.com/zodimo/go-compose/compose/effect"
	"github.com/zodimo/go-compose/lifecycle"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/state"
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

func (v *ViewModel) OnInit() {
}

func (v *ViewModel) MustContext() context.Context {
	if ctx, ok := v.CoroutineScope.Context(); ok {
		return ctx
	}
	panic("CoroutineScope not initialized, use SetViewModelScope")
}

func RememberViewModel[T lifecycle.ViewModel](c api.Composer, factory func() T) T {
	key := c.GenerateID()
	path := c.GetPath()

	viewModelPath := fmt.Sprintf("%d/%s/viewModel", key, path)

	currentViewModel := state.MustRemember(c, viewModelPath, func() T {
		return factory()
	}).Get()

	return currentViewModel
}
