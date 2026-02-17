package tokens

type TypographyTokenKey int

const (
	TypographyTokenKeyUnspecified TypographyTokenKey = iota
	TypographyTokenKeyBodyLarge
	TypographyTokenKeyBodyMedium
	TypographyTokenKeyBodySmall
	TypographyTokenKeyDisplayLarge
	TypographyTokenKeyDisplayMedium
	TypographyTokenKeyDisplaySmall
	TypographyTokenKeyHeadlineLarge
	TypographyTokenKeyHeadlineMedium
	TypographyTokenKeyHeadlineSmall
	TypographyTokenKeyLabelLarge
	TypographyTokenKeyLabelMedium
	TypographyTokenKeyLabelSmall
	TypographyTokenKeyTitleLarge
	TypographyTokenKeyTitleMedium
	TypographyTokenKeyTitleSmall
	TypographyTokenKeyBodyLargeEmphasized
	TypographyTokenKeyBodyMediumEmphasized
	TypographyTokenKeyBodySmallEmphasized
	TypographyTokenKeyDisplayLargeEmphasized
	TypographyTokenKeyDisplayMediumEmphasized
	TypographyTokenKeyDisplaySmallEmphasized
	TypographyTokenKeyHeadlineLargeEmphasized
	TypographyTokenKeyHeadlineMediumEmphasized
	TypographyTokenKeyHeadlineSmallEmphasized
	TypographyTokenKeyLabelLargeEmphasized
	TypographyTokenKeyLabelMediumEmphasized
	TypographyTokenKeyLabelSmallEmphasized
	TypographyTokenKeyTitleLargeEmphasized
	TypographyTokenKeyTitleMediumEmphasized
	TypographyTokenKeyTitleSmallEmphasized
)

func (k TypographyTokenKey) String() string {
	switch k {
	case TypographyTokenKeyBodyLarge:
		return "BodyLarge"
	case TypographyTokenKeyBodyMedium:
		return "BodyMedium"
	case TypographyTokenKeyBodySmall:
		return "BodySmall"
	case TypographyTokenKeyDisplayLarge:
		return "DisplayLarge"
	case TypographyTokenKeyDisplayMedium:
		return "DisplayMedium"
	case TypographyTokenKeyDisplaySmall:
		return "DisplaySmall"
	case TypographyTokenKeyHeadlineLarge:
		return "HeadlineLarge"
	case TypographyTokenKeyHeadlineMedium:
		return "HeadlineMedium"
	case TypographyTokenKeyHeadlineSmall:
		return "HeadlineSmall"
	case TypographyTokenKeyLabelLarge:
		return "LabelLarge"
	case TypographyTokenKeyLabelMedium:
		return "LabelMedium"
	case TypographyTokenKeyLabelSmall:
		return "LabelSmall"
	case TypographyTokenKeyTitleLarge:
		return "TitleLarge"
	case TypographyTokenKeyTitleMedium:
		return "TitleMedium"
	case TypographyTokenKeyTitleSmall:
		return "TitleSmall"
	case TypographyTokenKeyBodyLargeEmphasized:
		return "BodyLargeEmphasized"
	case TypographyTokenKeyBodyMediumEmphasized:
		return "BodyMediumEmphasized"
	case TypographyTokenKeyBodySmallEmphasized:
		return "BodySmallEmphasized"
	case TypographyTokenKeyDisplayLargeEmphasized:
		return "DisplayLargeEmphasized"
	case TypographyTokenKeyDisplayMediumEmphasized:
		return "DisplayMediumEmphasized"
	case TypographyTokenKeyDisplaySmallEmphasized:
		return "DisplaySmallEmphasized"
	case TypographyTokenKeyHeadlineLargeEmphasized:
		return "HeadlineLargeEmphasized"
	case TypographyTokenKeyHeadlineMediumEmphasized:
		return "HeadlineMediumEmphasized"
	case TypographyTokenKeyHeadlineSmallEmphasized:
		return "HeadlineSmallEmphasized"
	case TypographyTokenKeyLabelLargeEmphasized:
		return "LabelLargeEmphasized"
	case TypographyTokenKeyLabelMediumEmphasized:
		return "LabelMediumEmphasized"
	case TypographyTokenKeyLabelSmallEmphasized:
		return "LabelSmallEmphasized"
	case TypographyTokenKeyTitleLargeEmphasized:
		return "TitleLargeEmphasized"
	case TypographyTokenKeyTitleMediumEmphasized:
		return "TitleMediumEmphasized"
	case TypographyTokenKeyTitleSmallEmphasized:
		return "TitleSmallEmphasized"
	default:
		return "Unspecified"
	}
}

var TypographyKeyTokens = &TypographyKeyTokensData{
	BodyLarge:                TypographyTokenKeyBodyLarge,
	BodyMedium:               TypographyTokenKeyBodyMedium,
	BodySmall:                TypographyTokenKeyBodySmall,
	DisplayLarge:             TypographyTokenKeyDisplayLarge,
	DisplayMedium:            TypographyTokenKeyDisplayMedium,
	DisplaySmall:             TypographyTokenKeyDisplaySmall,
	HeadlineLarge:            TypographyTokenKeyHeadlineLarge,
	HeadlineMedium:           TypographyTokenKeyHeadlineMedium,
	HeadlineSmall:            TypographyTokenKeyHeadlineSmall,
	LabelLarge:               TypographyTokenKeyLabelLarge,
	LabelMedium:              TypographyTokenKeyLabelMedium,
	LabelSmall:               TypographyTokenKeyLabelSmall,
	TitleLarge:               TypographyTokenKeyTitleLarge,
	TitleMedium:              TypographyTokenKeyTitleMedium,
	TitleSmall:               TypographyTokenKeyTitleSmall,
	BodyLargeEmphasized:      TypographyTokenKeyBodyLargeEmphasized,
	BodyMediumEmphasized:     TypographyTokenKeyBodyMediumEmphasized,
	BodySmallEmphasized:      TypographyTokenKeyBodySmallEmphasized,
	DisplayLargeEmphasized:   TypographyTokenKeyDisplayLargeEmphasized,
	DisplayMediumEmphasized:  TypographyTokenKeyDisplayMediumEmphasized,
	DisplaySmallEmphasized:   TypographyTokenKeyDisplaySmallEmphasized,
	HeadlineLargeEmphasized:  TypographyTokenKeyHeadlineLargeEmphasized,
	HeadlineMediumEmphasized: TypographyTokenKeyHeadlineMediumEmphasized,
	HeadlineSmallEmphasized:  TypographyTokenKeyHeadlineSmallEmphasized,
	LabelLargeEmphasized:     TypographyTokenKeyLabelLargeEmphasized,
	LabelMediumEmphasized:    TypographyTokenKeyLabelMediumEmphasized,
	LabelSmallEmphasized:     TypographyTokenKeyLabelSmallEmphasized,
	TitleLargeEmphasized:     TypographyTokenKeyTitleLargeEmphasized,
	TitleMediumEmphasized:    TypographyTokenKeyTitleMediumEmphasized,
	TitleSmallEmphasized:     TypographyTokenKeyTitleSmallEmphasized,
}

type TypographyKeyTokensData struct {
	BodyLarge                TypographyTokenKey
	BodyMedium               TypographyTokenKey
	BodySmall                TypographyTokenKey
	DisplayLarge             TypographyTokenKey
	DisplayMedium            TypographyTokenKey
	DisplaySmall             TypographyTokenKey
	HeadlineLarge            TypographyTokenKey
	HeadlineMedium           TypographyTokenKey
	HeadlineSmall            TypographyTokenKey
	LabelLarge               TypographyTokenKey
	LabelMedium              TypographyTokenKey
	LabelSmall               TypographyTokenKey
	TitleLarge               TypographyTokenKey
	TitleMedium              TypographyTokenKey
	TitleSmall               TypographyTokenKey
	BodyLargeEmphasized      TypographyTokenKey
	BodyMediumEmphasized     TypographyTokenKey
	BodySmallEmphasized      TypographyTokenKey
	DisplayLargeEmphasized   TypographyTokenKey
	DisplayMediumEmphasized  TypographyTokenKey
	DisplaySmallEmphasized   TypographyTokenKey
	HeadlineLargeEmphasized  TypographyTokenKey
	HeadlineMediumEmphasized TypographyTokenKey
	HeadlineSmallEmphasized  TypographyTokenKey
	LabelLargeEmphasized     TypographyTokenKey
	LabelMediumEmphasized    TypographyTokenKey
	LabelSmallEmphasized     TypographyTokenKey
	TitleLargeEmphasized     TypographyTokenKey
	TitleMediumEmphasized    TypographyTokenKey
	TitleSmallEmphasized     TypographyTokenKey
}
