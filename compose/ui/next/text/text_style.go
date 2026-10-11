package text

import (
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/next/text/font"
	"github.com/zodimo/go-compose/compose/ui/next/text/intl"
	"github.com/zodimo/go-compose/compose/ui/next/text/style"
	"github.com/zodimo/go-compose/compose/ui/next/text/unit"
	"github.com/zodimo/go-sentinel-helper/sentinel/companions"
	"github.com/zodimo/go-sentinel-helper/sentinel/stringutils"
)

// textStyleUnspecified is the sentinel for TextStyle. It is deliberately
// package-private: the only public way to obtain it is TextStyleCompanion.Unspecified().
var textStyleUnspecified *TextStyle = &TextStyle{
	spanStyle:      spanStyleUnspecified,
	paragraphStyle: paragraphStyleUnspecified,
	platformStyle:  nil,
}

var _ TextStyleInterface = (*TextStyle)(nil)

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
	return TextStyleCompanion.IsSpecified(style)
}

func TakeOrElseTextStyle(style, defaultStyle *TextStyle) *TextStyle {
	return takeOrElseTextStyle(style, defaultStyle)
}

// Identity (2 ns)
func SameTextStyle(a, b *TextStyle) bool {
	return sameTextStyle(a, b)
}

// Semantic equality (field-by-field, 20 ns)
func SemanticEqualTextStyle(a, b *TextStyle) bool {
	return semanticEqualTextStyle(a, b)
}

func EqualTextStyle(a, b *TextStyle) bool {
	return equalTextStyle(a, b)
}

func MergeTextStyle(a, b *TextStyle) *TextStyle {
	return mergeTextStyle(a, b)
}

func CoalesceTextStyle(ptr, def *TextStyle) *TextStyle {
	return coalesceTextStyle(ptr, def)
}

// === Sentinel operations ===

func isSpecifiedTextStyle(style *TextStyle) bool {
	return style != nil && style != textStyleUnspecified
}

func isUnspecifiedTextStyle(style *TextStyle) bool {
	return !isSpecifiedTextStyle(style)
}

func takeOrElseTextStyle(style, defaultStyle *TextStyle) *TextStyle {
	if style == nil || style == textStyleUnspecified {
		return defaultStyle
	}
	return style
}

func coalesceTextStyle(ptr, def *TextStyle) *TextStyle {
	if ptr == nil {
		return def
	}
	return ptr
}

func mergeTextStyle(a, b *TextStyle) *TextStyle {
	a = coalesceTextStyle(a, textStyleUnspecified)
	b = coalesceTextStyle(b, textStyleUnspecified)

	if a == textStyleUnspecified {
		return b
	}
	if b == textStyleUnspecified {
		return a
	}

	// Both are custom: allocate new merged style. b wins where specified.
	return &TextStyle{
		spanStyle:      mergeSpanStyle(a.spanStyle, b.spanStyle),
		paragraphStyle: mergeParagraphStyle(a.paragraphStyle, b.paragraphStyle),
	}
}

func stringTextStyle(ts *TextStyle) string {
	if !isSpecifiedTextStyle(ts) {
		return "TextStyle{Unspecified}"
	}
	return "TextStyle{spanStyle=" + stringSpanStyle(ts.spanStyle) +
		", paragraphStyle=" + stringParagraphStyle(ts.paragraphStyle) + "}"
}

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

func semanticEqualTextStyle(a, b *TextStyle) bool {

	a = coalesceTextStyle(a, textStyleUnspecified)
	b = coalesceTextStyle(b, textStyleUnspecified)

	if a == textStyleUnspecified && b == textStyleUnspecified {
		return true
	}

	return semanticEqualSpanStyle(a.spanStyle, b.spanStyle) &&
		semanticEqualParagraphStyle(a.paragraphStyle, b.paragraphStyle)
}

func equalTextStyle(a, b *TextStyle) bool {
	if !sameTextStyle(a, b) {
		return semanticEqualTextStyle(a, b)
	}
	return true
}

// copyTextStyle returns a copy of ts; the source is never mutated.
func copyTextStyle(ts *TextStyle, options ...func(*TextStyle)) *TextStyle {
	cp := &TextStyle{}
	if isSpecifiedTextStyle(ts) {
		cp.spanStyle = ts.spanStyle
		cp.paragraphStyle = ts.paragraphStyle
		cp.platformStyle = ts.platformStyle
	}
	for _, option := range options {
		option(cp)
	}
	return cp
}

// TextStyleCompanion is the single exported sentinel-family symbol for TextStyle.
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

func TextStyleResolveDefaults(ts *TextStyle, direction unit.LayoutDirection) *TextStyle {
	ts = coalesceTextStyle(ts, textStyleUnspecified)

	return &TextStyle{
		spanStyle:      SpanStyleResolveDefaults(ts.spanStyle),
		paragraphStyle: ParagraphStyleResolveDefaults(ts.paragraphStyle, direction),
		platformStyle:  ts.platformStyle,
	}
}

func SpanStyleResolveDefaults(ss *SpanStyle) *SpanStyle {
	ss = coalesceSpanStyle(ss, spanStyleUnspecified)
	return &SpanStyle{
		textForegroundStyle:    style.TakeOrElseTextForegroundStyle(ss.textForegroundStyle, DefaultColorForegroundStyle),
		FontSize:               unit.TextUnitCompanion.TakeOrElse(ss.FontSize, DefaultFontSize),
		FontWeight:             font.FontWeightCompanion.TakeOrElse(ss.FontWeight, font.FontWeightNormal),
		FontStyle:              font.FontStyleCompanion.TakeOrElse(ss.FontStyle, font.FontStyleNormal),
		FontSynthesis:          font.TakeOrElseFontSynthesis(ss.FontSynthesis, font.FontSynthesisAll),
		FontFamily:             font.TakeOrElseFontFamily(ss.FontFamily, font.FontFamilyDefault),
		FontFeatureSettings:    stringutils.StringsCompanion.TakeOrElse(ss.FontFeatureSettings, ""),
		LetterSpacing:          unit.TextUnitCompanion.TakeOrElse(ss.LetterSpacing, DefaultLetterSpacing),
		BaselineShift:          style.TakeOrElseBaselineShift(ss.BaselineShift, style.BaselineShiftNone),
		TextGeometricTransform: style.TakeOrElseTextGeometricTransform(ss.TextGeometricTransform, style.TextGeometricTransformNone),
		LocaleList:             intl.LocaleListCompanion.TakeOrElse(ss.LocaleList, nil), //LocaleList.Current - local provider
		Background:             graphics.ColorCompanion.TakeOrElse(ss.Background, DefaultBackgroundColor),
		TextDecoration:         style.TakeOrElseTextDecoration(ss.TextDecoration, style.TextDecorationNone),
		Shadow:                 graphics.ShadowCompanion.TakeOrElse(ss.Shadow, graphics.ShadowNone),
		PlatformStyle:          ss.PlatformStyle,
		DrawStyle:              graphics.TakeOrElseDrawStyle(ss.DrawStyle, graphics.DrawStyleFill),
	}
}

func ParagraphStyleResolveDefaults(ps *ParagraphStyle, direction unit.LayoutDirection) *ParagraphStyle {
	ps = coalesceParagraphStyle(ps, paragraphStyleUnspecified)
	return &ParagraphStyle{
		TextAlign:       style.TextAlignCompanion.TakeOrElse(ps.TextAlign, style.TextAlignStart),
		TextDirection:   style.ResolveTextDirection(direction, ps.TextDirection),
		LineHeight:      unit.TextUnitCompanion.TakeOrElse(ps.LineHeight, unit.TextUnitCompanion.Unspecified()),
		TextIndent:      style.TakeOrElseTextIndent(ps.TextIndent, style.TextIndentNone),
		PlatformStyle:   ps.PlatformStyle,
		LineHeightStyle: ps.LineHeightStyle,
		LineBreak:       style.LineBreakCompanion.TakeOrElse(ps.LineBreak, style.LineBreakSimple),
		Hyphens:         style.HyphensCompanion.TakeOrElse(ps.Hyphens, style.HyphensNone),
		TextMotion:      style.TakeOrElseTextMotion(ps.TextMotion, style.TextMotionStatic),
	}
}
