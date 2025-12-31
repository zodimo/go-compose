package font

import (
	"fmt"

	giofont "gioui.org/font"
)

// ToGioFont converts go-compose font attributes to a gio font.Font.
func ToGioFont(f FontFamily, w FontWeight, s FontStyle) giofont.Font {
	return giofont.Font{
		Typeface: giofont.Typeface(ResolveGioTypeface(f)),
		Weight:   ToGioWeight(w),
		Style:    ToGioStyle(s),
	}
}

// ToGioWeight converts a go-compose FontWeight to a gio font.Weight.
// Gio weights are offset from 0 (Normal/400).
// 100 -> -300 (Thin)
// 400 -> 0 (Normal)
// 700 -> 300 (Bold)
func ToGioWeight(w FontWeight) giofont.Weight {

	w = w.TakeOrElse(FontWeightNormal)

	// Gio weight is (CSS weight - 400).
	// We first ensure it's clamped to valid CSS range [1, 1000] (which FontWeight already should be),
	// but mostly we care about the offset.
	// Since FontWeight is an int, we can just subtract 400.
	// However, we should handle Unspecified.
	if !w.IsFontWeight() {
		return giofont.Normal
	}
	return giofont.Weight(int(w) - 400)
}

// ToGioStyle converts a go-compose FontStyle to a gio font.Style.
func ToGioStyle(s FontStyle) giofont.Style {

	s = s.TakeOrElse(FontStyleNormal)

	switch s {
	case FontStyleNormal:
		return giofont.Regular
	case FontStyleItalic:
		return giofont.Italic
	default:
		panic(fmt.Sprintf("unhandled font style: %s", s))
	}
}

// ResolveGioTypeface resolves a FontFamily to a gio Typeface string.
func ResolveGioTypeface(f FontFamily) string {
	f = CoalesceFontFamily(f, FontFamilyDefault)

	switch family := f.(type) {
	case *GenericFontFamily:
		return family.Name()
	case *DefaultFontFamily:
		// empty string lets Gio fallback to its default (usually Go font)
		return ""
	case *LoadedFontFamily:
		// Best effort: if the loaded typeface has a family name, use it?
		// For now, we can't extract a string name from a generic Typeface interface easily
		// unless we cast it to something specific, which we don't know about here.
		// Fallback to empty.
		return ""
	case *FontListFontFamily:
		// Similar to loaded, we have a list of fonts. We can't easily turn this into a single string name
		// without more context or registration.
		// Fallback to empty.
		return ""
	default:
		return ""
	}
}
