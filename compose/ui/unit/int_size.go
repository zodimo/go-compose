package unit

import (
	"fmt"

	"github.com/zodimo/go-compose/pkg/sentinel"
)

// @TODO implement size fully
// IntSize represents a size with Int dimensions.
type IntSize struct {
	Width  int
	Height int
}

// NewDpSize creates a new DpSize.
func NewIntSize(width, height int) IntSize {
	return IntSize{Width: width, Height: height}
}

var IntSizeUnspecified = IntSize{Width: sentinel.IntUnspecified, Height: sentinel.IntUnspecified}
var IntSizeZero = IntSize{Width: 0, Height: 0}

func (s IntSize) Copy(width, height int) IntSize {
	return IntSize{Width: width, Height: height}
}

func (s IntSize) Subtract(other IntSize) IntSize {
	return IntSize{Width: s.Width - other.Width, Height: s.Height - other.Height}
}

func (s IntSize) Add(other IntSize) IntSize {
	return IntSize{Width: s.Width + other.Width, Height: s.Height + other.Height}
}

func (s IntSize) Times(other int) IntSize {
	return IntSize{Width: s.Width * other, Height: s.Height * other}
}

func (s IntSize) Div(other int) IntSize {
	return IntSize{Width: s.Width / other, Height: s.Height / other}
}

func (s IntSize) IsSpecified() bool {
	// Kotlin: packedValue != NaN_NaN. So if ANY part is not NaN, it is NOT Unspecified?
	// Wait, IsUnspecified is packedValue == NaN_NaN.
	// So IsSpecified is != NaN_NaN.
	// If one is NaN and other is 0, it is NOT IsUnspecified, so it IS IsSpecified?
	// No, DpSize is value class.
	// Let's assume strict checks: valid if both valid.
	return sentinel.IsSpecifiedIntValue(s.Width) && sentinel.IsSpecifiedIntValue(s.Height)
}

func (s IntSize) IsUnspecified() bool {
	return sentinel.IsUnspecifiedIntValue(s.Width) && sentinel.IsUnspecifiedIntValue(s.Height)
}

func (s IntSize) String() string {
	if s.IsSpecified() {
		return fmt.Sprintf("%v x %v", s.Width, s.Height)
	}
	return "IntSize.Unspecified"
}

func (s IntSize) TakeOrElse(block IntSize) IntSize {
	if s.IsSpecified() {
		return s
	}
	return block
}

func (s IntSize) Center() IntSize {
	return IntSize{
		Width:  s.Width / 2,
		Height: s.Height / 2,
	}
}

func LerpIntSize(start, stop IntSize, fraction float32) IntSize {
	return IntSize{
		Width:  LerpInt(start.Width, stop.Width, fraction),
		Height: LerpInt(start.Height, stop.Height, fraction),
	}
}

// LerpInt linearly interpolates between two Ints.
func LerpInt(start, stop int, fraction float32) int {
	return int(lerpBetween(float32(start), float32(stop), float64(fraction)))
}
