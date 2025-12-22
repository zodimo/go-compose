package floatutils

import "math"

func IsNaN(f float32) bool {
	return float64(f) == math.NaN()
}
