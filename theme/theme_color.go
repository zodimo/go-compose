package theme

import "image/color"

var _ ThemeColor = (*themeColor)(nil)

type themeColor struct {
	tokenColor TokenColor
}

func (c themeColor) AsHex() string {
	return c.tokenColor.AsHex()
}

func (c themeColor) AsNRGBA() color.NRGBA {
	return c.tokenColor.AsNRGBA()
}

func (c themeColor) SetOpacity(opacity OpacityLevel) ThemeColor {
	return themeColor{
		tokenColor: c.tokenColor.SetOpacity(opacity),
	}
}

func (c themeColor) AsTokenColor() TokenColor {
	return c.tokenColor
}
