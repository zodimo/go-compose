package material

import (
	"gioui.org/font/gofont"
	"gioui.org/text"
	"gioui.org/widget/material"
	"github.com/zodimo/go-compose/compose"
)

type ThemeInterface interface {
	GioMaterialTheme() *material.Theme
}

func Theme(c compose.Composer) ThemeInterface {
	return themeImpl{
		composer: c,
	}
}

type themeImpl struct {
	composer compose.Composer
}

func (t themeImpl) GioMaterialTheme() *material.Theme {
	return LocalGioMaterialTheme.Current(t.composer)
}

var LocalGioMaterialTheme = compose.CompositionLocalOf(func() *material.Theme {
	return defaultMaterialTheme()
})

func defaultMaterialTheme() *material.Theme {
	materialTheme := material.NewTheme()
	materialTheme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

	return materialTheme
}
