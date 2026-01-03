package scaffold

import (
	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/material3/surface"
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
		opts := DefaultScaffoldOptions(c)
		for _, option := range options {
			if option == nil {
				continue
			}
			option(&opts)
		}

		surfaceOpts := []surface.SurfaceOption{
			surface.WithColor(opts.ContainerColor),
			surface.WithContentColor(opts.ContentColor),
			surface.WithModifier(opts.Modifier),
		}

		return surface.Surface(
			func(c Composer) Composer {
				// We use a Box as the root to overlay FAB and Snackbars over the body
				return box.Box( // Root Stacking Context
					compose.Sequence(
						// Layer 1: The App Structure (TopBar -> Content -> BottomBar)
						column.Column(
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
						),
						// Layer 2: Floating Action Button
						c.When(
							opts.FloatingActionButton != nil,
							c.If(
								opts.FloatingActionButtonPosition == FabPositionCenter,
								box.Box(
									opts.FloatingActionButton,
									box.WithAlignment(layout.S), // Bottom Center
									// Add standard padding for FAB
									box.WithModifier(padding_modifier.All(16).
										// Wrapper must fill max to align FAB relative to screen
										Then(size.FillMax()),
									),
								),
								box.Box(
									opts.FloatingActionButton,
									box.WithAlignment(layout.SE), // Bottom End
									// Add standard padding for FAB
									box.WithModifier(padding_modifier.All(16).
										// Wrapper must fill max to align FAB relative to screen
										Then(size.FillMax()),
									),
								),
							),
						),
						// Layer 3: Snackbar Host
						c.When(
							opts.SnackbarHost != nil,
							opts.SnackbarHost,
						),
					),
					box.WithModifier(size.FillMax()),
				)(c)
			},
			surfaceOpts...,
		)(c)
	}
}
