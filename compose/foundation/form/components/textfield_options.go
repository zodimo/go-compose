package components

import (
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/pkg/sentinel"

	"github.com/zodimo/go-compose/internal/modifier"
)

// TextFieldComponentOptions configures TextFieldComponent.
type TextFieldComponentOptions struct {
	Modifier ui.Modifier

	// Label is the floating field label. Unspecified renders no label.
	Label string
	// HintText is shown beneath the field when valid. Unspecified renders no
	// supporting line.
	HintText string
	// Inline renders the field alongside surrounding content in a row instead of
	// as a full-width column entry.
	Inline bool
}

// TextFieldComponentOption configures a TextFieldComponent.
type TextFieldComponentOption func(o *TextFieldComponentOptions)

// DefaultTextFieldComponentOptions returns the default field options.
func DefaultTextFieldComponentOptions() TextFieldComponentOptions {
	return TextFieldComponentOptions{
		Modifier: modifier.EmptyModifier,
		Label:    sentinel.StringValueUnspecified,
		HintText: sentinel.StringValueUnspecified,
		Inline:   false,
	}
}

// TextFieldWithModifier appends a modifier to the field.
func TextFieldWithModifier(m ui.Modifier) TextFieldComponentOption {
	return func(o *TextFieldComponentOptions) {
		o.Modifier = o.Modifier.Then(m)
	}
}

// TextFieldWithLabel sets the floating field label.
func TextFieldWithLabel(label string) TextFieldComponentOption {
	return func(o *TextFieldComponentOptions) {
		o.Label = label
	}
}

// TextFieldWithHintText sets the supporting hint shown when the field is valid.
func TextFieldWithHintText(hint string) TextFieldComponentOption {
	return func(o *TextFieldComponentOptions) {
		o.HintText = hint
	}
}

// TextFieldWithInline renders the field alongside surrounding content in a row
// instead of as a full-width column entry.
func TextFieldWithInline(inline bool) TextFieldComponentOption {
	return func(o *TextFieldComponentOptions) {
		o.Inline = inline
	}
}
