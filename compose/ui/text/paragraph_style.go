package text

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/ui/text/style"
	"github.com/zodimo/go-compose/compose/ui/text/unit"
	"github.com/zodimo/go-compose/pkg/floatutils/lerp"
	"github.com/zodimo/go-sentinel-helper/sentinel/companions"
)

// paragraphStyleUnspecified is the sentinel for ParagraphStyle. It is deliberately
// package-private: the only public way to obtain it is ParagraphStyleCompanion.Unspecified().
var paragraphStyleUnspecified *ParagraphStyle = &ParagraphStyle{
	textAlign:     style.TextAlignCompanion.Unspecified(),
	textDirection: style.TextDirectionCompanion.Unspecified(),
	lineHeight:    unit.TextUnitCompanion.Unspecified(),
	lineBreak:     style.LineBreakCompanion.Unspecified(),
}

var _ ParaghStyleInterface = (*ParagraphStyle)(nil)

type ParagraphStyle struct {
	textAlign     style.TextAlign
	textDirection style.TextDirection
	lineHeight    unit.TextUnit
	lineBreak     style.LineBreak
}

type ParagraphStyleOption = func(*ParagraphStyle)

func ParagraphStyleWithTextAlign(textAlign style.TextAlign) ParagraphStyleOption {
	return func(opts *ParagraphStyle) {
		opts.textAlign = textAlign
	}
}

func ParagraphStyleWithTextDirection(textDirection style.TextDirection) ParagraphStyleOption {
	return func(opts *ParagraphStyle) {
		opts.textDirection = textDirection
	}
}

func ParagraphStyleWithLineHeight(lineHeight unit.TextUnit) ParagraphStyleOption {
	return func(opts *ParagraphStyle) {
		opts.lineHeight = lineHeight
	}
}

func ParagraphStyleWithLineBreak(lineBreak style.LineBreak) ParagraphStyleOption {
	return func(opts *ParagraphStyle) {
		opts.lineBreak = lineBreak
	}
}

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

func stringParagraphStyle(s *ParagraphStyle) string {
	s = coalesceParagraphStyle(s, paragraphStyleUnspecified)
	if !isSpecifiedParagraphStyle(s) {
		return "ParagraphStyle{Unspecified}"
	}
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

func isSpecifiedParagraphStyle(s *ParagraphStyle) bool {
	return s != nil && s != paragraphStyleUnspecified
}

func isUnspecifiedParagraphStyle(s *ParagraphStyle) bool {
	return !isSpecifiedParagraphStyle(s)
}

func takeOrElseParagraphStyle(s, def *ParagraphStyle) *ParagraphStyle {
	if !isSpecifiedParagraphStyle(s) {
		return def
	}
	return s
}

// Identity (2 ns)
func sameParagraphStyle(a, b *ParagraphStyle) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return b == paragraphStyleUnspecified
	}
	if b == nil {
		return a == paragraphStyleUnspecified
	}
	return a == b
}

// Semantic equality (field-by-field, 20 ns)
func semanticEqualParagraphStyle(a, b *ParagraphStyle) bool {

	a = coalesceParagraphStyle(a, paragraphStyleUnspecified)
	b = coalesceParagraphStyle(b, paragraphStyleUnspecified)

	return a.textAlign == b.textAlign &&
		a.textDirection == b.textDirection &&
		a.lineHeight == b.lineHeight &&
		a.lineBreak == b.lineBreak
}

func equalParagraphStyle(a, b *ParagraphStyle) bool {
	if sameParagraphStyle(a, b) {
		return true
	}
	return semanticEqualParagraphStyle(a, b)
}

// mergeParagraphStyle layers b onto a: a value from b wins only where b actually set it.
func mergeParagraphStyle(a, b *ParagraphStyle) *ParagraphStyle {
	a = coalesceParagraphStyle(a, paragraphStyleUnspecified)
	b = coalesceParagraphStyle(b, paragraphStyleUnspecified)

	if a == paragraphStyleUnspecified {
		return b
	}
	if b == paragraphStyleUnspecified {
		return a
	}

	return &ParagraphStyle{
		textAlign:     style.TextAlignCompanion.TakeOrElse(b.textAlign, a.textAlign),
		textDirection: style.TextDirectionCompanion.TakeOrElse(b.textDirection, a.textDirection),
		lineHeight:    unit.TextUnitCompanion.TakeOrElse(b.lineHeight, a.lineHeight),
		lineBreak:     style.LineBreakCompanion.TakeOrElse(b.lineBreak, a.lineBreak),
	}
}

func coalesceParagraphStyle(ptr, def *ParagraphStyle) *ParagraphStyle {
	if ptr == nil {
		return def
	}
	return ptr
}

// copyParagraphStyle returns a copy of s with the given bare options applied on top.
// By the time an option runs, the working copy of s has already been seeded from the
// source, so an option that only touches one field leaves the rest intact.
func copyParagraphStyle(s *ParagraphStyle, options ...func(*ParagraphStyle)) *ParagraphStyle {
	s = coalesceParagraphStyle(s, paragraphStyleUnspecified)
	changed := *paragraphStyleUnspecified
	for _, option := range options {
		option(&changed)
	}
	return mergeParagraphStyle(s, &changed)
}

// lerpParagraphStyle interpolates between two ParagraphStyles.
func lerpParagraphStyle(start, stop *ParagraphStyle, fraction float32) *ParagraphStyle {
	start = coalesceParagraphStyle(start, paragraphStyleUnspecified)
	stop = coalesceParagraphStyle(stop, paragraphStyleUnspecified)

	return &ParagraphStyle{
		textAlign:     lerp.LerpDiscrete(start.textAlign, stop.textAlign, fraction),
		textDirection: lerp.LerpDiscrete(start.textDirection, stop.textDirection, fraction),
		lineHeight:    unit.LerpTextUnitInheritable(start.lineHeight, stop.lineHeight, fraction),
		lineBreak:     lerp.LerpDiscrete(start.lineBreak, stop.lineBreak, fraction),
	}
}

func paragraphStyleResolveDefaults(s *ParagraphStyle, direction unit.LayoutDirection) *ParagraphStyle {
	s = coalesceParagraphStyle(s, paragraphStyleUnspecified)

	return &ParagraphStyle{
		textAlign:     style.TextAlignCompanion.TakeOrElse(s.textAlign, DefaultTextAlign),
		textDirection: style.ResolveTextDirection(direction, s.textDirection),
		lineHeight:    unit.TextUnitCompanion.TakeOrElse(s.lineHeight, DefaultLineHeight),
		lineBreak:     style.LineBreakCompanion.TakeOrElse(s.lineBreak, DefaultLineBreak),
	}
}

// ParagraphStyleCompanion is the only exported value of this file's sentinel family member.
var ParagraphStyleCompanion = companions.NewStructCompanion(
	func() *ParagraphStyle { return paragraphStyleUnspecified },
	stringParagraphStyle,
	isSpecifiedParagraphStyle,
	isUnspecifiedParagraphStyle,
	takeOrElseParagraphStyle,
	mergeParagraphStyle,
	sameParagraphStyle,
	semanticEqualParagraphStyle,
	equalParagraphStyle,
	coalesceParagraphStyle,
	copyParagraphStyle,
)

// ensureMutableParagraphStyle panics if ts is a sentinel value.
// Call this at the start of any function that mutates a SpanStyle to fail-fast on misuse.
func ensureMutableParagraphStyle(s *ParagraphStyle) {
	if ParagraphStyleCompanion.IsUnspecified(s) {
		panic("attempt to mutate sentinel; use XCompanion.Copy first")
	}
}

// SpanStyleFromOptions builds a SpanStyle from options applied over the sentinel.
func ParagraphStyleFromOptions(options ...ParagraphStyleOption) *ParagraphStyle {
	return ParagraphStyleCompanion.Copy(paragraphStyleUnspecified, options...)
}
