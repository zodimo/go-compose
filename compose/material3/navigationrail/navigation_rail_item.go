package navigationrail

import (
	"fmt"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/material3"
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	"github.com/zodimo/go-compose/compose/ui/text"
	"github.com/zodimo/go-compose/modifiers/clickable"
	"github.com/zodimo/go-compose/modifiers/clip"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"

	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/material3/surface"
	"github.com/zodimo/go-ternary"

	"github.com/zodimo/go-compose/compose/ui/unit"

	"gioui.org/widget"
)

// NavigationRailItem represents an item within a NavigationRail.
//
// Material 3 Specs:
// - Height: Can vary, usually minimum touch target.
// - Indicator: Pill shape around icon when selected.
// - Label: Text below icon.
func NavigationRailItem(
	selected bool,
	onClick func(),
	icon Composable,
	label Composable,
	modifier ui.Modifier,
	// colors, enabled, interactionSource... omitted for brevity
) Composable {
	return func(c Composer) Composer {
		theme := material3.Theme(c)
		// State for click interaction
		key := c.GenerateID()
		path := c.GetPath()
		clickStatePath := fmt.Sprintf("%d/%s/railitem_click", key, path)
		clickValue := c.State(clickStatePath, func() any { return &widget.Clickable{} })
		clickWidget := clickValue.Get().(*widget.Clickable)

		// Define indicator styling (pill shape)
		// Usually 56x32dp or similar for indicator.
		// We'll wrap the icon in a Surface that acts as the indicator.
		defaultIconsSizeProvider := compose.LocalTextStyle.Provides(
			text.CopyTextStyle(text.TextStyleUnspecified, text.WithFontSize(unit.NewTextUnit(24, unit.TextUnitTypeSp))),
		)

		return column.Column(
			func(c Composer) Composer {
				// Icon Container (Indicator)
				surface.Surface(
					compose.CompositionLocalProvider(
						[]api.ProvidedValue{
							defaultIconsSizeProvider,
						},
						icon,
					),
					surface.WithColor(ternary.Ternary(
						selected,
						theme.ColorScheme().SecondaryContainer,
						graphics.ColorTransparent,
					)),
					surface.WithShape(&shape.RoundedCornerShape{Radius: unit.Dp(16)}), // Pill shape (approx)
					surface.WithModifier(
						size.Height(32).
							Then(clip.Clip(&shape.RoundedCornerShape{Radius: unit.Dp(16)})). // Clip to pill shape
							Then(padding.Padding(16, 4, 16, 4)),                             // Padding inside indicator? Or just center icon.
					),
					surface.WithAlignment(box.Center),
				)(c)

				// Label
				if label != nil {
					spacer.Uniform(4)(c)
					label(c) // Should be styled with smaller font-size usually
				}
				return c
			},
			column.WithAlignment(column.Middle),
			column.WithModifier(
				modifier.
					Then(size.FillMaxWidth()).
					Then(clickable.OnClick(func() {
						if onClick != nil {
							onClick()
						}
					}, clickable.WithClickable(clickWidget))).
					Then(padding.Padding(0, 4, 0, 4)),
			),
		)(c)
	}
}
