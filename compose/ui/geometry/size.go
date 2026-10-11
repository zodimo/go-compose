package geometry

import (
	"fmt"
	"math"

	"github.com/zodimo/go-compose/pkg/floatutils"
	"github.com/zodimo/go-compose/pkg/floatutils/lerp"
	"github.com/zodimo/go-sentinel-helper/sentinel/companions"
)

// Sentinel Values, Values are Stack allocated == FAST

// sizeUnspecified Represents an unspecified [Size] value, usually a replacement for `null` when a
// primitive value is desired. It is deliberately package-private: the only public
// way to obtain it is SizeCompanion.Unspecified().
var sizeUnspecified = NewSize(floatutils.Float32Unspecified, floatutils.Float32Unspecified)

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
	if isSpecifiedSize(s) {
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
	if s == sizeUnspecified && other == sizeUnspecified {
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

// IsSpecified returns false when this is Size.Unspecified. It delegates to
// SizeCompanion so the companion remains the single source of truth for the
// predicate; the sentinel itself stays unreachable from outside the package.
func (s Size) IsSpecified() bool {
	return SizeCompanion.IsSpecified(s)
}

func isSpecifiedSize(s Size) bool {
	return s != sizeUnspecified
}

func isUnspecifiedSize(s Size) bool {
	return s == sizeUnspecified
}

func takeOrElseSize(s, def Size) Size {
	if isUnspecifiedSize(s) {
		return def
	}
	return s
}

// mergeSize layers b onto a: a value from b wins only where b is specified.
func mergeSize(a, b Size) Size {
	if isUnspecifiedSize(b) {
		return a
	}
	return b
}

func stringSize(s Size) string {
	if isUnspecifiedSize(s) {
		return "Size.Unspecified"
	}
	return fmt.Sprintf("Size(%.1f, %.1f)", s.Width(), s.Height())
}

func sameSize(a, b Size) bool {
	return a == b
}

// semanticEqualSize treats two unspecified sizes as equal; comparing the packed
// NaN representation with == would otherwise report them as different.
func semanticEqualSize(a, b Size) bool {
	if isUnspecifiedSize(a) && isUnspecifiedSize(b) {
		return true
	}
	return a == b
}

func equalSize(a, b Size) bool {
	return a.Equal(b)
}

func copySize(s Size) Size {
	return s
}

// SizeCompanion is the single exported sentinel-family symbol for Size.
var SizeCompanion = companions.NewScalarCompanion(
	func() Size { return sizeUnspecified },
	stringSize,
	isSpecifiedSize,
	isUnspecifiedSize,
	takeOrElseSize,
	mergeSize,
	sameSize,
	semanticEqualSize,
	equalSize,
	copySize,
)
