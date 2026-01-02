package text

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/text/font"
	"github.com/zodimo/go-compose/compose/ui/text/style"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

var TextStyleUnspecified *TextStyle = &TextStyle{
	spanStyle:      SpanStyleUnspecified,
	paragraphStyle: ParagraphStyleUnspecified,
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

func StringTextStyle(ts *TextStyle) string {
	if !IsSpecifiedTextStyle(ts) {
		return "TextStyle{Unspecified}"
	}
	return fmt.Sprintf("TextStyle{spanStyle: %s, paragraphStyle: %s}",
		StringSpanStyle(ts.spanStyle),
		StringParagraphStyle(ts.paragraphStyle),
	)
}

// ensureMutableTextStyle panics if ts is a sentinel value.
// Call this at the start of any function that mutates a TextStyle to fail-fast on misuse.
func ensureMutableTextStyle(ts *TextStyle) {
	if !IsSpecifiedTextStyle(ts) {
		panic("attempt to mutate sentinel TextStyleUnspecified; use CopyTextStyle first")
	}
}

func TextStyleResolveDefaults(ts *TextStyle, direction unit.LayoutDirection) *TextStyle {
	ts = CoalesceTextStyle(ts, TextStyleUnspecified)
	return &TextStyle{
		spanStyle:      SpanStyleResolveDefaults(ts.spanStyle),
		paragraphStyle: ParagraphStyleResolveDefaults(ts.paragraphStyle, direction),
	}
}

func IsSpecifiedTextStyle(s *TextStyle) bool {
	return s != nil && s != TextStyleUnspecified
}
func TakeOrElseTextStyle(s, def *TextStyle) *TextStyle {
	if !IsSpecifiedTextStyle(s) {
		return def
	}
	return s
}

func MergeTextStyle(a, b *TextStyle) *TextStyle {
	a = CoalesceTextStyle(a, TextStyleUnspecified)
	b = CoalesceTextStyle(b, TextStyleUnspecified)

	if a == TextStyleUnspecified {
		return b
	}
	if b == TextStyleUnspecified {
		return a
	}

	// Both are custom: allocate new merged style
	return &TextStyle{
		spanStyle:      MergeSpanStyle(a.spanStyle, b.spanStyle),
		paragraphStyle: MergeParagraphStyle(a.paragraphStyle, b.paragraphStyle),
	}
}

func CoalesceTextStyle(ptr, def *TextStyle) *TextStyle {
	if ptr == nil {
		return def
	}
	return ptr
}

func TextStyleFromOptions(options ...TextStyleOption) *TextStyle {
	return CopyTextStyle(TextStyleUnspecified, options...)
}

// Identity (2 ns)
func SameTextStyle(a, b *TextStyle) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return b == TextStyleUnspecified
	}
	if b == nil {
		return a == TextStyleUnspecified
	}
	return a == b
}

// Semantic equality (field-by-field, 20 ns)
func SemanticEqualTextStyle(a, b *TextStyle) bool {

	a = CoalesceTextStyle(a, TextStyleUnspecified)
	b = CoalesceTextStyle(b, TextStyleUnspecified)

	return EqualSpanStyle(a.spanStyle, b.spanStyle) &&
		EqualParagraphStyle(a.paragraphStyle, b.paragraphStyle)
}

func EqualTextStyle(a, b *TextStyle) bool {
	if !SameTextStyle(a, b) {
		return SemanticEqualTextStyle(a, b)
	}
	return true
}

func CopyTextStyle(ts *TextStyle, options ...TextStyleOption) *TextStyle {
	opt := *TextStyleUnspecified

	for _, option := range options {
		option(&opt)
	}
	return &TextStyle{
		spanStyle:      MergeSpanStyle(ts.spanStyle, opt.spanStyle),
		paragraphStyle: MergeParagraphStyle(ts.paragraphStyle, opt.paragraphStyle),
	}
}
