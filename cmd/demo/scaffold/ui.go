package main

import (
	"fmt"
	"image/color"

	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/material3/button"
	"github.com/zodimo/go-compose/compose/foundation/material3/scaffold"
	"github.com/zodimo/go-compose/compose/foundation/material3/snackbar"
	"github.com/zodimo/go-compose/compose/foundation/material3/surface"
	m3text "github.com/zodimo/go-compose/compose/foundation/material3/text"
	ftext "github.com/zodimo/go-compose/compose/foundation/text"
	padding_modifier "github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"
)

func UI(c api.Composer) api.LayoutNode {

	// Create snackbar host state outside the loop to persist state
	snackbarHostState := c.State("snackbarHostState", func() any { return snackbar.NewSnackbarHostState() }).Get().(*snackbar.SnackbarHostState)

	return scaffold.Scaffold(
		box.Box(
			func(c api.Composer) api.Composer {
				return column.Column(
					func(c api.Composer) api.Composer {
						m3text.Text("Scaffold Content Area", m3text.TypestyleBodyLarge)(c)
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
		scaffold.WithSnackbarHost(snackbar.SnackbarHost(snackbarHostState)),
	)(c).Build()
}
