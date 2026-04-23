package tokens

import (
	materialTokens "github.com/zodimo/go-compose/compose/material/tokens"
	"github.com/zodimo/go-compose/compose/ui/text/font"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

var TypefaceTokens = materialTokens.TypefaceTokens

var TypeScaleTokens = TypeScaleTokensData{
	BodyLargeFont:       TypefaceTokens.Plain,
	BodyLargeLineHeight: unit.Sp(24.0),
	BodyLargeSize:       unit.Sp(16),
	BodyLargeTracking:   unit.Sp(0.5),
	BodyLargeWeight:     TypefaceTokens.WeightRegular,

	BodyMediumFont:       TypefaceTokens.Plain,
	BodyMediumLineHeight: unit.Sp(20.0),
	BodyMediumSize:       unit.Sp(14),
	BodyMediumTracking:   unit.Sp(0.25),
	BodyMediumWeight:     TypefaceTokens.WeightRegular,

	BodySmallFont:       TypefaceTokens.Plain,
	BodySmallLineHeight: unit.Sp(16.0),
	BodySmallSize:       unit.Sp(12),
	BodySmallTracking:   unit.Sp(0.4),
	BodySmallWeight:     TypefaceTokens.WeightRegular,

	DisplayLargeFont:       TypefaceTokens.Brand,
	DisplayLargeLineHeight: unit.Sp(64.0),
	DisplayLargeSize:       unit.Sp(57),
	DisplayLargeTracking:   unit.Sp(-0.25),
	DisplayLargeWeight:     TypefaceTokens.WeightRegular,

	DisplayMediumFont:       TypefaceTokens.Brand,
	DisplayMediumLineHeight: unit.Sp(52.0),
	DisplayMediumSize:       unit.Sp(45),
	DisplayMediumTracking:   unit.Sp(0.0),
	DisplayMediumWeight:     TypefaceTokens.WeightRegular,

	DisplaySmallFont:       TypefaceTokens.Brand,
	DisplaySmallLineHeight: unit.Sp(44.0),
	DisplaySmallSize:       unit.Sp(36),
	DisplaySmallTracking:   unit.Sp(0.0),
	DisplaySmallWeight:     TypefaceTokens.WeightRegular,

	HeadlineLargeFont:       TypefaceTokens.Brand,
	HeadlineLargeLineHeight: unit.Sp(40.0),
	HeadlineLargeSize:       unit.Sp(32),
	HeadlineLargeTracking:   unit.Sp(0.0),
	HeadlineLargeWeight:     TypefaceTokens.WeightRegular,

	HeadlineMediumFont:       TypefaceTokens.Brand,
	HeadlineMediumLineHeight: unit.Sp(36.0),
	HeadlineMediumSize:       unit.Sp(28),
	HeadlineMediumTracking:   unit.Sp(0.0),
	HeadlineMediumWeight:     TypefaceTokens.WeightRegular,

	HeadlineSmallFont:       TypefaceTokens.Brand,
	HeadlineSmallLineHeight: unit.Sp(32.0),
	HeadlineSmallSize:       unit.Sp(24),
	HeadlineSmallTracking:   unit.Sp(0.0),
	HeadlineSmallWeight:     TypefaceTokens.WeightRegular,

	LabelLargeFont:       TypefaceTokens.Plain,
	LabelLargeLineHeight: unit.Sp(20.0),
	LabelLargeSize:       unit.Sp(14),
	LabelLargeTracking:   unit.Sp(0.1),
	LabelLargeWeight:     TypefaceTokens.WeightMedium,

	LabelMediumFont:       TypefaceTokens.Plain,
	LabelMediumLineHeight: unit.Sp(16.0),
	LabelMediumSize:       unit.Sp(12),
	LabelMediumTracking:   unit.Sp(0.5),
	LabelMediumWeight:     TypefaceTokens.WeightMedium,

	LabelSmallFont:       TypefaceTokens.Plain,
	LabelSmallLineHeight: unit.Sp(16.0),
	LabelSmallSize:       unit.Sp(11),
	LabelSmallTracking:   unit.Sp(0.5),
	LabelSmallWeight:     TypefaceTokens.WeightMedium,

	TitleLargeFont:       TypefaceTokens.Brand,
	TitleLargeLineHeight: unit.Sp(28.0),
	TitleLargeSize:       unit.Sp(22),
	TitleLargeTracking:   unit.Sp(0.0),
	TitleLargeWeight:     TypefaceTokens.WeightRegular,

	TitleMediumFont:       TypefaceTokens.Plain,
	TitleMediumLineHeight: unit.Sp(24.0),
	TitleMediumSize:       unit.Sp(16),
	TitleMediumTracking:   unit.Sp(0.15),
	TitleMediumWeight:     TypefaceTokens.WeightMedium,

	TitleSmallFont:       TypefaceTokens.Plain,
	TitleSmallLineHeight: unit.Sp(20.0),
	TitleSmallSize:       unit.Sp(14),
	TitleSmallTracking:   unit.Sp(0.1),
	TitleSmallWeight:     TypefaceTokens.WeightMedium,

	// Emphasized tokens
	BodyLargeEmphasizedFont:       TypefaceTokens.Plain,
	BodyLargeEmphasizedLineHeight: unit.Sp(24.0),
	BodyLargeEmphasizedSize:       unit.Sp(16),
	BodyLargeEmphasizedTracking:   unit.Sp(0.15),
	BodyLargeEmphasizedWeight:     TypefaceTokens.WeightMedium,

	BodyMediumEmphasizedFont:       TypefaceTokens.Plain,
	BodyMediumEmphasizedLineHeight: unit.Sp(20.0),
	BodyMediumEmphasizedSize:       unit.Sp(14),
	BodyMediumEmphasizedTracking:   unit.Sp(0.25),
	BodyMediumEmphasizedWeight:     TypefaceTokens.WeightMedium,

	BodySmallEmphasizedFont:       TypefaceTokens.Plain,
	BodySmallEmphasizedLineHeight: unit.Sp(16.0),
	BodySmallEmphasizedSize:       unit.Sp(12),
	BodySmallEmphasizedTracking:   unit.Sp(0.4),
	BodySmallEmphasizedWeight:     TypefaceTokens.WeightMedium,

	DisplayLargeEmphasizedFont:       TypefaceTokens.Brand,
	DisplayLargeEmphasizedLineHeight: unit.Sp(64.0),
	DisplayLargeEmphasizedSize:       unit.Sp(57),
	DisplayLargeEmphasizedTracking:   unit.Sp(0),
	DisplayLargeEmphasizedWeight:     TypefaceTokens.WeightMedium,

	DisplayMediumEmphasizedFont:       TypefaceTokens.Brand,
	DisplayMediumEmphasizedLineHeight: unit.Sp(52.0),
	DisplayMediumEmphasizedSize:       unit.Sp(45),
	DisplayMediumEmphasizedTracking:   unit.Sp(0),
	DisplayMediumEmphasizedWeight:     TypefaceTokens.WeightMedium,

	DisplaySmallEmphasizedFont:       TypefaceTokens.Brand,
	DisplaySmallEmphasizedLineHeight: unit.Sp(44.0),
	DisplaySmallEmphasizedSize:       unit.Sp(36),
	DisplaySmallEmphasizedTracking:   unit.Sp(0),
	DisplaySmallEmphasizedWeight:     TypefaceTokens.WeightMedium,

	HeadlineLargeEmphasizedFont:       TypefaceTokens.Brand,
	HeadlineLargeEmphasizedLineHeight: unit.Sp(40.0),
	HeadlineLargeEmphasizedSize:       unit.Sp(32),
	HeadlineLargeEmphasizedTracking:   unit.Sp(0),
	HeadlineLargeEmphasizedWeight:     TypefaceTokens.WeightMedium,

	HeadlineMediumEmphasizedFont:       TypefaceTokens.Brand,
	HeadlineMediumEmphasizedLineHeight: unit.Sp(36.0),
	HeadlineMediumEmphasizedSize:       unit.Sp(28),
	HeadlineMediumEmphasizedTracking:   unit.Sp(0),
	HeadlineMediumEmphasizedWeight:     TypefaceTokens.WeightMedium,

	HeadlineSmallEmphasizedFont:       TypefaceTokens.Brand,
	HeadlineSmallEmphasizedLineHeight: unit.Sp(32.0),
	HeadlineSmallEmphasizedSize:       unit.Sp(24),
	HeadlineSmallEmphasizedTracking:   unit.Sp(0),
	HeadlineSmallEmphasizedWeight:     TypefaceTokens.WeightMedium,

	LabelLargeEmphasizedFont:       TypefaceTokens.Plain,
	LabelLargeEmphasizedLineHeight: unit.Sp(20.0),
	LabelLargeEmphasizedSize:       unit.Sp(14),
	LabelLargeEmphasizedTracking:   unit.Sp(0.1),
	LabelLargeEmphasizedWeight:     TypefaceTokens.WeightBold,

	LabelMediumEmphasizedFont:       TypefaceTokens.Plain,
	LabelMediumEmphasizedLineHeight: unit.Sp(16.0),
	LabelMediumEmphasizedSize:       unit.Sp(12),
	LabelMediumEmphasizedTracking:   unit.Sp(0.5),
	LabelMediumEmphasizedWeight:     TypefaceTokens.WeightBold,

	LabelSmallEmphasizedFont:       TypefaceTokens.Plain,
	LabelSmallEmphasizedLineHeight: unit.Sp(16.0),
	LabelSmallEmphasizedSize:       unit.Sp(11),
	LabelSmallEmphasizedTracking:   unit.Sp(0.5),
	LabelSmallEmphasizedWeight:     TypefaceTokens.WeightBold,

	TitleLargeEmphasizedFont:       TypefaceTokens.Brand,
	TitleLargeEmphasizedLineHeight: unit.Sp(28.0),
	TitleLargeEmphasizedSize:       unit.Sp(22),
	TitleLargeEmphasizedTracking:   unit.Sp(0),
	TitleLargeEmphasizedWeight:     TypefaceTokens.WeightMedium,

	TitleMediumEmphasizedFont:       TypefaceTokens.Plain,
	TitleMediumEmphasizedLineHeight: unit.Sp(24.0),
	TitleMediumEmphasizedSize:       unit.Sp(16),
	TitleMediumEmphasizedTracking:   unit.Sp(0.15),
	TitleMediumEmphasizedWeight:     TypefaceTokens.WeightBold,

	TitleSmallEmphasizedFont:       TypefaceTokens.Plain,
	TitleSmallEmphasizedLineHeight: unit.Sp(20.0),
	TitleSmallEmphasizedSize:       unit.Sp(14),
	TitleSmallEmphasizedTracking:   unit.Sp(0.1),
	TitleSmallEmphasizedWeight:     TypefaceTokens.WeightBold,
}

type TypeScaleTokensData struct {
	BodyLargeFont       font.FontFamily
	BodyLargeLineHeight unit.TextUnit
	BodyLargeSize       unit.TextUnit
	BodyLargeTracking   unit.TextUnit
	BodyLargeWeight     font.FontWeight

	BodyMediumFont       font.FontFamily
	BodyMediumLineHeight unit.TextUnit
	BodyMediumSize       unit.TextUnit
	BodyMediumTracking   unit.TextUnit
	BodyMediumWeight     font.FontWeight

	BodySmallFont       font.FontFamily
	BodySmallLineHeight unit.TextUnit
	BodySmallSize       unit.TextUnit
	BodySmallTracking   unit.TextUnit
	BodySmallWeight     font.FontWeight

	DisplayLargeFont       font.FontFamily
	DisplayLargeLineHeight unit.TextUnit
	DisplayLargeSize       unit.TextUnit
	DisplayLargeTracking   unit.TextUnit
	DisplayLargeWeight     font.FontWeight

	DisplayMediumFont       font.FontFamily
	DisplayMediumLineHeight unit.TextUnit
	DisplayMediumSize       unit.TextUnit
	DisplayMediumTracking   unit.TextUnit
	DisplayMediumWeight     font.FontWeight

	DisplaySmallFont       font.FontFamily
	DisplaySmallLineHeight unit.TextUnit
	DisplaySmallSize       unit.TextUnit
	DisplaySmallTracking   unit.TextUnit
	DisplaySmallWeight     font.FontWeight

	HeadlineLargeFont       font.FontFamily
	HeadlineLargeLineHeight unit.TextUnit
	HeadlineLargeSize       unit.TextUnit
	HeadlineLargeTracking   unit.TextUnit
	HeadlineLargeWeight     font.FontWeight

	HeadlineMediumFont       font.FontFamily
	HeadlineMediumLineHeight unit.TextUnit
	HeadlineMediumSize       unit.TextUnit
	HeadlineMediumTracking   unit.TextUnit
	HeadlineMediumWeight     font.FontWeight

	HeadlineSmallFont       font.FontFamily
	HeadlineSmallLineHeight unit.TextUnit
	HeadlineSmallSize       unit.TextUnit
	HeadlineSmallTracking   unit.TextUnit
	HeadlineSmallWeight     font.FontWeight

	LabelLargeFont       font.FontFamily
	LabelLargeLineHeight unit.TextUnit
	LabelLargeSize       unit.TextUnit
	LabelLargeTracking   unit.TextUnit
	LabelLargeWeight     font.FontWeight

	LabelMediumFont       font.FontFamily
	LabelMediumLineHeight unit.TextUnit
	LabelMediumSize       unit.TextUnit
	LabelMediumTracking   unit.TextUnit
	LabelMediumWeight     font.FontWeight

	LabelSmallFont       font.FontFamily
	LabelSmallLineHeight unit.TextUnit
	LabelSmallSize       unit.TextUnit
	LabelSmallTracking   unit.TextUnit
	LabelSmallWeight     font.FontWeight

	TitleLargeFont       font.FontFamily
	TitleLargeLineHeight unit.TextUnit
	TitleLargeSize       unit.TextUnit
	TitleLargeTracking   unit.TextUnit
	TitleLargeWeight     font.FontWeight

	TitleMediumFont       font.FontFamily
	TitleMediumLineHeight unit.TextUnit
	TitleMediumSize       unit.TextUnit
	TitleMediumTracking   unit.TextUnit
	TitleMediumWeight     font.FontWeight

	TitleSmallFont       font.FontFamily
	TitleSmallLineHeight unit.TextUnit
	TitleSmallSize       unit.TextUnit
	TitleSmallTracking   unit.TextUnit
	TitleSmallWeight     font.FontWeight

	// Emphasized tokens
	BodyLargeEmphasizedFont       font.FontFamily
	BodyLargeEmphasizedLineHeight unit.TextUnit
	BodyLargeEmphasizedSize       unit.TextUnit
	BodyLargeEmphasizedTracking   unit.TextUnit
	BodyLargeEmphasizedWeight     font.FontWeight

	BodyMediumEmphasizedFont       font.FontFamily
	BodyMediumEmphasizedLineHeight unit.TextUnit
	BodyMediumEmphasizedSize       unit.TextUnit
	BodyMediumEmphasizedTracking   unit.TextUnit
	BodyMediumEmphasizedWeight     font.FontWeight

	BodySmallEmphasizedFont       font.FontFamily
	BodySmallEmphasizedLineHeight unit.TextUnit
	BodySmallEmphasizedSize       unit.TextUnit
	BodySmallEmphasizedTracking   unit.TextUnit
	BodySmallEmphasizedWeight     font.FontWeight

	DisplayLargeEmphasizedFont       font.FontFamily
	DisplayLargeEmphasizedLineHeight unit.TextUnit
	DisplayLargeEmphasizedSize       unit.TextUnit
	DisplayLargeEmphasizedTracking   unit.TextUnit
	DisplayLargeEmphasizedWeight     font.FontWeight

	DisplayMediumEmphasizedFont       font.FontFamily
	DisplayMediumEmphasizedLineHeight unit.TextUnit
	DisplayMediumEmphasizedSize       unit.TextUnit
	DisplayMediumEmphasizedTracking   unit.TextUnit
	DisplayMediumEmphasizedWeight     font.FontWeight

	DisplaySmallEmphasizedFont       font.FontFamily
	DisplaySmallEmphasizedLineHeight unit.TextUnit
	DisplaySmallEmphasizedSize       unit.TextUnit
	DisplaySmallEmphasizedTracking   unit.TextUnit
	DisplaySmallEmphasizedWeight     font.FontWeight

	HeadlineLargeEmphasizedFont       font.FontFamily
	HeadlineLargeEmphasizedLineHeight unit.TextUnit
	HeadlineLargeEmphasizedSize       unit.TextUnit
	HeadlineLargeEmphasizedTracking   unit.TextUnit
	HeadlineLargeEmphasizedWeight     font.FontWeight

	HeadlineMediumEmphasizedFont       font.FontFamily
	HeadlineMediumEmphasizedLineHeight unit.TextUnit
	HeadlineMediumEmphasizedSize       unit.TextUnit
	HeadlineMediumEmphasizedTracking   unit.TextUnit
	HeadlineMediumEmphasizedWeight     font.FontWeight

	HeadlineSmallEmphasizedFont       font.FontFamily
	HeadlineSmallEmphasizedLineHeight unit.TextUnit
	HeadlineSmallEmphasizedSize       unit.TextUnit
	HeadlineSmallEmphasizedTracking   unit.TextUnit
	HeadlineSmallEmphasizedWeight     font.FontWeight

	LabelLargeEmphasizedFont       font.FontFamily
	LabelLargeEmphasizedLineHeight unit.TextUnit
	LabelLargeEmphasizedSize       unit.TextUnit
	LabelLargeEmphasizedTracking   unit.TextUnit
	LabelLargeEmphasizedWeight     font.FontWeight

	LabelMediumEmphasizedFont       font.FontFamily
	LabelMediumEmphasizedLineHeight unit.TextUnit
	LabelMediumEmphasizedSize       unit.TextUnit
	LabelMediumEmphasizedTracking   unit.TextUnit
	LabelMediumEmphasizedWeight     font.FontWeight

	LabelSmallEmphasizedFont       font.FontFamily
	LabelSmallEmphasizedLineHeight unit.TextUnit
	LabelSmallEmphasizedSize       unit.TextUnit
	LabelSmallEmphasizedTracking   unit.TextUnit
	LabelSmallEmphasizedWeight     font.FontWeight

	TitleLargeEmphasizedFont       font.FontFamily
	TitleLargeEmphasizedLineHeight unit.TextUnit
	TitleLargeEmphasizedSize       unit.TextUnit
	TitleLargeEmphasizedTracking   unit.TextUnit
	TitleLargeEmphasizedWeight     font.FontWeight

	TitleMediumEmphasizedFont       font.FontFamily
	TitleMediumEmphasizedLineHeight unit.TextUnit
	TitleMediumEmphasizedSize       unit.TextUnit
	TitleMediumEmphasizedTracking   unit.TextUnit
	TitleMediumEmphasizedWeight     font.FontWeight

	TitleSmallEmphasizedFont       font.FontFamily
	TitleSmallEmphasizedLineHeight unit.TextUnit
	TitleSmallEmphasizedSize       unit.TextUnit
	TitleSmallEmphasizedTracking   unit.TextUnit
	TitleSmallEmphasizedWeight     font.FontWeight
}
