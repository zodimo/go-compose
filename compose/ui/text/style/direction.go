package style

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/ui/text/unit"
	"github.com/zodimo/go-sentinel-helper/sentinel/companions"
)

// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-text/src/commonMain/kotlin/androidx/compose/ui/next/text/style/TextDirection.kt

const (
	// textDirectionUnspecified represents an unset value, a usual replacement for "null"
	// when a primitive value is desired. It is deliberately package-private: the
	// only public way to obtain it is TextDirectionCompanion.Unspecified().
	textDirectionUnspecified TextDirection = 0

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
	case textDirectionUnspecified:
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
		return textDirectionUnspecified, fmt.Errorf(
			"the given value=%d is not recognized by TextDirection", value,
		)
	}
	return TextDirection(value), nil
}

// TextDirection returns true if this TextDirection is not TextDirectionUnspecified.
func (t TextDirection) IsSpecified() bool {
	return TextDirectionCompanion.IsSpecified(t)
}

// === Sentinel operations ===

func isSpecifiedTextDirection(t TextDirection) bool {
	return t != textDirectionUnspecified
}

func isUnspecifiedTextDirection(t TextDirection) bool {
	return t == textDirectionUnspecified
}

func takeOrElseTextDirection(t, def TextDirection) TextDirection {
	if isSpecifiedTextDirection(t) {
		return t
	}
	return def
}

func mergeTextDirection(a, b TextDirection) TextDirection {
	if isUnspecifiedTextDirection(b) {
		return a
	}
	return b
}

func stringTextDirection(t TextDirection) string {
	return t.String()
}

func sameTextDirection(a, b TextDirection) bool {
	return a == b
}

func semanticEqualTextDirection(a, b TextDirection) bool {
	return a == b
}

func equalTextDirection(a, b TextDirection) bool {
	return a == b
}

func copyTextDirection(t TextDirection) TextDirection {
	return t
}

// TextDirectionCompanion is the single exported sentinel-family symbol for TextDirection.
var TextDirectionCompanion = companions.NewScalarCompanion(
	func() TextDirection { return textDirectionUnspecified },
	stringTextDirection,
	isSpecifiedTextDirection,
	isUnspecifiedTextDirection,
	takeOrElseTextDirection,
	mergeTextDirection,
	sameTextDirection,
	semanticEqualTextDirection,
	equalTextDirection,
	copyTextDirection,
)

// ResolveTextDirection resolves the TextDirection, using the LayoutDirection if necessary.
// If the TextDirection is unspecified, it defaults to Ltr.
func ResolveTextDirection(layoutDirection unit.LayoutDirection, textDirection TextDirection) TextDirection {
	if isSpecifiedTextDirection(textDirection) {
		return textDirection
	}
	if layoutDirection == unit.LayoutDirectionRtl {
		return TextDirectionRtl
	}
	return TextDirectionLtr
}
