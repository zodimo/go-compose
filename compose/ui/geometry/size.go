package geometry

import (
	"fmt"
	"math"

	"github.com/zodimo/go-compose/pkg/floatutils"
)

// Size represents a 2D floating-point size.
// You can think of this as an Offset from the origin.
type Size struct {
	Width  float32
	Height float32
}

// SizeZero is an empty size, one with a zero width and a zero height.
var SizeZero = Size{Width: 0, Height: 0}

// SizeUnspecified Represents an unspecified [Size] value, usually a replacement for `null` when a
// primitive value is desired.
var SizeUnspecified = Size{Width: float32(math.NaN()), Height: float32(math.NaN())}

// NewSize constructs a Size from the given width and height.
func NewSize(width, height float32) Size {
	return Size{Width: width, Height: height}
}

// IsEmpty returns true if the size encloses a non-zero area.
// Negative areas are considered empty.
// Note: IsUnspecified check is included to match Kotlin's behavior.
func (s Size) IsEmpty() bool {
	return s.IsUnspecified() || s.Width <= 0 || s.Height <= 0
}

// Times returns a Size whose dimensions are the dimensions of the left-hand-side operand (a Size)
// multiplied by the scalar right-hand-side operand (a Float).
func (s Size) Times(operand float32) Size {
	return Size{Width: s.Width * operand, Height: s.Height * operand}
}

// Div returns a Size whose dimensions are the dimensions of the left-hand-side operand (a Size)
// divided by the scalar right-hand-side operand (a Float).
func (s Size) Div(operand float32) Size {
	return Size{Width: s.Width / operand, Height: s.Height / operand}
}

// MinDimension returns the lesser of the magnitudes of the width and the height.
func (s Size) MinDimension() float32 {
	return float32(math.Min(math.Abs(float64(s.Width)), math.Abs(float64(s.Height))))
}

// MaxDimension returns the greater of the magnitudes of the width and the height.
func (s Size) MaxDimension() float32 {
	return float32(math.Max(math.Abs(float64(s.Width)), math.Abs(float64(s.Height))))
}

// String returns a string representation of the object.
func (s Size) String() string {
	if s.IsSpecified() {
		return fmt.Sprintf("Size(%.1f, %.1f)", s.Width, s.Height)
	}
	return "Size.Unspecified"
}

// IsSpecified returns false when this is Size.Unspecified.
func (s Size) IsSpecified() bool {
	return !math.IsNaN(float64(s.Width)) && !math.IsNaN(float64(s.Height))
}

// IsUnspecified returns true when this is Size.Unspecified.
// Note: In Kotlin it checks for a specific packed NaN value.
// Here we check if either component is NaN, or following the Offset pattern:
// Offset checks both. Size.kt says:
// packedValue == 0x7fc00000_7fc00000L // NaN_NaN
// So both must be NaN.
func (s Size) IsUnspecified() bool {
	return math.IsNaN(float64(s.Width)) && math.IsNaN(float64(s.Height))
}

// TakeOrElse returns this size if Specified, otherwise executes the block and returns its result.
func (s Size) TakeOrElse(block func() Size) Size {
	if s.IsSpecified() {
		return s
	}
	return block()
}

// Center returns the Offset of the center of the rect from the point of [0, 0] with this Size.
func (s Size) Center() Offset {
	return Offset{X: s.Width / 2, Y: s.Height / 2}
}

// Equal checks equality with another Size.
func (s Size) Equal(other Size) bool {
	if s.IsUnspecified() && other.IsUnspecified() {
		return true
	}
	return floatutils.Float32Equals(s.Width, other.Width, floatutils.Float32EqualityThreshold) &&
		floatutils.Float32Equals(s.Height, other.Height, floatutils.Float32EqualityThreshold)
}

// LerpSize linearly interpolates between two sizes.
func LerpSize(start, stop Size, fraction float32) Size {
	return Size{
		Width:  lerpBetween(start.Width, stop.Width, float64(fraction)),
		Height: lerpBetween(start.Height, stop.Height, float64(fraction)),
	}
}
