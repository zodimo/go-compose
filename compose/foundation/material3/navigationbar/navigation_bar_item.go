package navigationbar

import (
	"fmt"
	"image/color"

	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/foundation/material3/surface"
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	"github.com/zodimo/go-compose/modifiers/clickable"
	"github.com/zodimo/go-compose/modifiers/clip"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/modifiers/weight"
	"github.com/zodimo/go-compose/theme"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
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
						surface.WithColor(func() theme.ColorDescriptor {
							if selected {
								return colors.IndicatorColor
							}
							return theme.ColorHelper.SpecificColor(color.NRGBA{A: 0}) // Transparent
						}()),
						surface.WithContentColor(func() theme.ColorDescriptor {
							if selected {
								return theme.ColorHelper.ColorSelector().SecondaryRoles.OnContainer
							}
							return colors.ContentColor
						}()),
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
