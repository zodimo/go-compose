package material

import (
	"image/color"

	"github.com/zodimo/go-compose/compose"

	gioMaterial "gioui.org/widget/material"
)

// BasicTheme is the framework-owned representation of a Material theme.
// It holds the palette colors used by theme consumers.
// Engine-bound components that need the underlying gioui material theme
// should convert at their boundary using the seam pattern.
type BasicTheme struct {
	Bg         color.NRGBA
	Fg         color.NRGBA
	ContrastBg color.NRGBA
	ContrastFg color.NRGBA
}

// ThemeInterface provides access to the framework-owned Material theme.
type ThemeInterface interface {
	GioMaterialTheme() *BasicTheme
}

// Theme returns a ThemeInterface for the given composer.
func Theme(c compose.Composer) ThemeInterface {
	return themeImpl{
		composer: c,
	}
}

type themeImpl struct {
	composer compose.Composer
}

func (t themeImpl) GioMaterialTheme() *BasicTheme {
	gioTheme := localGioMaterialTheme.Current(t.composer)
	shaper := compose.LocalTextShaper.Current(t.composer)
	gioTheme.Shaper = shaper.Shaper
	return &BasicTheme{
		Bg:         gioTheme.Bg,
		Fg:         gioTheme.Fg,
		ContrastBg: gioTheme.ContrastBg,
		ContrastFg: gioTheme.ContrastFg,
	}
}

// LocalBasicTheme is a CompositionLocal holding the framework-owned BasicTheme.
var LocalBasicTheme = compose.CompositionLocalOf(func() *BasicTheme {
	return &BasicTheme{}
})

// localGioMaterialTheme is the internal composition local holding the
// underlying gioui material theme. It is NOT exported in the public API.
var localGioMaterialTheme = compose.CompositionLocalOf(func() *gioMaterial.Theme {
	return defaultMaterialTheme()
})

// GioThemeForEngine returns the engine material theme for engine-bound
// components (classified Door 2, e.g. textfield). The returned value is the
// gioui *material.Theme; callers assert it at their engine boundary. The
// return type is opaque to keep gioui out of the public API surface.
func GioThemeForEngine(c compose.Composer) any {
	return localGioMaterialTheme.Current(c)
}

func defaultMaterialTheme() *gioMaterial.Theme {
	materialTheme := gioMaterial.NewTheme()
	return materialTheme
}
