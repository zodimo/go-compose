package overlay

import (
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/compose/ui/graphics"
)

type OverlayOptions struct {
	Modifier   ui.Modifier
	OnClick    func()
	ScrimColor graphics.Color
}

type OverlayOption func(*OverlayOptions)

func DefaultOverlayOptions() OverlayOptions {
	return OverlayOptions{
		Modifier:   ui.EmptyModifier,
		ScrimColor: graphics.ColorCompanion.Unspecified(),
	}
}

func WithModifier(m ui.Modifier) OverlayOption {
	return func(o *OverlayOptions) {
		o.Modifier = m
	}
}

func WithOnClick(f func()) OverlayOption {
	return func(o *OverlayOptions) {
		o.OnClick = f
	}
}

func WithScrimColor(c graphics.Color) OverlayOption {
	return func(o *OverlayOptions) {
		o.ScrimColor = c
	}
}
