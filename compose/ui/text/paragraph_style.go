package text

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/ui/text/style"
	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/pkg/floatutils/lerp"
	"github.com/zodimo/go-sentinel-helper/sentinel/companions"
)

// paragraphStyleUnspecified is the sentinel for ParagraphStyle. It is deliberately
// package-private: the only public way to obtain it is ParagraphStyleCompanion.Unspecified().
var paragraphStyleUnspecified *ParagraphStyle = &ParagraphStyle{
	textAlign:     style.TextAlignUnspecified,
	textDirection: style.TextDirectionUnspecified,
	lineHeight:    unit.TextUnitUnspecified,
	lineBreak:     style.LineBreakUnspecified,
}

var _ ParaghStyleInterface = (*ParagraphStyle)(nil)

type ParagraphStyle struct {
	textAlign     style.TextAlign
	textDirection style.TextDirection
	lineHeight    unit.TextUnit
	lineBreak     style.LineBreak
}

// ParagraphStyleOptions carries one field per ParagraphStyle property. A field left at
// its sentinel (or nil) means "leave the source value alone"; a field set to any other
// value overrides it. Options start from the sentinel set, so they defer by default.
type ParagraphStyleOptions struct {
	TextAlign     style.TextAlign
	TextDirection style.TextDirection
	LineHeight    unit.TextUnit
	LineBreak     style.LineBreak
}

func newParagraphStyleOptions() ParagraphStyleOptions {
	return ParagraphStyleOptions{
		TextAlign:     style.TextAlignUnspecified,
		TextDirection: style.TextDirectionUnspecified,
		LineHeight:    unit.TextUnitUnspecified,
		LineBreak:     style.LineBreakUnspecified,
	}
}

type ParagraphStyleOption = func(*ParagraphStyleOptions)

// asOption lets an already-populated ParagraphStyleOptions seed a
// ParagraphStyleOptions target, so the whole struct can be passed as a single
// ParagraphStyleOption.
func (o ParagraphStyleOptions) asOption(opts *ParagraphStyleOptions) {
	*opts = o
}

func ParagraphStyleWithTextAlign(textAlign style.TextAlign) ParagraphStyleOption {
	return func(opts *ParagraphStyleOptions) {
		opts.TextAlign = textAlign
	}
}

func ParagraphStyleWithTextDirection(textDirection style.TextDirection) ParagraphStyleOption {
	return func(opts *ParagraphStyleOptions) {
		opts.TextDirection = textDirection
	}
}

func ParagraphStyleWithLineHeight(lineHeight unit.TextUnit) ParagraphStyleOption {
	return func(opts *ParagraphStyleOptions) {
		opts.LineHeight = lineHeight
	}
}

func ParagraphStyleWithLineBreak(lineBreak style.LineBreak) ParagraphStyleOption {
	return func(opts *ParagraphStyleOptions) {
		opts.LineBreak = lineBreak
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

func stringParagraphStyle(s *ParagraphStyle) string {	s = coalesceParagraphStyle(s, paragraphStyleUnspecified)
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
		textAlign:     b.textAlign.TakeOrElse(a.textAlign),
		textDirection: b.textDirection.TakeOrElse(a.textDirection),
		lineHeight:    b.lineHeight.TakeOrElse(a.lineHeight),
		lineBreak:     b.lineBreak.TakeOrElse(a.lineBreak),
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
	work := *s
	for _, option := range options {
		option(&work)
	}
	return &work
}

// copyParagraphStyleWithOptions is the option-constructor form used by callers.
// Options are applied to a fresh ParagraphStyleOptions seeded with sentinels, then
// only the fields the caller actually set are overlaid onto s.
func copyParagraphStyleWithOptions(s *ParagraphStyle, options ...ParagraphStyleOption) *ParagraphStyle {
	s = coalesceParagraphStyle(s, paragraphStyleUnspecified)
	opts := newParagraphStyleOptions()
	for _, option := range options {
		option(&opts)
	}
	return &ParagraphStyle{
		textAlign:     opts.TextAlign.TakeOrElse(s.textAlign),
		textDirection: opts.TextDirection.TakeOrElse(s.textDirection),
		lineHeight:    opts.LineHeight.TakeOrElse(s.lineHeight),
		lineBreak:     opts.LineBreak.TakeOrElse(s.lineBreak),
	}
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
		textAlign:     s.textAlign.TakeOrElse(DefaultTextAlign),
		textDirection: style.ResolveTextDirection(direction, s.textDirection),
		lineHeight:    s.lineHeight.TakeOrElse(DefaultLineHeight),
		lineBreak:     s.lineBreak.TakeOrElse(DefaultLineBreak),
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
