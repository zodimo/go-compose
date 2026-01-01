package graphics

import (
	"github.com/zodimo/go-compose/compose/ui/graphics/colorspace"
	"github.com/zodimo/go-maybe"
)

type ColorCopyOptions struct {
	Alpha, Red, Green, Blue maybe.Maybe[float32]
}
type ColorCopyOption func(*ColorCopyOptions) ColorCopyOptions

func CopyWithAlpha(alpha float32) ColorCopyOption {
	return func(o *ColorCopyOptions) ColorCopyOptions {
		o.Alpha = maybe.Some(alpha)
		return *o
	}
}

func CopyWithRed(red float32) ColorCopyOption {
	return func(o *ColorCopyOptions) ColorCopyOptions {
		o.Red = maybe.Some(red)
		return *o
	}
}

func CopyWithGreen(green float32) ColorCopyOption {
	return func(o *ColorCopyOptions) ColorCopyOptions {
		o.Green = maybe.Some(green)
		return *o
	}
}

func CopyWithBlue(blue float32) ColorCopyOption {
	return func(o *ColorCopyOptions) ColorCopyOptions {
		o.Blue = maybe.Some(blue)
		return *o
	}
}

// Copy creates a new color with modified components.
func (c Color) Copy(opts ...ColorCopyOption) Color {
	var o ColorCopyOptions = ColorCopyOptions{
		Alpha: maybe.None[float32](),
		Red:   maybe.None[float32](),
		Green: maybe.None[float32](),
		Blue:  maybe.None[float32](),
	}
	for _, opt := range opts {
		opt(&o)
	}

	id := c.ColorSpaceId()
	space := colorspace.Get(id)
	return NewColor(
		o.Alpha.OrElse(c.Alpha()),
		o.Red.OrElse(c.Red()),
		o.Green.OrElse(c.Green()),
		o.Blue.OrElse(c.Blue()),
		space,
	)
}
