package textfield

import (
	"gioui.org/gesture"
	"gioui.org/layout"
	"gioui.org/widget"
)

const Material3TextFieldNodeID = "Material3TextField"

type HandlerWrapper struct {
	Func func(string)
}

type OnSubmitWrapper struct {
	Func func()
}

type TextFieldStateTracker struct {
	LastValue string
}

// TextField implements a Material Design 3 text field.
// It defaults to the Filled variant for backward compatibility,
// but users should prefer explicit Filled or Outlined calls.
func TextField(
	value string,
	onValueChange func(string),
	label string,
	options ...TextFieldOption,
) Composable {
	options = append(options, WithLabel(label))
	return Filled(value, onValueChange, options...)
}

// TextField implements the Material Design Text Field
// described here: https://material.io/components/text-fields
type TextFieldWidget struct {
	// Editor contains the edit buffer.
	widget.Editor
	// click detects when the mouse pointer clicks or hovers
	// within the textfield.
	click gesture.Click

	// Helper text to give additional context to a field.
	Helper string
	// CharLimit specifies the maximum number of characters the text input
	// will allow. Zero means "no limit".
	CharLimit uint
	// Prefix appears before the content of the text input.
	Prefix layout.Widget
	// Suffix appears after the content of the text input.
	Suffix layout.Widget

	// Animation state.
	state
	label  label
	border border
	helper helper
	anim   *Progress

	// errored tracks whether the input is in an errored state.
	// This is orthogonal to the other states: the input can be both errored
	// and inactive for example.
	errored bool
}
