package icon

import (
	"github.com/zodimo/go-maybe"
)

type IconOptions struct {
	Modifier  Modifier
	Color     maybe.Maybe[ColorDescriptor]
	LazyColor maybe.Maybe[func() ColorDescriptor]
}

type IconOption func(*IconOptions)

func DefaultIconOptions() IconOptions {
	return IconOptions{
		Modifier: EmptyModifier,
		// Default Fallback is black
		Color:     maybe.None[ColorDescriptor](),
		LazyColor: maybe.None[func() ColorDescriptor](),
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
func WithLazyColor(desc func() ColorDescriptor) IconOption {
	return func(o *IconOptions) {
		o.LazyColor = maybe.Some(desc)
	}
}
