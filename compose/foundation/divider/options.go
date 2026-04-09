package divider

import (
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/layout"
)

type DividerOptions struct {
	Modifier  ui.Modifier
	Thickness int
	Color     graphics.Color
	Axis      layout.Axis
}

type DividerOption func(*DividerOptions)

func WithModifier(m ui.Modifier) DividerOption {
	return func(o *DividerOptions) {
		o.Modifier = m
	}
}

func WithThickness(value int) DividerOption {
	return func(o *DividerOptions) {
		o.Thickness = value
	}
}

func WithColor(col graphics.Color) DividerOption {
	return func(o *DividerOptions) {
		o.Color = col
	}
}

func WithAxis(axis layout.Axis) DividerOption {
	return func(o *DividerOptions) {
		o.Axis = axis
	}
}
