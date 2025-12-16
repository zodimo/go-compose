package scaffold

import (
	"image/color"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/material3/surface"
	padding_modifier "github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/modifiers/weight"

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
			if option == nil {
				continue
			}
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
									// Top Bar
									c.When(opts.TopBar != nil,
										opts.TopBar,
									),

									// Content Body
									// This needs to expand to fill available space.
									box.Box(
										content,
										box.WithModifier(
											// Expand to fill remaining vertical space
											weight.Weight(1),
										),
									),

									// Bottom Bar
									c.When(opts.BottomBar != nil,
										opts.BottomBar,
									),
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
