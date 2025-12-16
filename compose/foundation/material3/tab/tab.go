package tab

import (
	"image/color"

	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/material3/surface"
	"github.com/zodimo/go-compose/modifiers/background"
	"github.com/zodimo/go-compose/modifiers/clickable"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/theme"
)

// TabRow contains a row of Tabs and displays an indicator underneath the currently selected Tab.
func TabRow(
	selectedTabIndex int,
	content Composable,
	options ...TabRowOption,
) Composable {
	return func(c Composer) Composer {
		opts := DefaultTabRowOptions()
		for _, option := range options {
			if option == nil {
				continue
			}
			option(&opts)
		}

		// Calculate indicator offset?
		// For now, simpler implementation: standard row.
		// Handling indicator animation requires knowing the positions of tabs.
		// This is hard with a simple Row.
		// We might need a custom layout or SubcomposeLayout equivalent if we want a sliding indicator.
		// For MVP, we can render the indicator INSIDE the selected Tab?
		// No, visual spec says it's a track at the bottom.

		// Alternative: Use a standard Row, but wrap it in a Surface.
		// The indicator handling is tricky without "SubcomposeLayout".
		// But in Gio, we can probably do it if we know the Width/Count.
		// If TabRow implies Equal Width tabs, the indicator position is just (Index * Width / Count).

		SurfaceOptions := []surface.SurfaceOption{}
		// Surface Options
		if opts.ContainerColor.IsSome() {
			SurfaceOptions = append(SurfaceOptions, surface.WithColor(theme.ColorHelper.SpecificColor(opts.ContainerColor.UnwrapUnsafe())))
		}
		if opts.ContentColor.IsSome() {
			SurfaceOptions = append(SurfaceOptions, surface.WithContentColor(theme.ColorHelper.SpecificColor(opts.ContentColor.UnwrapUnsafe())))
		}

		SurfaceOptions = append(SurfaceOptions, surface.WithModifier(opts.Modifier.Then(size.FillMaxWidth())))

		// Let's implement fixed tabs (filled width) for now.
		return surface.Surface(
			// Stack: Row of Tabs + Indicator
			// But Indicator needs to be top/bottom aligned?
			// Actually, TabRow container usually handles it.

			// Let's use a Column: Row(Tabs) then Indicator? No, indicator is overlay or below.
			// In M3, indicator is at the bottom of the container.

			// Let's use a Box to stack:
			box.Box(
				func(c Composer) Composer {
					// 1. The Row of Tabs
					row.Row(
						content,
						row.WithModifier(size.FillMaxWidth()),
					)(c)

					// 2. The Indicator (Absolute positioned? or just calculated)
					// For Fixed TabRow, we can simulate the indicator by rendering it
					// based on selectedTabIndex and total tabs (if we knew them).
					// But we don't know total tabs count here easily unless valid content is analyzed.
					// Or we require the User to pass the count? No.

					// Workaround:
					// If we only support "Indicator inside Tab" (active state), it's easier.
					// But M3 has a continuous track.

					// For now, let's skip the sliding indicator and just let Tabs render their own active state (border bottom?).
					// Wait, TabRowDefaults.Indicator is a composable.
					// If we want a proper indicator, we need strict layout.

					// Let's simplify: Standard TabRow just renders the content in a Row.
					// We can verify this basic behavior first.

					return c
				},
				box.WithModifier(size.FillMaxWidth()),
			),
			SurfaceOptions...,
		)(c)
	}
}

// Tab is a single tab in a TabRow.
func Tab(
	selected bool,
	onClick func(),
	content Composable, // Usually text and/or icon
	options ...TabOption,
) Composable {
	return func(c Composer) Composer {
		opts := DefaultTabOptions()
		for _, option := range options {
			if option == nil {
				continue
			}
			option(&opts)
		}

		// Tab is a semantic button
		// It needs to be clickable.
		// We can reuse Button or build from Surface + Clickable.
		// Given specific styling needs, Surface + Clickable is better.

		// Clickable state
		// We need a Clickable modifier or widget.
		// github.com/zodimo/go-compose usually has interaction handling.

		// Let's look at how Button does clicking.
		// It wraps `button.Button` from gio-mw.
		// We can use `surface.Surface` and `modifiers.Clickable`.
		// But `Clickable` modifier might not exist or be easy.
		// Using a transparent button could work.

		// For now, let's use a simple Column to structure the tab:
		// Icon (optional)
		// Text (optional)
		// Indicator (if selected and we render it here)

		contentColor := opts.UnselectedContentColor
		if selected {
			contentColor = opts.SelectedContentColor
		}

		return surface.Surface(
			func(c Composer) Composer {
				// Determine layout: Column?
				return column.Column(
					func(c Composer) Composer {
						if content != nil {
							content(c)
						}

						// If we render indicator here:
						if selected {
							// Render Indicator
							box.Box(
								func(c Composer) Composer { return c },
								box.WithModifier(
									size.FillMaxWidth().
										Then(size.Height(3)). // 3dp
										Then(background.Background(theme.ColorHelper.SpecificColor(TabRowDefaults.IndicatorColor()))),
								),
							)(c)
						} else {
							// Placeholder to keep height consistent?
							// Or just align bottom.
							box.Box(
								func(c Composer) Composer { return c },
								box.WithModifier(
									size.FillMaxWidth().
										Then(size.Height(3)).
										Then(background.Background(theme.ColorHelper.SpecificColor(color.NRGBA{}))), // Transparent
								),
							)(c)
						}
						return c
					},
					column.WithAlignment(column.Middle),                             // Center content
					column.WithSpacing(column.SpaceBetween),                         // Push content up, indicator down?
					column.WithModifier(size.FillMaxHeight().Then(padding.All(12))), // Padding
				)(c)
			},
			surface.WithModifier(
				opts.Modifier.
					Then(size.FillMaxHeight()).       // Fill row height
					Then(clickable.OnClick(onClick)), // Use clickable package
			),
			surface.WithColor(theme.ColorHelper.SpecificColor(color.NRGBA{})), // Transparent container
			surface.WithContentColor(theme.ColorHelper.SpecificColor(contentColor)),
		)(c)
	}
}
