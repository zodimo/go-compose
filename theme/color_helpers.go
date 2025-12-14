package theme

import "image/color"

func ToNRGBA(c color.Color) color.NRGBA {
	return c.(color.NRGBA)
}

func ToColorData(c color.Color) ThemeColorDescriptor {
	return ThemeColorDescriptor{
		color:   ToNRGBA(c),
		isColor: true,
	}
}
