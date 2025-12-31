package modifiers

import "github.com/zodimo/go-compose/compose/ui/next/text"

type TextSubstitutionValue struct {
	Original              text.AnnotatedString
	Substitution          text.AnnotatedString
	IsShowingSubstitution bool
	// LayoutCache *MultiParagraphLayoutCache
}
