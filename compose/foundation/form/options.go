package fform

import (
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/compose/ui/text"

	"github.com/zodimo/go-compose/internal/modifier"
)

type FormOptions struct {
	Modifier ui.Modifier

	TextStyle *text.TextStyle

	Orientation Orientation
}

type FormOption func(o *FormOptions)

func DefaultFormOptions() FormOptions {
	return FormOptions{
		Modifier:    modifier.EmptyModifier,
		TextStyle:   text.TextStyleUnspecified,
		Orientation: Horizontal,
	}
}

func WithModifier(m ui.Modifier) FormOption {
	return func(o *FormOptions) {
		o.Modifier = o.Modifier.Then(m)
	}
}

// DEAD: superseded by the form engine (FormNode/Control/Group/Array)
func WithTextStyle(style *text.TextStyle) FormOption {
	return func(o *FormOptions) {
		o.TextStyle = style
	}
}

// DEAD: superseded by the form engine (FormNode/Control/Group/Array)
func WithTextStyleOption(textStyleOption text.TextStyleOption) FormOption {
	return func(o *FormOptions) {
		o.TextStyle = text.CopyTextStyle(o.TextStyle, textStyleOption)
	}
}

// DEAD: superseded by the form engine (FormNode/Control/Group/Array)
func WithOrientation(orentation Orientation) FormOption {
	return func(o *FormOptions) {
		o.Orientation = orentation
	}
}
