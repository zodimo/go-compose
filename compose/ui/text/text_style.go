package text

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/text/font"
	"github.com/zodimo/go-compose/compose/ui/text/style"
	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-sentinel-helper/sentinel/companions"
)

// textStyleUnspecified is the sentinel for TextStyle. It is deliberately
// package-private: the only public way to obtain it is TextStyleCompanion.Unspecified().
var textStyleUnspecified *TextStyle = &TextStyle{
	spanStyle:      spanStyleUnspecified,
	paragraphStyle: paragraphStyleUnspecified,
}

var _ TextStyleInterface = (*TextStyle)(nil)

type TextStyle struct {
	spanStyle      *SpanStyle
	paragraphStyle *ParagraphStyle
}

func (ts TextStyle) Alpha() float32 {
	return ts.spanStyle.Color().Alpha()
}
func (ts TextStyle) Background() graphics.Color {
	return ts.spanStyle.Background()
}
func (ts TextStyle) Color() graphics.Color {
	return ts.spanStyle.Color()
}
func (ts TextStyle) FontFamily() font.FontFamily {
	return ts.spanStyle.FontFamily()
}
func (ts TextStyle) FontSize() unit.TextUnit {
	return ts.spanStyle.FontSize()
}
func (ts TextStyle) FontStyle() font.FontStyle {
	return ts.spanStyle.FontStyle()
}
func (ts TextStyle) FontWeight() font.FontWeight {
	return ts.spanStyle.FontWeight()
}
func (ts TextStyle) LetterSpacing() unit.TextUnit {
	return ts.spanStyle.LetterSpacing()
}

func (ts TextStyle) LineBreak() style.LineBreak {
	return ts.paragraphStyle.LineBreak()
}
func (ts TextStyle) LineHeight() unit.TextUnit {
	return ts.paragraphStyle.LineHeight()
}
func (ts TextStyle) TextAlign() style.TextAlign {
	return ts.paragraphStyle.TextAlign()
}
func (ts TextStyle) TextDecoration() *style.TextDecoration {
	return ts.spanStyle.TextDecoration()
}
func (ts TextStyle) TextDirection() style.TextDirection {
	return ts.paragraphStyle.TextDirection()
}

func stringTextStyle(ts *TextStyle) string {
	localTs := coalesceTextStyle(ts, textStyleUnspecified)
	if !isSpecifiedTextStyle(localTs) {
		return "TextStyle{Unspecified}"
	}
	return fmt.Sprintf("TextStyle{spanStyle: %s, paragraphStyle: %s}",
		stringSpanStyle(localTs.spanStyle),
		stringParagraphStyle(localTs.paragraphStyle),
	)
}

func textStyleResolveDefaults(ts *TextStyle, direction unit.LayoutDirection) *TextStyle {
	localTs := coalesceTextStyle(ts, textStyleUnspecified)
	return &TextStyle{
		spanStyle:      spanStyleResolveDefaults(localTs.spanStyle),
		paragraphStyle: paragraphStyleResolveDefaults(localTs.paragraphStyle, direction),
	}
}

func isSpecifiedTextStyle(s *TextStyle) bool {
	return s != nil && s != textStyleUnspecified
}

func isUnspecifiedTextStyle(s *TextStyle) bool {
	return !isSpecifiedTextStyle(s)
}

func takeOrElseTextStyle(s, def *TextStyle) *TextStyle {
	if !isSpecifiedTextStyle(s) {
		return def
	}
	return s
}

// mergeTextStyle layers b onto a: a value from b wins only where b actually set it.
func mergeTextStyle(a, b *TextStyle) *TextStyle {

	a = coalesceTextStyle(a, textStyleUnspecified)
	b = coalesceTextStyle(b, textStyleUnspecified)

	if a == textStyleUnspecified {
		return b
	}
	if b == textStyleUnspecified {
		return a
	}

	// Both are custom: allocate new merged style
	return &TextStyle{
		spanStyle:      mergeSpanStyle(a.spanStyle, b.spanStyle),
		paragraphStyle: mergeParagraphStyle(a.paragraphStyle, b.paragraphStyle),
	}
}

func coalesceTextStyle(ptr, def *TextStyle) *TextStyle {
	if ptr == nil {
		return def
	}
	return ptr
}

// copyTextStyle returns a copy of ts with the given bare options applied on top.
// Options override; every field they did not touch keeps ts's value. The source
// is never mutated, so a sentinel source is safe.

func copyTextStyle(ts *TextStyle, options ...TextStyleOption) *TextStyle {

	ts = coalesceTextStyle(ts, textStyleUnspecified)
	opt := *textStyleUnspecified

	for _, option := range options {
		option(&opt)
	}
	return &TextStyle{
		spanStyle:      mergeSpanStyle(ts.spanStyle, opt.spanStyle),
		paragraphStyle: mergeParagraphStyle(ts.paragraphStyle, opt.paragraphStyle),
	}
}

// ResolveTextStyleDefaults resolves an unspecified TextStyle against the framework defaults.
func ResolveTextStyleDefaults(ts *TextStyle, direction unit.LayoutDirection) *TextStyle {
	return textStyleResolveDefaults(ts, direction)
}

// ResolveSpanStyleDefaults resolves an unspecified SpanStyle against the framework defaults.
func ResolveSpanStyleDefaults(s *SpanStyle) *SpanStyle {
	return spanStyleResolveDefaults(s)
}

// ResolveParagraphStyleDefaults resolves an unspecified ParagraphStyle against the framework defaults.
func ResolveParagraphStyleDefaults(s *ParagraphStyle, direction unit.LayoutDirection) *ParagraphStyle {
	return paragraphStyleResolveDefaults(s, direction)
}

// Identity (2 ns)
func sameTextStyle(a, b *TextStyle) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return b == textStyleUnspecified
	}
	if b == nil {
		return a == textStyleUnspecified
	}
	return a == b
}

// Semantic equality (field-by-field, 20 ns)
func semanticEqualTextStyle(a, b *TextStyle) bool {

	a = coalesceTextStyle(a, textStyleUnspecified)
	b = coalesceTextStyle(b, textStyleUnspecified)

	return equalSpanStyle(a.spanStyle, b.spanStyle) &&
		equalParagraphStyle(a.paragraphStyle, b.paragraphStyle)
}

func equalTextStyle(a, b *TextStyle) bool {
	if sameTextStyle(a, b) {
		return true
	}
	return semanticEqualTextStyle(a, b)
}

// TextStyleCompanion is the only exported value of this file's sentinel family member.
var TextStyleCompanion = companions.NewStructCompanion(
	func() *TextStyle { return textStyleUnspecified },
	stringTextStyle,
	isSpecifiedTextStyle,
	isUnspecifiedTextStyle,
	takeOrElseTextStyle,
	mergeTextStyle,
	sameTextStyle,
	semanticEqualTextStyle,
	equalTextStyle,
	coalesceTextStyle,
	copyTextStyle,
)

// ensureMutableTextStyle panics if ts is a sentinel value.
// Call this at the start of any function that mutates a TextStyle to fail-fast on misuse.
func ensureMutableTextStyle(s *TextStyle) {
	if TextStyleCompanion.IsUnspecified(s) {
		panic("attempt to mutate sentinel; use XCompanion.Copy first")
	}
}

// TextStyleFromOptions builds a TextStyle from options applied over the sentinel.
func TextStyleFromOptions(options ...TextStyleOption) *TextStyle {
	return TextStyleCompanion.Copy(textStyleUnspecified, options...)
}
