package navigationdrawer

import (
	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/material3/surface"
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/theme"

	"gioui.org/unit"
)

// PermanentNavigationDrawer always shows the drawer.
// It is used for large screens (Expanded/Extra-large).
// The drawer is placed at the start of the content.
func PermanentNavigationDrawer(
	drawerContent Composable,
	content Composable,
	modifier Modifier,
) Composable {
	return func(c Composer) Composer {
		tm := theme.GetThemeManager()
		m3 := tm.GetMaterial3Theme()

		drawerContainerColor := m3.Scheme.SurfaceContainerLow.AsNRGBA()

		return row.Row(
			func(c Composer) Composer {
				// 1. Drawer Sheet
				surface.Surface(
					drawerContent,
					surface.WithColor(theme.ColorHelper.SpecificColor(drawerContainerColor)),
					// Standard drawer doesn't usually have rounded corners on the edge touching the content
					// unless it's a specific variant, but M3 defaults often show 0 radius or small radius.
					// We'll stick to a standard square edge or small radius if needed.
					// M3: "Permanent navigation drawers are co-planar with the content."
					surface.WithShape(shape.RoundedCornerShape{Radius: unit.Dp(0)}),
					surface.WithModifier(
						modifier.
							Then(size.Width(360)).
							Then(size.FillMaxHeight()),
					),
				)(c)

				// 2. Main Content
				box.Box(
					content,
					box.WithModifier(size.FillMax()),
				)(c)

				return c
			},
			row.WithModifier(size.FillMax()),
		)(c)
	}
}
