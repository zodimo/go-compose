package main

import (
	"context"
	"log"
	"time"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
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
		snackbarHostState := snackbar.RememberSnackbarHostState(c)

		// Define UI inline
		return box.Box(
			func(c compose.Composer) compose.Composer {
				// Content
				column.Column(
					c.Sequence(
						// Headline
						m3Text.HeadlineLarge("Snackbars!"),

						// Body
						m3Text.BodyLarge("This is a simple body message, click for debug info."),

						// Button
						button.Elevated(func() {
							log.Println("Showing Snackbar")
							snackbarHostState.ShowSnackbar("Hi!")
						}, "Say Hi!",
							button.WithModifier(padding.Vertical(20, 0)),
						),
						spacer.Height(16),
						button.Elevated(func() {
							log.Println("Showing Snackbar with long text")
							snackbarHostState.ShowSnackbar("Hi! This is a long message that will wrap around the snackbar.")
						}, "Say Long Hi!",
							button.WithModifier(padding.Vertical(20, 0)),
						),
						spacer.Height(16),
						button.Elevated(func() {
							log.Println("Showing Snackbar with long text")
							ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
							go func() {
								// Wait for the context to time out
								<-ctx.Done()
								cancel()
							}()

							snackbarHostState.ShowSnackbar("Hi! 1 second", snackbar.WithContext(ctx))
						}, "Say Hi! (custom context 1s)",
							button.WithModifier(padding.Vertical(20, 0)),
						),
						spacer.Height(16),
						button.Elevated(func() {
							log.Println("Showing Snackbar with long text")
							ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
							go func() {
								// Wait for the context to time out
								<-ctx.Done()
								cancel()
							}()

							snackbarHostState.ShowSnackbar("Hi! 200ms", snackbar.WithContext(ctx))
						}, "Say Hi! (custom context 200ms)",
							button.WithModifier(padding.Vertical(20, 0)),
						),
					),

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
