package unit

import (
	"fmt"
)

// DpSize represents a size with Dp dimensions.
type DpSize struct {
	Width  Dp
	Height Dp
}

// NewDpSize creates a new DpSize.
func NewDpSize(width, height Dp) DpSize {
	return DpSize{Width: width, Height: height}
}

// DpSizeUnspecified is a size whose width and height are unspecified.
// This is usually a replacement for `null` when a primitive value is desired.
// Access to Width or Height on an unspecified size should be avoided.
var DpSizeUnspecified = DpSize{Width: DpUnspecified, Height: DpUnspecified}

// DpSizeZero is a DpSize with 0 DP width and 0 DP height values.
var DpSizeZero = DpSize{Width: 0, Height: 0}

// Copy returns a copy of this DpSize instance optionally overriding the width or height parameter.
func (s DpSize) Copy(width, height Dp) DpSize {
	return DpSize{Width: width, Height: height}
}

// Subtract subtracts another DpSize from this one.
func (s DpSize) Subtract(other DpSize) DpSize {
	return DpSize{Width: s.Width - other.Width, Height: s.Height - other.Height}
}

// Add adds another DpSize to this one.
func (s DpSize) Add(other DpSize) DpSize {
	return DpSize{Width: s.Width + other.Width, Height: s.Height + other.Height}
}

// Times multiplies this DpSize by a scalar.
func (s DpSize) Times(other float32) DpSize {
	return DpSize{Width: s.Width.Times(other), Height: s.Height.Times(other)}
}

// TimesInt multiplies this DpSize by an integer scalar.
func (s DpSize) TimesInt(other int) DpSize {
	return DpSize{Width: s.Width.TimesInt(other), Height: s.Height.TimesInt(other)}
}

// Div divides this DpSize by a scalar.
func (s DpSize) Div(other float32) DpSize {
	return DpSize{Width: s.Width.Div(other), Height: s.Height.Div(other)}
}

// DivInt divides this DpSize by an integer scalar.
func (s DpSize) DivInt(other int) DpSize {
	return DpSize{Width: s.Width.DivInt(other), Height: s.Height.DivInt(other)}
}

// IsSpecified returns true when both width and height are specified.
func (s DpSize) IsSpecified() bool {
	return s.Width.IsSpecified() && s.Height.IsSpecified()
}

// IsUnspecified returns true when both width and height are unspecified.
func (s DpSize) IsUnspecified() bool {
	return s.Width.IsUnspecified() && s.Height.IsUnspecified()
}

// String returns the string representation.
func (s DpSize) String() string {
	if s.IsSpecified() {
		return fmt.Sprintf("%v x %v", s.Width, s.Height)
	}
	return "DpSize.Unspecified"
}

// TakeOrElse returns this DpSize if specified, otherwise returns the result of the block.
func (s DpSize) TakeOrElse(block DpSize) DpSize {
	if s.IsSpecified() {
		return s
	}
	return block
}

// Center returns the DpOffset of the center of the rect from the point of [0, 0] with this DpSize.
func (s DpSize) Center() DpOffset {
	return DpOffset{
		X: s.Width.Div(2),
		Y: s.Height.Div(2),
	}
}

// LerpDpSize linearly interpolates between two DpSizes.
//
// The fraction argument represents position on the timeline, with 0.0 meaning that the
// interpolation has not started, returning start, 1.0 meaning that the interpolation has
// finished, returning stop, and values in between meaning that the interpolation is at the
// relevant point on the timeline between start and stop. The interpolation can be extrapolated
// beyond 0.0 and 1.0, so negative values and values greater than 1.0 are valid.
func LerpDpSize(start, stop DpSize, fraction float32) DpSize {
	return DpSize{
		Width:  LerpDp(start.Width, stop.Width, fraction),
		Height: LerpDp(start.Height, stop.Height, fraction),
	}
}
