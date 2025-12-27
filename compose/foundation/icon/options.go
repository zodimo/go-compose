package icon

import (
	"github.com/zodimo/go-compose/compose/ui/graphics"
)

type IconOptions struct {
	Modifier Modifier
	Color    graphics.Color
}

type IconOption func(*IconOptions)

func DefaultIconOptions() IconOptions {
	return IconOptions{
		Modifier: EmptyModifier,
		Color:    graphics.ColorUnspecified,
	}
}

func WithModifier(m Modifier) IconOption {
	return func(o *IconOptions) {
		o.Modifier = m
	}
}

func WithColor(col graphics.Color) IconOption {
	return func(o *IconOptions) {
		o.Color = col
	}
}
