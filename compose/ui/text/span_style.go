package text

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/text/font"
	"github.com/zodimo/go-compose/compose/ui/text/style"
	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/pkg/floatutils/lerp"
	"github.com/zodimo/go-sentinel-helper/sentinel/companions"
)

// spanStyleUnspecified is the sentinel for SpanStyle. It is deliberately
// package-private: the only public way to obtain it is SpanStyleCompanion.Unspecified().
var spanStyleUnspecified *SpanStyle = &SpanStyle{
	color:          graphics.ColorUnspecified,
	fontSize:       unit.TextUnitUnspecified,
	fontWeight:     font.FontWeightUnspecified,
	fontStyle:      font.FontStyleUnspecified,
	fontFamily:     nil,
	letterSpacing:  unit.TextUnitUnspecified,
	background:     graphics.ColorUnspecified,
	textDecoration: nil,
	shadow:         nil,
}

var _ SpanStyleInterface = (*SpanStyle)(nil)

type SpanStyle struct {
	color          graphics.Color
	fontSize       unit.TextUnit
	fontWeight     font.FontWeight
	fontStyle      font.FontStyle
	fontFamily     font.FontFamily
	letterSpacing  unit.TextUnit
	background     graphics.Color
	textDecoration *style.TextDecoration
	shadow         *graphics.Shadow
}

// SpanStyleOptions carries one field per SpanStyle property. A field left at its
// sentinel (or nil) means "leave the source value alone"; anything else overrides it.
type SpanStyleOptions struct {
	Color          graphics.Color
	FontSize       unit.TextUnit
	FontWeight     font.FontWeight
	FontStyle      font.FontStyle
	FontFamily     font.FontFamily
	LetterSpacing  unit.TextUnit
	Background     graphics.Color
	TextDecoration *style.TextDecoration
	Shadow         *graphics.Shadow
}

func newSpanStyleOptions() SpanStyleOptions {
	return SpanStyleOptions{
		Color:          graphics.ColorUnspecified,
		FontSize:       unit.TextUnitUnspecified,
		FontWeight:     font.FontWeightUnspecified,
		FontStyle:      font.FontStyleUnspecified,
		FontFamily:     nil,
		LetterSpacing:  unit.TextUnitUnspecified,
		Background:     graphics.ColorUnspecified,
		TextDecoration: nil,
		Shadow:         nil,
	}
}

type SpanStyleOption = func(*SpanStyleOptions)

// asOption lets an already-populated SpanStyleOptions seed a SpanStyleOptions
// target, so the whole struct can be passed as a single SpanStyleOption.
func (o SpanStyleOptions) asOption(opts *SpanStyleOptions) {
	*opts = o
}

func SpanStyleWithColor(color graphics.Color) SpanStyleOption {
	return func(opts *SpanStyleOptions) {
		opts.Color = color
	}
}

func SpanStyleWithFontSize(size unit.TextUnit) SpanStyleOption {
	return func(opts *SpanStyleOptions) {
		opts.FontSize = size
	}
}

func SpanStyleWithFontWeight(weight font.FontWeight) SpanStyleOption {
	return func(opts *SpanStyleOptions) {
		opts.FontWeight = weight
	}
}

func SpanStyleWithFontStyle(fontStyle font.FontStyle) SpanStyleOption {
	return func(opts *SpanStyleOptions) {
		opts.FontStyle = fontStyle
	}
}

func SpanStyleWithFontFamily(family font.FontFamily) SpanStyleOption {
	return func(opts *SpanStyleOptions) {
		opts.FontFamily = family
	}
}

func SpanStyleWithLetterSpacing(spacing unit.TextUnit) SpanStyleOption {
	return func(opts *SpanStyleOptions) {
		opts.LetterSpacing = spacing
	}
}

func SpanStyleWithBackground(color graphics.Color) SpanStyleOption {
	return func(opts *SpanStyleOptions) {
		opts.Background = color
	}
}

func SpanStyleWithTextDecoration(decoration *style.TextDecoration) SpanStyleOption {
	return func(opts *SpanStyleOptions) {
		opts.TextDecoration = decoration
	}
}

func SpanStyleWithShadow(shadow *graphics.Shadow) SpanStyleOption {
	return func(opts *SpanStyleOptions) {
		opts.Shadow = shadow
	}
}

func (s SpanStyle) Color() graphics.Color {
	return s.color
}

func (s SpanStyle) FontSize() unit.TextUnit {
	return s.fontSize
}
func (s SpanStyle) FontWeight() font.FontWeight {
	return s.fontWeight
}
func (s SpanStyle) FontStyle() font.FontStyle {
	return s.fontStyle
}
func (s SpanStyle) FontFamily() font.FontFamily {
	return s.fontFamily
}
func (s SpanStyle) LetterSpacing() unit.TextUnit {
	return s.letterSpacing
}

func (s SpanStyle) Background() graphics.Color {
	return s.background
}
func (s SpanStyle) TextDecoration() *style.TextDecoration {
	return s.textDecoration
}
func (s SpanStyle) Shadow() *graphics.Shadow {
	return s.shadow
}

// copySpanStyle returns a copy of s with the given bare options applied on top.
// By the time an option runs, the working copy of s has already been seeded from the
// source, so an option that only touches one field leaves the rest intact.
func copySpanStyle(s *SpanStyle, options ...func(*SpanStyle)) *SpanStyle {
	s = coalesceSpanStyle(s, spanStyleUnspecified)
	work := *s
	for _, option := range options {
		option(&work)
	}
	return &work
}

// copySpanStyleWithOptions is the option-constructor form used by callers.
// Options are applied to a fresh SpanStyleOptions seeded with sentinels, then only
// the fields the caller actually set are overlaid onto s.
func copySpanStyleWithOptions(s *SpanStyle, options ...SpanStyleOption) *SpanStyle {
	s = coalesceSpanStyle(s, spanStyleUnspecified)
	opts := newSpanStyleOptions()
	for _, option := range options {
		option(&opts)
	}
	return &SpanStyle{
		color:          opts.Color.TakeOrElse(s.color),
		fontSize:       opts.FontSize.TakeOrElse(s.fontSize),
		fontWeight:     opts.FontWeight.TakeOrElse(s.fontWeight),
		fontStyle:      opts.FontStyle.TakeOrElse(s.fontStyle),
		fontFamily:     font.TakeOrElseFontFamily(opts.FontFamily, s.fontFamily),
		letterSpacing:  opts.LetterSpacing.TakeOrElse(s.letterSpacing),
		background:     opts.Background.TakeOrElse(s.background),
		textDecoration: style.TakeOrElseTextDecoration(opts.TextDecoration, s.textDecoration),
		shadow:         graphics.TakeOrElseShadow(opts.Shadow, s.shadow),
	}
}

func stringSpanStyle(s *SpanStyle) string {
	s = coalesceSpanStyle(s, spanStyleUnspecified)
	if !isSpecifiedSpanStyle(s) {
		return "SpanStyle{Unspecified}"
	}
	return fmt.Sprintf("SpanStyle("+
		"Color=%s, "+
		"FontSize=%s, "+
		"FontWeight=%s, "+
		"FontStyle=%s, "+
		"FontFamily=%s, "+
		"LetterSpacing=%s, "+
		"Background=%s, "+
		"TextDecoration=%s, "+
		"Shadow=%s)",
		s.color.String(),
		s.fontSize.String(),
		s.fontWeight.String(),
		s.fontStyle.String(),
		font.StringFontFamily(s.fontFamily),
		s.letterSpacing.String(),
		s.background.String(),
		style.StringTextDecoration(s.textDecoration),
		graphics.StringShadow(s.shadow))
}

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

// Identity (2 ns)
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

// Semantic equality (field-by-field, 20 ns)
func semanticEqualSpanStyle(a, b *SpanStyle) bool {

	a = coalesceSpanStyle(a, spanStyleUnspecified)
	b = coalesceSpanStyle(b, spanStyleUnspecified)

	return a.color == b.color &&
		a.fontSize == b.fontSize &&
		a.fontWeight == b.fontWeight &&
		a.fontStyle == b.fontStyle &&
		font.EqualFontFamily(a.fontFamily, b.fontFamily) &&
		a.letterSpacing == b.letterSpacing &&
		a.background == b.background &&
		style.EqualTextDecoration(a.textDecoration, b.textDecoration) &&
		graphics.EqualShadow(a.shadow, b.shadow)
}

func equalSpanStyle(a, b *SpanStyle) bool {
	if sameSpanStyle(a, b) {
		return true
	}
	return semanticEqualSpanStyle(a, b)
}

// mergeSpanStyle layers b onto a: a value from b wins only where b actually set it.
func mergeSpanStyle(a, b *SpanStyle) *SpanStyle {
	a = coalesceSpanStyle(a, spanStyleUnspecified)
	b = coalesceSpanStyle(b, spanStyleUnspecified)

	if a == spanStyleUnspecified {
		return b
	}
	if b == spanStyleUnspecified {
		return a
	}

	return &SpanStyle{
		color:          b.color.TakeOrElse(a.color),
		fontSize:       b.fontSize.TakeOrElse(a.fontSize),
		fontWeight:     b.fontWeight.TakeOrElse(a.fontWeight),
		fontStyle:      b.fontStyle.TakeOrElse(a.fontStyle),
		fontFamily:     font.TakeOrElseFontFamily(b.fontFamily, a.fontFamily),
		letterSpacing:  b.letterSpacing.TakeOrElse(a.letterSpacing),
		background:     b.background.TakeOrElse(a.background),
		textDecoration: style.TakeOrElseTextDecoration(b.textDecoration, a.textDecoration),
		shadow:         graphics.TakeOrElseShadow(b.shadow, a.shadow),
	}
}

func coalesceSpanStyle(ptr, def *SpanStyle) *SpanStyle {
	if ptr == nil {
		return def
	}
	return ptr
}

func lerpSpanStyle(start, stop *SpanStyle, fraction float32) *SpanStyle {
	start = coalesceSpanStyle(start, spanStyleUnspecified)
	stop = coalesceSpanStyle(stop, spanStyleUnspecified)

	return &SpanStyle{
		color:          graphics.LerpColor(start.color, stop.color, fraction),
		fontSize:       unit.LerpTextUnitInheritable(start.fontSize, stop.fontSize, fraction),
		fontWeight:     font.LerpFontWeight(start.fontWeight, stop.fontWeight, fraction),
		fontStyle:      lerp.LerpDiscrete(start.fontStyle, stop.fontStyle, fraction),
		fontFamily:     lerp.LerpDiscrete(start.fontFamily, stop.fontFamily, fraction),
		letterSpacing:  unit.LerpTextUnitInheritable(start.letterSpacing, stop.letterSpacing, fraction),
		background:     graphics.LerpColor(start.background, stop.background, fraction),
		textDecoration: lerp.LerpDiscrete(start.textDecoration, stop.textDecoration, fraction),
		shadow:         graphics.LerpShadow(start.shadow, stop.shadow, fraction),
	}

}

func spanStyleResolveDefaults(s *SpanStyle) *SpanStyle {
	s = coalesceSpanStyle(s, spanStyleUnspecified)
	return &SpanStyle{
		color:          s.color.TakeOrElse(DefaultColor),
		fontSize:       s.fontSize.TakeOrElse(DefaultFontSize),
		fontWeight:     s.fontWeight.TakeOrElse(DefaultFontWeight),
		fontStyle:      s.fontStyle.TakeOrElse(DefaultFontStyle),
		fontFamily:     font.TakeOrElseFontFamily(s.fontFamily, DefaultFontFamily),
		letterSpacing:  s.letterSpacing.TakeOrElse(DefaultLetterSpacing),
		background:     s.background.TakeOrElse(DefaultBackgroundColor),
		textDecoration: style.TakeOrElseTextDecoration(s.textDecoration, DefaultTextDecoration),
		shadow:         graphics.TakeOrElseShadow(s.shadow, DefaultShadow),
	}
}

// SpanStyleCompanion is the only exported value of this file's sentinel family member.
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
