package text

import (
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/text/font"
	"github.com/zodimo/go-compose/compose/ui/text/intl"
	"github.com/zodimo/go-compose/compose/ui/text/style"
	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/pkg/sentinel"
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

func (s SpanStyle) isAnnotation() {
}

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

type SpanStyleOptions struct {
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

func SpanStyleWithColor(color graphics.Color) SpanStyleOption {
	return func(opts *SpanStyleOptions) {
		opts.textForegroundStyle = &style.TextForegroundStyle{
			Color: color,
		}
	}
}

func SpanStyleWithBrush(brush graphics.Brush, alpha float32) SpanStyleOption {
	return func(opts *SpanStyleOptions) {
		opts.textForegroundStyle = &style.TextForegroundStyle{
			Brush: brush,
			Alpha: alpha,
		}
	}
}
func SpanStyleWithFontSize(fontSize unit.TextUnit) SpanStyleOption {
	return func(opts *SpanStyleOptions) {
		opts.FontSize = fontSize
	}
}
func SpanStyleWithFontWeight(fontWeight font.FontWeight) SpanStyleOption {
	return func(opts *SpanStyleOptions) {
		opts.FontWeight = fontWeight
	}
}
func SpanStyleWithFontStyle(fontStyle font.FontStyle) SpanStyleOption {
	return func(opts *SpanStyleOptions) {
		opts.FontStyle = fontStyle
	}
}
func SpanStyleWithFontSynthesis(fontSynthesis *font.FontSynthesis) SpanStyleOption {
	return func(opts *SpanStyleOptions) {
		opts.FontSynthesis = fontSynthesis
	}
}
func SpanStyleWithFontFamily(fontFamily font.FontFamily) SpanStyleOption {
	return func(opts *SpanStyleOptions) {
		opts.FontFamily = fontFamily
	}
}
func SpanStyleWithFontFeatureSettings(fontFeatureSettings string) SpanStyleOption {
	return func(opts *SpanStyleOptions) {
		opts.FontFeatureSettings = fontFeatureSettings
	}
}
func SpanStyleWithLetterSpacing(letterSpacing unit.TextUnit) SpanStyleOption {
	return func(opts *SpanStyleOptions) {
		opts.LetterSpacing = letterSpacing
	}
}
func SpanStyleWithBaselineShift(baselineShift style.BaselineShift) SpanStyleOption {
	return func(opts *SpanStyleOptions) {
		opts.BaselineShift = baselineShift
	}
}
func SpanStyleWithTextGeometricTransform(textGeometricTransform *style.TextGeometricTransform) SpanStyleOption {
	return func(opts *SpanStyleOptions) {
		opts.TextGeometricTransform = textGeometricTransform
	}
}
func SpanStyleWithLocaleList(localeList *intl.LocaleList) SpanStyleOption {
	return func(opts *SpanStyleOptions) {
		opts.LocaleList = localeList
	}
}
func SpanStyleWithBackground(background graphics.Color) SpanStyleOption {
	return func(opts *SpanStyleOptions) {
		opts.Background = background
	}
}
func SpanStyleWithTextDecoration(textDecoration *style.TextDecoration) SpanStyleOption {
	return func(opts *SpanStyleOptions) {
		opts.TextDecoration = textDecoration
	}
}
func SpanStyleWithShadow(shadow *graphics.Shadow) SpanStyleOption {
	return func(opts *SpanStyleOptions) {
		opts.Shadow = shadow
	}
}
func SpanStyleWithPlatformStyle(platformStyle *PlatformSpanStyle) SpanStyleOption {
	return func(opts *SpanStyleOptions) {
		opts.PlatformStyle = platformStyle
	}
}
func SpanStyleWithDrawStyle(drawStyle graphics.DrawStyle) SpanStyleOption {
	return func(opts *SpanStyleOptions) {
		opts.DrawStyle = drawStyle
	}
}

type SpanStyleOption = func(*SpanStyleOptions)

func (s SpanStyle) Copy(options ...SpanStyleOption) *SpanStyle {
	opts := &SpanStyleOptions{}
	for _, option := range options {
		option(opts)
	}
	return &SpanStyle{
		textForegroundStyle:    style.TakeOrElseTextForegroundStyle(opts.textForegroundStyle, s.textForegroundStyle),
		FontSize:               opts.FontSize.TakeOrElse(s.FontSize),
		FontWeight:             opts.FontWeight.TakeOrElse(s.FontWeight),
		FontStyle:              opts.FontStyle.TakeOrElse(s.FontStyle),
		FontSynthesis:          font.TakeOrElseFontSynthesis(opts.FontSynthesis, s.FontSynthesis),
		FontFamily:             font.TakeOrElseFontFamily(opts.FontFamily, s.FontFamily),
		FontFeatureSettings:    sentinel.TakeOrElseString(opts.FontFeatureSettings, s.FontFeatureSettings),
		LetterSpacing:          opts.LetterSpacing.TakeOrElse(s.LetterSpacing),
		BaselineShift:          style.TakeOrElseBaselineShift(opts.BaselineShift, s.BaselineShift),
		TextGeometricTransform: style.TakeOrElseTextGeometricTransform(opts.TextGeometricTransform, s.TextGeometricTransform),
		LocaleList:             intl.TakeOrElseLocaleList(opts.LocaleList, s.LocaleList),
		Background:             opts.Background.TakeOrElse(s.Background),
		TextDecoration:         style.TakeOrElseTextDecoration(opts.TextDecoration, s.TextDecoration),
		Shadow:                 graphics.TakeOrElseShadow(opts.Shadow, s.Shadow),
		PlatformStyle:          TakeOrElsePlatformSpanStyle(opts.PlatformStyle, s.PlatformStyle),
		DrawStyle:              graphics.TakeOrElseDrawStyle(opts.DrawStyle, s.DrawStyle),
	}
}

func StringSpanStyle(s *SpanStyle) string {
	panic("SpanStyle ToString not implemented")
}

func IsSpecifiedSpanStyle(s *SpanStyle) bool {
	return s != nil && s != SpanStyleUnspecified
}

func TakeOrElseSpanStyle(s, def *SpanStyle) *SpanStyle {
	if !IsSpecifiedSpanStyle(s) {
		return def
	}
	return s
}

// Identity (2 ns)
func SameSpanStyle(a, b *SpanStyle) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return b == SpanStyleUnspecified
	}
	if b == nil {
		return a == SpanStyleUnspecified
	}
	return a == b
}

// Semantic equality (field-by-field, 20 ns)
func SemanticEqualSpanStyle(a, b *SpanStyle) bool {

	a = CoalesceSpanStyle(a, SpanStyleUnspecified)
	b = CoalesceSpanStyle(b, SpanStyleUnspecified)

	return style.EqualTextForegroundStyle(a.textForegroundStyle, b.textForegroundStyle) &&
		a.FontSize == b.FontSize &&
		a.FontWeight == b.FontWeight &&
		a.FontStyle == b.FontStyle &&
		font.EqualFontSynthesis(a.FontSynthesis, b.FontSynthesis) &&
		font.EqualFontFamily(a.FontFamily, b.FontFamily) &&
		a.FontFeatureSettings == b.FontFeatureSettings &&
		a.LetterSpacing == b.LetterSpacing &&
		a.BaselineShift == b.BaselineShift &&
		style.EqualTextGeometricTransform(a.TextGeometricTransform, b.TextGeometricTransform) &&
		intl.EqualLocaleList(a.LocaleList, b.LocaleList) &&
		a.Background == b.Background &&
		style.EqualTextDecoration(a.TextDecoration, b.TextDecoration) &&
		graphics.EqualShadow(a.Shadow, b.Shadow) &&
		EqualPlatformSpanStyle(a.PlatformStyle, b.PlatformStyle) &&
		graphics.EqualDrawStyle(a.DrawStyle, b.DrawStyle)
}

func EqualSpanStyle(a, b *SpanStyle) bool {
	if !SameSpanStyle(a, b) {
		return SemanticEqualSpanStyle(a, b)
	}
	return true
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
