package style

import "fmt"

// LineBreak configuration for line breaking.
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-text/src/commonMain/kotlin/androidx/compose/ui/text/style/LineBreak.kt
type LineBreak int32

// LineBreakStrategy The strategy used for line breaking.
type LineBreakStrategy int32

const (
	// LineBreakStrategySimple Basic, fast break strategy. Hyphenation, if enabled, is done only for words that
	// don't fit on an entire line by themselves.
	LineBreakStrategySimple LineBreakStrategy = 1

	// LineBreakStrategyHighQuality Does whole paragraph optimization for more readable text, including hyphenation if
	// enabled.
	LineBreakStrategyHighQuality LineBreakStrategy = 2

	// LineBreakStrategyBalanced Attempts to balance the line lengths of the text, also applying automatic hyphenation
	// if enabled. Suitable for small screens.
	LineBreakStrategyBalanced LineBreakStrategy = 3

	// LineBreakStrategyUnspecified This represents an unset value, a usual replacement for "null" when a primitive value
	// is desired.
	LineBreakStrategyUnspecified LineBreakStrategy = 0
)

func (s LineBreakStrategy) String() string {
	switch s {
	case LineBreakStrategySimple:
		return "Strategy.Simple"
	case LineBreakStrategyHighQuality:
		return "Strategy.HighQuality"
	case LineBreakStrategyBalanced:
		return "Strategy.Balanced"
	case LineBreakStrategyUnspecified:
		return "Strategy.Unspecified"
	default:
		return "Invalid"
	}
}

// LineBreakStrictness Describes the strictness of line breaking, determining before which characters line breaks
// can be inserted. It is useful when working with CJK scripts.
type LineBreakStrictness int32

const (
	// LineBreakStrictnessDefault Default breaking rules for the locale, which may correspond to [Normal] or [Strict].
	LineBreakStrictnessDefault LineBreakStrictness = 1

	// LineBreakStrictnessLoose The least restrictive rules, suitable for short lines.
	//
	// For example, in Japanese it allows breaking before iteration marks, such as 々, 〻.
	LineBreakStrictnessLoose LineBreakStrictness = 2

	// LineBreakStrictnessNormal The most common rules for line breaking.
	//
	// For example, in Japanese it allows breaking before characters like small hiragana
	// (ぁ), small katakana (ァ), halfwidth variants (ｧ).
	LineBreakStrictnessNormal LineBreakStrictness = 3

	// LineBreakStrictnessStrict The most stringent rules for line breaking.
	//
	// For example, in Japanese it does not allow breaking before characters like small
	// hiragana (ぁ), small katakana (ァ), halfwidth variants (ｧ).
	LineBreakStrictnessStrict LineBreakStrictness = 4

	// LineBreakStrictnessUnspecified This represents an unset value, a usual replacement for "null" when a primitive value
	// is desired.
	LineBreakStrictnessUnspecified LineBreakStrictness = 0
)

func (s LineBreakStrictness) String() string {
	switch s {
	case LineBreakStrictnessDefault:
		return "Strictness.None"
	case LineBreakStrictnessLoose:
		return "Strictness.Loose"
	case LineBreakStrictnessNormal:
		return "Strictness.Normal"
	case LineBreakStrictnessStrict:
		return "Strictness.Strict"
	case LineBreakStrictnessUnspecified:
		return "Strictness.Unspecified"
	default:
		return "Invalid"
	}
}

// LineBreakWordBreak Describes how line breaks should be inserted within words.
type LineBreakWordBreak int32

const (
	// LineBreakWordBreakDefault Default word breaking rules for the locale. In latin scripts this means inserting
	// line breaks between words, while in languages that don't use whitespace (e.g.
	// Japanese) the line can break between characters.
	LineBreakWordBreakDefault LineBreakWordBreak = 1

	// LineBreakWordBreakPhrase Line breaking is based on phrases. In languages that don't use whitespace (e.g.
	// Japanese), line breaks are not inserted between characters that are part of the same
	// phrase unit. This is ideal for short text such as titles and UI labels.
	LineBreakWordBreakPhrase LineBreakWordBreak = 2

	// LineBreakWordBreakUnspecified This represents an unset value, a usual replacement for "null" when a primitive value
	// is desired.
	LineBreakWordBreakUnspecified LineBreakWordBreak = 0
)

func (w LineBreakWordBreak) String() string {
	switch w {
	case LineBreakWordBreakDefault:
		return "WordBreak.None"
	case LineBreakWordBreakPhrase:
		return "WordBreak.Phrase"
	case LineBreakWordBreakUnspecified:
		return "WordBreak.Unspecified"
	default:
		return "Invalid"
	}
}

// Basic, fast line breaking. Ideal for text input fields, as it will cause minimal text
// reflow when editing.
var LineBreakSimple = LineBreakOf(LineBreakStrategySimple, LineBreakStrictnessNormal, LineBreakWordBreakDefault)

// Looser breaking rules, suitable for short text such as titles or narrow newspaper
// columns. For longer lines of text, use [Paragraph] for improved readability.
var LineBreakHeading = LineBreakOf(LineBreakStrategyBalanced, LineBreakStrictnessLoose, LineBreakWordBreakPhrase)

// Slower, higher quality line breaking for improved readability. Suitable for larger
// amounts of text.
var LineBreakParagraph = LineBreakOf(LineBreakStrategyHighQuality, LineBreakStrictnessStrict, LineBreakWordBreakDefault)

// This represents an unset value, a usual replacement for "null" when a primitive value is
// desired.
const LineBreakUnspecified LineBreak = 0

// LineBreakOf creates a LineBreak from individual components.
func LineBreakOf(strategy LineBreakStrategy, strictness LineBreakStrictness, wordBreak LineBreakWordBreak) LineBreak {
	return LineBreak(packBytes(int32(strategy), int32(strictness), int32(wordBreak)))
}

// Strategy returns the strategy used for line breaking.
func (l LineBreak) Strategy() LineBreakStrategy {
	return LineBreakStrategy(unpackByte1(int32(l)))
}

// Strictness returns the strictness of line breaking.
func (l LineBreak) Strictness() LineBreakStrictness {
	return LineBreakStrictness(unpackByte2(int32(l)))
}

// WordBreak returns how line breaks should be inserted within words.
func (l LineBreak) WordBreak() LineBreakWordBreak {
	return LineBreakWordBreak(unpackByte3(int32(l)))
}

// Copy creates a copy of this LineBreak with the specified fields replaced.
func (l LineBreak) Copy(
	strategy LineBreakStrategy,
	strictness LineBreakStrictness,
	wordBreak LineBreakWordBreak,
) LineBreak {
	return LineBreakOf(
		strategy,
		strictness,
		wordBreak,
	)
}

func (l LineBreak) String() string {
	return fmt.Sprintf("LineBreak(strategy=%s, strictness=%s, wordBreak=%s)", l.Strategy(), l.Strictness(), l.WordBreak())
}

// Packs 3 bytes represented as Integers into a single Integer.
//
// A byte can represent any value between 0 and 256.
//
// Only the 8 least significant bits of any given value are packed into the returned Integer.
func packBytes(i1, i2, i3 int32) int32 {
	return i1 | (i2 << 8) | (i3 << 16)
}

func unpackByte1(mask int32) int32 {
	return 0x000000FF & mask
}

func unpackByte2(mask int32) int32 {
	return 0x000000FF & (mask >> 8)
}

func unpackByte3(mask int32) int32 {
	return 0x000000FF & (mask >> 16)
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
