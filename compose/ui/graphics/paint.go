package graphics

// Paint holds the style and color information about how to draw geometries, text and bitmaps.
type Paint struct {
	Alpha       float32
	Color       Color
	Shader      Shader
	BlendMode   BlendMode
	StrokeWidth float32
	Style       DrawStyle
	// Add other fields as needed
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
