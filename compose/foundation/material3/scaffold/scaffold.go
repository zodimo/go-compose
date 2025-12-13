package scaffold

import (
	"go-compose-dev/compose"
	"go-compose-dev/compose/foundation/layout/box"
	"go-compose-dev/compose/foundation/layout/column"
	"go-compose-dev/compose/foundation/material3/surface"
	padding_modifier "go-compose-dev/internal/modifiers/padding"
	"go-compose-dev/internal/modifiers/size"
	"go-compose-dev/internal/modifiers/weight"
	"image/color"

	"gioui.org/layout"
)

type FabPosition int

const (
	FabPositionCenter FabPosition = iota
	FabPositionEnd
)

// Scaffold implements the basic material design visual layout structure.
// This component provides API to put together several material components to construct your
// screen, by ensuring proper layout strategy for them and collecting necessary data so these
// components will work together correctly.
func Scaffold(content Composable, options ...ScaffoldOption) Composable {
	return func(c Composer) Composer {
		opts := DefaultScaffoldOptions()
		for _, option := range options {
			option(&opts)
		}

		// Prepare surface options
		var surfaceOpts []surface.SurfaceOption
		if opts.ContainerColor != (color.Color)(nil) && opts.ContainerColor != (color.NRGBA{}) {
			surfaceOpts = append(surfaceOpts, surface.WithColor(opts.ContainerColor))
		}
		if opts.ContentColor != (color.Color)(nil) && opts.ContentColor != (color.NRGBA{}) {
			surfaceOpts = append(surfaceOpts, surface.WithContentColor(opts.ContentColor))
		}
		// Apply the scaffold modifier to the root surface
		surfaceOpts = append(surfaceOpts, surface.WithModifier(opts.Modifier))

		return surface.Surface(
			func(c Composer) Composer {
				// We use a Box as the root to overlay FAB and Snackbars over the body
				return box.Box( // Root Stacking Context
					compose.Sequence(
						// Layer 1: The App Structure (TopBar -> Content -> BottomBar)
						func(c Composer) Composer {
							return column.Column(
								compose.Sequence(
									func(c Composer) Composer {
										// Top Bar
										if opts.TopBar != nil {
											opts.TopBar(c)
										}
										return c
									},
									func(c Composer) Composer {
										// Content Body
										// This needs to expand to fill available space.
										return box.Box(
											content,
											box.WithModifier(
												// Expand to fill remaining vertical space
												weight.Weight(1),
											),
											box.WithModifier(size.FillMaxWidth()),
										)(c)
									},
									func(c Composer) Composer {
										// Bottom Bar
										if opts.BottomBar != nil {
											opts.BottomBar(c)
										}
										return c
									},
								),
							)(c)
						},
						// Layer 2: Floating Action Button
						func(c Composer) Composer {
							if opts.FloatingActionButton != nil {
								alignment := layout.SE // Bottom End
								if opts.FloatingActionButtonPosition == FabPositionCenter {
									alignment = layout.S // Bottom Center
								}

								return box.Box(
									opts.FloatingActionButton,
									box.WithAlignment(alignment),
									// Add standard padding for FAB
									box.WithModifier(padding_modifier.All(16)),
									// Wrapper must fill max to align FAB relative to screen
									box.WithModifier(size.FillMax()),
								)(c)
							}
							return c
						},
						// Layer 3: Snackbar Host
						func(c Composer) Composer {
							if opts.SnackbarHost != nil {
								return box.Box(
									opts.SnackbarHost,
									box.WithAlignment(layout.S), // Bottom Center
									// Ensure it sits above navigation bars if possible,
									// currently just bottom aligned in the root Box.
									box.WithModifier(size.FillMax()),
								)(c)
							}
							return c
						},
					),
					box.WithModifier(size.FillMax()),
				)(c)
			},
			surfaceOpts...,
		)(c)
	}
}
