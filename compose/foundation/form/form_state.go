package fform

import (
	"fmt"

	"github.com/zodimo/go-compose/pkg/api"
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

func RememberFormFieldState[T any](c api.Composer, itemState state.MutableValueTyped[T]) *FormFieldState[T] {
	key := c.GenerateID()
	path := c.GetPath()

	fieldID := fmt.Sprintf("%d/%s/formField", key, path)

	formFieldState := state.MustRemember[*FormFieldState[T]](c, fieldID, func() *FormFieldState[T] {
		return NewFormFieldState(fieldID, itemState)
	})

	return formFieldState.Get()
}
