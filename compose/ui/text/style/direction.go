package style

import (
	"fmt"
)

// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-text/src/commonMain/kotlin/androidx/compose/ui/text/style/TextDirection.kt

const (
	// TextDirectionUnspecified represents an unset value, a usual replacement for "null"
	// when a primitive value is desired.
	TextDirectionUnspecified TextDirection = 0

	// TextDirectionLtr always sets the text direction to be Left to Right.
	TextDirectionLtr TextDirection = 1

	// TextDirectionRtl always sets the text direction to be Right to Left.
	TextDirectionRtl TextDirection = 2

	// TextDirectionContent indicates that the text direction depends on the first strong
	// directional character in the text according to the Unicode Bidirectional Algorithm.
	// If no strong directional character is present, then LayoutDirection is used to
	// resolve the final TextDirection.
	// * if used while creating a Paragraph object, LocaleList will be used to resolve
	//   the direction as a fallback instead of LayoutDirection.
	TextDirectionContent TextDirection = 3

	// TextDirectionContentOrLtr indicates that the text direction depends on the first
	// strong directional character in the text according to the Unicode Bidirectional
	// Algorithm. If no strong directional character is present, then Left to Right will
	// be used as the default direction.
	TextDirectionContentOrLtr TextDirection = 4

	// TextDirectionContentOrRtl indicates that the text direction depends on the first
	// strong directional character in the text according to the Unicode Bidirectional
	// Algorithm. If no strong directional character is present, then Right to Left will
	// be used as the default direction.
	TextDirectionContentOrRtl TextDirection = 5
)

// TextDirection defines the algorithm to be used while determining the text direction.
type TextDirection int

// String returns the string representation of the TextDirection.
func (t TextDirection) String() string {
	switch t {
	case TextDirectionLtr:
		return "Ltr"
	case TextDirectionRtl:
		return "Rtl"
	case TextDirectionContent:
		return "Content"
	case TextDirectionContentOrLtr:
		return "ContentOrLtr"
	case TextDirectionContentOrRtl:
		return "ContentOrRtl"
	case TextDirectionUnspecified:
		return "Unspecified"
	default:
		return "Invalid"
	}
}

// TextDirectionValues returns a list containing all possible values of TextDirection.
// Note: Does not include TextDirectionUnspecified.
func TextDirectionValues() []TextDirection {
	return []TextDirection{
		TextDirectionLtr,
		TextDirectionRtl,
		TextDirectionContent,
		TextDirectionContentOrLtr,
		TextDirectionContentOrRtl,
	}
}

// TextDirectionValueOf creates a TextDirection from the given integer value.
// This can be useful if you need to serialize/deserialize TextDirection values.
// Returns an error if the given value is not recognized by the preset TextDirection values.
func TextDirectionValueOf(value int) (TextDirection, error) {
	if value < 0 || value > 5 {
		return TextDirectionUnspecified, fmt.Errorf(
			"the given value=%d is not recognized by TextDirection", value,
		)
	}
	return TextDirection(value), nil
}

// TextDirection returns true if this TextDirection is not TextDirectionUnspecified.
func (t TextDirection) IsTextDirection() bool {
	return t != TextDirectionUnspecified
}

// TakeOrElse returns this TextDirection if IsSpecified() is true,
// otherwise executes the provided function and returns its result.
func (t TextDirection) TakeOrElse(block TextDirection) TextDirection {
	if t.IsTextDirection() {
		return t
	}
	return block
}
