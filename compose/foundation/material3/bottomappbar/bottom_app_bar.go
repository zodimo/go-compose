package bottomappbar

import (
	"image/color"

	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/material3/surface"
	padding_modifier "github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/modifiers/weight"
	"github.com/zodimo/go-compose/theme"

	"gioui.org/layout"
)

// BottomAppBar displays navigation and key actions at the bottom of the screen.
func BottomAppBar(
	actions Composable,
	options ...BottomAppBarOption,
) Composable {
	return func(c Composer) Composer {
		opts := DefaultBottomAppBarOptions()
		for _, option := range options {
			if option == nil {
				continue
			}
			option(&opts)
		}

		return surface.Surface(
			func(c Composer) Composer {
				return row.Row(
					func(c Composer) Composer {
						// Actions Section (pushes FAB to the end using weight)
						box.Box(
							func(c Composer) Composer {
								if actions != nil {
									return surface.Surface(
										actions,
										surface.WithContentColor(opts.ContentColor),
										surface.WithColor(theme.ColorHelper.SpecificColor(color.NRGBA{})), // Transparent background
									)(c)
								}
								return c
							},
							box.WithModifier(weight.Weight(1)), // Flex to fill available space
							box.WithAlignment(layout.W),        // Align content to start
						)(c)

						// Floating Action Button Section
						if opts.FloatingActionButton != nil {
							box.Box(
								opts.FloatingActionButton,
								box.WithAlignment(layout.E),
								// Add padding around FAB if needed?
								// Usually FAB has its own padding or shadow, but BottomAppBar might add some.
								// M3 spec says FAB has specific padding from end.
								// Let's assume content padding covers it, or add specific margin.
								// ContentPadding is applied to the Row, so the FAB sits inside that padding.
							)(c)
						}
						return c
					},
					row.WithModifier(size.FillMaxWidth()),
					row.WithModifier(size.Height(80)), // Standard height for BottomAppBar
					row.WithAlignment(row.Middle),     // Vertically center content
					row.WithModifier(padding_modifier.Padding(
						int(opts.ContentPadding.Start),
						int(opts.ContentPadding.Top),
						int(opts.ContentPadding.End),
						int(opts.ContentPadding.Bottom),
					)),
				)(c)
			},
			surface.WithModifier(opts.Modifier),
			surface.WithColor(opts.ContainerColor),
			surface.WithTonalElevation(opts.TonalElevation),
			// Check Surface implementation. It usually supports WithShadow (elevation).
			// If TonalElevation is just color adjustment, we might need a separate mechanism or it's handled by Surface logic?
			// For now let's skip explicit TonalElevation call if Surface doesn't have it, or use WithShadow(opts.TonalElevation).
		)(c)
	}
}
