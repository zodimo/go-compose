package appbar

import (
	"image/color"

	"go-compose-dev/compose/foundation/layout/box"
	"go-compose-dev/compose/foundation/layout/row"
	"go-compose-dev/compose/foundation/material3/surface"
	padding_modifier "go-compose-dev/internal/modifiers/padding"
	"go-compose-dev/internal/modifiers/size"
	"go-compose-dev/internal/modifiers/weight"

	"gioui.org/layout"
)

// SingleRowTopAppBar is an internal component to layout the TopAppBar content in a single row.
// It is used by SmallTopAppBar and CenterAlignedTopAppBar.
func SingleRowTopAppBar(
	modifier Modifier,
	title Composable,
	navigationIcon Composable,
	actions Composable,
	colors TopAppBarColors,
) Composable {
	return surface.Surface(
		func(c Composer) Composer {
			return row.Row(
				func(c Composer) Composer {
					// Navigation Icon
					if navigationIcon != nil {
						box.Box(
							surface.Surface(
								navigationIcon,
								surface.WithContentColor(colors.NavigationIconContentColor),
								surface.WithColor(color.NRGBA{}), // Transparent background
							),
							box.WithAlignment(layout.W),
							box.WithModifier(padding_modifier.Padding(4, 0, 0, 0)), // Start(4)
						)(c)
					}

					// Title
					box.Box(
						func(c Composer) Composer {
							return surface.Surface(
								title,
								surface.WithContentColor(colors.TitleContentColor),
								surface.WithColor(color.NRGBA{}), // Transparent
							)(c)
						},
						box.WithModifier(weight.Weight(1)),                    // Occupy remaining space
						box.WithAlignment(layout.W),                           // Align text to start
						box.WithModifier(padding_modifier.Horizontal(16, 16)), // Horizontal(16, 16)
					)(c)

					// Actions
					if actions != nil {
						row.Row(
							surface.Surface(
								actions,
								surface.WithContentColor(colors.ActionIconContentColor),
								surface.WithColor(color.NRGBA{}), // Transparent
							),
							row.WithAlignment(row.Middle),                          // Vertical alignment
							row.WithModifier(padding_modifier.Padding(0, 0, 4, 0)), // End(4)
						)(c)
					}
					return c
				},
				row.WithModifier(size.FillMaxWidth()),
				row.WithModifier(size.Height(64)), // Standard Height
				row.WithAlignment(row.Middle),     // Vertical Alignment
			)(c)
		},
		surface.WithModifier(modifier),
		surface.WithColor(colors.ContainerColor),
	)
}

// TopAppBar displays information and actions at the top of a screen.
// This is equivalent to SmallTopAppBar in Material 3.
func TopAppBar(
	title Composable,
	options ...TopAppBarOption,
) Composable {
	return func(c Composer) Composer {
		opts := DefaultTopAppBarOptions()
		for _, option := range options {
			option(&opts)
		}

		return SingleRowTopAppBar(
			opts.Modifier,
			title,
			opts.NavigationIcon,
			opts.Actions,
			opts.Colors,
		)(c)
	}
}
