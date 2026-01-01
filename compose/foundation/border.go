package foundation

import (
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

type BorderStroke struct {
	Width unit.Dp
	Brush graphics.Brush
}

// Helper to create a BorderStroke with a solid color
func BorderStrokeAndColor(width unit.Dp, color graphics.Color) BorderStroke {
	return BorderStroke{
		Width: width,
		Brush: graphics.SolidColor{Value: color},
	}
}
