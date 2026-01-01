package material3

import "github.com/zodimo/go-compose/compose/ui/text"

type TypographyOption func(*Typography)

func WithDisplaySmall(style *text.TextStyle) TypographyOption {
	if !text.IsSpecifiedTextStyle(style) {
		panic("DisplaySmall style must be specified")
	}
	return func(t *Typography) {
		t.DisplaySmall = style
	}
}

func WithDisplayMedium(style *text.TextStyle) TypographyOption {
	if !text.IsSpecifiedTextStyle(style) {
		panic("DisplayMedium style must be specified")
	}
	return func(t *Typography) {
		t.DisplayMedium = style
	}
}

func WithDisplayLarge(style *text.TextStyle) TypographyOption {
	if !text.IsSpecifiedTextStyle(style) {
		panic("DisplayLarge style must be specified")
	}
	return func(t *Typography) {
		t.DisplayLarge = style
	}
}

func WithHeadlineSmall(style *text.TextStyle) TypographyOption {
	if !text.IsSpecifiedTextStyle(style) {
		panic("HeadlineSmall style must be specified")
	}
	return func(t *Typography) {
		t.HeadlineSmall = style
	}
}

func WithHeadlineMedium(style *text.TextStyle) TypographyOption {
	if !text.IsSpecifiedTextStyle(style) {
		panic("HeadlineMedium style must be specified")
	}
	return func(t *Typography) {
		t.HeadlineMedium = style
	}
}

func WithHeadlineLarge(style *text.TextStyle) TypographyOption {
	if !text.IsSpecifiedTextStyle(style) {
		panic("HeadlineLarge style must be specified")
	}
	return func(t *Typography) {
		t.HeadlineLarge = style
	}
}

func WithTitleSmall(style *text.TextStyle) TypographyOption {
	if !text.IsSpecifiedTextStyle(style) {
		panic("TitleSmall style must be specified")
	}
	return func(t *Typography) {
		t.TitleSmall = style
	}
}

func WithTitleMedium(style *text.TextStyle) TypographyOption {
	if !text.IsSpecifiedTextStyle(style) {
		panic("TitleMedium style must be specified")
	}
	return func(t *Typography) {
		t.TitleMedium = style
	}
}

func WithTitleLarge(style *text.TextStyle) TypographyOption {
	if !text.IsSpecifiedTextStyle(style) {
		panic("TitleLarge style must be specified")
	}
	return func(t *Typography) {
		t.TitleLarge = style
	}
}

func WithBodySmall(style *text.TextStyle) TypographyOption {
	if !text.IsSpecifiedTextStyle(style) {
		panic("BodySmall style must be specified")
	}
	return func(t *Typography) {
		t.BodySmall = style
	}
}

func WithBodyMedium(style *text.TextStyle) TypographyOption {
	if !text.IsSpecifiedTextStyle(style) {
		panic("BodyMedium style must be specified")
	}
	return func(t *Typography) {
		t.BodyMedium = style
	}
}

func WithBodyLarge(style *text.TextStyle) TypographyOption {
	if !text.IsSpecifiedTextStyle(style) {
		panic("BodyLarge style must be specified")
	}
	return func(t *Typography) {
		t.BodyLarge = style
	}
}

func WithLabelSmall(style *text.TextStyle) TypographyOption {
	if !text.IsSpecifiedTextStyle(style) {
		panic("LabelSmall style must be specified")
	}
	return func(t *Typography) {
		t.LabelSmall = style
	}
}

func WithLabelMedium(style *text.TextStyle) TypographyOption {
	if !text.IsSpecifiedTextStyle(style) {
		panic("LabelMedium style must be specified")
	}
	return func(t *Typography) {
		t.LabelMedium = style
	}
}

func WithLabelLarge(style *text.TextStyle) TypographyOption {
	if !text.IsSpecifiedTextStyle(style) {
		panic("LabelLarge style must be specified")
	}
	return func(t *Typography) {
		t.LabelLarge = style
	}
}

func WithDisplaySmallEmphasized(style *text.TextStyle) TypographyOption {
	if !text.IsSpecifiedTextStyle(style) {
		panic("DisplaySmallEmphasized style must be specified")
	}
	return func(t *Typography) {
		t.DisplaySmallEmphasized = style
	}
}

func WithDisplayMediumEmphasized(style *text.TextStyle) TypographyOption {
	if !text.IsSpecifiedTextStyle(style) {
		panic("DisplayMediumEmphasized style must be specified")
	}
	return func(t *Typography) {
		t.DisplayMediumEmphasized = style
	}
}

func WithDisplayLargeEmphasized(style *text.TextStyle) TypographyOption {
	if !text.IsSpecifiedTextStyle(style) {
		panic("DisplayLargeEmphasized style must be specified")
	}
	return func(t *Typography) {
		t.DisplayLargeEmphasized = style
	}
}

func WithHeadlineSmallEmphasized(style *text.TextStyle) TypographyOption {
	if !text.IsSpecifiedTextStyle(style) {
		panic("HeadlineSmallEmphasized style must be specified")
	}
	return func(t *Typography) {
		t.HeadlineSmallEmphasized = style
	}
}

func WithHeadlineMediumEmphasized(style *text.TextStyle) TypographyOption {
	if !text.IsSpecifiedTextStyle(style) {
		panic("HeadlineMediumEmphasized style must be specified")
	}
	return func(t *Typography) {
		t.HeadlineMediumEmphasized = style
	}
}

func WithHeadlineLargeEmphasized(style *text.TextStyle) TypographyOption {
	if !text.IsSpecifiedTextStyle(style) {
		panic("HeadlineLargeEmphasized style must be specified")
	}
	return func(t *Typography) {
		t.HeadlineLargeEmphasized = style
	}
}

func WithTitleSmallEmphasized(style *text.TextStyle) TypographyOption {
	if !text.IsSpecifiedTextStyle(style) {
		panic("TitleSmallEmphasized style must be specified")
	}
	return func(t *Typography) {
		t.TitleSmallEmphasized = style
	}
}

func WithTitleMediumEmphasized(style *text.TextStyle) TypographyOption {
	if !text.IsSpecifiedTextStyle(style) {
		panic("TitleMediumEmphasized style must be specified")
	}
	return func(t *Typography) {
		t.TitleMediumEmphasized = style
	}
}

func WithTitleLargeEmphasized(style *text.TextStyle) TypographyOption {
	if !text.IsSpecifiedTextStyle(style) {
		panic("TitleLargeEmphasized style must be specified")
	}
	return func(t *Typography) {
		t.TitleLargeEmphasized = style
	}
}

func WithBodySmallEmphasized(style *text.TextStyle) TypographyOption {
	if !text.IsSpecifiedTextStyle(style) {
		panic("BodySmallEmphasized style must be specified")
	}
	return func(t *Typography) {
		t.BodySmallEmphasized = style
	}
}

func WithBodyMediumEmphasized(style *text.TextStyle) TypographyOption {
	if !text.IsSpecifiedTextStyle(style) {
		panic("BodyMediumEmphasized style must be specified")
	}
	return func(t *Typography) {
		t.BodyMediumEmphasized = style
	}
}

func WithBodyLargeEmphasized(style *text.TextStyle) TypographyOption {
	if !text.IsSpecifiedTextStyle(style) {
		panic("BodyLargeEmphasized style must be specified")
	}
	return func(t *Typography) {
		t.BodyLargeEmphasized = style
	}
}

func WithLabelSmallEmphasized(style *text.TextStyle) TypographyOption {
	if !text.IsSpecifiedTextStyle(style) {
		panic("LabelSmallEmphasized style must be specified")
	}
	return func(t *Typography) {
		t.LabelSmallEmphasized = style
	}
}

func WithLabelMediumEmphasized(style *text.TextStyle) TypographyOption {
	if !text.IsSpecifiedTextStyle(style) {
		panic("LabelMediumEmphasized style must be specified")
	}
	return func(t *Typography) {
		t.LabelMediumEmphasized = style
	}
}

func WithLabelLargeEmphasized(style *text.TextStyle) TypographyOption {
	if !text.IsSpecifiedTextStyle(style) {
		panic("LabelLargeEmphasized style must be specified")
	}
	return func(t *Typography) {
		t.LabelLargeEmphasized = style
	}
}
