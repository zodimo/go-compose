package main

import (
	"fmt"
	"go-compose-dev/compose/foundation/layout/box"
	"go-compose-dev/compose/foundation/layout/column"
	"go-compose-dev/compose/foundation/material3/button"
	"go-compose-dev/compose/foundation/material3/scaffold"
	"go-compose-dev/compose/foundation/material3/surface"
	m3text "go-compose-dev/compose/foundation/material3/text"
	ftext "go-compose-dev/compose/foundation/text"
	padding_modifier "go-compose-dev/internal/modifiers/padding"
	"go-compose-dev/internal/modifiers/size"
	"go-compose-dev/pkg/api"
	"image/color"
)

func UI(c api.Composer) api.LayoutNode {
	return scaffold.Scaffold(
		box.Box(
			func(c api.Composer) api.Composer {
				return column.Column(
					func(c api.Composer) api.Composer {
						m3text.Text("Scaffold Content Area", m3text.TypestyleBodyLarge)(c)
						button.Filled(func() {
							fmt.Println("Fab clicked!")
						}, "Click Me")(c)

						return c
					},
					column.WithAlignment(column.Middle),
				)(c)
			},
			box.WithAlignment(box.Center),
			box.WithModifier(size.FillMax()),
		),
		scaffold.WithTopBar(func(c api.Composer) api.Composer {
			return surface.Surface(
				func(c api.Composer) api.Composer {
					return box.Box(
						func(c api.Composer) api.Composer {
							m3text.Text("Top Bar", m3text.TypestyleTitleLarge,
								ftext.WithTextStyleOptions(ftext.StyleWithColor(color.NRGBA{R: 255, G: 255, B: 255, A: 255})),
							)(c)
							return c
						},
						box.WithModifier(padding_modifier.All(16)),
					)(c)
				},
				surface.WithColor(color.NRGBA{R: 98, G: 0, B: 238, A: 255}), // Purple 500
				surface.WithModifier(size.FillMaxWidth()),
				surface.WithModifier(size.Height(56)),
			)(c)
		}),
		scaffold.WithBottomBar(func(c api.Composer) api.Composer {
			return surface.Surface(
				func(c api.Composer) api.Composer {
					return box.Box(
						func(c api.Composer) api.Composer {
							m3text.Text("Bottom Bar", m3text.TypestyleBodyMedium)(c)
							return c
						},
						box.WithModifier(padding_modifier.All(16)),
						box.WithAlignment(box.Center),
					)(c)
				},
				surface.WithColor(color.NRGBA{R: 240, G: 240, B: 240, A: 255}), // Light Gray
				surface.WithModifier(size.FillMaxWidth()),
				surface.WithModifier(size.Height(80)),
			)(c)
		}),
		scaffold.WithFloatingActionButton(func(c api.Composer) api.Composer {
			return button.Filled(func() {
				fmt.Println("FAB Clicked")
			}, "+",
				button.WithModifier(size.Size(56, 56)),
				button.WithModifier(padding_modifier.All(0)), // Reset padding?
			)(c)
		}),
	)(c).Build()
}
