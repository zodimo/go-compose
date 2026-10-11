package text

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/text/font"
	"github.com/zodimo/go-compose/compose/ui/text/style"
	"github.com/zodimo/go-compose/compose/ui/text/unit"
	"github.com/zodimo/go-compose/pkg/floatutils/lerp"
	"github.com/zodimo/go-sentinel-helper/sentinel/companions"
)

// spanStyleUnspecified is the sentinel for SpanStyle. It is deliberately
// package-private: the only public way to obtain it is SpanStyleCompanion.Unspecified().
var spanStyleUnspecified *SpanStyle = &SpanStyle{
	color:          graphics.ColorCompanion.Unspecified(),
	fontSize:       unit.TextUnitCompanion.Unspecified(),
	fontWeight:     font.FontWeightCompanion.Unspecified(),
	fontStyle:      font.FontStyleCompanion.Unspecified(),
	fontFamily:     nil,
	letterSpacing:  unit.TextUnitCompanion.Unspecified(),
	background:     graphics.ColorCompanion.Unspecified(),
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

type SpanStyleOption = func(*SpanStyle)

func SpanStyleWithColor(color graphics.Color) SpanStyleOption {
	return func(opts *SpanStyle) {
		opts.color = color
	}
}

func SpanStyleWithFontSize(size unit.TextUnit) SpanStyleOption {
	return func(opts *SpanStyle) {
		opts.fontSize = size
	}
}

func SpanStyleWithFontWeight(weight font.FontWeight) SpanStyleOption {
	return func(opts *SpanStyle) {
		opts.fontWeight = weight
	}
}

func SpanStyleWithFontStyle(fontStyle font.FontStyle) SpanStyleOption {
	return func(opts *SpanStyle) {
		opts.fontStyle = fontStyle
	}
}

func SpanStyleWithFontFamily(family font.FontFamily) SpanStyleOption {
	return func(opts *SpanStyle) {
		opts.fontFamily = family
	}
}

func SpanStyleWithLetterSpacing(spacing unit.TextUnit) SpanStyleOption {
	return func(opts *SpanStyle) {
		opts.letterSpacing = spacing
	}
}

func SpanStyleWithBackground(color graphics.Color) SpanStyleOption {
	return func(opts *SpanStyle) {
		opts.background = color
	}
}

func SpanStyleWithTextDecoration(decoration *style.TextDecoration) SpanStyleOption {
	return func(opts *SpanStyle) {
		opts.textDecoration = decoration
	}
}

func SpanStyleWithShadow(shadow *graphics.Shadow) SpanStyleOption {
	return func(opts *SpanStyle) {
		opts.shadow = shadow
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
	localSs := coalesceSpanStyle(s, spanStyleUnspecified)
	changed := *spanStyleUnspecified
	for _, option := range options {
		option(&changed)
	}
	return &SpanStyle{
		color:          graphics.ColorCompanion.TakeOrElse(changed.color, localSs.color),
		fontSize:       unit.TextUnitCompanion.TakeOrElse(changed.fontSize, localSs.fontSize),
		fontWeight:     font.FontWeightCompanion.TakeOrElse(changed.fontWeight, localSs.fontWeight),
		fontStyle:      font.FontStyleCompanion.TakeOrElse(changed.fontStyle, localSs.fontStyle),
		fontFamily:     font.TakeOrElseFontFamily(changed.fontFamily, localSs.fontFamily),
		letterSpacing:  unit.TextUnitCompanion.TakeOrElse(changed.letterSpacing, localSs.letterSpacing),
		background:     graphics.ColorCompanion.TakeOrElse(changed.background, localSs.background),
		textDecoration: style.TakeOrElseTextDecoration(changed.textDecoration, localSs.textDecoration),
		shadow:         graphics.TakeOrElseShadow(changed.shadow, localSs.shadow),
	}

}

func stringSpanStyle(s *SpanStyle) string {
	localSs := coalesceSpanStyle(s, spanStyleUnspecified)
	if !isSpecifiedSpanStyle(localSs) {
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
	localA := coalesceSpanStyle(a, spanStyleUnspecified)
	localB := coalesceSpanStyle(b, spanStyleUnspecified)

	if localA == spanStyleUnspecified {
		return localB
	}
	if localB == spanStyleUnspecified {
		return localA
	}

	return &SpanStyle{
		color:          graphics.ColorCompanion.TakeOrElse(localB.color, localA.color),
		fontSize:       unit.TextUnitCompanion.TakeOrElse(localB.fontSize, localA.fontSize),
		fontWeight:     font.FontWeightCompanion.TakeOrElse(localB.fontWeight, localA.fontWeight),
		fontStyle:      font.FontStyleCompanion.TakeOrElse(localB.fontStyle, localA.fontStyle),
		fontFamily:     font.TakeOrElseFontFamily(localB.fontFamily, localA.fontFamily),
		letterSpacing:  unit.TextUnitCompanion.TakeOrElse(localB.letterSpacing, localA.letterSpacing),
		background:     graphics.ColorCompanion.TakeOrElse(localB.background, localA.background),
		textDecoration: style.TakeOrElseTextDecoration(localB.textDecoration, localA.textDecoration),
		shadow:         graphics.TakeOrElseShadow(localB.shadow, localA.shadow),
	}
}

func coalesceSpanStyle(ptr, def *SpanStyle) *SpanStyle {
	var zero *SpanStyle
	if ptr == zero {
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
	localSs := coalesceSpanStyle(s, spanStyleUnspecified)
	return &SpanStyle{
		color:          graphics.ColorCompanion.TakeOrElse(localSs.color, DefaultColor),
		fontSize:       unit.TextUnitCompanion.TakeOrElse(localSs.fontSize, DefaultFontSize),
		fontWeight:     font.FontWeightCompanion.TakeOrElse(localSs.fontWeight, DefaultFontWeight),
		fontStyle:      font.FontStyleCompanion.TakeOrElse(localSs.fontStyle, DefaultFontStyle),
		fontFamily:     font.TakeOrElseFontFamily(localSs.fontFamily, DefaultFontFamily),
		letterSpacing:  unit.TextUnitCompanion.TakeOrElse(localSs.letterSpacing, DefaultLetterSpacing),
		background:     graphics.ColorCompanion.TakeOrElse(localSs.background, DefaultBackgroundColor),
		textDecoration: style.TakeOrElseTextDecoration(localSs.textDecoration, DefaultTextDecoration),
		shadow:         graphics.TakeOrElseShadow(localSs.shadow, DefaultShadow),
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

// ensureMutableSpanStyle panics if ts is a sentinel value.
// Call this at the start of any function that mutates a SpanStyle to fail-fast on misuse.
func ensureMutableSpanStyle(s *SpanStyle) {
	if SpanStyleCompanion.IsUnspecified(s) {
		panic("attempt to mutate sentinel; use XCompanion.Copy first")
	}
}

// SpanStyleFromOptions builds a SpanStyle from options applied over the sentinel.
func SpanStyleFromOptions(options ...SpanStyleOption) *SpanStyle {
	return SpanStyleCompanion.Copy(spanStyleUnspecified, options...)
}
