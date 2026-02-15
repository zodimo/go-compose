package main

import (
	"context"
	"fmt"

	"github.com/zodimo/go-compose/lifecycle"
)

var _ lifecycle.ViewModel = (*InputViewModel)(nil)
var _ lifecycle.HasViewModelScope = (*InputViewModel)(nil)

type InputViewModel struct {
}

func NewInputViewModel() *InputViewModel {
	return &InputViewModel{}
}

func (vm *InputViewModel) OnCleared() {
	fmt.Println("InputViewModel cleared")
}

func (vm *InputViewModel) SetViewModelScope(ctx context.Context) {
	fmt.Println("InputViewModel scope set")

	go func() {
		fmt.Println("InputViewModel scope waiting to be cleared")
		<-ctx.Done()
		fmt.Println("InputViewModel scope cleared")
	}()
}
