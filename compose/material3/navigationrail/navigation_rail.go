package navigationrail

import (
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/material3"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"

	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/material3/surface"
)

// NavigationRail represents a side navigation component.
//
// Material 3 Specs:
// - Width: 80dp
// - Container Layout: Centered horizontally, Top or Center vertically (usually).
func NavigationRail(
	modifier Modifier,
	header Composable,
	content Composable,
) Composable {
	return func(c Composer) Composer {
		theme := material3.Theme(c)
		// Use theme role selectors directly
		containerColor := theme.ColorScheme().Surface.Color //theme.ColorHelper.ColorSelector().SurfaceRoles.Surface
		contentColor := theme.ColorScheme().Surface.OnColor //theme.ColorHelper.ColorSelector().SurfaceRoles.OnSurface

		return surface.Surface(
			func(c Composer) Composer {
				return column.Column(
					func(c Composer) Composer {
						if header != nil {
							header(c)
							spacer.Uniform(8)(c) // Spacing after header
						}

						// Content (Items)
						// M3: "Destinations are top-aligned in the rail by default, but can be centered."
						// For now, let's keep them top-aligned as default.
						content(c)
						return c
					},
					column.WithAlignment(column.Middle),
					column.WithModifier(
						// Add padding at top and bottom if needed
						padding.Padding(0, 12, 0, 12),
					),
				)(c)
			},
			surface.WithColor(containerColor),
			surface.WithContentColor(contentColor),
			surface.WithModifier(
				modifier.
					Then(size.Width(80)).
					Then(size.FillMaxHeight()),
			),
		)(c)
	}
}
