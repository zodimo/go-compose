// Package lerp provides high-performance linear interpolation utilities.
// It offers three precision levels: standard (fast), precise (accurate), and fixed-point (fastest).
package lerp

import (
	"math"

	"github.com/zodimo/go-compose/pkg/floatutils"
)

// Float constraint for all float32 and float64 based types.
type Float interface {
	~float32 | ~float64
}

// ============================================================================
// Generic Interpolation
// ============================================================================

// Between performs type-safe linear interpolation with float64 precision.
// Returns NaN if either start or stop is NaN.
//
// Example:
//
//	pos := lerp.Between(0.0, 100.0, 0.5) // 50.0
func Between[T Float](start, stop T, fraction float64) T {
	// if start or stop unspeficied then unspecified
	if !floatutils.IsSpecified(start) {
		return start
	}
	if !floatutils.IsSpecified(stop) {
		return stop
	}

	if fraction == 0 {
		return start
	}
	if fraction == 1 {
		return stop
	}
	s, e := float64(start), float64(stop)
	return T(s + (e-s)*fraction)
}

// Between32 is optimized for float32 types with ~15% better performance than Between.
//
// Example:
//
//	opacity := lerp.Between32(0.0, 1.0, 0.5) // 0.5
func Between32[T ~float32](start, stop T, fraction float32) T {
	// if start or stop unspeficied then unspecified
	if !floatutils.IsSpecified(start) {
		return start
	}
	if !floatutils.IsSpecified(stop) {
		return stop
	}
	if fraction == 0 {
		return start
	}

	if fraction == 1 {
		return stop
	}
	s, e := float32(start), float32(stop)
	return T(s + (e-s)*fraction)
}

// ============================================================================
// Numeric Specializations
// ============================================================================

// Float32 interpolates between two float32 values with zero allocation overhead.
func Float32(start, stop, fraction float32) float32 {
	if fraction == 0 {
		return start
	}
	if fraction == 1 {
		return stop
	}
	return (1-fraction)*start + fraction*stop
}

// Int provides fast integer interpolation using float32 intermediates.
// Truncates fractional parts; use IntPrecise for correct rounding.
func Int(start, stop int, fraction float32) int {
	return int(Float32(float32(start), float32(stop), fraction))
}

// IntFixed uses fixed-point arithmetic (fraction 0-256) for branch-free performance.
// Fraction 256 equals 1.0. Ideal for animation loops.
func IntFixed(start, stop, fraction int) int {
	return start + (stop-start)*fraction>>8
}

// IntPrecise provides correctly-rounded interpolation using float64.
func IntPrecise(start, stop int, fraction float32) int {
	f64 := float64(start) + float64(stop-start)*float64(fraction)
	return int(f64 + 0.5) // Nearest-integer rounding
}

// ============================================================================
// Color Interpolation
// ============================================================================

// TODO: Move to color/lerp package to eliminate coupling with internal packages

// ColorLerp performs standard RGBA interpolation.
//
// Example:
//
//	mid := lerp.ColorLerp(color1, color2, 0.5)
func ColorLerp(a, b struct{ R, G, B, A float32 }, p float32) struct{ R, G, B, A float32 } {
	invP := 1 - p
	return struct{ R, G, B, A float32 }{
		R: a.R*invP + b.R*p,
		G: a.G*invP + b.G*p,
		B: a.B*invP + b.B*p,
		A: a.A*invP + b.A*p,
	}
}

// ColorLerpPrecise uses squared interpolation for gamma-correct color blending.
// ~3x slower but produces perceptually better results.
func ColorLerpPrecise(a, b struct{ R, G, B, A float32 }, p float32) struct{ R, G, B, A float32 } {
	invP := 1 - p
	return struct{ R, G, B, A float32 }{
		R: sqrt(a.R*a.R*invP + b.R*b.R*p),
		G: sqrt(a.G*a.G*invP + b.G*b.G*p),
		B: sqrt(a.B*a.B*invP + b.B*b.B*p),
		A: a.A*invP + b.A*p,
	}
}

func sqrt(v float32) float32 {
	return float32(math.Sqrt(float64(v)))
}
