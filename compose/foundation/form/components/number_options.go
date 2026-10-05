package components

import (
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/internal/modifier"
	"github.com/zodimo/go-compose/pkg/sentinel"
)

// NumberFieldComponentOptions configures NumberFieldComponent.
type NumberFieldComponentOptions struct {
	Modifier ui.Modifier

	// Label is the floating field label. Unspecified renders no label.
	Label string
	// HintText is shown beneath the field when valid. Unspecified renders no
	// supporting line.
	HintText string

	// TextStateKey is the remember key for the raw text buffer. It must be
	// unique per field instance within a composition; an empty value falls back
	// to "numberFieldText", which is only safe when a single number field is
	// rendered under one path.
	TextStateKey string
}

// NumberFieldComponentOption configures a NumberFieldComponent.
type NumberFieldComponentOption func(o *NumberFieldComponentOptions)

// DefaultNumberFieldComponentOptions returns the default number field options.
func DefaultNumberFieldComponentOptions() NumberFieldComponentOptions {
	return NumberFieldComponentOptions{
		Modifier:     modifier.EmptyModifier,
		Label:        sentinel.StringValueUnspecified,
		HintText:     sentinel.StringValueUnspecified,
		TextStateKey: "",
	}
}

// NumberFieldWithModifier appends a modifier to the field.
func NumberFieldWithModifier(m ui.Modifier) NumberFieldComponentOption {
	return func(o *NumberFieldComponentOptions) {
		o.Modifier = o.Modifier.Then(m)
	}
}

// NumberFieldWithLabel sets the floating field label.
func NumberFieldWithLabel(label string) NumberFieldComponentOption {
	return func(o *NumberFieldComponentOptions) {
		o.Label = label
	}
}

// NumberFieldWithHintText sets the supporting hint shown when the field is valid.
func NumberFieldWithHintText(hint string) NumberFieldComponentOption {
	return func(o *NumberFieldComponentOptions) {
		o.HintText = hint
	}
}

// NumberFieldWithTextStateKey overrides the remember key for the raw text
// buffer. Use it when more than one number field is rendered under the same
// composition path.
func NumberFieldWithTextStateKey(key string) NumberFieldComponentOption {
	return func(o *NumberFieldComponentOptions) {
		o.TextStateKey = key
	}
}
