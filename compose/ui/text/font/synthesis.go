package font

import (
	"fmt"
)

// FontSynthesis specifies whether the system should fake bold or slanted glyphs
// when the FontFamily used does not contain bold or oblique Fonts.
type FontSynthesis struct {
	value int
}

const (
	// Internal flag constants matching Kotlin implementation
	synthesisAllFlags   = 0xffff
	synthesisWeightFlag = 0x1
	synthesisStyleFlag  = 0x2
)

var (
	FontSynthesisUnspecified = &FontSynthesis{value: -1}
	// FontSynthesisNone turns off font synthesis.
	// Neither bold nor slanted faces are synthesized.
	FontSynthesisNone = &FontSynthesis{value: 0}

	// FontSynthesisWeight synthesizes only bold font if not available.
	// Slanted fonts will not be synthesized.
	FontSynthesisWeight = &FontSynthesis{value: synthesisWeightFlag}

	// FontSynthesisStyle synthesizes only slanted font if not available.
	// Bold fonts will not be synthesized.
	FontSynthesisStyle = &FontSynthesis{value: synthesisStyleFlag}

	// FontSynthesisAll synthesizes both bold and slanted fonts if either
	// is not available in the FontFamily.
	FontSynthesisAll = &FontSynthesis{value: synthesisAllFlags}
)

// Value returns the underlying integer value.
func (f FontSynthesis) Value() int {
	return f.value
}

// IsWeightOn returns true if weight synthesis is enabled.
func (f FontSynthesis) IsWeightOn() bool {
	return f.value&synthesisWeightFlag != 0
}

// IsStyleOn returns true if style synthesis is enabled.
func (f FontSynthesis) IsStyleOn() bool {
	return f.value&synthesisStyleFlag != 0
}

// FontSynthesisValueOf creates a FontSynthesis from an integer value.
// Returns an error if the value is not recognized.
func FontSynthesisValueOf(value int) (*FontSynthesis, error) {
	if value != 0 && value != synthesisWeightFlag && value != synthesisStyleFlag && value != synthesisAllFlags {
		return nil, fmt.Errorf("the given value=%d is not recognized by FontSynthesis", value)
	}
	return &FontSynthesis{value: value}, nil
}

func StringFontSynthesis(f *FontSynthesis) string {
	if f == nil {
		return "FontSynthesis(nil)"
	}
	switch f.value {
	case -1:
		return "Unspecified"
	case 0:
		return "None"
	case synthesisWeightFlag:
		return "Weight"
	case synthesisStyleFlag:
		return "Style"
	case synthesisAllFlags:
		return "All"
	default:
		return "Invalid"
	}
}

func IsFontSynthesis(f *FontSynthesis) bool {
	return f != nil && f != FontSynthesisUnspecified
}

func TakeOrElseFontSynthesis(f, defaultFontSynthesis *FontSynthesis) *FontSynthesis {
	if !IsFontSynthesis(f) {
		return defaultFontSynthesis
	}
	return f
}

// Identity (2 ns)
func SameFontSynthesis(a, b *FontSynthesis) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return b == FontSynthesisUnspecified
	}
	if b == nil {
		return a == FontSynthesisUnspecified
	}
	return a == b
}

func EqualFontSynthesis(a, b *FontSynthesis) bool {
	if !SameFontSynthesis(a, b) {
		return SemanticEqualFontSynthesis(a, b)
	}
	return true
}

// Semantic equality (field-by-field, 20 ns)
func SemanticEqualFontSynthesis(a, b *FontSynthesis) bool {

	a = CoalesceFontSynthesis(a, FontSynthesisUnspecified)
	b = CoalesceFontSynthesis(b, FontSynthesisUnspecified)

	return a.value == b.value
}

func CoalesceFontSynthesis(ptr, def *FontSynthesis) *FontSynthesis {
	if ptr == nil {
		return def
	}
	return ptr
}
