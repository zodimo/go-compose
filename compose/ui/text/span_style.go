package text

import (
	"fmt"

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
	return style.CoalesceTextForegroundStyle(s.textForegroundStyle, style.TextForegroundStyleUnspecified).Color
}
func (s SpanStyle) Brush() graphics.Brush {
	return style.CoalesceTextForegroundStyle(s.textForegroundStyle, style.TextForegroundStyleUnspecified).Brush
}
func (s SpanStyle) Alpha() float32 {
	return style.CoalesceTextForegroundStyle(s.textForegroundStyle, style.TextForegroundStyleUnspecified).Alpha
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

	if s == nil {
		return "SpanStyle(nil)"
	}

	s = CoalesceSpanStyle(s, SpanStyleUnspecified)

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
	if b == nil {
		return a
	}
	if a == nil {
		return b
	}

	a = CoalesceSpanStyle(a, SpanStyleUnspecified)
	b = CoalesceSpanStyle(b, SpanStyleUnspecified)
	// b is kept as is because we merge INTO it

	return &SpanStyle{
		textForegroundStyle:    style.MergeTextForegroundStyle(a.textForegroundStyle, b.textForegroundStyle),
		FontSize:               b.FontSize.TakeOrElse(a.FontSize),
		FontWeight:             b.FontWeight.TakeOrElse(a.FontWeight),
		FontStyle:              b.FontStyle.TakeOrElse(a.FontStyle),
		FontSynthesis:          font.TakeOrElseFontSynthesis(b.FontSynthesis, a.FontSynthesis),
		FontFamily:             font.TakeOrElseFontFamily(b.FontFamily, a.FontFamily),
		FontFeatureSettings:    sentinel.TakeOrElseString(b.FontFeatureSettings, a.FontFeatureSettings),
		LetterSpacing:          b.LetterSpacing.TakeOrElse(a.LetterSpacing),
		BaselineShift:          style.TakeOrElseBaselineShift(b.BaselineShift, a.BaselineShift),
		TextGeometricTransform: style.TakeOrElseTextGeometricTransform(b.TextGeometricTransform, a.TextGeometricTransform),
		LocaleList:             intl.TakeOrElseLocaleList(b.LocaleList, a.LocaleList),
		Background:             b.Background.TakeOrElse(a.Background),
		TextDecoration:         style.TakeOrElseTextDecoration(b.TextDecoration, a.TextDecoration),
		Shadow:                 graphics.TakeOrElseShadow(b.Shadow, a.Shadow),
		PlatformStyle:          TakeOrElsePlatformSpanStyle(b.PlatformStyle, a.PlatformStyle),
		DrawStyle:              graphics.TakeOrElseDrawStyle(b.DrawStyle, a.DrawStyle),
	}
}

func CoalesceSpanStyle(ptr, def *SpanStyle) *SpanStyle {
	if ptr == nil {
		return def
	}
	return ptr
}

func LerpSpanStyle(width, start, stop *SpanStyle, fraction float32) *SpanStyle {
	start = CoalesceSpanStyle(start, SpanStyleUnspecified)
	stop = CoalesceSpanStyle(stop, SpanStyleUnspecified)

	return &SpanStyle{
		textForegroundStyle: style.LerpTextForegroundStyle(start.textForegroundStyle, stop.textForegroundStyle, fraction),
		FontSize:            LerpTextUnitInheritable(start.FontSize, stop.FontSize, fraction),
		FontWeight:          font.LerpFontWeight(start.FontWeight, stop.FontWeight, fraction),
		FontStyle:           lerpDiscrete(start.FontStyle, stop.FontStyle, fraction),
		FontSynthesis:       lerpDiscrete(start.FontSynthesis, stop.FontSynthesis, fraction),
		FontFamily:          lerpDiscrete(start.FontFamily, stop.FontFamily, fraction),
		FontFeatureSettings: lerpDiscrete(start.FontFeatureSettings, stop.FontFeatureSettings, fraction),
		LetterSpacing:       LerpTextUnitInheritable(start.LetterSpacing, stop.LetterSpacing, fraction),
		BaselineShift:       style.LerpBaselineShift(start.BaselineShift, stop.BaselineShift, fraction),
		TextGeometricTransform: ptr(style.LerpGeometricTransform(
			style.TakeOrElseTextGeometricTransform(start.TextGeometricTransform, style.TextGeometricTransformNone),
			style.TakeOrElseTextGeometricTransform(stop.TextGeometricTransform, style.TextGeometricTransformNone),
			fraction)),
		LocaleList:     lerpDiscrete(start.LocaleList, stop.LocaleList, fraction),
		Background:     graphics.LerpColor(start.Background, stop.Background, fraction),
		TextDecoration: lerpDiscrete(start.TextDecoration, stop.TextDecoration, fraction),
		Shadow:         graphics.LerpShadow(start.Shadow, stop.Shadow, fraction),
		PlatformStyle:  LerpPlatformSpanStyle(start.PlatformStyle, stop.PlatformStyle, fraction),
		DrawStyle:      lerpDiscrete(start.DrawStyle, stop.DrawStyle, fraction),
	}

}

func LerpTextUnitInheritable(a, b unit.TextUnit, t float32) unit.TextUnit {
	if a.IsUnspecified() || b.IsUnspecified() {
		return lerpDiscrete(a, b, t)
	}
	return unit.LerpTextUnit(a, b, t)
}

func lerpDiscrete[T any](a, b T, fraction float32) T {
	if fraction < 0.5 {
		return a
	}
	return b
}

func ptr[T any](v T) *T {
	return &v
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
	s = CoalesceSpanStyle(s, SpanStyleUnspecified)
	return &SpanStyle{
		textForegroundStyle:    style.TakeOrElseTextForegroundStyle(s.textForegroundStyle, DefaultColorForegroundStyle),
		FontSize:               s.FontSize.TakeOrElse(DefaultFontSize),
		FontWeight:             s.FontWeight.TakeOrElse(font.FontWeightNormal),
		FontStyle:              s.FontStyle.TakeOrElse(font.FontStyleNormal),
		FontSynthesis:          font.TakeOrElseFontSynthesis(s.FontSynthesis, font.FontSynthesisAll),
		FontFamily:             font.TakeOrElseFontFamily(s.FontFamily, font.FontFamilyDefault),
		FontFeatureSettings:    sentinel.TakeOrElseString(s.FontFeatureSettings, ""),
		LetterSpacing:          s.LetterSpacing.TakeOrElse(DefaultLetterSpacing),
		BaselineShift:          style.TakeOrElseBaselineShift(s.BaselineShift, style.BaselineShiftNone),
		TextGeometricTransform: style.TakeOrElseTextGeometricTransform(s.TextGeometricTransform, style.TextGeometricTransformNone), // TODO: check sentinel for nil
		LocaleList:             intl.TakeOrElseLocaleList(s.LocaleList, nil),                                                       // TODO: verify default
		Background:             s.Background.TakeOrElse(graphics.ColorTransparent),                                                 // Should be transparent?
		TextDecoration:         style.TakeOrElseTextDecoration(s.TextDecoration, nil),                                              // TODO verify
		Shadow:                 graphics.TakeOrElseShadow(s.Shadow, graphics.ShadowNone),
		PlatformStyle:          TakeOrElsePlatformSpanStyle(s.PlatformStyle, nil), // TODO
		DrawStyle:              graphics.TakeOrElseDrawStyle(s.DrawStyle, nil),
	}
}
