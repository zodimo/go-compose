package components

import (
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/compose/ui/text"
	"github.com/zodimo/go-compose/pkg/sentinel"

	"github.com/zodimo/go-compose/internal/modifier"
)

type TextFieldComponentOptions struct {
	Modifier ui.Modifier

	// TextStyle is reserved: the Material 3 textfield option takes a different
	// TextStyle type (foundation/next/text) than the form's ui/text style, so it
	// is not forwarded to the field. The form-level style is provided through
	// Form(WithTextStyle(...)) instead.
	TextStyle *text.TextStyle

	Label    string
	Inline   bool
	HintText string
}

type TextFieldComponentOption func(o *TextFieldComponentOptions)

func DefaultTextFieldComponentOptions() TextFieldComponentOptions {
	return TextFieldComponentOptions{
		Modifier:  modifier.EmptyModifier,
		TextStyle: text.TextStyleUnspecified,
		Label:     sentinel.StringValueUnspecified,
		Inline:    false,
		HintText:  sentinel.StringValueUnspecified,
	}
}

// TextFieldWithModifier appends a modifier to the field.
func TextFieldWithModifier(m ui.Modifier) TextFieldComponentOption {
	return func(o *TextFieldComponentOptions) {
		o.Modifier = o.Modifier.Then(m)
	}
}

// TextFieldWithTextStyle sets the field's text style.
func TextFieldWithTextStyle(style *text.TextStyle) TextFieldComponentOption {
	return func(o *TextFieldComponentOptions) {
		o.TextStyle = style
	}
}

// TextFieldWithTextStyleOption merges a text style option into the field's style.
func TextFieldWithTextStyleOption(textStyleOption text.TextStyleOption) TextFieldComponentOption {
	return func(o *TextFieldComponentOptions) {
		o.TextStyle = text.CopyTextStyle(o.TextStyle, textStyleOption)
	}
}

func TextFieldWithLabel(label string) TextFieldComponentOption {
	return func(o *TextFieldComponentOptions) {
		o.Label = label
	}
}

// TextFieldWithInline renders the field alongside surrounding content in a row
// instead of as a full-width column entry.
func TextFieldWithInline(inline bool) TextFieldComponentOption {
	return func(o *TextFieldComponentOptions) {
		o.Inline = inline
	}
}

func TextFieldWithHintText(hint string) TextFieldComponentOption {
	return func(o *TextFieldComponentOptions) {
		o.HintText = hint
	}
}
