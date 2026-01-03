package floatutils

import "math"

// floatEquals compares two float32 values with absolute epsilon tolerance.
func Float32Equals(a, b, epsilon float32) bool {
	return math.Abs(float64(a-b)) <= float64(epsilon)
}

// floatEquals compares two float32 values with absolute epsilon tolerance.
func Float64Equals(a, b, epsilon float64) bool {
	return math.Abs(a-b) <= epsilon
}

func IsInfinite[T Float](f T) bool {
	if !IsSpecified(f) {
		return false
	}
	return math.IsInf(float64(f), 0)
}

func IsSpecified[T Float](f T) bool {
	return !math.IsNaN(float64(f))
}

func IsUnspecified[T Float](f T) bool {
	return math.IsNaN(float64(f))
}

// Deprecated
func IsNaN[T Float](f T) bool {
	return math.IsNaN(float64(f))
}
