package theme

import (
	"image/color"
)

type ColorReader = func(themeColor ThemeColor) color.Color

type ThemeColorSetter struct {
	ThemeColor ColorReader
}
