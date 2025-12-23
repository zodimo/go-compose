package floatutils

import "math"

type Float interface {
	~float32 | ~float64
}

var FloatUnspecified = math.NaN()

func IsSpecified[T Float](f T) bool {
	//same as !math.IsNaN()
	return float64(f) != FloatUnspecified
}
