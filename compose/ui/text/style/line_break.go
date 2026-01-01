package style

import (
	"fmt"

	gioText "gioui.org/text"
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

// This represents an unset value, a usual replacement for "null" when a primitive value is
// desired.
const LineBreakUnspecified LineBreak = 99

func (l LineBreak) String() string {
	var description string
	switch l {
	case LineBreakSimple:
		description = "LineBreakSimple/WrapGraphemes"
	case LineBreakHeading:
		description = "LineBreakHeading/WrapWords"
	case LineBreakParagraph:
		description = "LineBreakParagraph/WrapHeuristically"
	case LineBreakUnspecified:
		description = "LineBreakUnspecified"
	default:
		panic(fmt.Sprintf("unknown line break strategy: %d", l))
	}
	return fmt.Sprintf("LineBreak(strategy=%s)", description)
}

// IsSpecified returns true if it is not [LineBreakUnspecified].
func (l LineBreak) IsSpecified() bool {
	return l != LineBreakUnspecified
}

func (l LineBreak) TakeOrElse(other LineBreak) LineBreak {
	if l.IsSpecified() {
		return l
	}
	return other
}
