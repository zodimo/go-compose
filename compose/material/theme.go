package material

import (
	gioMaterial "gioui.org/widget/material"
	"github.com/zodimo/go-compose/compose"
)

type ThemeInterface interface {
	GioMaterialTheme() *gioMaterial.Theme
}

func Theme(c compose.Composer) ThemeInterface {
	return themeImpl{
		composer: c,
	}
}

type themeImpl struct {
	composer compose.Composer
}

func (t themeImpl) GioMaterialTheme() *gioMaterial.Theme {
	shaper := compose.LocalTextShaper.Current(t.composer)
	theme := LocalGioMaterialTheme.Current(t.composer)
	theme.Shaper = shaper.Shaper
	return theme
}

var LocalGioMaterialTheme = compose.CompositionLocalOf(func() *gioMaterial.Theme {
	return defaultMaterialTheme()
})

func defaultMaterialTheme() *gioMaterial.Theme {
	materialTheme := gioMaterial.NewTheme()
	return materialTheme
}
