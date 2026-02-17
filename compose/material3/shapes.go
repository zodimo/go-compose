package material3

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/material3/tokens"
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
)

type Shapes struct {
	ExtraExtraLarge     shape.Shape
	ExtraLarge          shape.Shape
	ExtraLargeIncreased shape.Shape
	ExtraLargeTop       shape.Shape
	ExtraSmall          shape.Shape
	ExtraSmallTop       shape.Shape
	Full                shape.Shape
	Large               shape.Shape
	LargeEnd            shape.Shape
	LargeIncreased      shape.Shape
	LargeStart          shape.Shape
	LargeTop            shape.Shape
	Medium              shape.Shape
	None                shape.Shape
	Small               shape.Shape
}

var DefaultShapes = &Shapes{
	ExtraExtraLarge:     tokens.ShapeTokens.CornerExtraExtraLarge,
	ExtraLarge:          tokens.ShapeTokens.CornerExtraLarge,
	ExtraLargeIncreased: tokens.ShapeTokens.CornerExtraLargeIncreased,
	ExtraLargeTop:       tokens.ShapeTokens.CornerExtraLargeTop,
	ExtraSmall:          tokens.ShapeTokens.CornerExtraSmall,
	ExtraSmallTop:       tokens.ShapeTokens.CornerExtraSmallTop,
	Full:                tokens.ShapeTokens.CornerFull,
	Large:               tokens.ShapeTokens.CornerLarge,
	LargeEnd:            tokens.ShapeTokens.CornerLargeEnd,
	LargeIncreased:      tokens.ShapeTokens.CornerLargeIncreased,
	LargeStart:          tokens.ShapeTokens.CornerLargeStart,
	LargeTop:            tokens.ShapeTokens.CornerLargeTop,
	Medium:              tokens.ShapeTokens.CornerMedium,
	None:                tokens.ShapeTokens.CornerNone,
	Small:               tokens.ShapeTokens.CornerSmall,
}

func (c Shapes) FromToken(value tokens.ShapeTokenKey) shape.Shape {
	switch value {
	case tokens.ShapeTokenKeyCornerExtraExtraLarge:
		return c.ExtraExtraLarge
	case tokens.ShapeTokenKeyCornerExtraLarge:
		return c.ExtraLarge
	case tokens.ShapeTokenKeyCornerExtraLargeIncreased:
		return c.ExtraLargeIncreased
	case tokens.ShapeTokenKeyCornerExtraLargeTop:
		return c.ExtraLargeTop
	case tokens.ShapeTokenKeyCornerExtraSmall:
		return c.ExtraSmall
	case tokens.ShapeTokenKeyCornerExtraSmallTop:
		return c.ExtraSmallTop
	case tokens.ShapeTokenKeyCornerFull:
		return c.Full
	case tokens.ShapeTokenKeyCornerLarge:
		return c.Large
	case tokens.ShapeTokenKeyCornerLargeEnd:
		return c.LargeEnd
	case tokens.ShapeTokenKeyCornerLargeIncreased:
		return c.LargeIncreased
	case tokens.ShapeTokenKeyCornerLargeStart:
		return c.LargeStart
	case tokens.ShapeTokenKeyCornerLargeTop:
		return c.LargeTop
	case tokens.ShapeTokenKeyCornerMedium:
		return c.Medium
	case tokens.ShapeTokenKeyCornerNone:
		return c.None
	case tokens.ShapeTokenKeyCornerSmall:
		return c.Small
	default:
		panic(fmt.Sprintf("unknown shape token key: %s", value.String()))
	}
}
