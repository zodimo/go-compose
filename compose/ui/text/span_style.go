package text

import (
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/text/font"
	"github.com/zodimo/go-compose/compose/ui/text/intl"
	"github.com/zodimo/go-compose/compose/ui/text/style"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

var SpanStyleUnspecified *SpanStyle = &SpanStyle{
	textForegroundStyle:    nil,
	FontSize:               unit.TextUnitUnspecified,
	FontWeight:             font.FontWeightUnspecified,
	FontStyle:              font.FontStyleUnspecified,
	FontSynthesis:          nil,
	FontFamily:             nil,
	FontFeatureSettings:    "",
	LetterSpacing:          unit.TextUnitUnspecified,
	BaselineShift:          style.BaselineShiftUnspecified,
	TextGeometricTransform: nil,
	LocaleList:             nil,
	Background:             graphics.ColorUnspecified,
	TextDecoration:         nil,
	Shadow:                 nil,
	PlatformStyle:          nil,
	DrawStyle:              nil,
}

var _ Annotation = (*SpanStyle)(nil)

type SpanStyle struct {
	textForegroundStyle    *style.TextForegroundStyle
	FontSize               unit.TextUnit
	FontWeight             font.FontWeight
	FontStyle              font.FontStyle
	FontSynthesis          *font.FontSynthesis
	FontFamily             font.FontFamily
	FontFeatureSettings    string
	LetterSpacing          unit.TextUnit
	BaselineShift          style.BaselineShift
	TextGeometricTransform *style.TextGeometricTransform
	LocaleList             *intl.LocaleList
	Background             graphics.Color
	TextDecoration         *style.TextDecoration
	Shadow                 *graphics.Shadow
	PlatformStyle          *PlatformSpanStyle
	DrawStyle              graphics.DrawStyle
}

func (s SpanStyle) isAnnotation() {}

// Props
func (s SpanStyle) Color() graphics.Color {
	return s.textForegroundStyle.Color
}
func (s SpanStyle) Brush() graphics.Brush {
	return s.textForegroundStyle.Brush
}
func (s SpanStyle) Alpha() float32 {
	return s.textForegroundStyle.Alpha
}

func (s SpanStyle) Copy() *SpanStyle {
	panic("SpanStyle Copy not implemented")
}

func (s SpanStyle) ToString() string {
	panic("SpanStyle ToString not implemented")
}

func MergeSpanStyle(a, b *SpanStyle) *SpanStyle {
	a = CoalesceSpanStyle(a, SpanStyleUnspecified)
	b = CoalesceSpanStyle(b, SpanStyleUnspecified)

	if a == SpanStyleUnspecified {
		return b
	}
	if b == SpanStyleUnspecified {
		return a
	}

	// Both are custom: allocate new merged style
	return &SpanStyle{
		textForegroundStyle: style.MergeTextForegroundStyle(a.textForegroundStyle, b.textForegroundStyle),
	}
}

func CoalesceSpanStyle(ptr, def *SpanStyle) *SpanStyle {
	if ptr == nil {
		return def
	}
	return ptr
}
