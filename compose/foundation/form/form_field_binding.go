package fform

import (
	"fmt"

	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/state"
)

// FormFieldBinding adapts a *Control[T] for UI consumption.
//
// It reads Value/Errors/IsTouched/HasErrors directly from the control and writes
// through SetValue, which sets the value and marks the control touched.
type FormFieldBinding[T any] struct {
	control *Control[T]
}

// NewFormFieldBinding wraps control in a FormFieldBinding.
func NewFormFieldBinding[T any](control *Control[T]) *FormFieldBinding[T] {
	return &FormFieldBinding[T]{control: control}
}

// Value returns the control's current value.
func (b *FormFieldBinding[T]) Value() T {
	return b.control.Value()
}

// SetValue writes v to the control and marks it touched.
func (b *FormFieldBinding[T]) SetValue(v T) {
	b.control.Set(v)
	b.control.MarkAsTouched()
}

// Errors returns the control's validation errors keyed by path.
func (b *FormFieldBinding[T]) Errors() map[string]string {
	return b.control.Errors()
}

// HasErrors reports whether the control currently has validation errors.
func (b *FormFieldBinding[T]) HasErrors() bool {
	return b.control.HasErrors()
}

// IsTouched reports whether the control has been touched.
func (b *FormFieldBinding[T]) IsTouched() bool {
	return b.control.IsTouched()
}

// Control returns the underlying control.
func (b *FormFieldBinding[T]) Control() *Control[T] {
	return b.control
}

// Deprecated: use RememberFormFieldBinding.
func RememberFormFieldState[T any](c api.Composer, itemState state.MutableValueTyped[T]) *FormFieldBinding[T] {
	return RememberFormFieldBinding(c, itemState)
}

// RememberFormFieldBinding remembers a standalone FormFieldBinding backed by the
// given itemState. This is the standalone (non-tree-integrated) path; the
// tree-integrated path is FormScope.Field.
func RememberFormFieldBinding[T any](c api.Composer, itemState state.MutableValueTyped[T]) *FormFieldBinding[T] {
	fieldID := fmt.Sprintf("%d/%s/formField", c.GenerateID(), c.GetPath())

	binding := state.MustRemember[*FormFieldBinding[T]](c, fieldID, func() *FormFieldBinding[T] {
		control := NewControl(
			NewMutableValueValueStore(itemState),
			itemState.Get(),
		)
		return NewFormFieldBinding(control)
	})

	return binding.Get()
}
