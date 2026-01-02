package floatutils

import "math"

const (
	Float64EqualityThreshold float64 = 1e-9
	Float32EqualityThreshold float32 = 1e-6
)

var Float64Unspecified = math.NaN()
var Float32Unspecified = float32(math.NaN())

var Float32Infinite = float32(math.Inf(1))
var FloatInfinite = math.Inf(1)
