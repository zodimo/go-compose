package style

import (
	"fmt"

	"github.com/zodimo/go-compose/pkg/floatutils"
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
var BaselineShiftUnspecified = BaselineShift(floatutils.Float32Unspecified)

// IsBaselineShift returns true if this baseline shift is not Unspecified (i.e., not NaN)
func IsBaselineShift(bs BaselineShift) bool {
	return floatutils.IsSpecified(bs)
}

// TakeOrElseBaselineShift returns this BaselineShift if it's specified, otherwise returns the result of block.
// This is equivalent to Kotlin's inline function with lambda.
func TakeOrElseBaselineShift(bs, block BaselineShift) BaselineShift {
	if IsBaselineShift(bs) {
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

// MergeBaselineShift returns the first argument if it is specified, otherwise the second.
func MergeBaselineShift(a, b BaselineShift) BaselineShift {
	if IsBaselineShift(a) {
		return a
	}
	return b
}

// StringBaselineShift returns a string representation of the BaselineShift
func StringBaselineShift(bs BaselineShift) string {
	switch {
	case !IsBaselineShift(bs):
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

// String returns a string representation of the BaselineShift
func (bs BaselineShift) String() string {
	return StringBaselineShift(bs)
}

// LerpBaselineShift linearly interpolates between two BaselineShifts.
// If either start or stop is Unspecified (NaN), the result will be Unspecified.
func LerpBaselineShift(start, stop BaselineShift, fraction float32) BaselineShift {
	// Kotlin implementation: return BaselineShift(lerp(start.multiplier, stop.multiplier, fraction))
	// If one is NaN, the arithmetic naturally produces NaN.
	return lerp.Between32(start, stop, fraction)
}
