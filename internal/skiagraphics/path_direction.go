package skiagraphics

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-skia-support/skia/enums"
)

func SkiaPathDirectionToGraphicsPathDirection(d enums.PathDirection) graphics.PathDirection {
	switch d {
	case enums.PathDirectionCCW:
		return graphics.PathDirectionCounterClockwise
	case enums.PathDirectionCW:
		return graphics.PathDirectionClockwise
	default:
		panic(fmt.Sprintf("Unknown PathDirection: %s", d.String()))
	}
}

func GraphicsPathDirectionToSkiaPathDirection(d graphics.PathDirection) enums.PathDirection {
	switch d {
	case graphics.PathDirectionCounterClockwise:
		return enums.PathDirectionCCW
	case graphics.PathDirectionClockwise:
		return enums.PathDirectionCW
	default:
		panic(fmt.Sprintf("Unknown PathDirection: %s", d.String()))
	}
}
