// Package components provides bound Material 3 field widgets driven by the form
// engine: text fields, checkboxes, switches, and selects.
//
// Every component takes a *fform.FormFieldBinding[T] built over a control in the
// form tree, reads its value/errors/touched state, and writes back through
// SetValue. Error display follows the Material 3 convention and is gated on the
// control being touched, so invalid-but-unvisited fields render without error
// styling.
package components

import (
	fform "github.com/zodimo/go-compose/compose/foundation/form"
)

// fieldError returns the touched-gated error message for a control and whether
// the field is currently in error. The message is empty and hasError false when
// the control is valid, disabled, or has not been touched. It is the single
// place the touched-gating rule lives, shared by every field component.
func fieldError[T any](binding *fform.FormFieldBinding[T]) (message string, hasError bool) {
	if binding == nil || !binding.HasErrors() || !binding.IsTouched() {
		return "", false
	}
	return binding.ErrorMessage(), true
}

// supportingText picks the text to show beneath a field: the error message when
// the field is in error, otherwise the caller's hint text. This keeps the hint
// visible for valid fields while letting errors take over the supporting slot,
// matching Material 3's WithSupportingText convention.
func supportingText(hint string, message string, hasError bool) string {
	if hasError && message != "" {
		return message
	}
	return hint
}
