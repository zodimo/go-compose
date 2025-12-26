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

func ParagraphStyleWithTextIndent(textIndent *style.TextIndent) ParagraphStyleOption {
	return func(opts *ParagraphStyleOptions) {
		opts.TextIndent = textIndent
	}
}

func ParagraphStyleWithPlatformStyle(platformStyle *PlatformParagraphStyle) ParagraphStyleOption {
	return func(opts *ParagraphStyleOptions) {
		opts.PlatformStyle = platformStyle
	}
}

func ParagraphStyleWithLineHeightStyle(lineHeightStyle *style.LineHeightStyle) ParagraphStyleOption {
	return func(opts *ParagraphStyleOptions) {
		opts.LineHeightStyle = lineHeightStyle
	}
}

func ParagraphStyleWithLineBreak(lineBreak style.LineBreak) ParagraphStyleOption {
	return func(opts *ParagraphStyleOptions) {
		opts.LineBreak = lineBreak
	}
}

func ParagraphStyleWithHyphens(hyphens style.Hyphens) ParagraphStyleOption {
	return func(opts *ParagraphStyleOptions) {
		opts.Hyphens = hyphens
	}
}

func ParagraphStyleWithTextMotion(textMotion *style.TextMotion) ParagraphStyleOption {
	return func(opts *ParagraphStyleOptions) {
		opts.TextMotion = textMotion
	}
}

type ParagraphStyleOption = func(*ParagraphStyleOptions)

func (s ParagraphStyle) Copy(options ...ParagraphStyleOption) *ParagraphStyle {
	opts := &ParagraphStyleOptions{
		TextAlign:       style.TextAlignUnspecified,
		TextDirection:   style.TextDirectionUnspecified,
		LineHeight:      unit.TextUnitUnspecified,
		TextIndent:      nil,
		PlatformStyle:   nil,
		LineHeightStyle: nil,
		LineBreak:       style.LineBreakUnspecified,
		Hyphens:         style.HyphensUnspecified,
		TextMotion:      nil,
	}
	for _, option := range options {
		option(opts)
	}
	return &ParagraphStyle{
		TextAlign:       opts.TextAlign.TakeOrElse(s.TextAlign),
		TextDirection:   opts.TextDirection.TakeOrElse(s.TextDirection),
		LineHeight:      opts.LineHeight.TakeOrElse(s.LineHeight),
		TextIndent:      style.TakeOrElseTextIndent(opts.TextIndent, s.TextIndent),
		PlatformStyle:   TakeOrElsePlatformParagraphStyle(opts.PlatformStyle, s.PlatformStyle),
		LineHeightStyle: style.MergeLineHeightStyle(opts.LineHeightStyle, s.LineHeightStyle),
		LineBreak:       opts.LineBreak.TakeOrElse(s.LineBreak),
		Hyphens:         opts.Hyphens.TakeOrElse(s.Hyphens),
		TextMotion:      style.TakeOrElseTextMotion(opts.TextMotion, s.TextMotion),
	}
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

	return a.TextAlign == b.TextAlign &&
		a.TextDirection == b.TextDirection &&
		a.LineHeight == b.LineHeight &&
		style.EqualTextIndent(a.TextIndent, b.TextIndent) &&
		style.EqualLineHeightStyle(a.LineHeightStyle, b.LineHeightStyle) &&
		a.LineBreak == b.LineBreak &&
		a.Hyphens == b.Hyphens &&
		style.EqualTextMotion(a.TextMotion, b.TextMotion)

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
