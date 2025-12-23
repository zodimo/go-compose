package colorspace

import (
	"github.com/zodimo/go-compose/compose/ui/util"
)

// WhitePoint represents a white point in the CIE xy chromaticity diagram.
type WhitePoint struct {
	X float32
	Y float32
}

// NewWhitePoint creates a new WhitePoint.
func NewWhitePoint(x, y float32) WhitePoint {
	return WhitePoint{X: x, Y: y}
}

// ToXyz converts this white point to reduced XYZ values (luminance Y = 1).
func (wp WhitePoint) ToXyz() []float32 {
	return []float32{wp.X / wp.Y, 1.0, (1.0 - wp.X - wp.Y) / wp.Y}
}

// Compare compares two WhitePoints with a precision of 1e-3.
func CompareWhitePoint(a, b WhitePoint) bool {
	if a == b {
		return true
	}
	return util.Abs(a.X-b.X) < 1e-3 && util.Abs(a.Y-b.Y) < 1e-3
}
