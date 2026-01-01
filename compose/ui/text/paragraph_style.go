package text

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/ui/next/text/style"
	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/pkg/floatutils/lerp"
)

var ParagraphStyleUnspecified *ParagraphStyle = &ParagraphStyle{
	textAlign:     style.TextAlignUnspecified,
	textDirection: style.TextDirectionUnspecified,
	lineHeight:    unit.TextUnitUnspecified,
	lineBreak:     style.LineBreakUnspecified,
}

type ParaghStyleInterface interface {
	TextAlign() style.TextAlign
	TextDirection() style.TextDirection
	LineHeight() unit.TextUnit
	LineBreak() style.LineBreak
}

var _ ParaghStyleInterface = (*ParagraphStyle)(nil)

type ParagraphStyle struct {
	textAlign     style.TextAlign
	textDirection style.TextDirection
	lineHeight    unit.TextUnit
	lineBreak     style.LineBreak
}

type ParagraphStyleOption func(*ParagraphStyle)

func (ps ParagraphStyle) TextAlign() style.TextAlign {
	return ps.textAlign
}
func (ps ParagraphStyle) TextDirection() style.TextDirection {
	return ps.textDirection
}
func (ps ParagraphStyle) LineHeight() unit.TextUnit {
	return ps.lineHeight
}
func (ps ParagraphStyle) LineBreak() style.LineBreak {
	return ps.lineBreak
}

func ParagraphStyleCopy(s *ParagraphStyle, options ...ParagraphStyleOption) *ParagraphStyle {
	copy := *s
	for _, option := range options {
		option(&copy)
	}
	return &copy
}

func StringParagraphStyle(s *ParagraphStyle) string {
	s = CoalesceParagraphStyle(s, ParagraphStyleUnspecified)
	return fmt.Sprintf("ParagraphStyle("+
		"TextAlign=%s, "+
		"TextDirection=%s, "+
		"LineHeight=%s, "+
		"LineBreak=%s)",
		s.textAlign.String(),
		s.textDirection.String(),
		s.lineHeight.String(),
		s.lineBreak.String())
}

func IsSpecifiedParagraphStyle(s *ParagraphStyle) bool {
	return s != nil && s != ParagraphStyleUnspecified
}

func TakeOrElseParagraphStyle(s, def *ParagraphStyle) *ParagraphStyle {
	if !IsSpecifiedParagraphStyle(s) {
		return def
	}
	return s
}

// Identity (2 ns)
func SameParagraphStyle(a, b *ParagraphStyle) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return b == ParagraphStyleUnspecified
	}
	if b == nil {
		return a == ParagraphStyleUnspecified
	}
	return a == b
}

// Semantic equality (field-by-field, 20 ns)
func SemanticEqualParagraphStyle(a, b *ParagraphStyle) bool {

	a = CoalesceParagraphStyle(a, ParagraphStyleUnspecified)
	b = CoalesceParagraphStyle(b, ParagraphStyleUnspecified)

	return a.textAlign == b.textAlign &&
		a.textDirection == b.textDirection &&
		a.lineHeight == b.lineHeight &&
		a.lineBreak == b.lineBreak
}

func EqualParagraphStyle(a, b *ParagraphStyle) bool {
	if !SameParagraphStyle(a, b) {
		return SemanticEqualParagraphStyle(a, b)
	}
	return true
}

func MergeParagraphStyle(a, b *ParagraphStyle) *ParagraphStyle {
	if b == nil {
		return a
	}
	if a == nil {
		return b
	}

	a = CoalesceParagraphStyle(a, ParagraphStyleUnspecified)
	b = CoalesceParagraphStyle(b, ParagraphStyleUnspecified)

	return &ParagraphStyle{
		textAlign:     b.textAlign.TakeOrElse(a.textAlign),
		textDirection: b.textDirection.TakeOrElse(a.textDirection),
		lineHeight:    b.lineHeight.TakeOrElse(a.lineHeight),
		lineBreak:     b.lineBreak.TakeOrElse(a.lineBreak),
	}
}

func CoalesceParagraphStyle(ptr, def *ParagraphStyle) *ParagraphStyle {
	if ptr == nil {
		return def
	}
	return ptr
}

func LerpParagraphStyle(width, start, stop *ParagraphStyle, fraction float32) *ParagraphStyle {
	start = CoalesceParagraphStyle(start, ParagraphStyleUnspecified)
	stop = CoalesceParagraphStyle(stop, ParagraphStyleUnspecified)

	return &ParagraphStyle{
		textAlign:     lerp.LerpDiscrete(start.textAlign, stop.textAlign, fraction),
		textDirection: lerp.LerpDiscrete(start.textDirection, stop.textDirection, fraction),
		lineHeight:    style.LerpTextUnitInheritable(start.lineHeight, stop.lineHeight, fraction),
		lineBreak:     lerp.LerpDiscrete(start.lineBreak, stop.lineBreak, fraction),
	}

}
