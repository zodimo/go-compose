package surface

import (
	// This import is still needed for DefaultSurfaceOptions and WithBorder, but the instruction implies it should be removed.
	// The provided snippet for imports is malformed. I will try to infer the correct imports.
	// Based on the struct changes, "image/color" will be replaced by "github.com/zodimo/go-compose/theme".
	// "github.com/zodimo/go-compose/modifiers/shadow" is added in the snippet, but not directly used in the struct definition.
	// "github.com/zodimo/go-compose/internal/modifier" is still needed for Modifier type.

	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	"github.com/zodimo/go-compose/internal/modifier"
	"github.com/zodimo/go-compose/theme" // Added based on instruction
)

type SurfaceOptions struct {
	Modifier        Modifier
	Shape           Shape
	Color           theme.ColorDescriptor // Changed from color.Color
	ContentColor    theme.ColorDescriptor // Changed from color.Color
	TonalElevation  Dp
	ShadowElevation Dp
	BorderWidth     Dp
	BorderColor     theme.ColorDescriptor // Changed from color.Color
	Alignment       box.Direction         // Optional alignment for content inside surface
}

type SurfaceOption func(*SurfaceOptions)

func DefaultSurfaceOptions() SurfaceOptions {
	colorHelper := theme.ColorHelper
	return SurfaceOptions{
		Modifier:        modifier.EmptyModifier,
		Shape:           shape.ShapeRectangle,
		Color:           colorHelper.UnspecifiedColor(),
		ContentColor:    colorHelper.UnspecifiedColor(),
		TonalElevation:  0,
		ShadowElevation: 0,
		BorderWidth:     0,
		BorderColor:     colorHelper.UnspecifiedColor(),
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

func WithColor(colorDesc theme.ColorDescriptor) SurfaceOption {
	return func(o *SurfaceOptions) {
		o.Color = colorDesc
	}
}

func WithContentColor(colorDesc theme.ColorDescriptor) SurfaceOption {
	return func(o *SurfaceOptions) {
		o.ContentColor = colorDesc
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

func WithBorder(width Dp, colorDesc theme.ColorDescriptor) SurfaceOption {
	return func(o *SurfaceOptions) {
		o.BorderWidth = width
		o.BorderColor = colorDesc
	}
}

func WithAlignment(alignment box.Direction) SurfaceOption {
	return func(o *SurfaceOptions) {
		o.Alignment = alignment
	}
}
