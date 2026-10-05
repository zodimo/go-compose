package fform

import (
	"fmt"
	"sort"
	"strings"

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

// IsEnabled reports whether the control is enabled (neither own- nor
// inherited-disabled). Field components use it to render disabled styling.
func (b *FormFieldBinding[T]) IsEnabled() bool {
	return b.control.IsEnabled()
}

// IsDirty reports whether the control's value differs from its initial value.
func (b *FormFieldBinding[T]) IsDirty() bool {
	return b.control.IsDirty()
}

// Status returns the control's validation/lifecycle status.
func (b *FormFieldBinding[T]) Status() Status {
	return b.control.Status()
}

// ErrorMessage returns every validation error joined into one supporting-text
// string, or "" when the control is valid or disabled. It is not touched-gated;
// components gate on IsTouched before rendering it.
func (b *FormFieldBinding[T]) ErrorMessage() string {
	return errorMessage(b.control.Errors())
}

// Control returns the underlying control.
func (b *FormFieldBinding[T]) Control() *Control[T] {
	return b.control
}

// errorMessage joins the messages of errs, sorted for determinism, into one
// supporting-text string. A control exposes a single entry keyed by its path
// whose value is the errors.Join aggregation of every validator failure.
func errorMessage(errs map[string]string) string {
	if len(errs) == 0 {
		return ""
	}
	msgs := make([]string, 0, len(errs))
	for _, msg := range errs {
		msgs = append(msgs, msg)
	}
	sort.Strings(msgs)
	return strings.Join(msgs, "\n")
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
