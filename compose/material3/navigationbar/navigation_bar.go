package navigationbar

import (
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/material3/surface"
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
		opts := DefaultNavigationBarOptions(c)
		for _, option := range options {
			if option == nil {
				continue
			}
			option(&opts)
		}
		return surface.Surface(
			row.Row(
				func(c Composer) Composer {
					if content != nil {
						content(c)
					}
					return c
				},
				row.WithModifier(size.FillMax()),
				// Items utilize weight to distribute space evenly
				row.WithAlignment(row.Middle), // Vertically centered
			),
			surface.WithModifier(
				opts.Modifier.
					Then(size.FillMaxWidth()).
					Then(size.Height(int(opts.Height))),
			),
			surface.WithColor(opts.ContainerColor),
			surface.WithContentColor(opts.ContentColor),
			surface.WithTonalElevation(opts.TonalElevation),
		)(c)
	}
}
