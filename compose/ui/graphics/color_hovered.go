package graphics

import (
	"image/color"
)

// Hovered Logic Copied from gioui.org@v0.10.0/internal/f32color/rgba.go

// Hovered blends dark colors towards white, and light colors towards
// black. It is approximate because it operates in non-linear sRGB space.
func Hovered(baseColor Color) Color {
	if baseColor.Alpha() == 0 {
		// Provide a reasonable default for transparent widgets.
		col := color.NRGBA{A: 0x44, R: 0x88, G: 0x88, B: 0x88}
		return FromNRGBA(col)
	}
	// 0 - 255
	const ratio = 32
	m := ColorWhite.Copy(CopyWithAlpha(baseColor.Alpha()))
	if baseColor.Luminance() > 0.5 {
		m = ColorBlack.Copy(CopyWithAlpha(baseColor.Alpha()))
	}
	return m.MixWithWeight(baseColor, ratio)
}

func Selected(c Color) Color {
	return Hovered(c)
}
