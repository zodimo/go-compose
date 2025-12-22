package font

import "fmt"

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
	// FontSynthesisNone turns off font synthesis.
	// Neither bold nor slanted faces are synthesized.
	FontSynthesisNone = FontSynthesis{value: 0}

	// FontSynthesisWeight synthesizes only bold font if not available.
	// Slanted fonts will not be synthesized.
	FontSynthesisWeight = FontSynthesis{value: synthesisWeightFlag}

	// FontSynthesisStyle synthesizes only slanted font if not available.
	// Bold fonts will not be synthesized.
	FontSynthesisStyle = FontSynthesis{value: synthesisStyleFlag}

	// FontSynthesisAll synthesizes both bold and slanted fonts if either
	// is not available in the FontFamily.
	FontSynthesisAll = FontSynthesis{value: synthesisAllFlags}
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

// String returns a string representation of the FontSynthesis.
func (f FontSynthesis) String() string {
	switch f.value {
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

// Equals checks if two FontSynthesis values are equal.
func (f FontSynthesis) Equals(other FontSynthesis) bool {
	return f.value == other.value
}

// FontSynthesisValueOf creates a FontSynthesis from an integer value.
// Returns an error if the value is not recognized.
func FontSynthesisValueOf(value int) (FontSynthesis, error) {
	if value != 0 && value != synthesisWeightFlag && value != synthesisStyleFlag && value != synthesisAllFlags {
		return FontSynthesis{}, fmt.Errorf("the given value=%d is not recognized by FontSynthesis", value)
	}
	return FontSynthesis{value: value}, nil
}
