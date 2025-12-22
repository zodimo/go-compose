package unit

import (
	"github.com/zodimo/go-compose/pkg/floatutils"
	"github.com/zodimo/go-compose/pkg/floatutils/lerp"
)

var lerpBetween = lerp.Between[float32]
var float32Equals = floatutils.Float32Equals
var float32EqualityThreshold = floatutils.Float32EqualityThreshold
