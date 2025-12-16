package main

import (
	"fmt"
	"image/color"

	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/material3/button"
	"github.com/zodimo/go-compose/compose/foundation/material3/surface"
	m3text "github.com/zodimo/go-compose/compose/foundation/material3/text"
	"github.com/zodimo/go-compose/compose/foundation/text"
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"

	"gioui.org/unit"
	"git.sr.ht/~schnwalter/gio-mw/token"
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
				row.WithModifier(size.FillMaxWidth()),
				row.WithModifier(padding.All(16)),
			),
			// Content Area
			// Level 0: Surface (White/Background)
			surface.Surface(
				c.Sequence(
					column.Column(
						c.Sequence(
							m3text.Text("Aloha!", token.TypestyleHeadlineLarge, text.WithTextStyleOptions(text.StyleWithColor(color.NRGBA{R: 28, G: 27, B: 31, A: 255}))), // OnSurface
							m3text.Text("Level 1: 1", token.TypestyleBodySmall, text.WithTextStyleOptions(text.StyleWithColor(color.NRGBA{R: 28, G: 27, B: 31, A: 255}))),
							// Start Recursion from Level 2
							RecursiveSurface(2, active),
						),
						column.WithModifier(padding.All(16)),
					),
				),
				surface.WithColor(color.NRGBA{R: 255, G: 251, B: 254, A: 255}), // Surface
				surface.WithShape(shape.CutCornerShape{Radius: unit.Dp(28)}),   // TopEnd chamfer in demo usually, here simplified to all cut
				surface.WithShadowElevation(0),                                 // Root usually flat?
				surface.WithModifier(size.FillMax()),
			),
		),
		column.WithModifier(size.FillMax()),
		column.WithModifier(padding.All(8)),
	)(c)

	return c.Build()
}

// RecursiveSurface generates nested surfaces with alternating colors and shapes
func RecursiveSurface(level int, active bool) api.Composable {
	return func(c api.Composer) api.Composer {
		if level > 15 {
			return c
		}

		// Colors (Baseline M3 Light)
		primaryContainer := color.NRGBA{R: 234, G: 221, B: 255, A: 255}
		onPrimaryContainer := color.NRGBA{R: 33, G: 0, B: 93, A: 255}

		secondaryContainer := color.NRGBA{R: 232, G: 222, B: 248, A: 255}
		onSecondaryContainer := color.NRGBA{R: 29, G: 25, B: 43, A: 255}

		tertiaryContainer := color.NRGBA{R: 255, G: 216, B: 228, A: 255}
		onTertiaryContainer := color.NRGBA{R: 49, G: 17, B: 29, A: 255}

		inverseSurface := color.NRGBA{R: 49, G: 48, B: 51, A: 255}
		inverseOnSurface := color.NRGBA{R: 244, G: 239, B: 244, A: 255}

		var bgColor, textColor color.Color
		var currentShape surface.Shape
		var elevation unit.Dp

		// Cycle: Primary -> Secondary -> Tertiary -> Inverse
		// Level 2 is index 0 in our logic roughly?
		// Demo: Level 1 (Surface) -> Level 2 (Primary? or Secondary?)
		// Let's just cycle.
		cycle := (level) % 4
		switch cycle {
		case 0: // Primary Container
			bgColor = primaryContainer
			textColor = onPrimaryContainer
			currentShape = shape.CutCornerShape{Radius: unit.Dp(20)}
			if active {
				elevation = unit.Dp(1)
			}
		case 1: // Secondary Container
			bgColor = secondaryContainer
			textColor = onSecondaryContainer
			currentShape = shape.RoundedCornerShape{Radius: unit.Dp(16)}
			if active {
				elevation = unit.Dp(3)
			}
		case 2: // Tertiary Container
			bgColor = tertiaryContainer
			textColor = onTertiaryContainer
			currentShape = shape.CutCornerShape{Radius: unit.Dp(12)}
			if active {
				elevation = unit.Dp(6)
			}
		case 3: // Inverse Surface
			bgColor = inverseSurface
			textColor = inverseOnSurface
			currentShape = shape.RoundedCornerShape{Radius: unit.Dp(8)}
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
						m3text.Text(fmt.Sprintf("Level %d: %d", level, cycle), token.TypestyleBodyMedium, text.WithTextStyleOptions(text.StyleWithColor(toNRGBA(textColor)))),
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

func toNRGBA(c color.Color) color.NRGBA {
	return c.(color.NRGBA)
}
