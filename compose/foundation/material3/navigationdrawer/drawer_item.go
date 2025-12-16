package navigationdrawer

import (
	"image/color"

	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/foundation/material3/surface"
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	"github.com/zodimo/go-compose/modifiers/clickable"
	"github.com/zodimo/go-compose/modifiers/clip"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/theme"

	"gioui.org/unit"
	"gioui.org/widget"
)

// NavigationDrawerItem is a destination item within a Navigation Drawer.
// Material 3 Specs:
// - Height: 56dp
// - Shape: Full rounded (Stadium) or RoundedRect 100dp.
// - Layout: Icon (start) + Label.
func NavigationDrawerItem(
	selected bool,
	onClick func(),
	icon Composable,
	label Composable,
	m Modifier,
) Composable {
	return func(c Composer) Composer {
		tm := theme.GetThemeManager()
		m3 := tm.GetMaterial3Theme()

		// Colors
		var containerColor color.Color
		if selected {
			containerColor = m3.Scheme.SecondaryContainer.Color.AsNRGBA()
		} else {
			containerColor = color.NRGBA{A: 0} // Transparent
		}

		// Click interaction state
		key := c.GenerateID()
		path := c.GetPath()
		clickStatePath := key.String() + "/" + path.String() + "/draweritem_click"
		clickWidget := c.State(clickStatePath, func() any { return &widget.Clickable{} }).Get().(*widget.Clickable)

		return surface.Surface(
			func(c Composer) Composer {
				return row.Row(
					func(c Composer) Composer {
						// Icon
						// Icon color is handled by caller or implicit defaults if we had them.
						// Drawer items usually have 24dp icons, 12dp padding start.
						// We'll rely on the caller sizing the icon or simple wrapping.
						// We just layout children.
						if icon != nil {
							icon(c)
							spacer.Width(12)(c)
						}

						// Label
						if label != nil {
							box.Box(
								func(c Composer) Composer {
									label(c)
									return c
								},
								box.WithAlignment(box.Center),
								box.WithModifier(size.FillMaxHeight()),
							)(c)
						}
						return c
					},

					row.WithAlignment(row.Middle), // Vertically center content
					row.WithModifier(size.FillMax()),
					row.WithModifier(
						// Padding inside the item
						padding.Padding(16, 0, 24, 0), // 16dp start, 24dp end
					),
				)(c)
			},
			surface.WithColor(theme.ColorHelper.SpecificColor(containerColor)),
			surface.WithShape(shape.RoundedCornerShape{Radius: unit.Dp(28)}), // Stadium shape (height 56 / 2)
			surface.WithModifier(
				m.
					Then(size.FillMaxWidth()).
					Then(size.Height(56)).
					Then(clip.Clip(shape.RoundedCornerShape{Radius: unit.Dp(28)})).
					Then(clickable.OnClick(func() {
						if onClick != nil {
							onClick()
						}
					}, clickable.WithClickable(clickWidget))),
				// M3 says item connects to edge with 12dp horizontal padding usually.
				// Or the container has padding.
			),
		)(c)
	}
}
