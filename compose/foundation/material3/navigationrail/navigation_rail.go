package navigationrail

import (
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/theme"

	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/foundation/material3/surface"
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
		tm := theme.GetThemeManager()
		m3 := tm.GetMaterial3Theme()

		containerColor := m3.Scheme.Surface.Color.AsNRGBA()
		contentColor := m3.Scheme.Surface.OnColor.AsNRGBA()

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
			surface.WithColor(theme.ColorHelper.SpecificColor(containerColor)),
			surface.WithContentColor(theme.ColorHelper.SpecificColor(contentColor)),
			surface.WithModifier(
				modifier.
					Then(size.Width(80)).
					Then(size.FillMaxHeight()),
			),
		)(c)
	}
}
