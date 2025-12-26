package text

import (
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/text/font"
	"github.com/zodimo/go-compose/compose/ui/text/intl"
	"github.com/zodimo/go-compose/compose/ui/text/style"
)

var _ Annotation = (*SpanStyle)(nil)

type SpanStyle struct {
	textForegroundStyle    style.TextForegroundStyle
	FontSize               FontSize
	FontWeight             font.FontWeight
	FontStyle              font.FontStyle
	FontSynthesis          *font.FontSynthesis
	FontFamily             font.FontFamily
	FontFeatureSettings    string
	LetterSpacing          uiTextUnit
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
	return s.textForegroundStyle.Color()
}
func (s SpanStyle) Brush() graphics.Brush {
	return s.textForegroundStyle.Brush()
}
func (s SpanStyle) Alpha() float32 {
	return s.textForegroundStyle.Alpha()
}

func (s SpanStyle) Merge(other *SpanStyle) *SpanStyle {
	panic("SpanStyle Merge not implemented")
}

func (s SpanStyle) Copy() *SpanStyle {
	panic("SpanStyle Copy not implemented")
}

func (s SpanStyle) ToString() string {
	panic("SpanStyle ToString not implemented")
}
