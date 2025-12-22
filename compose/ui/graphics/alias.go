package graphics

import (
	"github.com/zodimo/go-compose/compose/ui/geometry"
	"github.com/zodimo/go-compose/compose/ui/utils/lerp"
	"github.com/zodimo/go-compose/pkg/floatutils"
)

type Offset = geometry.Offset

var ZeroOffset = geometry.OffsetZero

var lerpBetween = lerp.Between[float32]
var float32Equals = floatutils.Float32Equals
var float32EqualityThreshold = floatutils.Float32EqualityThreshold
