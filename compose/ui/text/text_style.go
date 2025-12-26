package text

import (
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/text/font"
	"github.com/zodimo/go-compose/compose/ui/text/intl"
	"github.com/zodimo/go-compose/compose/ui/text/style"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

var TextStyleUnspecified *TextStyle = &TextStyle{
	spanStyle:      SpanStyleUnspecified,
	paragraphStyle: ParagraphStyleUnspecified,
	platformStyle:  nil,
}

type TextStyle struct {
	spanStyle      *SpanStyle
	paragraphStyle *ParagraphStyle
	platformStyle  *PlatformTextStyle
}

func (ts TextStyle) ToSpanStyle() *SpanStyle {
	return ts.spanStyle
}

func (ts TextStyle) ToParagraphStyle() *ParagraphStyle {
	return ts.paragraphStyle
}

func (ts TextStyle) ToPlatformTextStyle() *PlatformTextStyle {
	return ts.platformStyle
}

func (ts TextStyle) MergeSpanStyle(other *SpanStyle) *TextStyle {
	return &TextStyle{
		spanStyle:      MergeSpanStyle(ts.spanStyle, other),
		paragraphStyle: ts.paragraphStyle,
	}
}

func (ts TextStyle) MergeParagraphStyle(other *ParagraphStyle) *TextStyle {
	return &TextStyle{
		spanStyle:      ts.spanStyle,
		paragraphStyle: MergeParagraphStyle(ts.paragraphStyle, other),
	}
}

func (ts *TextStyle) Plus(other *TextStyle) *TextStyle {
	//breaking the rules here
	return MergeTextStyle(ts, other)
}

func (ts TextStyle) PlusSpanStyle(other *SpanStyle) *TextStyle {
	return ts.MergeSpanStyle(other)
}

func (ts TextStyle) PlusParagraphStyle(other *ParagraphStyle) *TextStyle {
	return ts.MergeParagraphStyle(other)
}

func (ts TextStyle) Copy() *TextStyle {
	panic("TextStyle copy not implemented")
}

// Props on Kotlin TextStyle Class

func (ts TextStyle) Brush() graphics.Brush {
	return (*ts.spanStyle).Brush()
}
func (ts TextStyle) Color() graphics.Color {
	return (*ts.spanStyle).Color()
}

func (ts TextStyle) Alpha() float32 {
	return (*ts.spanStyle).Alpha()
}

func (ts TextStyle) FontSize() unit.TextUnit {
	return ts.spanStyle.FontSize
}

func (ts TextStyle) FontWeight() font.FontWeight {
	return ts.spanStyle.FontWeight
}

func (ts TextStyle) FontStyle() font.FontStyle {
	return ts.spanStyle.FontStyle
}

func (ts TextStyle) FontSynthesis() *font.FontSynthesis {
	return ts.spanStyle.FontSynthesis
}

func (ts TextStyle) FontFamily() font.FontFamily {
	return ts.spanStyle.FontFamily
}

func (ts TextStyle) FontFeatureSettings() string {
	return ts.spanStyle.FontFeatureSettings
}

func (ts TextStyle) LetterSpacing() unit.TextUnit {
	return ts.spanStyle.LetterSpacing
}

func (ts TextStyle) BaselineShift() style.BaselineShift {
	return ts.spanStyle.BaselineShift
}

func (ts TextStyle) TextGeometricTransform() *style.TextGeometricTransform {
	return ts.spanStyle.TextGeometricTransform
}

func (ts TextStyle) LocaleList() *intl.LocaleList {
	return ts.spanStyle.LocaleList
}

func (ts TextStyle) Background() graphics.Color {
	return ts.spanStyle.Background
}

func (ts TextStyle) TextDecoration() *style.TextDecoration {
	return ts.spanStyle.TextDecoration
}

func (ts TextStyle) Shadow() *graphics.Shadow {
	return ts.spanStyle.Shadow
}

func (ts TextStyle) DrawStyle() graphics.DrawStyle {
	return ts.spanStyle.DrawStyle
}

func (ts TextStyle) TextAlign() style.TextAlign {
	return ts.paragraphStyle.TextAlign
}

func (ts TextStyle) TextDirection() style.TextDirection {
	return ts.paragraphStyle.TextDirection
}

func (ts TextStyle) LineHeight() unit.TextUnit {
	return ts.paragraphStyle.LineHeight
}

func (ts TextStyle) TextIndent() *style.TextIndent {
	return ts.paragraphStyle.TextIndent
}

func (ts TextStyle) LineHeightStyle() *style.LineHeightStyle {
	return ts.paragraphStyle.LineHeightStyle
}

func (ts TextStyle) Hyphens() style.Hyphens {
	return ts.paragraphStyle.Hyphens
}

func (ts TextStyle) LineBreak() style.LineBreak {
	return ts.paragraphStyle.LineBreak
}

func (ts TextStyle) TextMotion() *style.TextMotion {
	return ts.paragraphStyle.TextMotion
}

func (s TextStyle) ToString() string {
	panic("TextStyle ToString not implemented")
}

func IsSpecifiedTextStyle(style *TextStyle) bool {
	return style != nil && style != TextStyleUnspecified
}

func TakeOrElseTextStyle(style, defaultStyle *TextStyle) *TextStyle {
	if style == nil || style == TextStyleUnspecified {
		return defaultStyle
	}
	return style
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

	return SemanticEqualSpanStyle(a.spanStyle, b.spanStyle) &&
		SemanticEqualParagraphStyle(a.paragraphStyle, b.paragraphStyle)
}

func EqualTextStyle(a, b *TextStyle) bool {
	if !SameTextStyle(a, b) {
		return SemanticEqualTextStyle(a, b)
	}
	return true
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
