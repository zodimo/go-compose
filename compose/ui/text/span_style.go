package text

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/text/font"
	"github.com/zodimo/go-compose/compose/ui/text/style"
	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/pkg/floatutils/lerp"
)

var SpanStyleUnspecified *SpanStyle = &SpanStyle{
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

type SpanStyleInterface interface {
	Color() graphics.Color
	FontSize() unit.TextUnit
	FontWeight() font.FontWeight
	FontStyle() font.FontStyle
	FontFamily() font.FontFamily
	LetterSpacing() unit.TextUnit
	Background() graphics.Color
	TextDecoration() *style.TextDecoration
	Shadow() *graphics.Shadow
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

type SpanStyleOption func(*SpanStyle)

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

func CopySpanStyle(s *SpanStyle, options ...SpanStyleOption) *SpanStyle {
	copy := *s
	for _, option := range options {
		option(&copy)
	}
	return &copy
}

func StringSpanStyle(s *SpanStyle) string {
	s = CoalesceSpanStyle(s, SpanStyleUnspecified)
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

func SpanStyleResolveDefaults(s *SpanStyle) *SpanStyle {
	s = CoalesceSpanStyle(s, SpanStyleUnspecified)
	return &SpanStyle{
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
