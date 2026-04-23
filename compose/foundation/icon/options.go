package icon

import (
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

type IconOptions struct {
	Modifier ui.Modifier
	Color    graphics.Color
	Size     unit.Dp // Unified size for both IconBytes and SymbolName

}

type IconOption func(*IconOptions)

func DefaultIconOptions() IconOptions {
	return IconOptions{
		Modifier: ui.EmptyModifier,
		Color:    graphics.ColorUnspecified,
		Size:     unit.DpUnspecified,
	}
}

func WithModifier(m ui.Modifier) IconOption {
	return func(o *IconOptions) {
		o.Modifier = m
	}
}

func WithColor(col graphics.Color) IconOption {
	return func(o *IconOptions) {
		o.Color = col
	}
}

func WithSize(size unit.Dp) IconOption {
	return func(o *IconOptions) {
		o.Size = size
	}
}
