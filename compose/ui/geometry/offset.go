package geometry

import (
	"fmt"
	"math"

	"github.com/zodimo/go-compose/pkg/floatutils"
	"github.com/zodimo/go-compose/pkg/floatutils/lerp"
	"github.com/zodimo/go-sentinel-helper/sentinel/companions"
)

// offsetUnspecified Represents an unspecified [Offset] value, usually a replacement for `null` when a
// primitive value is desired. It is deliberately package-private: the only public
// way to obtain it is OffsetCompanion.Unspecified().
var offsetUnspecified = NewOffset(floatutils.Float32Unspecified, floatutils.Float32Unspecified)

// Offset is an immutable 2D floating-point offset.
// It is a packed value where the x coordinate is in the high 32 bits and the y coordinate is in the low 32 bits.
type Offset int64

// OffsetZero is an offset with zero magnitude.
var OffsetZero = NewOffset(0, 0)

// OffsetInfinite is an offset with infinite x and y components.
var OffsetInfinite = NewOffset(floatutils.Float32Infinite, floatutils.Float32Infinite)

// NewOffset constructs an Offset from the given x and y values.
func NewOffset(x, y float32) Offset {
	return Offset(floatutils.PackFloats(x, y))
}

// X returns the x component of the offset.
func (o Offset) X() float32 {
	return floatutils.UnpackFloat1(int64(o))
}

// Y returns the y component of the offset.
func (o Offset) Y() float32 {
	return floatutils.UnpackFloat2(int64(o))
}

// IsValid returns true if x and y are not NaN.
func (o Offset) IsValid() bool {
	return !math.IsNaN(float64(o.X())) && !math.IsNaN(float64(o.Y()))
}

// GetDistance returns the magnitude of the offset.
func (o Offset) GetDistance() float32 {
	x, y := o.X(), o.Y()
	return float32(math.Sqrt(float64(x*x + y*y)))
}

// GetDistanceSquared returns the square of the magnitude of the offset.
func (o Offset) GetDistanceSquared() float32 {
	x, y := o.X(), o.Y()
	return x*x + y*y
}

// UnaryMinus returns an offset with the coordinates negated.
func (o Offset) UnaryMinus() Offset {
	return NewOffset(-o.X(), -o.Y())
}

// Minus returns an offset whose x value is the left-hand-side operand's x minus the
// right-hand-side operand's x and whose y value is the left-hand-side operand's y minus
// the right-hand-side operand's y.
func (o Offset) Minus(other Offset) Offset {
	return NewOffset(o.X()-other.X(), o.Y()-other.Y())
}

// Plus returns an offset whose x value is the sum of the x values of the two operands, and whose
// y value is the sum of the y values of the two operands.
func (o Offset) Plus(other Offset) Offset {
	return NewOffset(o.X()+other.X(), o.Y()+other.Y())
}

// Times returns an offset whose coordinates are the coordinates of the left-hand-side operand (an
// Offset) multiplied by the scalar right-hand-side operand (a Float).
func (o Offset) Times(operand float32) Offset {
	return NewOffset(o.X()*operand, o.Y()*operand)
}

// Div returns an offset whose coordinates are the coordinates of the left-hand-side operand (an
// Offset) divided by the scalar right-hand-side operand (a Float).
func (o Offset) Div(operand float32) Offset {
	return NewOffset(o.X()/operand, o.Y()/operand)
}

// Rem returns an offset whose coordinates are the remainder of dividing the coordinates of the
// left-hand-side operand (an Offset) by the scalar right-hand-side operand (a Float).
func (o Offset) Rem(operand float32) Offset {
	return NewOffset(
		float32(math.Mod(float64(o.X()), float64(operand))),
		float32(math.Mod(float64(o.Y()), float64(operand))),
	)
}

// String returns a string representation of the object.
func (o Offset) String() string {
	if isSpecifiedOffset(o) {
		return fmt.Sprintf("Offset(%.1f, %.1f)", o.X(), o.Y())
	}
	return "Offset.Unspecified"
}

// LerpOffset linearly interpolates between two offsets.
func LerpOffset(start, stop Offset, fraction float32) Offset {
	return NewOffset(
		lerp.Between32(start.X(), stop.X(), fraction),
		lerp.Between32(start.Y(), stop.Y(), fraction),
	)
}

// Equal checks equality with another Offset.
func (o Offset) Equal(other Offset) bool {
	if o == other {
		return true
	}
	if o == offsetUnspecified && other == offsetUnspecified {
		return true
	}
	return floatutils.Float32Equals(o.X(), other.X(), floatutils.Float32EqualityThreshold) &&
		floatutils.Float32Equals(o.Y(), other.Y(), floatutils.Float32EqualityThreshold)
}

// IsFinite returns true if both x and y values of the Offset are finite.
// NaN values are not considered finite.
func (o Offset) IsFinite() bool {
	x, y := o.X(), o.Y()
	return !math.IsInf(float64(x), 0) && !math.IsNaN(float64(x)) &&
		!math.IsInf(float64(y), 0) && !math.IsNaN(float64(y))
}

// IsInfinite returns true if either x or y value of the Offset is infinite.
func (o Offset) IsInfinite() bool {
	x, y := o.X(), o.Y()
	return math.IsInf(float64(x), 0) || math.IsInf(float64(y), 0)
}

// === Sentinel operations ===
//
// These follow the go-sentinel-helper contract: sentinel logic lives in the
// package-private ops below and is surfaced only through OffsetCompanion, so
// there is a single public source of truth for "unspecified".

func isSpecifiedOffset(o Offset) bool {
	return o != offsetUnspecified
}

func isUnspecifiedOffset(o Offset) bool {
	return o == offsetUnspecified
}

func takeOrElseOffset(o, def Offset) Offset {
	if isUnspecifiedOffset(o) {
		return def
	}
	return o
}

// mergeOffset layers b onto a: a value from b wins only where b is specified.
func mergeOffset(a, b Offset) Offset {
	if isUnspecifiedOffset(b) {
		return a
	}
	return b
}

func stringOffset(o Offset) string {
	if isUnspecifiedOffset(o) {
		return "Offset.Unspecified"
	}
	return fmt.Sprintf("Offset(%.1f, %.1f)", o.X(), o.Y())
}

func sameOffset(a, b Offset) bool {
	return a == b
}

// semanticEqualOffset treats two unspecified offsets as equal; comparing the
// packed NaN representation with == would otherwise report them as different.
func semanticEqualOffset(a, b Offset) bool {
	if isUnspecifiedOffset(a) && isUnspecifiedOffset(b) {
		return true
	}
	return a == b
}

func equalOffset(a, b Offset) bool {
	return a.Equal(b)
}

func copyOffset(o Offset) Offset {
	return o
}

// OffsetCompanion is the single exported sentinel-family symbol for Offset.
var OffsetCompanion = companions.NewScalarCompanion(
	func() Offset { return offsetUnspecified },
	stringOffset,
	isSpecifiedOffset,
	isUnspecifiedOffset,
	takeOrElseOffset,
	mergeOffset,
	sameOffset,
	semanticEqualOffset,
	equalOffset,
	copyOffset,
)
