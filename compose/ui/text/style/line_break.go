package style

import (
	"fmt"

	gioText "gioui.org/text"
	"github.com/zodimo/go-sentinel-helper/sentinel/companions"
)

// LineBreak configures strategies for choosing where to break lines of text for line
type LineBreak gioText.WrapPolicy

// Basic, fast line breaking. Ideal for text input fields, as it will cause minimal text
// reflow when editing.
var LineBreakSimple LineBreak = LineBreak(gioText.WrapGraphemes)

// Looser breaking rules, suitable for short text such as titles or narrow newspaper
// columns. For longer lines of text, use [Paragraph] for improved readability.
var LineBreakHeading LineBreak = LineBreak(gioText.WrapWords)

// Slower, higher quality line breaking for improved readability. Suitable for larger
// amounts of text.
var LineBreakParagraph LineBreak = LineBreak(gioText.WrapHeuristically)

// lineBreakUnspecified represents an unset value, a usual replacement for "null" when a
// primitive value is desired. It is deliberately package-private: the only public way
// to obtain it is LineBreakCompanion.Unspecified().
const lineBreakUnspecified LineBreak = 99

func (l LineBreak) String() string {
	var description string
	switch l {
	case LineBreakSimple:
		description = "LineBreakSimple/WrapGraphemes"
	case LineBreakHeading:
		description = "LineBreakHeading/WrapWords"
	case LineBreakParagraph:
		description = "LineBreakParagraph/WrapHeuristically"
	case lineBreakUnspecified:
		description = "LineBreakUnspecified"
	default:
		panic(fmt.Sprintf("unknown line break strategy: %d", l))
	}
	return fmt.Sprintf("LineBreak(strategy=%s)", description)
}

// IsSpecified returns true if it is not [LineBreakUnspecified].
func (l LineBreak) IsSpecified() bool {
	return isSpecifiedLineBreak(l)
}

// === Sentinel operations ===

func isSpecifiedLineBreak(l LineBreak) bool {
	return l != lineBreakUnspecified
}

func isUnspecifiedLineBreak(l LineBreak) bool {
	return l == lineBreakUnspecified
}

func takeOrElseLineBreak(l, def LineBreak) LineBreak {
	if isSpecifiedLineBreak(l) {
		return l
	}
	return def
}

func mergeLineBreak(a, b LineBreak) LineBreak {
	if isUnspecifiedLineBreak(b) {
		return a
	}
	return b
}

func stringLineBreak(l LineBreak) string {
	if isUnspecifiedLineBreak(l) {
		return "LineBreakUnspecified"
	}
	return l.String()
}

func sameLineBreak(a, b LineBreak) bool {
	return a == b
}

func semanticEqualLineBreak(a, b LineBreak) bool {
	return a == b
}

func equalLineBreak(a, b LineBreak) bool {
	return a == b
}

func copyLineBreak(l LineBreak) LineBreak {
	return l
}

// LineBreakCompanion is the single exported sentinel-family symbol for LineBreak.
var LineBreakCompanion = companions.NewScalarCompanion(
	func() LineBreak { return lineBreakUnspecified },
	stringLineBreak,
	isSpecifiedLineBreak,
	isUnspecifiedLineBreak,
	takeOrElseLineBreak,
	mergeLineBreak,
	sameLineBreak,
	semanticEqualLineBreak,
	equalLineBreak,
	copyLineBreak,
)
