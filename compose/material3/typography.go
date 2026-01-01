package material3

import "github.com/zodimo/go-compose/compose/ui/text"

type Typography struct {
	DisplaySmall  *text.TextStyle
	DisplayMedium *text.TextStyle
	DisplayLarge  *text.TextStyle

	HeadlineSmall  *text.TextStyle
	HeadlineMedium *text.TextStyle
	HeadlineLarge  *text.TextStyle

	TitleSmall  *text.TextStyle
	TitleMedium *text.TextStyle
	TitleLarge  *text.TextStyle

	BodySmall  *text.TextStyle
	BodyMedium *text.TextStyle
	BodyLarge  *text.TextStyle

	LabelSmall  *text.TextStyle
	LabelMedium *text.TextStyle
	LabelLarge  *text.TextStyle

	DisplaySmallEmphasized  *text.TextStyle
	DisplayMediumEmphasized *text.TextStyle
	DisplayLargeEmphasized  *text.TextStyle

	HeadlineSmallEmphasized  *text.TextStyle
	HeadlineMediumEmphasized *text.TextStyle
	HeadlineLargeEmphasized  *text.TextStyle

	TitleSmallEmphasized  *text.TextStyle
	TitleMediumEmphasized *text.TextStyle
	TitleLargeEmphasized  *text.TextStyle

	BodySmallEmphasized  *text.TextStyle
	BodyMediumEmphasized *text.TextStyle
	BodyLargeEmphasized  *text.TextStyle

	LabelSmallEmphasized  *text.TextStyle
	LabelMediumEmphasized *text.TextStyle
	LabelLargeEmphasized  *text.TextStyle
}

func (t *Typography) Copy(options ...TypographyOption) *Typography {
	copy := *t
	for _, option := range options {
		option(&copy)
	}
	return &copy
}
