package graphics

import (
	"github.com/zodimo/go-compose/compose/ui/geometry"
	"github.com/zodimo/go-compose/pkg/floatutils"
	"github.com/zodimo/go-compose/theme"
)

type Color = theme.ColorDescriptor
type Offset = geometry.Offset

var ZeroOffset = geometry.OffsetZero()

var floatEquals = func(a, b float32) bool {
	return floatutils.Float32Equals(a, b, floatutils.Float32EqualityThreshold)
}
