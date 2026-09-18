package components

import (
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/compose/ui/text"
	"github.com/zodimo/go-compose/pkg/sentinel"

	"github.com/zodimo/go-compose/internal/modifier"
)

type TextFieldComponentOptions struct {
	Modifier  ui.Modifier
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

// DEAD: superseded by the form engine (FormNode/Control/Group/Array)
func TextFieldWithModifier(m ui.Modifier) TextFieldComponentOption {
	return func(o *TextFieldComponentOptions) {
		o.Modifier = o.Modifier.Then(m)
	}
}

// DEAD: superseded by the form engine (FormNode/Control/Group/Array)
func TextFieldWithTextStyle(style *text.TextStyle) TextFieldComponentOption {
	return func(o *TextFieldComponentOptions) {
		o.TextStyle = style
	}
}

// DEAD: superseded by the form engine (FormNode/Control/Group/Array)
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

// DEAD: superseded by the form engine (FormNode/Control/Group/Array)
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
