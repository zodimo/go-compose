package text

import (
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/text/font"
	"github.com/zodimo/go-compose/compose/ui/text/style"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

var TextStyleUnspecified *TextStyle = &TextStyle{
	lineBreak: style.LineBreakUnspecified,
}

var _ TextStyleInterface = (*TextStyle)(nil)

type TextStyle struct {
	font       font.Font
	fontSize   unit.TextUnit
	fontWeight font.FontWeight
	fontStyle  font.FontStyle

	fontFamily font.FontFamily

	textDecoration *style.TextDecoration
	letterSpacing  unit.TextUnit

	textAlign     style.TextAlign
	textDirection style.TextDirection

	color      graphics.Color
	background graphics.Color
	lineHeight unit.TextUnit
	lineBreak  style.LineBreak
}

func (ts TextStyle) Alpha() float32 {
	return ts.color.Alpha()
}
func (ts TextStyle) Background() graphics.Color {
	return ts.background
}
func (ts TextStyle) Color() graphics.Color {
	return ts.color
}
func (ts TextStyle) FontFamily() font.FontFamily {
	return ts.fontFamily
}
func (ts TextStyle) FontFeatureSettings() string {
	panic("FontFeatureSettings not implemented")
}
func (ts TextStyle) FontSize() unit.TextUnit {
	return ts.fontSize
}
func (ts TextStyle) FontStyle() font.FontStyle {
	return ts.fontStyle
}
func (ts TextStyle) FontSynthesis() *font.FontSynthesis {
	panic("FontSynthesis not implemented")
}
func (ts TextStyle) FontWeight() font.FontWeight {
	return ts.fontWeight
}
func (ts TextStyle) LetterSpacing() unit.TextUnit {
	return ts.letterSpacing
}

func (ts TextStyle) LineBreak() style.LineBreak {
	return ts.lineBreak
}
func (ts TextStyle) LineHeight() unit.TextUnit {
	return ts.lineHeight
}
func (ts TextStyle) TextAlign() style.TextAlign {
	return ts.textAlign
}
func (ts TextStyle) TextDecoration() *style.TextDecoration {
	return ts.textDecoration
}
func (ts TextStyle) TextDirection() style.TextDirection {
	return ts.textDirection
}
func (ts TextStyle) ToString() string {
	panic("ToString not implemented")
}

func TextStyleResolveDefaults(ts *TextStyle, direction unit.LayoutDirection) *TextStyle {
	ts = CoalesceTextStyle(ts, TextStyleUnspecified)
	return &TextStyle{

		fontSize:   ts.FontSize().TakeOrElse(DefaultFontSize),
		fontWeight: ts.FontWeight().TakeOrElse(font.FontWeightNormal),
		fontStyle:  ts.FontStyle().TakeOrElse(font.FontStyleNormal),
		// FontSynthesis:          font.TakeOrElseFontSynthesis(ss.FontSynthesis, font.FontSynthesisAll),
		fontFamily: font.TakeOrElseFontFamily(ts.FontFamily(), font.FontFamilyDefault),
		// FontFeatureSettings:    sentinel.TakeOrElseString(ss.FontFeatureSettings, ""),
		letterSpacing: ts.LetterSpacing().TakeOrElse(DefaultLetterSpacing),
		// BaselineShift:          style.TakeOrElseBaselineShift(ss.BaselineShift, style.BaselineShiftNone),
		// TextGeometricTransform: style.TakeOrElseTextGeometricTransform(ss.TextGeometricTransform, style.TextGeometricTransformNone),
		// LocaleList:             intl.TakeOrElseLocaleList(ss.LocaleList, nil), //LocaleList.Current - local provider
		background:     ts.Background().TakeOrElse(DefaultBackgroundColor),
		textDecoration: style.TakeOrElseTextDecoration(ts.TextDecoration(), style.TextDecorationNone),
		// Shadow:                 graphics.TakeOrElseShadow(ss.Shadow, graphics.ShadowNone),
		// PlatformStyle:          ss.PlatformStyle,
		// DrawStyle:              graphics.TakeOrElseDrawStyle(ss.DrawStyle, graphics.DrawStyleFill),

		textAlign:     ts.TextAlign().TakeOrElse(style.TextAlignStart),
		textDirection: style.ResolveTextDirection(direction, ts.TextDirection()),
		lineHeight:    ts.LineHeight().TakeOrElse(unit.TextUnitUnspecified),
		// TextIndent:      style.TakeOrElseTextIndent(ps.TextIndent, style.TextIndentNone),
		// PlatformStyle:   ps.PlatformStyle,
		// LineHeightStyle: ps.LineHeightStyle,
		lineBreak: ts.LineBreak().TakeOrElse(style.LineBreakSimple),
		// Hyphens:         ps.Hyphens.TakeOrElse(style.HyphensNone),
		// TextMotion:      style.TakeOrElseTextMotion(ps.TextMotion, style.TextMotionStatic),
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

func CoalesceTextStyle(ptr, def *TextStyle) *TextStyle {
	if ptr == nil {
		return def
	}
	return ptr
}
