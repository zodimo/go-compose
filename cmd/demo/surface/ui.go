package main

import (
	"fmt"
	"image/color"

	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/text"
	"github.com/zodimo/go-compose/compose/material3"
	"github.com/zodimo/go-compose/compose/material3/button"
	"github.com/zodimo/go-compose/compose/material3/surface"
	m3text "github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	iuText "github.com/zodimo/go-compose/compose/ui/text"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"

	"github.com/zodimo/go-compose/compose/ui/unit"
)

func UI(c api.Composer) api.LayoutNode {
	// State
	activeState := c.State("elevation_active", func() any {
		return false
	})
	active := activeState.Get().(bool)

	c = column.Column(
		c.Sequence(
			// Controls
			row.Row(
				c.Sequence(
					button.Filled(
						func() {
							activeState.Set(!active)
							fmt.Printf("Toggle Elevation: %v\n", !active)
						},
						"Toggle Elevation",
						button.WithModifier(padding.All(8)),
					),
				),
				row.WithModifier(size.FillMaxWidth().Then(padding.All(16))),
			),
			// Content Area
			// Level 0: Surface (White/Background)
			surface.Surface(
				c.Sequence(
					column.Column(
						c.Sequence(
							m3text.HeadlineLarge("Aloha!", text.WithTextStyleOptions(iuText.WithColor(
								material3.Theme(c).ColorScheme().Surface.OnColor,
							))),
							m3text.BodySmall("Level 1: 1", text.WithTextStyleOptions(iuText.WithColor(
								material3.Theme(c).ColorScheme().Surface.OnColor,
							))),
							// Start Recursion from Level 2
							RecursiveSurface(2, active),
						),
						column.WithModifier(padding.All(16)),
					),
				),
				surface.WithColor(graphics.FromNRGBA(color.NRGBA{R: 255, G: 251, B: 254, A: 255})), // Surface
				surface.WithShape(&shape.CutCornerShape{Radius: unit.Dp(28)}),                      // TopEnd chamfer in demo usually, here simplified to all cut
				surface.WithShadowElevation(0),                                                     // Root usually flat?
				surface.WithModifier(size.FillMax()),
			),
		),
		column.WithModifier(size.FillMax().Then(padding.All(8))),
	)(c)

	return c.Build()
}

// RecursiveSurface generates nested surfaces with alternating colors and shapes
func RecursiveSurface(level int, active bool) api.Composable {
	return func(c api.Composer) api.Composer {
		if level > 15 {
			return c
		}

		var bgColor graphics.Color
		var textColor graphics.Color
		var currentShape surface.Shape
		var elevation unit.Dp

		// Cycle: Primary -> Secondary -> Tertiary -> Inverse
		// Level 2 is index 0 in our logic roughly?
		// Demo: Level 1 (Surface) -> Level 2 (Primary? or Secondary?)
		// Let's just cycle.
		cycle := (level) % 4
		switch cycle {
		case 0: // Primary Container
			bgColor = material3.Theme(c).ColorScheme().PrimaryContainer.Color     //theme.ColorHelper.ColorSelector().PrimaryRoles.Container
			textColor = material3.Theme(c).ColorScheme().PrimaryContainer.OnColor // theme.ColorHelper.ColorSelector().PrimaryRoles.OnContainer
			currentShape = &shape.CutCornerShape{Radius: unit.Dp(20)}
			if active {
				elevation = unit.Dp(1)
			}
		case 1: // Secondary Container
			bgColor = material3.Theme(c).ColorScheme().SecondaryContainer.Color     //theme.ColorHelper.ColorSelector().SecondaryRoles.Container
			textColor = material3.Theme(c).ColorScheme().SecondaryContainer.OnColor // theme.ColorHelper.ColorSelector().SecondaryRoles.OnContainer
			currentShape = &shape.RoundedCornerShape{Radius: unit.Dp(16)}
			if active {
				elevation = unit.Dp(3)
			}
		case 2: // Tertiary Container
			bgColor = material3.Theme(c).ColorScheme().TertiaryContainer.Color     //theme.ColorHelper.ColorSelector().TertiaryRoles.Container
			textColor = material3.Theme(c).ColorScheme().TertiaryContainer.OnColor // theme.ColorHelper.ColorSelector().TertiaryRoles.OnContainer
			currentShape = &shape.CutCornerShape{Radius: unit.Dp(12)}
			if active {
				elevation = unit.Dp(6)
			}
		case 3: // Inverse Surface
			bgColor = material3.Theme(c).ColorScheme().InverseSurface.Color     //theme.ColorHelper.ColorSelector().InverseRoles.Surface
			textColor = material3.Theme(c).ColorScheme().InverseSurface.OnColor // theme.ColorHelper.ColorSelector().InverseRoles.OnSurface
			currentShape = &shape.RoundedCornerShape{Radius: unit.Dp(8)}
			if active {
				elevation = unit.Dp(8)
			}
		}

		// Adjust shape radius based on level logic if needed?
		// The demo shrinks sizes effectively by padding.

		return surface.Surface(
			c.Sequence(
				// Padding/Margin inside the surface
				column.Column(
					c.Sequence(
						m3text.BodyMedium(
							fmt.Sprintf("Level %d: %d", level, cycle),
							text.WithTextStyleOptions(
								iuText.WithColor(textColor),
							// text.StyleWithColor(textColor),
							)),
						RecursiveSurface(level+1, active),
					),
					// Tighter padding to create tunnel effect
					column.WithModifier(padding.All(4)), // changed from 20 to 4
				),
			),
			surface.WithColor(bgColor),
			surface.WithShape(currentShape),
			surface.WithShadowElevation(elevation),
			// surface.WithBorder(unit.Dp(1), toNRGBA(textColor)), // Optional border for contrast
			surface.WithModifier(padding.All(8)), // External margin? Or internal padding?
			// In Compose, Surface modifier is external. Content is inside.
			// To nest them: outer Surface -> padding(margin) -> inner Surface.
			// But here we apply padding to the Inner Surface via WithModifier.
		)(c)
	}
}
