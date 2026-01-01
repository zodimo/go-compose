package navigationbar

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/material3/surface"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	"github.com/zodimo/go-compose/modifiers/clickable"
	"github.com/zodimo/go-compose/modifiers/clip"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/modifiers/weight"
	"github.com/zodimo/go-compose/theme"
	"github.com/zodimo/go-ternary"

	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

// NavigationBarItem represents an item within a NavigationBar.
//
// Material 3 Specs:
// - Height: Fill container (80dp).
// - Indicator: Pill shape (64x32dp) around icon when selected.
// - Label: Text below icon.
func NavigationBarItem(
	selected bool,
	onClick func(),
	icon Composable,
	label Composable,
	options ...NavigationBarItemOption,
) Composable {
	return func(c Composer) Composer {

		opts := DefaultNavigationBarItemOptions()
		for _, option := range options {
			if option == nil {
				continue
			}
			option(&opts)
		}

		// State for click interaction
		key := c.GenerateID()
		path := c.GetPath()
		clickStatePath := fmt.Sprintf("%d/%s/navitem_click", key, path)
		clickValue := c.State(clickStatePath, func() any { return &widget.Clickable{} })
		clickWidget := clickValue.Get().(*widget.Clickable)

		// Defaults
		colors := NavigationBarDefaults.Colors()

		return box.Box(
			column.Column(
				c.Sequence(
					// Top padding
					spacer.Uniform(12),
					// Icon Container (Indicator)
					surface.Surface(
						func(c Composer) Composer {
							return box.Box(
								icon,
								box.WithAlignment(layout.Center),
							)(c)
						},
						surface.WithColor(ternary.Ternary(
							selected,
							colors.IndicatorColor,
							theme.ColorHelper.SpecificColor(graphics.ColorTransparent), // Transparent
						)),
						surface.WithContentColor(ternary.Ternary(
							selected,
							theme.ColorHelper.ColorSelector().SecondaryRoles.OnContainer,
							colors.ContentColor,
						)),
						surface.WithShape(shape.RoundedCornerShape{Radius: unit.Dp(16)}),
						surface.WithModifier(
							EmptyModifier.
								Then(size.Width(64)).
								Then(size.Height(32)).
								Then(clip.Clip(shape.RoundedCornerShape{Radius: unit.Dp(16)})),
						),
						surface.WithAlignment(layout.Center),
					),
					// Label
					c.When(
						label != nil,
						c.Sequence(
							spacer.Uniform(4),
							label,
						),
					),
					// Bottom padding
					spacer.Uniform(12),
				),
				column.WithAlignment(column.Middle), // Center children horizontally
			),
			box.WithModifier(
				opts.Modifier.
					Then(weight.Weight(1)).     // Fill share of the Row
					Then(size.FillMaxHeight()). // Fill height of the Row
					Then(clickable.OnClick(func() {
						if onClick != nil {
							onClick()
						}
					}, clickable.WithClickable(clickWidget))),
			),
			box.WithAlignment(layout.Center), // Center the Column within the allocated slot
		)(c)
	}
}
