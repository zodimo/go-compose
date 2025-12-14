package theme

import (
	"gioui.org/font/gofont"
	"gioui.org/text"
	"gioui.org/widget/material"
)

func defaultMaterialTheme() *material.Theme {
	materialTheme := material.NewTheme()
	materialTheme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

	return materialTheme
}
