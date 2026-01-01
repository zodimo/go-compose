package material3

import (
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

type Shapes struct {
	ExtraSmall shape.Shape
	Small      shape.Shape
	Medium     shape.Shape
	Large      shape.Shape
	ExtraLarge shape.Shape
}

// Default shapes for Material 3
var (
	ShapeExtraSmall = shape.RoundedCornerShape{Radius: unit.Dp(4)}
	ShapeSmall      = shape.RoundedCornerShape{Radius: unit.Dp(8)}
	ShapeMedium     = shape.RoundedCornerShape{Radius: unit.Dp(12)}
	ShapeLarge      = shape.RoundedCornerShape{Radius: unit.Dp(16)}
	ShapeExtraLarge = shape.RoundedCornerShape{Radius: unit.Dp(28)}
)

func DefaultShapes() *Shapes {
	return &Shapes{
		ExtraSmall: ShapeExtraSmall,
		Small:      ShapeSmall,
		Medium:     ShapeMedium,
		Large:      ShapeLarge,
		ExtraLarge: ShapeExtraLarge,
	}
}
