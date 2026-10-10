package material3

import (
	"fmt"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/material3/tokens"
	"github.com/zodimo/go-compose/compose/ui/text"
)

var TypographyUnspecified = &Typography{
	BodyLarge:  text.TextStyleCompanion.Unspecified(),
	BodyMedium: text.TextStyleCompanion.Unspecified(),
	BodySmall:  text.TextStyleCompanion.Unspecified(),

	DisplayLarge:  text.TextStyleCompanion.Unspecified(),
	DisplayMedium: text.TextStyleCompanion.Unspecified(),
	DisplaySmall:  text.TextStyleCompanion.Unspecified(),

	HeadlineLarge:  text.TextStyleCompanion.Unspecified(),
	HeadlineMedium: text.TextStyleCompanion.Unspecified(),
	HeadlineSmall:  text.TextStyleCompanion.Unspecified(),

	LabelLarge:  text.TextStyleCompanion.Unspecified(),
	LabelMedium: text.TextStyleCompanion.Unspecified(),
	LabelSmall:  text.TextStyleCompanion.Unspecified(),

	TitleLarge:  text.TextStyleCompanion.Unspecified(),
	TitleMedium: text.TextStyleCompanion.Unspecified(),
	TitleSmall:  text.TextStyleCompanion.Unspecified(),

	// Emphasized
	BodyLargeEmphasized:  text.TextStyleCompanion.Unspecified(),
	BodyMediumEmphasized: text.TextStyleCompanion.Unspecified(),
	BodySmallEmphasized:  text.TextStyleCompanion.Unspecified(),

	DisplayLargeEmphasized:  text.TextStyleCompanion.Unspecified(),
	DisplayMediumEmphasized: text.TextStyleCompanion.Unspecified(),
	DisplaySmallEmphasized:  text.TextStyleCompanion.Unspecified(),

	HeadlineLargeEmphasized:  text.TextStyleCompanion.Unspecified(),
	HeadlineMediumEmphasized: text.TextStyleCompanion.Unspecified(),
	HeadlineSmallEmphasized:  text.TextStyleCompanion.Unspecified(),

	LabelLargeEmphasized:  text.TextStyleCompanion.Unspecified(),
	LabelMediumEmphasized: text.TextStyleCompanion.Unspecified(),
	LabelSmallEmphasized:  text.TextStyleCompanion.Unspecified(),

	TitleLargeEmphasized:  text.TextStyleCompanion.Unspecified(),
	TitleMediumEmphasized: text.TextStyleCompanion.Unspecified(),
	TitleSmallEmphasized:  text.TextStyleCompanion.Unspecified(),
}

type Typography struct {
	BodyLarge  *text.TextStyle
	BodyMedium *text.TextStyle
	BodySmall  *text.TextStyle

	DisplayLarge  *text.TextStyle
	DisplayMedium *text.TextStyle
	DisplaySmall  *text.TextStyle

	HeadlineLarge  *text.TextStyle
	HeadlineMedium *text.TextStyle
	HeadlineSmall  *text.TextStyle

	LabelLarge  *text.TextStyle
	LabelMedium *text.TextStyle
	LabelSmall  *text.TextStyle

	TitleLarge  *text.TextStyle
	TitleMedium *text.TextStyle
	TitleSmall  *text.TextStyle

	// Emphasized
	BodyLargeEmphasized  *text.TextStyle
	BodyMediumEmphasized *text.TextStyle
	BodySmallEmphasized  *text.TextStyle

	DisplayLargeEmphasized  *text.TextStyle
	DisplayMediumEmphasized *text.TextStyle
	DisplaySmallEmphasized  *text.TextStyle

	HeadlineLargeEmphasized  *text.TextStyle
	HeadlineMediumEmphasized *text.TextStyle
	HeadlineSmallEmphasized  *text.TextStyle

	LabelLargeEmphasized  *text.TextStyle
	LabelMediumEmphasized *text.TextStyle
	LabelSmallEmphasized  *text.TextStyle

	TitleLargeEmphasized  *text.TextStyle
	TitleMediumEmphasized *text.TextStyle
	TitleSmallEmphasized  *text.TextStyle
}

func (t *Typography) FromToken(value tokens.TypographyTokenKey) *text.TextStyle {
	switch value {
	case tokens.TypographyTokenKeyUnspecified:
		return text.TextStyleCompanion.Unspecified()
	case tokens.TypographyTokenKeyBodyLarge:
		return t.BodyLarge
	case tokens.TypographyTokenKeyBodyMedium:
		return t.BodyMedium
	case tokens.TypographyTokenKeyBodySmall:
		return t.BodySmall
	case tokens.TypographyTokenKeyDisplayLarge:
		return t.DisplayLarge
	case tokens.TypographyTokenKeyDisplayMedium:
		return t.DisplayMedium
	case tokens.TypographyTokenKeyDisplaySmall:
		return t.DisplaySmall
	case tokens.TypographyTokenKeyHeadlineLarge:
		return t.HeadlineLarge
	case tokens.TypographyTokenKeyHeadlineMedium:
		return t.HeadlineMedium
	case tokens.TypographyTokenKeyHeadlineSmall:
		return t.HeadlineSmall
	case tokens.TypographyTokenKeyLabelLarge:
		return t.LabelLarge
	case tokens.TypographyTokenKeyLabelMedium:
		return t.LabelMedium
	case tokens.TypographyTokenKeyLabelSmall:
		return t.LabelSmall
	case tokens.TypographyTokenKeyTitleLarge:
		return t.TitleLarge
	case tokens.TypographyTokenKeyTitleMedium:
		return t.TitleMedium
	case tokens.TypographyTokenKeyTitleSmall:
		return t.TitleSmall
	case tokens.TypographyTokenKeyBodyLargeEmphasized:
		return t.BodyLargeEmphasized
	case tokens.TypographyTokenKeyBodyMediumEmphasized:
		return t.BodyMediumEmphasized
	case tokens.TypographyTokenKeyBodySmallEmphasized:
		return t.BodySmallEmphasized
	case tokens.TypographyTokenKeyDisplayLargeEmphasized:
		return t.DisplayLargeEmphasized
	case tokens.TypographyTokenKeyDisplayMediumEmphasized:
		return t.DisplayMediumEmphasized
	case tokens.TypographyTokenKeyDisplaySmallEmphasized:
		return t.DisplaySmallEmphasized
	case tokens.TypographyTokenKeyHeadlineLargeEmphasized:
		return t.HeadlineLargeEmphasized
	case tokens.TypographyTokenKeyHeadlineMediumEmphasized:
		return t.HeadlineMediumEmphasized
	case tokens.TypographyTokenKeyHeadlineSmallEmphasized:
		return t.HeadlineSmallEmphasized
	case tokens.TypographyTokenKeyLabelLargeEmphasized:
		return t.LabelLargeEmphasized
	case tokens.TypographyTokenKeyLabelMediumEmphasized:
		return t.LabelMediumEmphasized
	case tokens.TypographyTokenKeyLabelSmallEmphasized:
		return t.LabelSmallEmphasized
	case tokens.TypographyTokenKeyTitleLargeEmphasized:
		return t.TitleLargeEmphasized
	case tokens.TypographyTokenKeyTitleMediumEmphasized:
		return t.TitleMediumEmphasized
	case tokens.TypographyTokenKeyTitleSmallEmphasized:
		return t.TitleSmallEmphasized
	default:
		panic(fmt.Sprintf("unknown typography token: %s", value))
	}
}

var LocalTypography = compose.CompositionLocalOf(func() *Typography {
	return DefaultTypography
})
