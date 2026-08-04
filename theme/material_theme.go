package theme

import (
	"github.com/zodimo/go-compose/assets/fonts"

	"gioui.org/text"
	"gioui.org/widget/material"
)

func defaultMaterialTheme() *BasicTheme {
	mt := material.NewTheme()
	mt.Shaper = text.NewShaper(text.WithCollection(fonts.Collection()))
	return &BasicTheme{
		Bg:         mt.Bg,
		Fg:         mt.Fg,
		ContrastBg: mt.ContrastBg,
		ContrastFg: mt.ContrastFg,
	}
}
