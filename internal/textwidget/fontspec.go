package textwidget

import (
	gioFont "gioui.org/font"
)

// FontSpec describes font attributes for text shaping, without gioui types.
type FontSpec struct {
	// Typeface is the name of the font family (e.g., "sans-serif").
	Typeface string
	// Weight is the CSS weight value [1-1000].
	Weight int
	// Italic selects the italic variant.
	Italic bool
}

// ToGioFont converts a FontSpec to a gioui font.Font.
func ToGioFont(f FontSpec) gioFont.Font {
	w := f.Weight
	if w < 1 {
		w = 400
	}
	style := gioFont.Regular
	if f.Italic {
		style = gioFont.Italic
	}
	return gioFont.Font{
		Typeface: gioFont.Typeface(f.Typeface),
		Weight:   gioFont.Weight(w - 400),
		Style:    style,
	}
}

// DefaultFontSpec returns the default font spec (sans-serif, normal weight, upright).
func DefaultFontSpec() FontSpec {
	return FontSpec{
		Typeface: "sans-serif",
		Weight:   400,
		Italic:   false,
	}
}
