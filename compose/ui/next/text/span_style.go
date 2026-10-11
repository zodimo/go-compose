package text

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/next/text/font"
	"github.com/zodimo/go-compose/compose/ui/next/text/intl"
	"github.com/zodimo/go-compose/compose/ui/next/text/style"
	"github.com/zodimo/go-compose/compose/ui/next/text/unit"
	"github.com/zodimo/go-compose/pkg/floatutils/lerp"
	"github.com/zodimo/go-sentinel-helper/sentinel/companions"
	"github.com/zodimo/go-sentinel-helper/sentinel/stringutils"
)

// spanStyleUnspecified is the sentinel for SpanStyle. It is deliberately
// package-private: the only public way to obtain it is SpanStyleCompanion.Unspecified().
var spanStyleUnspecified *SpanStyle = &SpanStyle{
	textForegroundStyle:    nil,
	FontSize:               unit.TextUnitCompanion.Unspecified(),
	FontWeight:             font.FontWeightCompanion.Unspecified(),
	FontStyle:              font.FontStyleCompanion.Unspecified(),
	FontSynthesis:          nil,
	FontFamily:             nil,
	FontFeatureSettings:    "",
	LetterSpacing:          unit.TextUnitCompanion.Unspecified(),
	BaselineShift:          style.BaselineShiftCompanion.Unspecified(),
	TextGeometricTransform: nil,
	LocaleList:             nil,
	Background:             graphics.ColorCompanion.Unspecified(),
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
	return style.CoalesceTextForegroundStyle(s.textForegroundStyle, style.TextForegroundStyleCompanion.Unspecified()).Color
}
func (s SpanStyle) Brush() graphics.Brush {
	return style.CoalesceTextForegroundStyle(s.textForegroundStyle, style.TextForegroundStyleCompanion.Unspecified()).Brush
}
func (s SpanStyle) Alpha() float32 {
	return style.CoalesceTextForegroundStyle(s.textForegroundStyle, style.TextForegroundStyleCompanion.Unspecified()).Alpha
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
		FontSize:               unit.TextUnitCompanion.TakeOrElse(opts.FontSize, s.FontSize),
		FontWeight:             font.FontWeightCompanion.TakeOrElse(opts.FontWeight, s.FontWeight),
		FontStyle:              font.FontStyleCompanion.TakeOrElse(opts.FontStyle, s.FontStyle),
		FontSynthesis:          font.TakeOrElseFontSynthesis(opts.FontSynthesis, s.FontSynthesis),
		FontFamily:             font.TakeOrElseFontFamily(opts.FontFamily, s.FontFamily),
		FontFeatureSettings:    stringutils.StringsCompanion.TakeOrElse(opts.FontFeatureSettings, s.FontFeatureSettings),
		LetterSpacing:          unit.TextUnitCompanion.TakeOrElse(opts.LetterSpacing, s.LetterSpacing),
		BaselineShift:          style.TakeOrElseBaselineShift(opts.BaselineShift, s.BaselineShift),
		TextGeometricTransform: style.TakeOrElseTextGeometricTransform(opts.TextGeometricTransform, s.TextGeometricTransform),
		LocaleList:             intl.LocaleListCompanion.TakeOrElse(opts.LocaleList, s.LocaleList),
		Background:             graphics.ColorCompanion.TakeOrElse(opts.Background, s.Background),
		TextDecoration:         style.TextDecorationCompanion.TakeOrElse(opts.TextDecoration, s.TextDecoration),
		Shadow:                 graphics.ShadowCompanion.TakeOrElse(opts.Shadow, s.Shadow),
		PlatformStyle:          takeOrElsePlatformSpanStyle(opts.PlatformStyle, s.PlatformStyle),
		DrawStyle:              graphics.TakeOrElseDrawStyle(opts.DrawStyle, s.DrawStyle),
	}
}

func StringSpanStyle(s *SpanStyle) string {

	if s == nil {
		return "SpanStyle(nil)"
	}

	s = coalesceSpanStyle(s, spanStyleUnspecified)

	return "SpanStyle(" +
		"color=" + s.Color().String() + ", " +
		"brush=" + fmt.Sprintf("%v", s.Brush()) + ", " +
		"alpha=" + fmt.Sprintf("%.2f", s.Alpha()) + ", " +
		"fontSize=" + s.FontSize.String() + ", " +
		"fontWeight=" + s.FontWeight.String() + ", " +
		"fontStyle=" + s.FontStyle.String() + ", " +
		"fontSynthesis=" + font.StringFontSynthesis(s.FontSynthesis) + ", " +
		"fontFamily=" + font.StringFontFamily(s.FontFamily) + ", " +
		"fontFeatureSettings=" + s.FontFeatureSettings + ", " +
		"letterSpacing=" + s.LetterSpacing.String() + ", " +
		"baselineShift=" + s.BaselineShift.String() + ", " +
		"textGeometricTransform=" + style.StringTextGeometricTransform(s.TextGeometricTransform) + ", " +
		"localeList=" + fmt.Sprintf("%v", s.LocaleList) + ", " +
		"background=" + s.Background.String() + ", " +
		"textDecoration=" + style.StringTextDecoration(s.TextDecoration) + ", " +
		"shadow=" + graphics.StringShadow(s.Shadow) + ", " +
		"platformStyle=" + fmt.Sprintf("%v", s.PlatformStyle) + ", " +
		"drawStyle=" + fmt.Sprintf("%v", s.DrawStyle) +
		")"
}

func IsSpecifiedSpanStyle(s *SpanStyle) bool {
	return SpanStyleCompanion.IsSpecified(s)
}

func TakeOrElseSpanStyle(s, def *SpanStyle) *SpanStyle {
	return takeOrElseSpanStyle(s, def)
}

// Identity (2 ns)
func SameSpanStyle(a, b *SpanStyle) bool {
	return sameSpanStyle(a, b)
}

// Semantic equality (field-by-field, 20 ns)
func SemanticEqualSpanStyle(a, b *SpanStyle) bool {
	return semanticEqualSpanStyle(a, b)
}

func EqualSpanStyle(a, b *SpanStyle) bool {
	return equalSpanStyle(a, b)
}

func MergeSpanStyle(a, b *SpanStyle) *SpanStyle {
	return mergeSpanStyle(a, b)
}

func CoalesceSpanStyle(ptr, def *SpanStyle) *SpanStyle {
	return coalesceSpanStyle(ptr, def)
}

// === Sentinel operations ===

func isSpecifiedSpanStyle(s *SpanStyle) bool {
	return s != nil && s != spanStyleUnspecified
}

func isUnspecifiedSpanStyle(s *SpanStyle) bool {
	return !isSpecifiedSpanStyle(s)
}

func takeOrElseSpanStyle(s, def *SpanStyle) *SpanStyle {
	if !isSpecifiedSpanStyle(s) {
		return def
	}
	return s
}

func coalesceSpanStyle(ptr, def *SpanStyle) *SpanStyle {
	if ptr == nil {
		return def
	}
	return ptr
}

func mergeSpanStyle(a, b *SpanStyle) *SpanStyle {
	if b == nil {
		return a
	}
	if a == nil {
		return b
	}

	a = coalesceSpanStyle(a, spanStyleUnspecified)
	b = coalesceSpanStyle(b, spanStyleUnspecified)
	// b is kept as is because we merge INTO it

	return &SpanStyle{
		textForegroundStyle:    style.MergeTextForegroundStyle(a.textForegroundStyle, b.textForegroundStyle),
		FontSize:               unit.TextUnitCompanion.TakeOrElse(b.FontSize, a.FontSize),
		FontWeight:             font.FontWeightCompanion.TakeOrElse(b.FontWeight, a.FontWeight),
		FontStyle:              font.FontStyleCompanion.TakeOrElse(b.FontStyle, a.FontStyle),
		FontSynthesis:          font.TakeOrElseFontSynthesis(b.FontSynthesis, a.FontSynthesis),
		FontFamily:             font.TakeOrElseFontFamily(b.FontFamily, a.FontFamily),
		FontFeatureSettings:    stringutils.StringsCompanion.TakeOrElse(b.FontFeatureSettings, a.FontFeatureSettings),
		LetterSpacing:          unit.TextUnitCompanion.TakeOrElse(b.LetterSpacing, a.LetterSpacing),
		BaselineShift:          style.BaselineShiftCompanion.TakeOrElse(b.BaselineShift, a.BaselineShift),
		TextGeometricTransform: style.TextGeometricTransformCompanion.TakeOrElse(b.TextGeometricTransform, a.TextGeometricTransform),
		LocaleList:             intl.LocaleListCompanion.TakeOrElse(b.LocaleList, a.LocaleList),
		Background:             graphics.ColorCompanion.TakeOrElse(b.Background, a.Background),
		TextDecoration:         style.TakeOrElseTextDecoration(b.TextDecoration, a.TextDecoration),
		Shadow:                 graphics.ShadowCompanion.TakeOrElse(b.Shadow, a.Shadow),
		PlatformStyle:          takeOrElsePlatformSpanStyle(b.PlatformStyle, a.PlatformStyle),
		DrawStyle:              graphics.TakeOrElseDrawStyle(b.DrawStyle, a.DrawStyle),
	}
}

func stringSpanStyle(s *SpanStyle) string {
	return StringSpanStyle(s)
}

func sameSpanStyle(a, b *SpanStyle) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return b == spanStyleUnspecified
	}
	if b == nil {
		return a == spanStyleUnspecified
	}
	return a == b
}

func semanticEqualSpanStyle(a, b *SpanStyle) bool {

	a = coalesceSpanStyle(a, spanStyleUnspecified)
	b = coalesceSpanStyle(b, spanStyleUnspecified)

	if a == spanStyleUnspecified && b == spanStyleUnspecified {
		return true
	}

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
		equalPlatformSpanStyle(a.PlatformStyle, b.PlatformStyle) &&
		graphics.EqualDrawStyle(a.DrawStyle, b.DrawStyle)
}

func equalSpanStyle(a, b *SpanStyle) bool {
	if !sameSpanStyle(a, b) {
		return semanticEqualSpanStyle(a, b)
	}
	return true
}

// copySpanStyle returns a copy of s; the source is never mutated.
func copySpanStyle(s *SpanStyle, options ...func(*SpanStyle)) *SpanStyle {
	cp := &SpanStyle{}
	if isSpecifiedSpanStyle(s) {
		*cp = *s
	}
	for _, option := range options {
		option(cp)
	}
	return cp
}

// SpanStyleCompanion is the single exported sentinel-family symbol for SpanStyle.
var SpanStyleCompanion = companions.NewStructCompanion(
	func() *SpanStyle { return spanStyleUnspecified },
	stringSpanStyle,
	isSpecifiedSpanStyle,
	isUnspecifiedSpanStyle,
	takeOrElseSpanStyle,
	mergeSpanStyle,
	sameSpanStyle,
	semanticEqualSpanStyle,
	equalSpanStyle,
	coalesceSpanStyle,
	copySpanStyle,
)

func LerpSpanStyle(width, start, stop *SpanStyle, fraction float32) *SpanStyle {
	start = CoalesceSpanStyle(start, spanStyleUnspecified)
	stop = CoalesceSpanStyle(stop, spanStyleUnspecified)

	return &SpanStyle{
		textForegroundStyle: style.LerpTextForegroundStyle(start.textForegroundStyle, stop.textForegroundStyle, fraction),
		FontSize:            unit.LerpTextUnitInheritable(start.FontSize, stop.FontSize, fraction),
		FontWeight:          font.LerpFontWeight(start.FontWeight, stop.FontWeight, fraction),
		FontStyle:           lerp.LerpDiscrete(start.FontStyle, stop.FontStyle, fraction),
		FontSynthesis:       lerp.LerpDiscrete(start.FontSynthesis, stop.FontSynthesis, fraction),
		FontFamily:          lerp.LerpDiscrete(start.FontFamily, stop.FontFamily, fraction),
		FontFeatureSettings: lerp.LerpDiscrete(start.FontFeatureSettings, stop.FontFeatureSettings, fraction),
		LetterSpacing:       unit.LerpTextUnitInheritable(start.LetterSpacing, stop.LetterSpacing, fraction),
		BaselineShift:       style.LerpBaselineShift(start.BaselineShift, stop.BaselineShift, fraction),
		TextGeometricTransform: style.LerpGeometricTransform(
			style.TakeOrElseTextGeometricTransform(start.TextGeometricTransform, style.TextGeometricTransformNone),
			style.TakeOrElseTextGeometricTransform(stop.TextGeometricTransform, style.TextGeometricTransformNone),
			fraction),
		LocaleList:     lerp.LerpDiscrete(start.LocaleList, stop.LocaleList, fraction),
		Background:     graphics.LerpColor(start.Background, stop.Background, fraction),
		TextDecoration: lerp.LerpDiscrete(start.TextDecoration, stop.TextDecoration, fraction),
		Shadow:         graphics.LerpShadow(start.Shadow, stop.Shadow, fraction),
		PlatformStyle:  LerpPlatformSpanStyle(start.PlatformStyle, stop.PlatformStyle, fraction),
		DrawStyle:      lerp.LerpDiscrete(start.DrawStyle, stop.DrawStyle, fraction),
	}

}

func LerpPlatformSpanStyle(start, stop *PlatformSpanStyle, fraction float32) *PlatformSpanStyle {
	if start == nil && stop == nil {
		return nil
	}
	if fraction < 0.5 {
		return start
	}
	return stop
}

func ResolveSpanStyleDefaults(s *SpanStyle) *SpanStyle {
	s = coalesceSpanStyle(s, spanStyleUnspecified)
	return &SpanStyle{
		textForegroundStyle:    style.TakeOrElseTextForegroundStyle(s.textForegroundStyle, DefaultColorForegroundStyle),
		FontSize:               unit.TextUnitCompanion.TakeOrElse(s.FontSize, DefaultFontSize),
		FontWeight:             font.FontWeightCompanion.TakeOrElse(s.FontWeight, font.FontWeightNormal),
		FontStyle:              font.FontStyleCompanion.TakeOrElse(s.FontStyle, font.FontStyleNormal),
		FontSynthesis:          font.TakeOrElseFontSynthesis(s.FontSynthesis, font.FontSynthesisAll),
		FontFamily:             font.TakeOrElseFontFamily(s.FontFamily, font.FontFamilyDefault),
		FontFeatureSettings:    stringutils.StringsCompanion.TakeOrElse(s.FontFeatureSettings, ""),
		LetterSpacing:          unit.TextUnitCompanion.TakeOrElse(s.LetterSpacing, DefaultLetterSpacing),
		BaselineShift:          style.TakeOrElseBaselineShift(s.BaselineShift, style.BaselineShiftNone),
		TextGeometricTransform: style.TakeOrElseTextGeometricTransform(s.TextGeometricTransform, style.TextGeometricTransformNone), // TODO: check sentinel for nil
		LocaleList:             intl.LocaleListCompanion.TakeOrElse(s.LocaleList, nil),                                             // TODO: verify default
		Background:             graphics.ColorCompanion.TakeOrElse(s.Background, graphics.ColorTransparent),                        // Should be transparent?
		TextDecoration:         style.TakeOrElseTextDecoration(s.TextDecoration, nil),                                              // TODO verify
		Shadow:                 graphics.TakeOrElseShadow(s.Shadow, graphics.ShadowNone),
		PlatformStyle:          takeOrElsePlatformSpanStyle(s.PlatformStyle, nil), // TODO
		DrawStyle:              graphics.TakeOrElseDrawStyle(s.DrawStyle, nil),
	}
}
