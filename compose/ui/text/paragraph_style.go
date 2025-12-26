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

func (s ParagraphStyle) isAnnotation() {}

func (s ParagraphStyle) Merge(other *ParagraphStyle) *ParagraphStyle {
	panic("ParagraphStyle Merge not implemented")
}
func (s ParagraphStyle) Copy() *ParagraphStyle {
	panic("ParagraphStyle Copy not implemented")
}
func (s ParagraphStyle) ToString() string {
	panic("ParagraphStyle ToString not implemented")
}
