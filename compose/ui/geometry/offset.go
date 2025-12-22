package geometry

import (
	"fmt"
	"math"

	"github.com/zodimo/go-compose/pkg/floatutils"
)

// Offset is an immutable 2D floating-point offset.
type Offset struct {
	X float32
	Y float32
}

// OffsetZero is an offset with zero magnitude.
var OffsetZero = Offset{X: 0, Y: 0}

// OffsetInfinite is an offset with infinite x and y components.
var OffsetInfinite = Offset{X: float32(math.Inf(1)), Y: float32(math.Inf(1))}

// OffsetUnspecified Represents an unspecified [Offset] value, usually a replacement for `null` when a
// primitive value is desired.
var OffsetUnspecified = Offset{X: float32(math.NaN()), Y: float32(math.NaN())}

// NewOffset constructs an Offset from the given x and y values.
func NewOffset(x, y float32) Offset {
	return Offset{X: x, Y: y}
}

// Copy returns a copy of this Offset instance optionally overriding the x or y parameter.
// Since Go doesn't support default arguments effectively for this, we provide a helper
// that takes pointers or just use struct literal if needed.
// However, to match Kotlin's copy, we can just return a new struct.
// For idiomatic Go, usually we just create a new struct.
// We will provide a helper that accepts optional values? No, Keep it simple.
// Users can just do geometry.Offset{X: o.X, Y: newY}.
// But we'll provide a functional one for chaining if needed, but for now NewOffset is enough.

// IsValid returns true if x and y are not NaN and not Infinite.
// Wait, Kotlin's isValid() checks if it's NOT NaN.
// Kotlin:
// - False if x or y is NaN
// - True if x or y is infinite
// - True otherwise
// Wait, Kotlin source says:
// isValid():
// "Take the unsigned packed floats and see if they are > InfinityBase (any NaN)"
// So it returns true if it is NOT NaN. Infinite IS valid.
func (o Offset) IsValid() bool {
	return !math.IsNaN(float64(o.X)) && !math.IsNaN(float64(o.Y))
}

// GetDistance returns the magnitude of the offset.
func (o Offset) GetDistance() float32 {
	return float32(math.Sqrt(float64(o.X*o.X + o.Y*o.Y)))
}

// GetDistanceSquared returns the square of the magnitude of the offset.
func (o Offset) GetDistanceSquared() float32 {
	return o.X*o.X + o.Y*o.Y
}

// UnaryMinus returns an offset with the coordinates negated.
func (o Offset) UnaryMinus() Offset {
	return Offset{X: -o.X, Y: -o.Y}
}

// Minus returns an offset whose x value is the left-hand-side operand's x minus the
// right-hand-side operand's x and whose y value is the left-hand-side operand's y minus
// the right-hand-side operand's y.
func (o Offset) Minus(other Offset) Offset {
	return Offset{X: o.X - other.X, Y: o.Y - other.Y}
}

// Plus returns an offset whose x value is the sum of the x values of the two operands, and whose
// y value is the sum of the y values of the two operands.
func (o Offset) Plus(other Offset) Offset {
	return Offset{X: o.X + other.X, Y: o.Y + other.Y}
}

// Times returns an offset whose coordinates are the coordinates of the left-hand-side operand (an
// Offset) multiplied by the scalar right-hand-side operand (a Float).
func (o Offset) Times(operand float32) Offset {
	return Offset{X: o.X * operand, Y: o.Y * operand}
}

// Div returns an offset whose coordinates are the coordinates of the left-hand-side operand (an
// Offset) divided by the scalar right-hand-side operand (a Float).
func (o Offset) Div(operand float32) Offset {
	return Offset{X: o.X / operand, Y: o.Y / operand}
}

// Rem returns an offset whose coordinates are the remainder of dividing the coordinates of the
// left-hand-side operand (an Offset) by the scalar right-hand-side operand (a Float).
func (o Offset) Rem(operand float32) Offset {
	return Offset{
		X: float32(math.Mod(float64(o.X), float64(operand))),
		Y: float32(math.Mod(float64(o.Y), float64(operand))),
	}
}

// String returns a string representation of the object.
func (o Offset) String() string {
	if o.IsSpecified() {
		return fmt.Sprintf("Offset(%.1f, %.1f)", o.X, o.Y)
	}
	return "Offset.Unspecified"
}

// LerpOffset linearly interpolates between two offsets.
func LerpOffset(start, stop Offset, fraction float32) Offset {
	return Offset{
		X: lerpBetween(start.X, stop.X, float64(fraction)),
		Y: lerpBetween(start.Y, stop.Y, float64(fraction)),
	}
}

// IsFinite returns true if both x and y values of the Offset are finite.
// NaN values are not considered finite.
func (o Offset) IsFinite() bool {
	return !math.IsInf(float64(o.X), 0) && !math.IsNaN(float64(o.X)) &&
		!math.IsInf(float64(o.Y), 0) && !math.IsNaN(float64(o.Y))
}

// IsSpecified returns true if this is not Offset.Unspecified.
func (o Offset) IsSpecified() bool {
	return !math.IsNaN(float64(o.X)) || !math.IsNaN(float64(o.Y))
}

// IsUnspecified returns true if this is Offset.Unspecified.
func (o Offset) IsUnspecified() bool {
	return math.IsNaN(float64(o.X)) && math.IsNaN(float64(o.Y))
}

// TakeOrElse returns this offset if Specified, otherwise executes the block and returns its result.
func (o Offset) TakeOrElse(block func() Offset) Offset {
	if o.IsSpecified() {
		return o
	}
	return block()
}

// Equal checks equality with another Offset.
func (o Offset) Equal(other Offset) bool {
	if o.IsUnspecified() && other.IsUnspecified() {
		return true
	}
	return floatutils.Float32Equals(o.X, other.X, floatutils.Float32EqualityThreshold) &&
		floatutils.Float32Equals(o.Y, other.Y, floatutils.Float32EqualityThreshold)
}
