package fform

import (
	"context"
	"fmt"

	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/pkg/stateful"
	"github.com/zodimo/go-compose/state"
)

type FormState struct{}

// RememberFormState creates a FormState that is remembered across compositions.
func RememberFormState(c api.Composer) *FormState {

	// key := c.GenerateID()
	// path := c.GetPath()

	// expandedPath := fmt.Sprintf("%d/%s/expanded", key, path)
	// selectedPath := fmt.Sprintf("%d/%s/selected", key, path)

	// expanded := state.MustRemember(c, expandedPath, func() map[any]bool {
	// 	return make(map[any]bool)
	// })
	// selected := state.MustRemember(c, selectedPath, func() map[any]bool {
	// 	return make(map[any]bool)
	// })

	return &FormState{
		// expandedItems: expanded,
		// selectedItems: selected,
	}
}

func RememberFormItemState(c api.Composer) *FormItemState {
	key := c.GenerateID()
	path := c.GetPath()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fieldID := fmt.Sprintf("%d/%s/fieldID", key, path)
	onForgottenPath := fmt.Sprintf("%d/%s/onForgottenPath", key, path)
	actorPath := fmt.Sprintf("%d/%s/actor", key, path)

	onForgotten := state.MustRemember(c, onForgottenPath, func() func() {
		fmt.Println("FormItemState forgotten, calling cancel")
		return cancel
	})

	actor := state.MustRemember(
		c,
		actorPath,
		func() *stateful.Actor[*FormFieldState, stateful.Stateful[*FormFieldState]] {
			return stateful.NewActor(ctx, stateful.NewStateful(NewFormFieldState(fieldID), func(s *FormFieldState) *FormFieldState { return s.Clone() }))
		},
		state.WithTypedOnForgotten[*stateful.Actor[*FormFieldState, stateful.Stateful[*FormFieldState]]](onForgotten.Get()),
	)

	return NewFormItemState(actor.Get())
}
