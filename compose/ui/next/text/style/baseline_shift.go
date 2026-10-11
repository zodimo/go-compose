package style

import (
	"fmt"

	"github.com/zodimo/go-compose/pkg/floatutils"
	"github.com/zodimo/go-compose/pkg/floatutils/lerp"
	"github.com/zodimo/go-sentinel-helper/sentinel/companions"
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

// baselineShiftUnspecified represents an unset baseline shift (NaN value). It is
// deliberately package-private: the only public way to obtain it is
// BaselineShiftCompanion.Unspecified().
var baselineShiftUnspecified = BaselineShift(floatutils.Float32Unspecified)

// IsBaselineShift returns true if this baseline shift is not Unspecified (i.e., not NaN)
func IsBaselineShift(bs BaselineShift) bool {
	return isSpecifiedBaselineShift(bs)
}

// TakeOrElseBaselineShift returns this BaselineShift if it's specified, otherwise returns the result of block.
// This is equivalent to Kotlin's inline function with lambda.
func TakeOrElseBaselineShift(bs, block BaselineShift) BaselineShift {
	return takeOrElseBaselineShift(bs, block)
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

// === Sentinel operations ===

func isSpecifiedBaselineShift(bs BaselineShift) bool {
	return floatutils.IsSpecified(bs)
}

func isUnspecifiedBaselineShift(bs BaselineShift) bool {
	return !isSpecifiedBaselineShift(bs)
}

func takeOrElseBaselineShift(bs, def BaselineShift) BaselineShift {
	if isSpecifiedBaselineShift(bs) {
		return bs
	}
	return def
}

// mergeBaselineShift layers b onto a: a value from b wins only where b is specified.
func mergeBaselineShift(a, b BaselineShift) BaselineShift {
	if isUnspecifiedBaselineShift(b) {
		return a
	}
	return b
}

func stringBaselineShift(bs BaselineShift) string {
	return StringBaselineShift(bs)
}

func sameBaselineShift(a, b BaselineShift) bool {
	if isUnspecifiedBaselineShift(a) && isUnspecifiedBaselineShift(b) {
		return true
	}
	return a == b
}

func semanticEqualBaselineShift(a, b BaselineShift) bool {
	if isUnspecifiedBaselineShift(a) && isUnspecifiedBaselineShift(b) {
		return true
	}
	return floatutils.Float32Equals(float32(a), float32(b), floatutils.Float32EqualityThreshold)
}

func equalBaselineShift(a, b BaselineShift) bool {
	return semanticEqualBaselineShift(a, b)
}

func copyBaselineShift(bs BaselineShift) BaselineShift {
	return bs
}

// BaselineShiftCompanion is the single exported sentinel-family symbol for BaselineShift.
var BaselineShiftCompanion = companions.NewScalarCompanion(
	func() BaselineShift { return baselineShiftUnspecified },
	stringBaselineShift,
	isSpecifiedBaselineShift,
	isUnspecifiedBaselineShift,
	takeOrElseBaselineShift,
	mergeBaselineShift,
	sameBaselineShift,
	semanticEqualBaselineShift,
	equalBaselineShift,
	copyBaselineShift,
)
