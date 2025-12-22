package font

// FontStyle defines whether the font is Italic or Normal.
type FontStyle int

const (
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
	default:
		return "Invalid"
	}
}

// Value returns the underlying integer value.
func (s FontStyle) Value() int {
	return int(s)
}

// FontStyleValues returns a list of possible FontStyle values.
func FontStyleValues() []FontStyle {
	return []FontStyle{FontStyleNormal, FontStyleItalic}
}
