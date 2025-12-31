package style

import (
	"fmt"

	gioText "gioui.org/text"
)

// TextAlign defines how to align text horizontally.
// TextAlign controls how text aligns in the space it appears.
type TextAlign gioText.Alignment

// TextAlign constants
const (
	// TextAlignUnspecified represents an unset value, a usual replacement for "null"
	// when a primitive value is desired.
	TextAlignUnspecified TextAlign = 99

	// TextAlignStart aligns the text on the leading edge of the container.
	// For Left to Right text, this is the left edge.
	// For Right to Left text, like Arabic, this is the right edge.
	TextAlignStart TextAlign = TextAlign(gioText.Start)

	// TextAlignEnd aligns the text on the trailing edge of the container.
	// For Left to Right text, this is the right edge.
	// For Right to Left text, like Arabic, this is the left edge.
	TextAlignEnd TextAlign = TextAlign(gioText.End)

	// TextAlignCenter aligns the text in the center of the container.
	TextAlignMiddle TextAlign = TextAlign(gioText.Middle)
)

// String returns the string representation of the TextAlign.
func (t TextAlign) String() string {
	switch t {
	case TextAlignStart:
		return "Start"
	case TextAlignEnd:
		return "End"
	case TextAlignMiddle:
		return "Middle"
	case TextAlignUnspecified:
		return "Unspecified"
	default:
		return "Invalid"
	}
}

// TextAlignValues returns a list containing all possible values of TextAlign.
func TextAlignValues() []TextAlign {
	return []TextAlign{
		TextAlignStart,
		TextAlignEnd,
		TextAlignMiddle,
	}
}

// TextAlignValueOf creates a TextAlign from the given integer value.
// This can be useful if you need to serialize/deserialize TextAlign values.
// Returns an error if the given value is not recognized by the preset TextAlign values.
func TextAlignValueOf(value int) (TextAlign, error) {
	switch value {
	case 0:
		return TextAlignStart, nil
	case 1:
		return TextAlignEnd, nil
	case 2:
		return TextAlignMiddle, nil
	case 99:
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
