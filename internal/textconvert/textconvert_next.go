package textconvert

import (
	"fmt"

	gioFont "gioui.org/font"
	nextFont "github.com/zodimo/go-compose/compose/ui/next/text/font"
)

// ToGioFontNext converts go-compose next/font attributes to a gio font.Font.
func ToGioFontNext(f nextFont.FontFamily, w nextFont.FontWeight, s nextFont.FontStyle) gioFont.Font {
	return gioFont.Font{
		Typeface: gioFont.Typeface(ResolveGioTypefaceNext(f)),
		Weight:   ToGioWeightNext(w),
		Style:    ToGioStyleNext(s),
	}
}

// ToGioWeightNext converts a go-compose next/font FontWeight to a gio font.Weight.
func ToGioWeightNext(w nextFont.FontWeight) gioFont.Weight {
	w = w.TakeOrElse(nextFont.FontWeightNormal)
	if !w.IsFontWeight() {
		return gioFont.Normal
	}
	return gioFont.Weight(int(w) - 400)
}

// ToGioStyleNext converts a go-compose next/font FontStyle to a gio font.Style.
func ToGioStyleNext(s nextFont.FontStyle) gioFont.Style {
	s = s.TakeOrElse(nextFont.FontStyleNormal)
	switch s {
	case nextFont.FontStyleNormal:
		return gioFont.Regular
	case nextFont.FontStyleItalic:
		return gioFont.Italic
	default:
		panic(fmt.Sprintf("unhandled font style: %s", s))
	}
}

// ResolveGioTypefaceNext resolves a next/font FontFamily to a gio Typeface string.
func ResolveGioTypefaceNext(f nextFont.FontFamily) string {
	f = nextFont.CoalesceFontFamily(f, nextFont.FontFamilyDefault)
	switch family := f.(type) {
	case *nextFont.GenericFontFamily:
		return family.Name()
	case *nextFont.DefaultFontFamily:
		return ""
	case *nextFont.LoadedFontFamily:
		return ""
	case *nextFont.FontListFontFamily:
		return ""
	default:
		return ""
	}
}
