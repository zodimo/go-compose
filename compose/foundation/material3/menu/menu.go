package menu

import (
	"go-compose-dev/compose/foundation/layout/column"
	m3Card "go-compose-dev/compose/foundation/material3/card"
	"go-compose-dev/internal/modifiers/padding"
	"go-compose-dev/internal/modifiers/size"
	"go-compose-dev/pkg/api"
)

// DropdownMenu Composable
// Currently renders inline. User should place it in a Box or Overlay.
func DropdownMenu(
	expanded bool,
	onDismissRequest func(),
	content api.Composable,
) api.Composable {
	return func(c api.Composer) api.Composer {
		if !expanded {
			return c
		}

		return m3Card.Elevated(
			m3Card.CardContents(
				m3Card.Content(
					column.Column(
						content,
						column.WithModifier(padding.Vertical(8, 8)),
						column.WithModifier(size.WrapContentWidth()),
					),
				),
			),
			// Use modifiers on the card if possible, but m3Card.Elevated signature might not support varargs modifiers in the first arg
			// Check m3Card signature.
		)(c)
	}
}
