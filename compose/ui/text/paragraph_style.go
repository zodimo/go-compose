package text

import (
	"github.com/zodimo/go-compose/compose/ui/text/style"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

var ParagraphStyleUnspecified *ParagraphStyle = &ParagraphStyle{
	TextAlign:       style.TextAlignUnspecified,
	TextDirection:   style.TextDirectionUnspecified,
	LineHeight:      unit.TextUnitUnspecified,
	TextIndent:      nil,
	LineHeightStyle: nil,
	LineBreak:       style.LineBreakUnspecified,
	Hyphens:         style.HyphensUnspecified,
	TextMotion:      nil,
}

var _ Annotation = (*ParagraphStyle)(nil)

// ParagraphStyle configuration for a paragraph.
type ParagraphStyle struct {
	TextAlign       style.TextAlign
	TextDirection   style.TextDirection
	LineHeight      unit.TextUnit
	TextIndent      *style.TextIndent
	PlatformStyle   *PlatformParagraphStyle
	LineHeightStyle *style.LineHeightStyle
	LineBreak       style.LineBreak
	Hyphens         style.Hyphens
	TextMotion      *style.TextMotion
}

func (s ParagraphStyle) isAnnotation() {}

type ParagraphStyleOptions struct {
}

type ParagraphStyleOption = func(*ParagraphStyleOptions)

func (s ParagraphStyle) Copy(options ...ParagraphStyleOption) *ParagraphStyle {
	panic("ParagraphStyle Copy not implemented")
}

func StringParagraphStyle(s *ParagraphStyle) string {
	panic("ParagraphStyle ToString not implemented")
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

	panic("ParagraphStyle SemanticEqualParagraphStyle not implemented")
}

func EqualParagraphStyle(a, b *ParagraphStyle) bool {
	if !SameParagraphStyle(a, b) {
		return SemanticEqualParagraphStyle(a, b)
	}
	return true
}

func MergeParagraphStyle(a, b *ParagraphStyle) *ParagraphStyle {
	a = CoalesceParagraphStyle(a, ParagraphStyleUnspecified)
	b = CoalesceParagraphStyle(b, ParagraphStyleUnspecified)

	if a == ParagraphStyleUnspecified {
		return b
	}
	if b == ParagraphStyleUnspecified {
		return a
	}

	// Both are custom: allocate new merged style
	return &ParagraphStyle{
		TextAlign:       b.TextAlign.TakeOrElse(a.TextAlign),
		TextDirection:   b.TextDirection.TakeOrElse(a.TextDirection),
		LineHeight:      b.LineHeight.TakeOrElse(a.LineHeight),
		TextIndent:      style.TakeOrElseTextIndent(b.TextIndent, a.TextIndent),
		PlatformStyle:   TakeOrElsePlatformParagraphStyle(b.PlatformStyle, a.PlatformStyle),
		LineHeightStyle: style.MergeLineHeightStyle(a.LineHeightStyle, b.LineHeightStyle),
		LineBreak:       b.LineBreak.TakeOrElse(a.LineBreak),
		Hyphens:         b.Hyphens.TakeOrElse(a.Hyphens),
		TextMotion:      style.TakeOrElseTextMotion(b.TextMotion, a.TextMotion),
	}
}

func CoalesceParagraphStyle(ptr, def *ParagraphStyle) *ParagraphStyle {
	if ptr == nil {
		return def
	}
	return ptr
}
