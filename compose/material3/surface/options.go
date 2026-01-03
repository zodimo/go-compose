package surface

import (
	// This import is still needed for DefaultSurfaceOptions and WithBorder, but the instruction implies it should be removed.
	// The provided snippet for imports is malformed. I will try to infer the correct imports.
	// Based on the struct changes, "image/color" will be replaced by "github.com/zodimo/go-compose/theme".
	// "github.com/zodimo/go-compose/modifiers/shadow" is added in the snippet, but not directly used in the struct definition.
	// "github.com/zodimo/go-compose/internal/modifier" is still needed for Modifier type.

	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	"github.com/zodimo/go-compose/internal/modifier"
	// Added based on instruction
)

type SurfaceOptions struct {
	Modifier        Modifier
	Shape           Shape
	Color           graphics.Color
	ContentColor    graphics.Color
	TonalElevation  Dp
	ShadowElevation Dp
	BorderWidth     Dp
	BorderColor     graphics.Color
	Alignment       box.Direction // Optional alignment for content inside surface
}

type SurfaceOption func(*SurfaceOptions)

func DefaultSurfaceOptions() SurfaceOptions {
	return SurfaceOptions{
		Modifier:        modifier.EmptyModifier,
		Shape:           shape.ShapeRectangle,
		Color:           graphics.ColorUnspecified,
		ContentColor:    graphics.ColorUnspecified,
		TonalElevation:  0,
		ShadowElevation: 0,
		BorderWidth:     0,
		BorderColor:     graphics.ColorUnspecified,
		Alignment:       box.NW,
	}
}

func WithModifier(m Modifier) SurfaceOption {
	return func(o *SurfaceOptions) {
		o.Modifier = m
	}
}

func WithShape(s Shape) SurfaceOption {
	return func(o *SurfaceOptions) {
		o.Shape = s
	}
}

func WithColor(col graphics.Color) SurfaceOption {
	return func(o *SurfaceOptions) {
		o.Color = col
	}
}

func WithContentColor(col graphics.Color) SurfaceOption {
	return func(o *SurfaceOptions) {
		o.ContentColor = col
	}
}

func WithTonalElevation(elevation Dp) SurfaceOption {
	return func(o *SurfaceOptions) {
		o.TonalElevation = elevation
	}
}

func WithShadowElevation(elevation Dp) SurfaceOption {
	return func(o *SurfaceOptions) {
		o.ShadowElevation = elevation
	}
}

func WithBorder(width Dp, color graphics.Color) SurfaceOption {
	return func(o *SurfaceOptions) {
		o.BorderWidth = width
		o.BorderColor = color
	}
}

func WithAlignment(alignment box.Direction) SurfaceOption {
	return func(o *SurfaceOptions) {
		o.Alignment = alignment
	}
}
