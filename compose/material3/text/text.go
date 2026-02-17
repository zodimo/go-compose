package text

import (
	"github.com/zodimo/go-compose/compose/foundation/text"
	"github.com/zodimo/go-compose/compose/material3"
	"github.com/zodimo/go-compose/compose/material3/tokens"

	uiText "github.com/zodimo/go-compose/compose/ui/text"

	"github.com/zodimo/go-compose/pkg/api"
)

func Text(value string, options ...text.TextOption) api.Composable {
	return func(c api.Composer) api.Composer {
		contentColor := material3.LocalContentColor.Current(c)
		finalOptions := append(
			[]text.TextOption{
				text.WithTextStyleOptions(uiText.WithColor(contentColor)),
			},
			options...,
		)
		return text.Text(
			value,
			finalOptions...,
		)(c)
	}
}

// textWithStyle displays text with the given style from the Material Theme.
// It retrieves the theme from the layout context at runtime.
func textWithStyle(value string, tokenStyle tokens.TypographyTokenKey, options ...text.TextOption) api.Composable {
	return func(c api.Composer) api.Composer {
		theme := material3.Theme(c)
		typography := theme.Typography()
		contentColor := material3.LocalContentColor.Current(c)

		baseStyle := typography.FromToken(tokenStyle)

		if contentColor.IsSpecified() {
			baseStyle = uiText.CopyTextStyle(baseStyle, uiText.WithColor(contentColor))
		}

		finalOptions := append([]text.TextOption{text.WithTextStyle(baseStyle)}, options...)

		return text.Text(
			value,
			finalOptions...,
		)(c)
	}
}
