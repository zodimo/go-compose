package components

import (
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/internal/modifier"
	"github.com/zodimo/go-compose/pkg/sentinel"
)

// CheckboxComponentOptions configures CheckboxComponent.
type CheckboxComponentOptions struct {
	Modifier ui.Modifier

	// Label is rendered inline beside the checkbox. Unspecified renders no label.
	Label string
	// HintText is shown beneath the checkbox when valid. Unspecified renders no
	// supporting line.
	HintText string

	// ErrorIndent indents the error supporting text so it visually aligns with
	// the label rather than the checkbox glyph.
	ErrorIndent int
	// ErrorSpacing is the vertical gap before the error supporting text.
	ErrorSpacing int
}

// CheckboxComponentOption configures a CheckboxComponent.
type CheckboxComponentOption func(o *CheckboxComponentOptions)

// DefaultCheckboxComponentOptions returns the default checkbox options.
func DefaultCheckboxComponentOptions() CheckboxComponentOptions {
	return CheckboxComponentOptions{
		Modifier:     modifier.EmptyModifier,
		Label:        sentinel.StringValueUnspecified,
		HintText:     sentinel.StringValueUnspecified,
		ErrorIndent:  0,
		ErrorSpacing: 4,
	}
}

// CheckboxWithModifier appends a modifier to the checkbox.
func CheckboxWithModifier(m ui.Modifier) CheckboxComponentOption {
	return func(o *CheckboxComponentOptions) {
		o.Modifier = o.Modifier.Then(m)
	}
}

// CheckboxWithLabel sets the inline label.
func CheckboxWithLabel(label string) CheckboxComponentOption {
	return func(o *CheckboxComponentOptions) {
		o.Label = label
	}
}

// CheckboxWithHintText sets the supporting hint shown when the field is valid.
func CheckboxWithHintText(hint string) CheckboxComponentOption {
	return func(o *CheckboxComponentOptions) {
		o.HintText = hint
	}
}

// CheckboxWithErrorIndent sets the indent of the error supporting text.
func CheckboxWithErrorIndent(indent int) CheckboxComponentOption {
	return func(o *CheckboxComponentOptions) {
		o.ErrorIndent = indent
	}
}
