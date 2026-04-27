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

// NewPaint creates a new Paint instance with default values.
func NewPaint() *Paint {
	return &Paint{
		Alpha: 1.0,
		// Default BlendMode is usually SrcOver
		BlendMode: BlendModeSrcOver,
	}
}

func (p *Paint) ApplyStyle(style DrawStyle) {
	p.Style = style
}
