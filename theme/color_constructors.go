package theme

import "image/color"

func SpecificColor(color color.Color) ColorDescriptor {
	return colorDescriptor{
		color:   toNRGBA(color),
		isColor: true,
	}
}

func ThemeRoleColor(colorRole ColorRole) ColorDescriptor {
	return colorDescriptor{
		colorRole: colorRole,
	}
}

func ThemeColorFromColor(color color.Color) ThemeColor {
	return themeColor{
		tokenColor: TokenColor(toNRGBA(color)),
	}
}
func ThemeColorFromNRGBA(color color.NRGBA) ThemeColor {
	return themeColor{
		tokenColor: TokenColor(color),
	}
}

func ThemeColorFromTokenColor(color TokenColor) ThemeColor {
	return themeColor{
		tokenColor: color,
	}
}
