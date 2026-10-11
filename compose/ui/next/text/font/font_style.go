package font

import (
	"github.com/zodimo/go-sentinel-helper/sentinel/companions"
)

// FontStyle defines whether the font is Italic or Normal.
type FontStyle int

const (
	// fontStyleUnspecified is the sentinel for FontStyle. It is deliberately
	// package-private: the only public way to obtain it is
	// FontStyleCompanion.Unspecified().
	fontStyleUnspecified = -1
	// FontStyleNormal uses the upright glyphs
	FontStyleNormal FontStyle = 0
	// FontStyleItalic uses glyphs designed for slanting
	FontStyleItalic FontStyle = 1
)

// String returns a string representation of the FontStyle.
func (s FontStyle) String() string {
	switch s {
	case FontStyleNormal:
		return "Normal"
	case FontStyleItalic:
		return "Italic"
	case fontStyleUnspecified:
		return "Unspecified"
	default:
		return "Invalid"
	}
}

// Value returns the underlying integer value.
func (s FontStyle) Value() int {
	return int(s)
}

// Equals checks if two FontStyles are equal.
func (s FontStyle) Equals(other FontStyle) bool {
	return s == other
}

func (s FontStyle) IsSpecified() bool {
	return FontStyleCompanion.IsSpecified(s)
}

// === Sentinel operations ===

func isSpecifiedFontStyle(s FontStyle) bool {
	return s != fontStyleUnspecified
}

func isUnspecifiedFontStyle(s FontStyle) bool {
	return s == fontStyleUnspecified
}

func takeOrElseFontStyle(s, def FontStyle) FontStyle {
	if isSpecifiedFontStyle(s) {
		return s
	}
	return def
}

func mergeFontStyle(a, b FontStyle) FontStyle {
	if isUnspecifiedFontStyle(b) {
		return a
	}
	return b
}

func stringFontStyle(s FontStyle) string {
	if isUnspecifiedFontStyle(s) {
		return "Unspecified"
	}
	return s.String()
}

func sameFontStyle(a, b FontStyle) bool {
	return a == b
}

func semanticEqualFontStyle(a, b FontStyle) bool {
	return a == b
}

func equalFontStyle(a, b FontStyle) bool {
	return a == b
}

func copyFontStyle(s FontStyle) FontStyle {
	return s
}

// FontStyleCompanion is the single exported sentinel-family symbol for FontStyle.
var FontStyleCompanion = companions.NewScalarCompanion(
	func() FontStyle { return fontStyleUnspecified },
	stringFontStyle,
	isSpecifiedFontStyle,
	isUnspecifiedFontStyle,
	takeOrElseFontStyle,
	mergeFontStyle,
	sameFontStyle,
	semanticEqualFontStyle,
	equalFontStyle,
	copyFontStyle,
)

// FontStyleValues returns a list of possible FontStyle values.
func FontStyleValues() []FontStyle {
	return []FontStyle{FontStyleNormal, FontStyleItalic}
}
