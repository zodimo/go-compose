package surface

import (
	"image/color"

	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	"github.com/zodimo/go-compose/internal/modifier"
	"github.com/zodimo/go-compose/theme"
)

type SurfaceOptions struct {
	Modifier        Modifier
	Shape           Shape
	Color           theme.ColorDescriptor
	ContentColor    theme.ColorDescriptor
	TonalElevation  Dp
	ShadowElevation Dp
	BorderWidth     Dp
	BorderColor     theme.ColorDescriptor
	Alignment       box.Direction // Optional alignment for content inside surface
}

type SurfaceOption func(*SurfaceOptions)

func DefaultSurfaceOptions() SurfaceOptions {
	colorHelper := theme.ColorHelper
	return SurfaceOptions{
		Modifier:        modifier.EmptyModifier,
		Shape:           shape.ShapeRectangle,
		Color:           colorHelper.ColorSelector().SurfaceRoles.Surface,
		ContentColor:    colorHelper.ColorSelector().SurfaceRoles.OnSurface,
		TonalElevation:  0,
		ShadowElevation: 0,
		BorderWidth:     0,
		BorderColor:     colorHelper.SpecificColor(color.NRGBA{A: 0}), // Transparent
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
