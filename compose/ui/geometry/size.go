package geometry

import (
	"fmt"
	"math"

	"github.com/zodimo/go-compose/pkg/floatutils"
	"github.com/zodimo/go-compose/pkg/floatutils/lerp"
)

// Sentinel Values, Values are Stack allocated == FAST

// SizeUnspecified Represents an unspecified [Size] value, usually a replacement for `null` when a
// primitive value is desired.
var SizeUnspecified = NewSize(floatutils.Float32Unspecified, floatutils.Float32Unspecified)

// Size represents a 2D floating-point size.
// You can think of this as an Offset from the origin.
// It is a packed value where the width is in the high 32 bits and the height is in the low 32 bits.
type Size int64

// SizeZero is an empty size, one with a zero width and a zero height.
var SizeZero = NewSize(0, 0)

// NewSize constructs a Size from the given width and height.
func NewSize(width, height float32) Size {
	return Size(floatutils.PackFloats(width, height))
}

// Width returns the width component of the size.
func (s Size) Width() float32 {
	return floatutils.UnpackFloat1(int64(s))
}

// Height returns the height component of the size.
func (s Size) Height() float32 {
	return floatutils.UnpackFloat2(int64(s))
}

// IsEmpty returns true if the Size is empty (width <= 0 or height <= 0).
func (s Size) IsEmpty() bool {
	return s.Width() <= 0 || s.Height() <= 0
}

// Times returns a Size whose dimensions are the dimensions of the left-hand-side operand (a Size)
// multiplied by the scalar right-hand-side operand (a Float).
func (s Size) Times(operand float32) Size {
	return NewSize(s.Width()*operand, s.Height()*operand)
}

// Div returns a Size whose dimensions are the dimensions of the left-hand-side operand (a Size)
// divided by the scalar right-hand-side operand (a Float).
func (s Size) Div(operand float32) Size {
	return NewSize(s.Width()/operand, s.Height()/operand)
}

// MinDimension returns the lesser of the magnitudes of the width and the height.
func (s Size) MinDimension() float32 {
	return float32(math.Min(math.Abs(float64(s.Width())), math.Abs(float64(s.Height()))))
}

// MaxDimension returns the greater of the magnitudes of the width and the height.
func (s Size) MaxDimension() float32 {
	return float32(math.Max(math.Abs(float64(s.Width())), math.Abs(float64(s.Height()))))
}

// String returns a string representation of the object.
func (s Size) String() string {
	if s.IsSpecified() {
		return fmt.Sprintf("Size(%.1f, %.1f)", s.Width(), s.Height())
	}
	return "Size.Unspecified"
}

// Center returns the Offset of the center of the rect from the point of [0, 0] with this Size.
func (s Size) Center() Offset {
	return NewOffset(s.Width()/2, s.Height()/2)
}

// Equal checks equality with another Size.
func (s Size) Equal(other Size) bool {
	if s == SizeUnspecified && other == SizeUnspecified {
		return true
	}
	return floatutils.Float32Equals(s.Width(), other.Width(), floatutils.Float32EqualityThreshold) &&
		floatutils.Float32Equals(s.Height(), other.Height(), floatutils.Float32EqualityThreshold)
}

// LerpSize linearly interpolates between two sizes.
func LerpSize(start, stop Size, fraction float32) Size {
	return NewSize(
		lerp.Between32(start.Width(), stop.Width(), fraction),
		lerp.Between32(start.Height(), stop.Height(), fraction),
	)
}

// IsSpecified returns false when this is Size.Unspecified.
func (s Size) IsSpecified() bool {
	return s != SizeUnspecified
}

// IsUnspecified returns true when this is Size.Unspecified.
func (s Size) IsUnspecified() bool {
	return s == SizeUnspecified
}

//////////////// Sentinel HELPERS //////////////////////////

func (s Size) TakeOrElse(block Size) Size {
	if s == SizeUnspecified {
		return block
	}
	return s
}
