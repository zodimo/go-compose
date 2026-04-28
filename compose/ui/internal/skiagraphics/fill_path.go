package skiagraphics

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-skia-support/skia/enums"
)

func SkiaFillPathTypeToGraphicsFillType(p enums.PathFillType) graphics.PathFillType {
	switch p {
	case enums.PathFillTypeWinding:
		return graphics.PathFillTypeNonZero
	case enums.PathFillTypeEvenOdd:
		return graphics.PathFillTypeEvenOdd
	default:
		panic(fmt.Sprintf("Unknown PathFillType: %s", p.String()))
	}
}

func GraphicsFillPathTypeToSkiaFillPathType(p graphics.PathFillType) enums.PathFillType {
	switch p {
	case graphics.PathFillTypeNonZero:
		return enums.PathFillTypeWinding
	case graphics.PathFillTypeEvenOdd:
		return enums.PathFillTypeEvenOdd
	default:
		panic(fmt.Sprintf("Unknown PathFillType: %s", p.String()))
	}
}
