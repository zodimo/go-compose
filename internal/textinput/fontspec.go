package textinput

import (
	"github.com/zodimo/go-compose/compose/ui/next/text/font"
	nextFont "github.com/zodimo/go-compose/compose/ui/text/font"
	tw "github.com/zodimo/go-compose/internal/textwidget"
)

// ToFontSpec converts go-compose text/font attributes to a textwidget.FontSpec.
func ToFontSpec(f nextFont.FontFamily, w nextFont.FontWeight, s nextFont.FontStyle) tw.FontSpec {
	f = nextFont.CoalesceFontFamily(f, nextFont.FontFamilyDefault)
	typeface := ""
	weight := 400
	italic := false

	switch family := f.(type) {
	case *nextFont.GenericFontFamily:
		typeface = family.Name()
	case *nextFont.DefaultFontFamily, *nextFont.LoadedFontFamily, *nextFont.FontListFontFamily:
		typeface = "sans-serif"
	}

	if w.IsFontWeight() {
		weight = int(w)
	}

	switch s {
	case nextFont.FontStyleItalic:
		italic = true
	}

	return tw.FontSpec{
		Typeface: typeface,
		Weight:   weight,
		Italic:   italic,
	}
}

// ToFontSpecNext converts go-compose next/font attributes to a textwidget.FontSpec.
func ToFontSpecNext(f font.FontFamily, w font.FontWeight, s font.FontStyle) tw.FontSpec {
	f = font.CoalesceFontFamily(f, font.FontFamilyDefault)
	typeface := ""
	weight := 400
	italic := false

	switch family := f.(type) {
	case *font.GenericFontFamily:
		typeface = family.Name()
	case *font.DefaultFontFamily, *font.LoadedFontFamily, *font.FontListFontFamily:
		typeface = "sans-serif"
	}

	if w.IsFontWeight() {
		weight = int(w)
	}

	switch s {
	case font.FontStyleItalic:
		italic = true
	}

	return tw.FontSpec{
		Typeface: typeface,
		Weight:   weight,
		Italic:   italic,
	}
}
