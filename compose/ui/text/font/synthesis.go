package font

import (
	"fmt"

	"github.com/zodimo/go-sentinel-helper/sentinel/companions"
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
	// fontSynthesisUnspecified is the sentinel for FontSynthesis. It is
	// deliberately package-private: the only public way to obtain it is
	// FontSynthesisCompanion.Unspecified().
	fontSynthesisUnspecified = &FontSynthesis{value: -1}
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
	return f != nil && f != fontSynthesisUnspecified
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
		return b == fontSynthesisUnspecified
	}
	if b == nil {
		return a == fontSynthesisUnspecified
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

	a = CoalesceFontSynthesis(a, fontSynthesisUnspecified)
	b = CoalesceFontSynthesis(b, fontSynthesisUnspecified)

	return a.value == b.value
}

func CoalesceFontSynthesis(ptr, def *FontSynthesis) *FontSynthesis {
	if ptr == nil {
		return def
	}
	return ptr
}

// === Sentinel operations ===
//
// These package-private operations back FontSynthesisCompanion, the single
// exported sentinel-family symbol for FontSynthesis. Callers reach the
// reserved "unspecified" value only through FontSynthesisCompanion.Unspecified().

func isSpecifiedFontSynthesis(f *FontSynthesis) bool {
	return IsFontSynthesis(f)
}

func isUnspecifiedFontSynthesis(f *FontSynthesis) bool {
	return !IsFontSynthesis(f)
}

func takeOrElseFontSynthesis(f, def *FontSynthesis) *FontSynthesis {
	return TakeOrElseFontSynthesis(f, def)
}

func mergeFontSynthesis(a, b *FontSynthesis) *FontSynthesis {
	if isUnspecifiedFontSynthesis(b) {
		return coalesceFontSynthesis(a, fontSynthesisUnspecified)
	}
	return b
}

func stringFontSynthesis(f *FontSynthesis) string {
	return StringFontSynthesis(f)
}

func sameFontSynthesis(a, b *FontSynthesis) bool {
	return SameFontSynthesis(a, b)
}

func semanticEqualFontSynthesis(a, b *FontSynthesis) bool {
	return SemanticEqualFontSynthesis(a, b)
}

func equalFontSynthesis(a, b *FontSynthesis) bool {
	return EqualFontSynthesis(a, b)
}

func coalesceFontSynthesis(ptr, def *FontSynthesis) *FontSynthesis {
	return CoalesceFontSynthesis(ptr, def)
}

// copyFontSynthesis returns a copy of s. The source is never mutated.
func copyFontSynthesis(s *FontSynthesis, options ...func(*FontSynthesis)) *FontSynthesis {
	cp := &FontSynthesis{}
	if s != nil {
		cp.value = s.value
	}
	for _, option := range options {
		option(cp)
	}
	return cp
}

// FontSynthesisCompanion is the single exported sentinel-family symbol for FontSynthesis.
var FontSynthesisCompanion = companions.NewStructCompanion(
	func() *FontSynthesis { return fontSynthesisUnspecified },
	stringFontSynthesis,
	isSpecifiedFontSynthesis,
	isUnspecifiedFontSynthesis,
	takeOrElseFontSynthesis,
	mergeFontSynthesis,
	sameFontSynthesis,
	semanticEqualFontSynthesis,
	equalFontSynthesis,
	coalesceFontSynthesis,
	copyFontSynthesis,
)
