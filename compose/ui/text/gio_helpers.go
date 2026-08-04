package text

import (
	gioFont "gioui.org/font"

	"github.com/zodimo/go-compose/compose/ui/text/font"
	"github.com/zodimo/go-compose/internal/textconvert"
)

// TextStyleFromGioFont creates a TextStyle from a framework font.Font descriptor.
// This is useful for integrating with font resources.
func TextStyleFromGioFont(f font.Font) *TextStyle {
	return TextStyleFromOptions(
		WithFontFamily(font.ToFontFamily(f)),
		WithFontWeight(f.Weight()),
		WithFontStyle(f.Style()),
	)
}

// toGioFont converts a TextStyle to a gio font.Font.
// This is seam-only; it returns a gioui type.
func toGioFont(ts *TextStyle) gioFont.Font {
	ts = CoalesceTextStyle(ts, TextStyleUnspecified)

	return textconvert.ToGioFont(
		ts.FontFamily(),
		ts.FontWeight(),
		ts.FontStyle(),
	)
}
