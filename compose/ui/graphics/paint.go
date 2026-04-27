package graphics

import (
	skiaiface "github.com/zodimo/go-skia-support/skia/interfaces"
)

// Paint holds the style and color information about how to draw geometries, text and bitmaps.
type Paint struct {
	Alpha       float32
	Color       Color
	Shader      Shader
	BlendMode   BlendMode
	StrokeWidth float32
	Style       DrawStyle
	// Color/Mask filters mapped to Skia's Paint state (aliases to Skia interfaces)
	ColorFilter skiaiface.ColorFilter
	MaskFilter  skiaiface.MaskFilter
}
type PaintOption func(*Paint)

func PaintWithColor(c Color) PaintOption {
	return func(p *Paint) {
		p.Color = c
	}
}
func PaintWithAlpha(alpha float32) PaintOption {
	return func(p *Paint) {
		p.Alpha = alpha
	}
}
func PaintWithBlendMode(blendMode BlendMode) PaintOption {
	return func(p *Paint) {
		p.BlendMode = blendMode
	}
}
func PaintWithStrokeWidth(strokeWidth float32) PaintOption {
	return func(p *Paint) {
		p.StrokeWidth = strokeWidth
	}
}
func PaintWithStyle(style DrawStyle) PaintOption {
	return func(p *Paint) {
		p.Style = style
	}
}

// NewPaint creates a new Paint instance with default values.
func NewPaint(opts ...PaintOption) *Paint {
	p := &Paint{
		Alpha: 1.0,
		// Default BlendMode is usually SrcOver
		BlendMode: BlendModeSrcOver,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(p)
		}
	}
	return p
}

func (p *Paint) ApplyStyle(style DrawStyle) {
	p.Style = style
}
