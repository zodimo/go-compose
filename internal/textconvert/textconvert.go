// Package textconvert provides conversion functions between go-compose text/font types
// and gioui.org text/font types. These functions were moved from compose/ui/text/font
// and compose/ui/text/style to keep gioui.org types out of the public API surface.
package textconvert

import (
	"fmt"

	gioFont "gioui.org/font"
	gioText "gioui.org/text"
	"github.com/zodimo/go-compose/compose/ui/text/font"
	"github.com/zodimo/go-compose/compose/ui/text/style"
)

// --- Font conversions (from compose/ui/text/font) ---

// ToGioFont converts go-compose font attributes to a gio font.Font.
func ToGioFont(f font.FontFamily, w font.FontWeight, s font.FontStyle) gioFont.Font {
	return gioFont.Font{
		Typeface: gioFont.Typeface(ResolveGioTypeface(f)),
		Weight:   ToGioWeight(w),
		Style:    ToGioStyle(s),
	}
}

// ToGioWeight converts a go-compose FontWeight to a gio font.Weight.
func ToGioWeight(w font.FontWeight) gioFont.Weight {
	w = w.TakeOrElse(font.FontWeightNormal)
	if !w.IsFontWeight() {
		return gioFont.Normal
	}
	return gioFont.Weight(int(w) - 400)
}

// ToGioStyle converts a go-compose FontStyle to a gio font.Style.
func ToGioStyle(s font.FontStyle) gioFont.Style {
	s = s.TakeOrElse(font.FontStyleNormal)
	switch s {
	case font.FontStyleNormal:
		return gioFont.Regular
	case font.FontStyleItalic:
		return gioFont.Italic
	default:
		panic(fmt.Sprintf("unhandled font style: %s", s))
	}
}

// ResolveGioTypeface resolves a FontFamily to a gio Typeface string.
func ResolveGioTypeface(f font.FontFamily) string {
	f = font.CoalesceFontFamily(f, font.FontFamilyDefault)
	switch family := f.(type) {
	case *font.GenericFontFamily:
		return family.Name()
	case *font.DefaultFontFamily:
		return ""
	case *font.LoadedFontFamily:
		return ""
	case *font.FontListFontFamily:
		return ""
	default:
		return ""
	}
}

// FromGioFont converts a gio font.Font to go-compose font attributes.
func FromGioFont(gf gioFont.Font) (font.FontFamily, font.FontWeight, font.FontStyle) {
	return FromGioTypeface(gf.Typeface), FromGioWeight(gf.Weight), FromGioStyle(gf.Style)
}

// FromGioWeight converts a gio font.Weight to a go-compose FontWeight.
func FromGioWeight(w gioFont.Weight) font.FontWeight {
	cssWeight := int(w) + 400
	if cssWeight < 1 {
		cssWeight = 1
	} else if cssWeight > 1000 {
		cssWeight = 1000
	}
	return font.FontWeight(cssWeight)
}

// FromGioStyle converts a gio font.Style to a go-compose FontStyle.
func FromGioStyle(s gioFont.Style) font.FontStyle {
	switch s {
	case gioFont.Regular:
		return font.FontStyleNormal
	case gioFont.Italic:
		return font.FontStyleItalic
	default:
		return font.FontStyleNormal
	}
}

// FromGioTypeface converts a gio font.Typeface to a go-compose FontFamily.
func FromGioTypeface(t gioFont.Typeface) font.FontFamily {
	name := string(t)
	if name == "" {
		return font.FontFamilyDefault
	}
	switch name {
	case "sans-serif":
		return font.FontFamilySansSerif
	case "serif":
		return font.FontFamilySerif
	case "monospace":
		return font.FontFamilyMonospace
	case "cursive":
		return font.FontFamilyCursive
	default:
		return font.NewGenericFontFamily(name, name)
	}
}

// --- Style conversions (from compose/ui/text/style) ---

// TextAlignToGioTextAlignment converts a go-compose TextAlign to a gio text.Alignment.
func TextAlignToGioTextAlignment(t style.TextAlign) gioText.Alignment {
	align := t.TakeOrElse(style.TextAlignStart)
	return gioText.Alignment(align)
}

// FromGioTextAlign converts a gio text.Alignment to a go-compose TextAlign.
func FromGioTextAlign(gioAlign gioText.Alignment) style.TextAlign {
	return style.TextAlign(gioAlign)
}

// LineBreakToGioWrapPolicy converts a go-compose LineBreak to a gio text.WrapPolicy.
func LineBreakToGioWrapPolicy(l style.LineBreak) gioText.WrapPolicy {
	linebreak := l.TakeOrElse(style.LineBreakParagraph)
	return gioText.WrapPolicy(linebreak)
}

// GioWrapPolicyToLineBreak converts a gio text.WrapPolicy to a go-compose LineBreak.
func GioWrapPolicyToLineBreak(gioWrapPolicy gioText.WrapPolicy) style.LineBreak {
	return style.LineBreak(gioWrapPolicy)
}
