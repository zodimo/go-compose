package text

import (
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

func CopyTextStyle(ts *TextStyle, options ...TextStyleOption) *TextStyle {
	copy := *ts
	for _, option := range options {
		option(&copy)
	}
	return &copy
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
