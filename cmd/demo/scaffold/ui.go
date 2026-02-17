package main

import (
	"fmt"
	"image/color"

	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/material3/button"
	"github.com/zodimo/go-compose/compose/material3/scaffold"
	"github.com/zodimo/go-compose/compose/material3/snackbar"
	"github.com/zodimo/go-compose/compose/material3/surface"
	m3text "github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	padding_modifier "github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"
)

func UI() api.Composable {
	return func(c api.Composer) api.Composer {

		// Create snackbar host state outside the loop to persist state
		snackbarHostState := snackbar.RememberSnackbarHostState(c)

		return scaffold.Scaffold(
			box.Box(
				func(c api.Composer) api.Composer {
					return column.Column(
						func(c api.Composer) api.Composer {
							m3text.BodyLarge("Scaffold Content Area")(c)
							button.Filled(func() {
								snackbarHostState.ShowSnackbar("Content area button clicked!")
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
							m3text.TitleLarge("Top Bar"),
							box.WithModifier(padding_modifier.All(16)),
						)(c)
					},
					surface.WithColor(graphics.FromNRGBA(color.NRGBA{R: 98, G: 0, B: 238, A: 255})), // Purple 500
					surface.WithModifier(size.FillMaxWidth().
						Then(size.Height(56)),
					),
				)(c)
			}),
			scaffold.WithBottomBar(func(c api.Composer) api.Composer {
				return surface.Surface(
					func(c api.Composer) api.Composer {
						return box.Box(
							func(c api.Composer) api.Composer {
								m3text.BodyMedium("Bottom Bar")(c)
								return c
							},
							box.WithModifier(padding_modifier.All(16)),
							box.WithAlignment(box.Center),
						)(c)
					},
					surface.WithColor(graphics.FromNRGBA(color.NRGBA{R: 240, G: 240, B: 240, A: 255})), // Light Gray
					surface.WithModifier(size.FillMaxWidth().
						Then(size.Height(80)),
					),
				)(c)
			}),
			scaffold.WithFloatingActionButton(func(c api.Composer) api.Composer {
				return button.Filled(func() {
					fmt.Println("FAB Clicked")
				}, "+",
					button.WithModifier(size.Size(56, 56).
						Then(padding_modifier.All(0)), // Reset padding?
					),
				)(c)
			}),
			scaffold.WithSnackbarHost(snackbar.SnackbarHost(snackbarHostState)),
		)(c)

	}
}
