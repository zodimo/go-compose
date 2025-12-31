package unit

import (
	"math"

	"github.com/zodimo/go-compose/compose/ui/geometry"
)

// Density provides information about the density of the display.
// Used for the conversions between pixels, [Dp], [int] and [TextUnit].
type Density interface {
	// Density returns the logical density of the display.
	// This is a scaling factor for the [Dp] unit.
	Density() float32

	// FontScale returns the current user preference for the scaling factor for fonts.
	FontScale() float32

	// DpToPx converts [Dp] to pixels. Pixels are used to paint to Canvas.
	DpToPx(dp Dp) float32

	// DpRoundToPx converts [Dp] to [int] by rounding.
	DpRoundToPx(dp Dp) int

	// TextUnitToPx converts Sp to pixels. Pixels are used to paint to Canvas.
	// Panics if TextUnit other than SP unit is specified.
	TextUnitToPx(tu TextUnit) float32

	// TextUnitRoundToPx converts Sp to [int] by rounding.
	TextUnitRoundToPx(tu TextUnit) int

	// IntToDp converts an [int] pixel value to [Dp].
	IntToDp(px int) Dp

	// IntToSp converts an [int] pixel value to Sp.
	IntToSp(px int) TextUnit

	// FloatToDp converts a [float32] pixel value to a [Dp].
	FloatToDp(px float32) Dp

	// FloatToSp converts a [float32] pixel value to a Sp.
	FloatToSp(px float32) TextUnit

	// DpRectToRect converts a [DpRect] to a [Rect].
	DpRectToRect(rect DpRect) geometry.Rect

	// DpSizeToSize converts a [DpSize] to a [Size].
	DpSizeToSize(size DpSize) geometry.Size

	// SizeToDpSize converts a [Size] to a [DpSize].
	SizeToDpSize(size geometry.Size) DpSize
}

// NewDensity creates a new Density.
//
// density: The logical density of the display.
// fontScale: Current user preference for the scaling factor for fonts. Default should be 1.0 if unknown.
func NewDensity(density, fontScale float32) Density {
	return &densityImpl{
		density:   density,
		fontScale: fontScale,
	}
}

type densityImpl struct {
	density   float32
	fontScale float32
}

func (d *densityImpl) Density() float32 {
	return d.density
}

func (d *densityImpl) FontScale() float32 {
	return d.fontScale
}

func (d *densityImpl) DpToPx(dp Dp) float32 {
	return dp.Value() * d.density
}

func (d *densityImpl) DpRoundToPx(dp Dp) int {
	px := d.DpToPx(dp)
	if math.IsInf(float64(px), 0) {
		return Infinity
	}
	return int(math.Round(float64(px)))
}

func (d *densityImpl) TextUnitToPx(tu TextUnit) float32 {
	if tu.Type() != TextUnitTypeSp {
		panic("Only Sp can convert to Px")
	}
	// Sp -> Dp -> Px
	// Dp = Sp * fontScale
	// Px = Dp * density
	return tu.Value() * d.fontScale * d.density
}

func (d *densityImpl) TextUnitRoundToPx(tu TextUnit) int {
	return int(math.Round(float64(d.TextUnitToPx(tu))))
}

func (d *densityImpl) IntToDp(px int) Dp {
	return NewDp(float32(px) / d.density)
}

func (d *densityImpl) IntToSp(px int) TextUnit {
	// Px -> Dp -> Sp
	// Dp = Px / density
	// Sp = Dp / fontScale
	return Sp((float32(px) / d.density) / d.fontScale)
}

func (d *densityImpl) FloatToDp(px float32) Dp {
	return NewDp(px / d.density)
}

func (d *densityImpl) FloatToSp(px float32) TextUnit {
	return Sp((px / d.density) / d.fontScale)
}

func (d *densityImpl) DpRectToRect(rect DpRect) geometry.Rect {
	return geometry.NewRect(
		d.DpToPx(rect.Left),
		d.DpToPx(rect.Top),
		d.DpToPx(rect.Right),
		d.DpToPx(rect.Bottom),
	)
}

func (d *densityImpl) DpSizeToSize(size DpSize) geometry.Size {
	if size.IsSpecified() {
		return geometry.NewSize(d.DpToPx(size.Width), d.DpToPx(size.Height))
	}
	return geometry.SizeUnspecified
}

func (d *densityImpl) SizeToDpSize(size geometry.Size) DpSize {
	if size.IsSpecified() {
		return NewDpSize(d.FloatToDp(size.Width()), d.FloatToDp(size.Height()))
	}
	return DpSizeUnspecified
}
