package text

import (
	gioFont "gioui.org/font"
	"github.com/zodimo/go-compose/internal/textconvert"
)

// TextStyleFromGioFont converts a gio font.Font to a TextStyle.
// This is useful for integrating with Gio's font system.
func TextStyleFromGioFont(gf gioFont.Font) *TextStyle {
	fontFamily, fontWeight, fontStyle := textconvert.FromGioFont(gf)

	return TextStyleFromOptions(
		WithFontFamily(fontFamily),
		WithFontWeight(fontWeight),
		WithFontStyle(fontStyle),
	)
}

// ToGioFont converts a TextStyle to a gio font.Font.
// This extracts the font family, weight, and style from the TextStyle
// and converts them to Gio's font representation.
func ToGioFont(ts *TextStyle) gioFont.Font {
	ts = CoalesceTextStyle(ts, TextStyleUnspecified)

	return textconvert.ToGioFont(
		ts.FontFamily(),
		ts.FontWeight(),
		ts.FontStyle(),
	)
}
