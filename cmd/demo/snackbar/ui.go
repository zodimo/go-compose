package main

import (
	"log"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/material3/button"
	"github.com/zodimo/go-compose/compose/material3/snackbar"
	m3Text "github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"
)

func UI() api.Composable {
	return func(c api.Composer) api.Composer {

		// Create snackbar host state outside the loop to persist state
		snackbarHostState := snackbar.RemeberSnackbarHostState(c)

		// Define UI inline
		return box.Box(
			func(c compose.Composer) compose.Composer {
				// Content
				column.Column(
					func(c compose.Composer) compose.Composer {
						// Headline
						m3Text.TextWithStyle("Snackbars!", m3Text.TypestyleHeadlineLarge)(c)

						// Body
						m3Text.TextWithStyle("This is a simple body message, click for debug info.", m3Text.TypestyleBodyLarge)(c)

						// Button
						button.Elevated(func() {
							log.Println("Showing Snackbar")
							snackbarHostState.ShowSnackbarAsync("Hi!")
						}, "Say Hi!",
							button.WithModifier(padding.Vertical(20, 0)),
						)(c)

						return c
					},
					column.WithModifier(padding.All(16)),
				)(c)

				// SnackbarHost overlay
				// Since Box stacks children, this will be drawn on top (last child)
				snackbar.SnackbarHost(snackbarHostState)(c)

				return c
			},
			box.WithModifier(
				size.FillMax(),
			),
		)(c)
	}
}
