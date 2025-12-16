package icon

import (
	"github.com/zodimo/go-maybe"
)

type IconOptions struct {
	Modifier Modifier
	Color    maybe.Maybe[ColorDescriptor]
}

type IconOption func(*IconOptions)

func DefaultIconOptions() IconOptions {
	return IconOptions{
		Modifier: EmptyModifier,
		Color:    maybe.Some(colorHelper.ColorSelector().SurfaceRoles.OnSurface),
	}
}

func WithModifier(m Modifier) IconOption {
	return func(o *IconOptions) {
		o.Modifier = m
	}
}

func WithColor(desc ColorDescriptor) IconOption {
	return func(o *IconOptions) {
		o.Color = maybe.Some(desc)
	}
}
