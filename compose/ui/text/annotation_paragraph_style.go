package text

import (
	"github.com/zodimo/go-compose/compose/ui/text/style"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

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

// NewParagraphStyle creates a new ParagraphStyle.
func NewParagraphStyle() ParagraphStyle {
	return ParagraphStyle{
		TextAlign:     style.TextAlignUnspecified, // Need to ensure this exists or use appropriate default
		TextDirection: style.TextDirectionUnspecified,
		LineHeight:    unit.TextUnitUnspecified,
		LineBreak:     style.LineBreakUnspecified,
		Hyphens:       style.HyphensUnspecified,
	}
}

func (s ParagraphStyle) isAnnotation() {}

func (s ParagraphStyle) Merge(other ParagraphStyle) ParagraphStyle {
	return ParagraphStyle{
		TextAlign:       takeTextAlignOrElse(s.TextAlign, other.TextAlign),
		TextDirection:   takeTextDirectionOrElse(s.TextDirection, other.TextDirection),
		LineHeight:      s.LineHeight.TakeOrElse(other.LineHeight),
		TextIndent:      takeTextIndentOrElse(s.TextIndent, other.TextIndent),
		PlatformStyle:   s.PlatformStyle.Merge(other.PlatformStyle),
		LineHeightStyle: takeLineHeightStyleOrElse(s.LineHeightStyle, other.LineHeightStyle),
		LineBreak:       takeLineBreakOrElse(s.LineBreak, other.LineBreak),
		Hyphens:         takeHyphensOrElse(s.Hyphens, other.Hyphens),
		TextMotion:      takeTextMotionOrElse(s.TextMotion, other.TextMotion),
	}
}

// Equals returns true if the other ParagraphStyle is equal to this one.
func (s ParagraphStyle) Equals(other ParagraphStyle) bool {
	return s.TextAlign == other.TextAlign &&
		s.TextDirection == other.TextDirection &&
		s.LineHeight.Equals(other.LineHeight) &&
		textIndentEquals(s.TextIndent, other.TextIndent) &&
		s.PlatformStyle.Equals(other.PlatformStyle) &&
		lineHeightStyleEquals(s.LineHeightStyle, other.LineHeightStyle) &&
		s.LineBreak == other.LineBreak &&
		s.Hyphens == other.Hyphens &&
		textMotionEquals(s.TextMotion, other.TextMotion)
}

func textIndentEquals(a, b *style.TextIndent) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Equals(*b)
}

func lineHeightStyleEquals(a, b *style.LineHeightStyle) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Equals(*b)
}

func textMotionEquals(a, b *style.TextMotion) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

// Helpers

func takeTextAlignOrElse(a, b style.TextAlign) style.TextAlign {
	if a != style.TextAlignUnspecified {
		return a
	}
	return b
}

func takeTextDirectionOrElse(a, b style.TextDirection) style.TextDirection {
	if a != style.TextDirectionUnspecified {
		return a
	}
	return b
}

func takeTextIndentOrElse(a, b *style.TextIndent) *style.TextIndent {
	if a != nil {
		return a
	}
	return b
}

func takeLineHeightStyleOrElse(a, b *style.LineHeightStyle) *style.LineHeightStyle {
	if a != nil {
		return a
	}
	return b
}

func takeLineBreakOrElse(a, b style.LineBreak) style.LineBreak {
	if a != style.LineBreakUnspecified {
		return a
	}
	return b
}

func takeHyphensOrElse(a, b style.Hyphens) style.Hyphens {
	if a != style.HyphensUnspecified {
		return a
	}
	return b
}

func takeTextMotionOrElse(a, b *style.TextMotion) *style.TextMotion {
	if a != nil {
		return a
	}
	return b
}

// LerpParagraphStyle interpolates between two ParagraphStyles.
func LerpParagraphStyle(start, stop ParagraphStyle, fraction float32) ParagraphStyle {
	return ParagraphStyle{
		TextAlign:       LerpDiscrete(start.TextAlign, stop.TextAlign, fraction),
		TextDirection:   LerpDiscrete(start.TextDirection, stop.TextDirection, fraction),
		LineHeight:      unit.LerpTextUnit(start.LineHeight, stop.LineHeight, fraction),
		TextIndent:      lerpTextIndent(start.TextIndent, stop.TextIndent, fraction),
		PlatformStyle:   lerpPlatformParagraphStyle(start.PlatformStyle, stop.PlatformStyle, fraction),
		LineHeightStyle: LerpDiscrete(start.LineHeightStyle, stop.LineHeightStyle, fraction),
		LineBreak:       LerpDiscrete(start.LineBreak, stop.LineBreak, fraction),
		Hyphens:         LerpDiscrete(start.Hyphens, stop.Hyphens, fraction),
		TextMotion:      LerpDiscrete(start.TextMotion, stop.TextMotion, fraction),
	}
}

func lerpTextIndent(start, stop *style.TextIndent, fraction float32) *style.TextIndent {
	// If both are nil, return nil
	if start == nil && stop == nil {
		return nil
	}
	// If one is nil, use default/empty or fall back to discrete
	// Compose implementation typically reconstructs if one is missing, assuming unit.Sp(0).
	// Let's rely on style.LerpTextIndent but we need to dereference carefully.

	s := style.TextIndentNone
	if start != nil {
		s = *start
	}
	e := style.TextIndentNone
	if stop != nil {
		e = *stop
	}
	result := style.LerpTextIndent(s, e, fraction)
	return &result
}
