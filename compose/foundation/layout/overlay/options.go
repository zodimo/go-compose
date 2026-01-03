package overlay

import (
	"github.com/zodimo/go-compose/compose/ui/graphics"
)

type OverlayOptions struct {
	Modifier   Modifier
	OnDismiss  func()
	ScrimColor graphics.Color
}

type OverlayOption func(*OverlayOptions)

func DefaultOverlayOptions() OverlayOptions {
	return OverlayOptions{
		Modifier:   EmptyModifier,
		ScrimColor: graphics.ColorUnspecified,
	}
}

func WithModifier(m Modifier) OverlayOption {
	return func(o *OverlayOptions) {
		o.Modifier = m
	}
}

func WithOnDismiss(f func()) OverlayOption {
	return func(o *OverlayOptions) {
		o.OnDismiss = f
	}
}

func WithScrimColor(c graphics.Color) OverlayOption {
	return func(o *OverlayOptions) {
		o.ScrimColor = c
	}
}
