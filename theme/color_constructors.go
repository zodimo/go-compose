package theme

import "image/color"

func SpecificColor(color color.Color) ThemeColorDescriptor {
	return ThemeColorDescriptor{
		color:   ToNRGBA(color),
		isColor: true,
	}
}

func ThemeRoleColor(colorRole ColorRole) ThemeColorDescriptor {
	return ThemeColorDescriptor{
		colorRole: colorRole,
	}
}

func ThemeColorFromColor(color color.Color) ThemeColor {
	return themeColor{
		tokenColor: TokenColor(ToNRGBA(color)),
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
