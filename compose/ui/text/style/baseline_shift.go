package style

import (
	"fmt"
	"math"

	"github.com/zodimo/go-compose/pkg/floatutils/lerp"
)

// BaselineShift represents the amount by which text is shifted up or down from the current baseline.
// The shift is calculated as: multiplier * (baseline - ascent)
// This is a 32-bit float value class equivalent to Kotlin's BaselineShift.
type BaselineShift float32

// Predefined baseline shift constants
const (
	// Superscript shifts text upward by 0.5 * (baseline - ascent)
	BaselineShiftSuperscript BaselineShift = 0.5
	// Subscript shifts text downward by -0.5 * (baseline - ascent)
	BaselineShiftSubscript BaselineShift = -0.5
	// None applies no baseline shift
	BaselineShiftNone BaselineShift = 0.0
)

// BaselineShiftUnspecified represents an unset baseline shift (NaN value)
var BaselineShiftUnspecified = BaselineShift(math.NaN())

// IsSpecified returns true if this baseline shift is not Unspecified (i.e., not NaN)
func (bs BaselineShift) IsSpecified() bool {
	return !math.IsNaN(float64(bs))
}

// TakeOrElse returns this BaselineShift if it's specified, otherwise returns the result of block.
// This is equivalent to Kotlin's inline function with lambda.
func (bs BaselineShift) TakeOrElse(block BaselineShift) BaselineShift {
	if bs.IsSpecified() {
		return bs
	}
	return block
}

// Multiplier returns the underlying multiplier value for use in calculations
func (bs BaselineShift) Multiplier() float32 {
	return float32(bs)
}

// NewBaselineShift creates a new BaselineShift with the given multiplier
func NewBaselineShift(multiplier float32) BaselineShift {
	return BaselineShift(multiplier)
}

// String returns a string representation of the BaselineShift
func (bs BaselineShift) String() string {
	switch {
	case !bs.IsSpecified():
		return "BaselineShift.Unspecified"
	case bs == BaselineShiftSuperscript:
		return "BaselineShift.Superscript"
	case bs == BaselineShiftSubscript:
		return "BaselineShift.Subscript"
	case bs == BaselineShiftNone:
		return "BaselineShift.None"
	default:
		return fmt.Sprintf("BaselineShift(multiplier=%v)", bs.Multiplier())
	}
}

// LerpBaselineShift linearly interpolates between two BaselineShifts.
// If either start or stop is Unspecified (NaN), the result will be Unspecified.
func LerpBaselineShift(start, stop BaselineShift, fraction float32) BaselineShift {
	// Kotlin implementation: return BaselineShift(lerp(start.multiplier, stop.multiplier, fraction))
	// If one is NaN, the arithmetic naturally produces NaN.
	return lerp.Between32(start, stop, fraction)
}
