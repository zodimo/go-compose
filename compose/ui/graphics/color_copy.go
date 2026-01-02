package graphics

import (
	"github.com/zodimo/go-compose/compose/ui/graphics/colorspace"
	"github.com/zodimo/go-compose/pkg/floatutils"
)

type ColorCopyOptions struct {
	Alpha, Red, Green, Blue float32
}
type ColorCopyOption func(*ColorCopyOptions) ColorCopyOptions

func CopyWithAlpha(alpha float32) ColorCopyOption {
	return func(o *ColorCopyOptions) ColorCopyOptions {
		o.Alpha = alpha
		return *o
	}
}

func CopyWithRed(red float32) ColorCopyOption {
	return func(o *ColorCopyOptions) ColorCopyOptions {
		o.Red = red
		return *o
	}
}

func CopyWithGreen(green float32) ColorCopyOption {
	return func(o *ColorCopyOptions) ColorCopyOptions {
		o.Green = green
		return *o
	}
}

func CopyWithBlue(blue float32) ColorCopyOption {
	return func(o *ColorCopyOptions) ColorCopyOptions {
		o.Blue = blue
		return *o
	}
}

// Copy creates a new color with modified components.
func (c Color) Copy(opts ...ColorCopyOption) Color {
	var o ColorCopyOptions = ColorCopyOptions{
		Alpha: floatutils.Float32Unspecified,
		Red:   floatutils.Float32Unspecified,
		Green: floatutils.Float32Unspecified,
		Blue:  floatutils.Float32Unspecified,
	}
	for _, opt := range opts {
		opt(&o)
	}

	id := c.ColorSpaceId()
	space := colorspace.Get(id)
	return NewColor(
		floatutils.TakeOrElse(o.Red, c.Red()),
		floatutils.TakeOrElse(o.Green, c.Green()),
		floatutils.TakeOrElse(o.Blue, c.Blue()),
		floatutils.TakeOrElse(o.Alpha, c.Alpha()),
		space,
	)
}
