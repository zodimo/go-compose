package style

import (
	"fmt"
)

// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-text/src/commonMain/kotlin/androidx/compose/ui/text/style/TextAlign.kt

// TextAlign defines how to align text horizontally.
// TextAlign controls how text aligns in the space it appears.
type TextAlign int

// TextAlign constants
const (
	// TextAlignUnspecified represents an unset value, a usual replacement for "null"
	// when a primitive value is desired.
	TextAlignUnspecified TextAlign = 0

	// TextAlignLeft aligns the text on the left edge of the container.
	TextAlignLeft TextAlign = 1

	// TextAlignRight aligns the text on the right edge of the container.
	TextAlignRight TextAlign = 2

	// TextAlignCenter aligns the text in the center of the container.
	TextAlignCenter TextAlign = 3

	// TextAlignJustify stretches lines of text that end with a soft line break
	// to fill the width of the container.
	// Lines that end with hard line breaks are aligned towards the Start edge.
	TextAlignJustify TextAlign = 4

	// TextAlignStart aligns the text on the leading edge of the container.
	// For Left to Right text, this is the left edge.
	// For Right to Left text, like Arabic, this is the right edge.
	TextAlignStart TextAlign = 5

	// TextAlignEnd aligns the text on the trailing edge of the container.
	// For Left to Right text, this is the right edge.
	// For Right to Left text, like Arabic, this is the left edge.
	TextAlignEnd TextAlign = 6
)

// String returns the string representation of the TextAlign.
func (t TextAlign) String() string {
	switch t {
	case TextAlignLeft:
		return "Left"
	case TextAlignRight:
		return "Right"
	case TextAlignCenter:
		return "Center"
	case TextAlignJustify:
		return "Justify"
	case TextAlignStart:
		return "Start"
	case TextAlignEnd:
		return "End"
	case TextAlignUnspecified:
		return "Unspecified"
	default:
		return "Invalid"
	}
}

// TextAlignValues returns a list containing all possible values of TextAlign.
func TextAlignValues() []TextAlign {
	return []TextAlign{
		TextAlignLeft,
		TextAlignRight,
		TextAlignCenter,
		TextAlignJustify,
		TextAlignStart,
		TextAlignEnd,
	}
}

// TextAlignValueOf creates a TextAlign from the given integer value.
// This can be useful if you need to serialize/deserialize TextAlign values.
// Returns an error if the given value is not recognized by the preset TextAlign values.
func TextAlignValueOf(value int) (TextAlign, error) {
	switch value {
	case 1:
		return TextAlignLeft, nil
	case 2:
		return TextAlignRight, nil
	case 3:
		return TextAlignCenter, nil
	case 4:
		return TextAlignJustify, nil
	case 5:
		return TextAlignStart, nil
	case 6:
		return TextAlignEnd, nil
	case 0:
		return TextAlignUnspecified, nil
	default:
		return TextAlignUnspecified, fmt.Errorf("the given value=%d is not recognized by TextAlign", value)
	}
}

// IsSpecified returns true if this TextAlign is not TextAlignUnspecified.
func (t TextAlign) IsSpecified() bool {
	return t != TextAlignUnspecified
}

// TakeOrElse returns this TextAlign if IsSpecified() is true,
// otherwise executes the provided function and returns its result.
func (t TextAlign) TakeOrElse(other TextAlign) TextAlign {
	if t.IsSpecified() {
		return t
	}
	return other
}
