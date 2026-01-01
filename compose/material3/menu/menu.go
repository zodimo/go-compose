package menu

import (
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/material3/surface"
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	"github.com/zodimo/go-compose/compose/ui/window" // Import internal modifier
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"

	"github.com/zodimo/go-compose/compose/ui/unit"
)

// DropdownMenu Composable
// Uses a Popup to display content on top of other content.
func DropdownMenu(
	expanded bool,
	onDismissRequest func(),
	content api.Composable,
	options ...DropdownMenuOption,
) api.Composable {
	return func(c api.Composer) api.Composer {
		if !expanded {
			return c
		}

		opts := DefaultDropdownMenuOptions()
		for _, option := range options {
			if option == nil {
				continue
			}
			option(&opts)
		}

		return window.Popup(
			surface.Surface(
				func(c api.Composer) api.Composer {
					return column.Column(
						content,
						column.WithModifier(
							padding.Vertical(8, 8),
						),
					)(c)
				},
				// M3 Menu Spec
				// Shape: Extra Small (4dp)
				surface.WithShape(shape.RoundedCornerShape{Radius: unit.Dp(4)}),
				// Container Color: Surface (default)
				// Elevation: Level 2 (3.dp)
				surface.WithShadowElevation(3),
				surface.WithModifier(
					// M3 Specs:
					// Min width: 112.dp
					// Max width: 280.dp
					// TODO: Implement MinWidth/MaxWidth modifiers in size package.
					// For now, we depend on content width.
					size.
						WrapContentWidth().
						Then(opts.Modifier),
				),
			),
			window.WithOffset(opts.OffsetX, opts.OffsetY),
		)(c)
	}
}
