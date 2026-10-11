package geometry

import (
	"fmt"

	"github.com/zodimo/go-compose/pkg/floatutils"
	"github.com/zodimo/go-compose/pkg/floatutils/lerp"
	"github.com/zodimo/go-sentinel-helper/sentinel/companions"
)

// CornerRadius represents the radii for corners of a rounded rectangle.
// It is a packed value where the x radius is in the high 32 bits and the y radius is in the low 32 bits.
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-graphics/src/commonMain/kotlin/androidx/compose/ui/graphics/CornerRadius.kt
type CornerRadius int64

// CornerRadiusZero is a CornerRadius with both x and y radii equal to zero (sharp corners).
var CornerRadiusZero = NewCornerRadius(0, 0)

// cornerRadiusUnspecified represents an unspecified corner radius. It is
// deliberately package-private: the only public way to obtain it is
// CornerRadiusCompanion.Unspecified().
var cornerRadiusUnspecified = NewCornerRadius(floatutils.Float32Unspecified, floatutils.Float32Unspecified)

// NewCornerRadius creates a CornerRadius from the given x and y radii.
// x is the radius of the corners along the x-axis.
// y is the radius of the corners along the y-axis (defaults to x if you want circular corners).
func NewCornerRadius(x, y float32) CornerRadius {
	return CornerRadius(floatutils.PackFloats(x, y))
}

// NewCircularCornerRadius creates a CornerRadius with equal x and y radii (circular corners).
func NewCircularCornerRadius(radius float32) CornerRadius {
	return NewCornerRadius(radius, radius)
}

// X returns the x component of the corner radius.
func (c CornerRadius) X() float32 {
	return floatutils.UnpackFloat1(int64(c))
}

// Y returns the y component of the corner radius.
func (c CornerRadius) Y() float32 {
	return floatutils.UnpackFloat2(int64(c))
}

// String returns a string representation of the CornerRadius.
func (c CornerRadius) String() string {
	if c == CornerRadiusZero {
		return "CornerRadius.Zero"
	}
	if c == cornerRadiusUnspecified {
		return "CornerRadius.Unspecified"
	}
	return fmt.Sprintf("CornerRadius(%.1f, %.1f)", c.X(), c.Y())
}

// Equal checks equality with another CornerRadius.
func (c CornerRadius) Equal(other CornerRadius) bool {
	if c == other {
		return true
	}
	return floatutils.Float32Equals(c.X(), other.X(), floatutils.Float32EqualityThreshold) &&
		floatutils.Float32Equals(c.Y(), other.Y(), floatutils.Float32EqualityThreshold)
}

// IsSpecified returns true if this is not CornerRadius.Unspecified.
func (c CornerRadius) IsSpecified() bool {
	return isSpecifiedCornerRadius(c)
}

func isSpecifiedCornerRadius(c CornerRadius) bool {
	return c != cornerRadiusUnspecified
}

func isUnspecifiedCornerRadius(c CornerRadius) bool {
	return c == cornerRadiusUnspecified
}

func takeOrElseCornerRadius(c, def CornerRadius) CornerRadius {
	if isUnspecifiedCornerRadius(c) {
		return def
	}
	return c
}

// mergeCornerRadius layers b onto a: a value from b wins only where b is specified.
func mergeCornerRadius(a, b CornerRadius) CornerRadius {
	if isUnspecifiedCornerRadius(b) {
		return a
	}
	return b
}

func stringCornerRadius(c CornerRadius) string {
	if c == CornerRadiusZero {
		return "CornerRadius.Zero"
	}
	if isUnspecifiedCornerRadius(c) {
		return "CornerRadius.Unspecified"
	}
	return fmt.Sprintf("CornerRadius(%.1f, %.1f)", c.X(), c.Y())
}

func sameCornerRadius(a, b CornerRadius) bool {
	return a == b
}

// semanticEqualCornerRadius treats two unspecified radii as equal; comparing the
// packed NaN representation with == would otherwise report them as different.
func semanticEqualCornerRadius(a, b CornerRadius) bool {
	if isUnspecifiedCornerRadius(a) && isUnspecifiedCornerRadius(b) {
		return true
	}
	return a == b
}

func equalCornerRadius(a, b CornerRadius) bool {
	return a.Equal(b)
}

func copyCornerRadius(c CornerRadius) CornerRadius {
	return c
}

// CornerRadiusCompanion is the single exported sentinel-family symbol for CornerRadius.
var CornerRadiusCompanion = companions.NewScalarCompanion(
	func() CornerRadius { return cornerRadiusUnspecified },
	stringCornerRadius,
	isSpecifiedCornerRadius,
	isUnspecifiedCornerRadius,
	takeOrElseCornerRadius,
	mergeCornerRadius,
	sameCornerRadius,
	semanticEqualCornerRadius,
	equalCornerRadius,
	copyCornerRadius,
)

// Plus returns a CornerRadius with radii that are the sum of both corner radii.
func (c CornerRadius) Plus(other CornerRadius) CornerRadius {
	return NewCornerRadius(c.X()+other.X(), c.Y()+other.Y())
}

// Minus returns a CornerRadius with radii that are the difference of both corner radii.
func (c CornerRadius) Minus(other CornerRadius) CornerRadius {
	return NewCornerRadius(c.X()-other.X(), c.Y()-other.Y())
}

// Times returns a CornerRadius with radii scaled by the given factor.
func (c CornerRadius) Times(operand float32) CornerRadius {
	return NewCornerRadius(c.X()*operand, c.Y()*operand)
}

// Div returns a CornerRadius with radii divided by the given factor.
func (c CornerRadius) Div(operand float32) CornerRadius {
	return NewCornerRadius(c.X()/operand, c.Y()/operand)
}

// LerpCornerRadius linearly interpolates between two corner radii.
func LerpCornerRadius(start, stop CornerRadius, fraction float32) CornerRadius {
	return NewCornerRadius(
		lerp.Between32(start.X(), stop.X(), fraction),
		lerp.Between32(start.Y(), stop.Y(), fraction),
	)
}
