package fform

import "github.com/zodimo/go-compose/pkg/stateful"

type FormItemState struct {
	actor *stateful.Actor[*FormFieldState, stateful.Stateful[*FormFieldState]]
}

func NewFormItemState(a *stateful.Actor[*FormFieldState, stateful.Stateful[*FormFieldState]]) *FormItemState {
	return &FormItemState{
		actor: a,
	}
}
