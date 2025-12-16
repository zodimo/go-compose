package navigationbar

import (
	"github.com/zodimo/go-compose/theme"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/material3/surface"
	"github.com/zodimo/go-compose/modifiers/size"
)

// NavigationBar displays a navigation bar at the bottom of the screen.
//
// Material 3 Specs:
// - Height: 80dp
// - Layout: Items are equally distributed.
func NavigationBar(
	content Composable,
	options ...NavigationBarOption,
) Composable {
	return func(c Composer) Composer {
		opts := DefaultNavigationBarOptions()
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
						if content != nil {
							content(c)
						}
						return c
					},
					row.WithModifier(size.FillMax()),
					// Items utilize weight to distribute space evenly
					row.WithAlignment(row.Middle), // Vertically centered
				)(c)
			},
			surface.WithModifier(
				opts.Modifier.
					Then(size.FillMaxWidth()).
					Then(size.Height(int(opts.Height))),
			),
			surface.WithColor(theme.ColorHelper.SpecificColor(opts.ContainerColor)),
			surface.WithContentColor(theme.ColorHelper.SpecificColor(opts.ContentColor)),
			surface.WithTonalElevation(opts.TonalElevation),
		)(c)
	}
}
