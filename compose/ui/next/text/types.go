package text

import (
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/next/text/font"
	"github.com/zodimo/go-compose/compose/ui/next/text/intl"
	"github.com/zodimo/go-compose/compose/ui/next/text/style"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

type TextStyleInterface interface {
	Alpha() float32
	Background() graphics.Color
	BaselineShift() style.BaselineShift
	Brush() graphics.Brush
	Color() graphics.Color
	Copy() *TextStyle
	DrawStyle() graphics.DrawStyle
	FontFamily() font.FontFamily
	FontFeatureSettings() string
	FontSize() unit.TextUnit
	FontStyle() font.FontStyle
	FontSynthesis() *font.FontSynthesis
	FontWeight() font.FontWeight
	Hyphens() style.Hyphens
	LetterSpacing() unit.TextUnit
	LineBreak() style.LineBreak
	LineHeight() unit.TextUnit
	LineHeightStyle() *style.LineHeightStyle
	LocaleList() *intl.LocaleList
	MergeParagraphStyle(other *ParagraphStyle) *TextStyle
	MergeSpanStyle(other *SpanStyle) *TextStyle
	Plus(other *TextStyle) *TextStyle
	PlusParagraphStyle(other *ParagraphStyle) *TextStyle
	PlusSpanStyle(other *SpanStyle) *TextStyle
	Shadow() *graphics.Shadow
	TextAlign() style.TextAlign
	TextDecoration() *style.TextDecoration
	TextDirection() style.TextDirection
	TextGeometricTransform() *style.TextGeometricTransform
	TextIndent() *style.TextIndent
	TextMotion() *style.TextMotion
	ToParagraphStyle() *ParagraphStyle
	ToPlatformTextStyle() *PlatformTextStyle
	ToSpanStyle() *SpanStyle
	ToString() string
}
