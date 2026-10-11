package text

import (
	"github.com/zodimo/go-compose/compose/ui/next/text/style"
	"github.com/zodimo/go-compose/compose/ui/next/text/unit"
	"github.com/zodimo/go-compose/pkg/floatutils/lerp"
	"github.com/zodimo/go-sentinel-helper/sentinel/companions"
)

// paragraphStyleUnspecified is the sentinel for ParagraphStyle. It is
// deliberately package-private: the only public way to obtain it is
// ParagraphStyleCompanion.Unspecified().
var paragraphStyleUnspecified *ParagraphStyle = &ParagraphStyle{
	TextAlign:       style.TextAlignCompanion.Unspecified(),
	TextDirection:   style.TextDirectionCompanion.Unspecified(),
	LineHeight:      unit.TextUnitCompanion.Unspecified(),
	TextIndent:      nil,
	LineHeightStyle: nil,
	LineBreak:       style.LineBreakCompanion.Unspecified(),
	Hyphens:         style.HyphensCompanion.Unspecified(),
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
		TextAlign:       style.TextAlignCompanion.Unspecified(),
		TextDirection:   style.TextDirectionCompanion.Unspecified(),
		LineHeight:      unit.TextUnitCompanion.Unspecified(),
		TextIndent:      nil,
		PlatformStyle:   nil,
		LineHeightStyle: nil,
		LineBreak:       style.LineBreakCompanion.Unspecified(),
		Hyphens:         style.HyphensCompanion.Unspecified(),
		TextMotion:      nil,
	}
	for _, option := range options {
		option(opts)
	}
	return &ParagraphStyle{
		TextAlign:       style.TextAlignCompanion.TakeOrElse(opts.TextAlign, s.TextAlign),
		TextDirection:   style.TextDirectionCompanion.TakeOrElse(opts.TextDirection, s.TextDirection),
		LineHeight:      unit.TextUnitCompanion.TakeOrElse(opts.LineHeight, s.LineHeight),
		TextIndent:      style.TakeOrElseTextIndent(opts.TextIndent, s.TextIndent),
		PlatformStyle:   TakeOrElsePlatformParagraphStyle(opts.PlatformStyle, s.PlatformStyle),
		LineHeightStyle: style.MergeLineHeightStyle(opts.LineHeightStyle, s.LineHeightStyle),
		LineBreak:       style.LineBreakCompanion.TakeOrElse(opts.LineBreak, s.LineBreak),
		Hyphens:         style.HyphensCompanion.TakeOrElse(opts.Hyphens, s.Hyphens),
		TextMotion:      style.TakeOrElseTextMotion(opts.TextMotion, s.TextMotion),
	}
}

func StringParagraphStyle(s *ParagraphStyle) string {
	return "ParagraphStyle(" +
		"TextAlign=" + s.TextAlign.String() + ", " +
		"TextDirection=" + s.TextDirection.String() + ", " +
		"LineHeight=" + s.LineHeight.String() + ", " +
		"TextIndent=" + style.StringTextIndent(s.TextIndent) + ", " +
		"PlatformStyle=" + StringPlatformParagraphStyle(s.PlatformStyle) + ", " +
		"LineHeightStyle=" + style.StringLineHeightStyle(s.LineHeightStyle) + ", " +
		"LineBreak=" + s.LineBreak.String() + ", " +
		"Hyphens=" + s.Hyphens.String() + ", " +
		"TextMotion=" + style.StringTextMotion(s.TextMotion) +
		")"
}

func IsSpecifiedParagraphStyle(s *ParagraphStyle) bool {
	return ParagraphStyleCompanion.IsSpecified(s)
}

func TakeOrElseParagraphStyle(s, def *ParagraphStyle) *ParagraphStyle {
	return takeOrElseParagraphStyle(s, def)
}

// Identity (2 ns)
func SameParagraphStyle(a, b *ParagraphStyle) bool {
	return sameParagraphStyle(a, b)
}

// Semantic equality (field-by-field, 20 ns)
func SemanticEqualParagraphStyle(a, b *ParagraphStyle) bool {
	return semanticEqualParagraphStyle(a, b)
}

func EqualParagraphStyle(a, b *ParagraphStyle) bool {
	return equalParagraphStyle(a, b)
}

func MergeParagraphStyle(a, b *ParagraphStyle) *ParagraphStyle {
	return mergeParagraphStyle(a, b)
}

func CoalesceParagraphStyle(ptr, def *ParagraphStyle) *ParagraphStyle {
	return coalesceParagraphStyle(ptr, def)
}

// === Sentinel operations ===

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

func coalesceParagraphStyle(ptr, def *ParagraphStyle) *ParagraphStyle {
	if ptr == nil {
		return def
	}
	return ptr
}

func mergeParagraphStyle(a, b *ParagraphStyle) *ParagraphStyle {
	a = coalesceParagraphStyle(a, paragraphStyleUnspecified)
	b = coalesceParagraphStyle(b, paragraphStyleUnspecified)

	if a == paragraphStyleUnspecified {
		return b
	}
	if b == paragraphStyleUnspecified {
		return a
	}

	// Both are custom: allocate new merged style. b wins where specified.
	return &ParagraphStyle{
		TextAlign:       style.TextAlignCompanion.TakeOrElse(b.TextAlign, a.TextAlign),
		TextDirection:   style.TextDirectionCompanion.TakeOrElse(b.TextDirection, a.TextDirection),
		LineHeight:      unit.TextUnitCompanion.TakeOrElse(b.LineHeight, a.LineHeight),
		TextIndent:      style.TextIndentCompanion.TakeOrElse(b.TextIndent, a.TextIndent),
		PlatformStyle:   takeOrElsePlatformParagraphStyle(b.PlatformStyle, a.PlatformStyle),
		LineHeightStyle: style.MergeLineHeightStyle(a.LineHeightStyle, b.LineHeightStyle),
		LineBreak:       style.LineBreakCompanion.TakeOrElse(b.LineBreak, a.LineBreak),
		Hyphens:         style.HyphensCompanion.TakeOrElse(b.Hyphens, a.Hyphens),
		TextMotion:      style.TakeOrElseTextMotion(b.TextMotion, a.TextMotion),
	}
}

func stringParagraphStyle(s *ParagraphStyle) string {
	return StringParagraphStyle(s)
}

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

func semanticEqualParagraphStyle(a, b *ParagraphStyle) bool {
	a = coalesceParagraphStyle(a, paragraphStyleUnspecified)
	b = coalesceParagraphStyle(b, paragraphStyleUnspecified)

	if a == paragraphStyleUnspecified && b == paragraphStyleUnspecified {
		return true
	}

	return a.TextAlign == b.TextAlign &&
		a.TextDirection == b.TextDirection &&
		a.LineHeight == b.LineHeight &&
		style.EqualTextIndent(a.TextIndent, b.TextIndent) &&
		style.EqualLineHeightStyle(a.LineHeightStyle, b.LineHeightStyle) &&
		a.LineBreak == b.LineBreak &&
		a.Hyphens == b.Hyphens &&
		style.EqualTextMotion(a.TextMotion, b.TextMotion)
}

func equalParagraphStyle(a, b *ParagraphStyle) bool {
	if !sameParagraphStyle(a, b) {
		return semanticEqualParagraphStyle(a, b)
	}
	return true
}

// copyParagraphStyle returns a copy of s; the source is never mutated.
func copyParagraphStyle(s *ParagraphStyle, options ...func(*ParagraphStyle)) *ParagraphStyle {
	cp := &ParagraphStyle{}
	if isSpecifiedParagraphStyle(s) {
		*cp = *s
	}
	for _, option := range options {
		option(cp)
	}
	return cp
}

// ParagraphStyleCompanion is the single exported sentinel-family symbol for ParagraphStyle.
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

// LerpParagraphStyle interpolates between two ParagraphStyles.
func LerpParagraphStyle(start, stop *ParagraphStyle, fraction float32) *ParagraphStyle {
	return &ParagraphStyle{
		TextAlign:       lerp.LerpDiscrete(start.TextAlign, stop.TextAlign, fraction),
		TextDirection:   lerp.LerpDiscrete(start.TextDirection, stop.TextDirection, fraction),
		LineHeight:      style.LerpTextUnitInheritable(start.LineHeight, stop.LineHeight, fraction),
		TextIndent:      style.LerpTextIndent(start.TextIndent, stop.TextIndent, fraction),
		PlatformStyle:   lerpPlatformParagraphStyle(start.PlatformStyle, stop.PlatformStyle, fraction),
		LineHeightStyle: lerp.LerpDiscrete(start.LineHeightStyle, stop.LineHeightStyle, fraction),
		LineBreak:       lerp.LerpDiscrete(start.LineBreak, stop.LineBreak, fraction),
		Hyphens:         lerp.LerpDiscrete(start.Hyphens, stop.Hyphens, fraction),
		TextMotion:      lerp.LerpDiscrete(start.TextMotion, stop.TextMotion, fraction),
	}
}

func lerpPlatformParagraphStyle(start, stop *PlatformParagraphStyle, fraction float32) *PlatformParagraphStyle {
	if start == nil && stop == nil {
		return nil
	}
	startNonNull := TakeOrElsePlatformParagraphStyle(start, platformParagraphStyleUnspecified)
	stopNonNull := TakeOrElsePlatformParagraphStyle(stop, platformParagraphStyleUnspecified)
	return lerp.LerpDiscrete(startNonNull, stopNonNull, fraction)
}

func ResolveParagraphStyleDefaults(s *ParagraphStyle, direction unit.LayoutDirection) *ParagraphStyle {
	s = CoalesceParagraphStyle(s, paragraphStyleUnspecified)

	return &ParagraphStyle{
		TextAlign:       style.TextAlignCompanion.TakeOrElse(s.TextAlign, DefaultTextAlign),
		TextDirection:   style.ResolveTextDirection(direction, s.TextDirection),
		LineHeight:      unit.TextUnitCompanion.TakeOrElse(s.LineHeight, DefaultLineHeight),
		TextIndent:      style.TakeOrElseTextIndent(s.TextIndent, DefaultTextIndent),
		PlatformStyle:   s.PlatformStyle,
		LineHeightStyle: s.LineHeightStyle,
		LineBreak:       style.LineBreakCompanion.TakeOrElse(s.LineBreak, DefaultLineBreak),
		Hyphens:         style.HyphensCompanion.TakeOrElse(s.Hyphens, DefaultHyphens),
		TextMotion:      style.TakeOrElseTextMotion(s.TextMotion, DefaultTextMotion),
	}
}
