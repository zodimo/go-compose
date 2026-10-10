package tokens

import (
	"github.com/zodimo/go-compose/compose/ui/text"
)

var defaultTextStyle *text.TextStyle = text.TextStyleCompanion.Unspecified()

type TypographyTokensData struct {
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

var TypographyTokens = TypographyTokensData{
	BodyLarge: text.CopyTextStyleWithOptions(defaultTextStyle,
		text.WithFontFamily(TypeScaleTokens.BodyLargeFont),
		text.WithFontWeight(TypeScaleTokens.BodyLargeWeight),
		text.WithFontSize(TypeScaleTokens.BodyLargeSize),
		text.WithLineHeight(TypeScaleTokens.BodyLargeLineHeight),
		text.WithLetterSpacing(TypeScaleTokens.BodyLargeTracking),
	),

	BodyMedium: text.CopyTextStyleWithOptions(defaultTextStyle,
		text.WithFontFamily(TypeScaleTokens.BodyMediumFont),
		text.WithFontWeight(TypeScaleTokens.BodyMediumWeight),
		text.WithFontSize(TypeScaleTokens.BodyMediumSize),
		text.WithLineHeight(TypeScaleTokens.BodyMediumLineHeight),
		text.WithLetterSpacing(TypeScaleTokens.BodyMediumTracking),
	),

	BodySmall: text.CopyTextStyleWithOptions(defaultTextStyle,
		text.WithFontFamily(TypeScaleTokens.BodySmallFont),
		text.WithFontWeight(TypeScaleTokens.BodySmallWeight),
		text.WithFontSize(TypeScaleTokens.BodySmallSize),
		text.WithLineHeight(TypeScaleTokens.BodySmallLineHeight),
		text.WithLetterSpacing(TypeScaleTokens.BodySmallTracking),
	),

	DisplayLarge: text.CopyTextStyleWithOptions(defaultTextStyle,
		text.WithFontFamily(TypeScaleTokens.DisplayLargeFont),
		text.WithFontWeight(TypeScaleTokens.DisplayLargeWeight),
		text.WithFontSize(TypeScaleTokens.DisplayLargeSize),
		text.WithLineHeight(TypeScaleTokens.DisplayLargeLineHeight),
		text.WithLetterSpacing(TypeScaleTokens.DisplayLargeTracking),
	),

	DisplayMedium: text.CopyTextStyleWithOptions(defaultTextStyle,
		text.WithFontFamily(TypeScaleTokens.DisplayMediumFont),
		text.WithFontWeight(TypeScaleTokens.DisplayMediumWeight),
		text.WithFontSize(TypeScaleTokens.DisplayMediumSize),
		text.WithLineHeight(TypeScaleTokens.DisplayMediumLineHeight),
		text.WithLetterSpacing(TypeScaleTokens.DisplayMediumTracking),
	),

	DisplaySmall: text.CopyTextStyleWithOptions(defaultTextStyle,
		text.WithFontFamily(TypeScaleTokens.DisplaySmallFont),
		text.WithFontWeight(TypeScaleTokens.DisplaySmallWeight),
		text.WithFontSize(TypeScaleTokens.DisplaySmallSize),
		text.WithLineHeight(TypeScaleTokens.DisplaySmallLineHeight),
		text.WithLetterSpacing(TypeScaleTokens.DisplaySmallTracking),
	),

	HeadlineLarge: text.CopyTextStyleWithOptions(defaultTextStyle,
		text.WithFontFamily(TypeScaleTokens.HeadlineLargeFont),
		text.WithFontWeight(TypeScaleTokens.HeadlineLargeWeight),
		text.WithFontSize(TypeScaleTokens.HeadlineLargeSize),
		text.WithLineHeight(TypeScaleTokens.HeadlineLargeLineHeight),
		text.WithLetterSpacing(TypeScaleTokens.HeadlineLargeTracking),
	),
	HeadlineMedium: text.CopyTextStyleWithOptions(defaultTextStyle,
		text.WithFontFamily(TypeScaleTokens.HeadlineMediumFont),
		text.WithFontWeight(TypeScaleTokens.HeadlineMediumWeight),
		text.WithFontSize(TypeScaleTokens.HeadlineMediumSize),
		text.WithLineHeight(TypeScaleTokens.HeadlineMediumLineHeight),
		text.WithLetterSpacing(TypeScaleTokens.HeadlineMediumTracking),
	),
	HeadlineSmall: text.CopyTextStyleWithOptions(defaultTextStyle,
		text.WithFontFamily(TypeScaleTokens.HeadlineSmallFont),
		text.WithFontWeight(TypeScaleTokens.HeadlineSmallWeight),
		text.WithFontSize(TypeScaleTokens.HeadlineSmallSize),
		text.WithLineHeight(TypeScaleTokens.HeadlineSmallLineHeight),
		text.WithLetterSpacing(TypeScaleTokens.HeadlineSmallTracking),
	),

	LabelLarge: text.CopyTextStyleWithOptions(defaultTextStyle,
		text.WithFontFamily(TypeScaleTokens.LabelLargeFont),
		text.WithFontWeight(TypeScaleTokens.LabelLargeWeight),
		text.WithFontSize(TypeScaleTokens.LabelLargeSize),
		text.WithLineHeight(TypeScaleTokens.LabelLargeLineHeight),
		text.WithLetterSpacing(TypeScaleTokens.LabelLargeTracking),
	),
	LabelMedium: text.CopyTextStyleWithOptions(defaultTextStyle,
		text.WithFontFamily(TypeScaleTokens.LabelMediumFont),
		text.WithFontWeight(TypeScaleTokens.LabelMediumWeight),
		text.WithFontSize(TypeScaleTokens.LabelMediumSize),
		text.WithLineHeight(TypeScaleTokens.LabelMediumLineHeight),
		text.WithLetterSpacing(TypeScaleTokens.LabelMediumTracking),
	),
	LabelSmall: text.CopyTextStyleWithOptions(defaultTextStyle,
		text.WithFontFamily(TypeScaleTokens.LabelSmallFont),
		text.WithFontWeight(TypeScaleTokens.LabelSmallWeight),
		text.WithFontSize(TypeScaleTokens.LabelSmallSize),
		text.WithLineHeight(TypeScaleTokens.LabelSmallLineHeight),
		text.WithLetterSpacing(TypeScaleTokens.LabelSmallTracking),
	),

	TitleLarge: text.CopyTextStyleWithOptions(defaultTextStyle,
		text.WithFontFamily(TypeScaleTokens.TitleLargeFont),
		text.WithFontWeight(TypeScaleTokens.TitleLargeWeight),
		text.WithFontSize(TypeScaleTokens.TitleLargeSize),
		text.WithLineHeight(TypeScaleTokens.TitleLargeLineHeight),
		text.WithLetterSpacing(TypeScaleTokens.TitleLargeTracking),
	),
	TitleMedium: text.CopyTextStyleWithOptions(defaultTextStyle,
		text.WithFontFamily(TypeScaleTokens.TitleMediumFont),
		text.WithFontWeight(TypeScaleTokens.TitleMediumWeight),
		text.WithFontSize(TypeScaleTokens.TitleMediumSize),
		text.WithLineHeight(TypeScaleTokens.TitleMediumLineHeight),
		text.WithLetterSpacing(TypeScaleTokens.TitleMediumTracking),
	),
	TitleSmall: text.CopyTextStyleWithOptions(defaultTextStyle,
		text.WithFontFamily(TypeScaleTokens.TitleSmallFont),
		text.WithFontWeight(TypeScaleTokens.TitleSmallWeight),
		text.WithFontSize(TypeScaleTokens.TitleSmallSize),
		text.WithLineHeight(TypeScaleTokens.TitleSmallLineHeight),
		text.WithLetterSpacing(TypeScaleTokens.TitleSmallTracking),
	),

	//Emphasized

	BodyLargeEmphasized: text.CopyTextStyleWithOptions(defaultTextStyle,
		text.WithFontFamily(TypeScaleTokens.BodyLargeEmphasizedFont),
		text.WithFontWeight(TypeScaleTokens.BodyLargeEmphasizedWeight),
		text.WithFontSize(TypeScaleTokens.BodyLargeEmphasizedSize),
		text.WithLineHeight(TypeScaleTokens.BodyLargeEmphasizedLineHeight),
		text.WithLetterSpacing(TypeScaleTokens.BodyLargeEmphasizedTracking),
	),
	BodyMediumEmphasized: text.CopyTextStyleWithOptions(defaultTextStyle,
		text.WithFontFamily(TypeScaleTokens.BodyMediumEmphasizedFont),
		text.WithFontWeight(TypeScaleTokens.BodyMediumEmphasizedWeight),
		text.WithFontSize(TypeScaleTokens.BodyMediumEmphasizedSize),
		text.WithLineHeight(TypeScaleTokens.BodyMediumEmphasizedLineHeight),
		text.WithLetterSpacing(TypeScaleTokens.BodyMediumEmphasizedTracking),
	),
	BodySmallEmphasized: text.CopyTextStyleWithOptions(defaultTextStyle,
		text.WithFontFamily(TypeScaleTokens.BodySmallEmphasizedFont),
		text.WithFontWeight(TypeScaleTokens.BodySmallEmphasizedWeight),
		text.WithFontSize(TypeScaleTokens.BodySmallEmphasizedSize),
		text.WithLineHeight(TypeScaleTokens.BodySmallEmphasizedLineHeight),
		text.WithLetterSpacing(TypeScaleTokens.BodySmallEmphasizedTracking),
	),

	DisplayLargeEmphasized: text.CopyTextStyleWithOptions(defaultTextStyle,
		text.WithFontFamily(TypeScaleTokens.DisplayLargeEmphasizedFont),
		text.WithFontWeight(TypeScaleTokens.DisplayLargeEmphasizedWeight),
		text.WithFontSize(TypeScaleTokens.DisplayLargeEmphasizedSize),
		text.WithLineHeight(TypeScaleTokens.DisplayLargeEmphasizedLineHeight),
		text.WithLetterSpacing(TypeScaleTokens.DisplayLargeEmphasizedTracking),
	),
	DisplayMediumEmphasized: text.CopyTextStyleWithOptions(defaultTextStyle,
		text.WithFontFamily(TypeScaleTokens.DisplayMediumEmphasizedFont),
		text.WithFontWeight(TypeScaleTokens.DisplayMediumEmphasizedWeight),
		text.WithFontSize(TypeScaleTokens.DisplayMediumEmphasizedSize),
		text.WithLineHeight(TypeScaleTokens.DisplayMediumEmphasizedLineHeight),
		text.WithLetterSpacing(TypeScaleTokens.DisplayMediumEmphasizedTracking),
	),
	DisplaySmallEmphasized: text.CopyTextStyleWithOptions(defaultTextStyle,
		text.WithFontFamily(TypeScaleTokens.DisplaySmallEmphasizedFont),
		text.WithFontWeight(TypeScaleTokens.DisplaySmallEmphasizedWeight),
		text.WithFontSize(TypeScaleTokens.DisplaySmallEmphasizedSize),
		text.WithLineHeight(TypeScaleTokens.DisplaySmallEmphasizedLineHeight),
		text.WithLetterSpacing(TypeScaleTokens.DisplaySmallEmphasizedTracking),
	),

	HeadlineLargeEmphasized: text.CopyTextStyleWithOptions(defaultTextStyle,
		text.WithFontFamily(TypeScaleTokens.HeadlineLargeEmphasizedFont),
		text.WithFontWeight(TypeScaleTokens.HeadlineLargeEmphasizedWeight),
		text.WithFontSize(TypeScaleTokens.HeadlineLargeEmphasizedSize),
		text.WithLineHeight(TypeScaleTokens.HeadlineLargeEmphasizedLineHeight),
		text.WithLetterSpacing(TypeScaleTokens.HeadlineLargeEmphasizedTracking),
	),
	HeadlineMediumEmphasized: text.CopyTextStyleWithOptions(defaultTextStyle,
		text.WithFontFamily(TypeScaleTokens.HeadlineMediumEmphasizedFont),
		text.WithFontWeight(TypeScaleTokens.HeadlineMediumEmphasizedWeight),
		text.WithFontSize(TypeScaleTokens.HeadlineMediumEmphasizedSize),
		text.WithLineHeight(TypeScaleTokens.HeadlineMediumEmphasizedLineHeight),
		text.WithLetterSpacing(TypeScaleTokens.HeadlineMediumEmphasizedTracking),
	),
	HeadlineSmallEmphasized: text.CopyTextStyleWithOptions(defaultTextStyle,
		text.WithFontFamily(TypeScaleTokens.HeadlineSmallEmphasizedFont),
		text.WithFontWeight(TypeScaleTokens.HeadlineSmallEmphasizedWeight),
		text.WithFontSize(TypeScaleTokens.HeadlineSmallEmphasizedSize),
		text.WithLineHeight(TypeScaleTokens.HeadlineSmallEmphasizedLineHeight),
		text.WithLetterSpacing(TypeScaleTokens.HeadlineSmallEmphasizedTracking),
	),

	LabelLargeEmphasized: text.CopyTextStyleWithOptions(defaultTextStyle,
		text.WithFontFamily(TypeScaleTokens.LabelLargeEmphasizedFont),
		text.WithFontWeight(TypeScaleTokens.LabelLargeEmphasizedWeight),
		text.WithFontSize(TypeScaleTokens.LabelLargeEmphasizedSize),
		text.WithLineHeight(TypeScaleTokens.LabelLargeEmphasizedLineHeight),
		text.WithLetterSpacing(TypeScaleTokens.LabelLargeEmphasizedTracking),
	),
	LabelMediumEmphasized: text.CopyTextStyleWithOptions(defaultTextStyle,
		text.WithFontFamily(TypeScaleTokens.LabelMediumEmphasizedFont),
		text.WithFontWeight(TypeScaleTokens.LabelMediumEmphasizedWeight),
		text.WithFontSize(TypeScaleTokens.LabelMediumEmphasizedSize),
		text.WithLineHeight(TypeScaleTokens.LabelMediumEmphasizedLineHeight),
		text.WithLetterSpacing(TypeScaleTokens.LabelMediumEmphasizedTracking),
	),
	LabelSmallEmphasized: text.CopyTextStyleWithOptions(defaultTextStyle,
		text.WithFontFamily(TypeScaleTokens.LabelSmallEmphasizedFont),
		text.WithFontWeight(TypeScaleTokens.LabelSmallEmphasizedWeight),
		text.WithFontSize(TypeScaleTokens.LabelSmallEmphasizedSize),
		text.WithLineHeight(TypeScaleTokens.LabelSmallEmphasizedLineHeight),
		text.WithLetterSpacing(TypeScaleTokens.LabelSmallEmphasizedTracking),
	),

	TitleLargeEmphasized: text.CopyTextStyleWithOptions(defaultTextStyle,
		text.WithFontFamily(TypeScaleTokens.TitleLargeEmphasizedFont),
		text.WithFontWeight(TypeScaleTokens.TitleLargeEmphasizedWeight),
		text.WithFontSize(TypeScaleTokens.TitleLargeEmphasizedSize),
		text.WithLineHeight(TypeScaleTokens.TitleLargeEmphasizedLineHeight),
		text.WithLetterSpacing(TypeScaleTokens.TitleLargeEmphasizedTracking),
	),
	TitleMediumEmphasized: text.CopyTextStyleWithOptions(defaultTextStyle,
		text.WithFontFamily(TypeScaleTokens.TitleMediumEmphasizedFont),
		text.WithFontWeight(TypeScaleTokens.TitleMediumEmphasizedWeight),
		text.WithFontSize(TypeScaleTokens.TitleMediumEmphasizedSize),
		text.WithLineHeight(TypeScaleTokens.TitleMediumEmphasizedLineHeight),
		text.WithLetterSpacing(TypeScaleTokens.TitleMediumEmphasizedTracking),
	),
	TitleSmallEmphasized: text.CopyTextStyleWithOptions(defaultTextStyle,
		text.WithFontFamily(TypeScaleTokens.TitleSmallEmphasizedFont),
		text.WithFontWeight(TypeScaleTokens.TitleSmallEmphasizedWeight),
		text.WithFontSize(TypeScaleTokens.TitleSmallEmphasizedSize),
		text.WithLineHeight(TypeScaleTokens.TitleSmallEmphasizedLineHeight),
		text.WithLetterSpacing(TypeScaleTokens.TitleSmallEmphasizedTracking),
	),
}
