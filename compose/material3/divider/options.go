package divider

import (
	"github.com/zodimo/go-compose/compose/ui/graphics"
)

type DividerOptions struct {
	Modifier  Modifier
	Thickness float64
	Color     graphics.Color
}

type DividerOption func(*DividerOptions)

func WithModifier(m Modifier) DividerOption {
	return func(o *DividerOptions) {
		o.Modifier = m
	}
}

func WithThickness(value int) DividerOption {
	return func(o *DividerOptions) {
		o.Thickness = float64(value)
	}
}

func WithColor(col graphics.Color) DividerOption {
	return func(o *DividerOptions) {
		o.Color = col
	}
}
