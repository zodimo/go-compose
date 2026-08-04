package text

import "github.com/zodimo/go-compose/compose/ui/next/text/style"

// Alignment defines how to align text horizontally within its container.
// Re-exports the framework-owned style.TextAlign type.
type Alignment = style.TextAlign

const (
	// Start aligns the text on the leading edge of the container.
	Start Alignment = style.TextAlignStart
	// End aligns the text on the trailing edge of the container.
	End Alignment = style.TextAlignEnd
	// Middle aligns the text in the center of the container.
	Middle Alignment = style.TextAlignCenter
)

// WrapPolicy configures strategies for choosing where to break lines of text for line
// wrapping.
// Re-exports the framework-owned style.LineBreak type.
type WrapPolicy = style.LineBreak

var (
	// WrapHeuristically tries to minimize breaking within words (UAX#14 text segments)
	// while also ensuring that text fits within the given MaxWidth. It will only break
	// a line within a word (on a UAX#29 grapheme cluster boundary) when that word cannot
	// fit on a line by itself. Additionally, when the final word of a line is being
	// truncated, this policy will preserve as many symbols of that word as
	// possible before the truncator.
	WrapHeuristically = style.LineBreakParagraph
	// WrapWords does not permit words (UAX#14 text segments) to be broken across lines.
	// This means that sometimes long words will exceed the MaxWidth they are wrapped with.
	WrapWords = style.LineBreakHeading
	// WrapGraphemes will maximize the amount of text on each line at the expense of readability,
	// breaking any word across lines on UAX#29 grapheme cluster boundaries to maximize the number of
	// grapheme clusters on each line.
	WrapGraphemes = style.LineBreakSimple
)
