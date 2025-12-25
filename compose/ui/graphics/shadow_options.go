package graphics

import (
	"github.com/zodimo/go-compose/compose/ui/geometry"
	"github.com/zodimo/go-compose/pkg/floatutils"
)

var ShadowOptionsDefault = ShadowOptions{
	Color:      ColorUnspecified,
	Offset:     geometry.OffsetUnspecified,
	BlurRadius: floatutils.Float32Unspecified,
}

type ShadowOptions struct {
	Color      Color
	Offset     geometry.Offset
	BlurRadius float32
}

type ShadowOption func(*ShadowOptions)

func WithColor(color Color) ShadowOption {
	return func(o *ShadowOptions) {
		o.Color = color
	}
}

func WithOffset(offset geometry.Offset) ShadowOption {
	return func(o *ShadowOptions) {
		o.Offset = offset
	}
}

func WithBlurRadius(blurRadius float32) ShadowOption {
	return func(o *ShadowOptions) {
		o.BlurRadius = blurRadius
	}
}
