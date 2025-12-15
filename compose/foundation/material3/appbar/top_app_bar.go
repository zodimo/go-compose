package appbar

import (
	"image/color"

	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/foundation/material3/surface"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/modifiers/weight"

	"gioui.org/layout"
)

// SingleRowTopAppBar is an internal component to layout the TopAppBar content in a single row.
// It is used by SmallTopAppBar and CenterAlignedTopAppBar.
func SingleRowTopAppBar(
	modifier Modifier,
	title Composable,
	navigationIcon Composable,
	actions []Composable,
	colors TopAppBarColors,
) Composable {
	return func(c Composer) Composer {
		return surface.Surface(
			row.Row(
				c.Sequence(
					// Navigation Icon
					c.When(
						navigationIcon != nil,
						box.Box(
							surface.Surface(
								navigationIcon,
								surface.WithContentColor(colors.NavigationIconContentColor),
								surface.WithColor(color.NRGBA{}), // Transparent background
							),
							box.WithAlignment(layout.W),
							box.WithModifier(padding.Padding(4, 0, 0, 0)), // Start(4)
						),
					),

					// Title
					box.Box(
						c.When(
							title != nil,
							surface.Surface(
								title,
								surface.WithContentColor(colors.TitleContentColor),
								surface.WithColor(color.NRGBA{}), // Transparent
							),
						),
						box.WithModifier(weight.Weight(1)),           // Occupy remaining space
						box.WithAlignment(layout.W),                  // Align text to start
						box.WithModifier(padding.Horizontal(16, 16)), // Horizontal(16, 16)
					),
					spacer.SpacerWeight(1),
					// Actions
					c.When(
						len(actions) > 0,
						row.Row(
							c.Sequence(actions...),
							row.WithAlignment(row.Middle),
							row.WithModifier(padding.Padding(0, 0, 4, 0)), // End(4)
						),
					),
				),
				row.WithModifier(size.FillMaxWidth()),
				row.WithModifier(size.Height(64)), // Standard Height
				row.WithAlignment(row.Middle),     // Vertical Alignment
			),
			surface.WithModifier(modifier),
			surface.WithColor(colors.ContainerColor),
		)(c)
	}
}

// TopAppBar displays information and actions at the top of a screen.
// This is equivalent to SmallTopAppBar in Material 3.
func TopAppBar(
	title Composable,
	options ...TopAppBarOption,
) Composable {
	return func(c Composer) Composer {
		opts := DefaultTopAppBarOptions()
		for _, option := range options {
			if option == nil {
				continue
			}
			option(&opts)
		}

		return SingleRowTopAppBar(
			opts.Modifier,
			title,
			opts.NavigationIcon,
			opts.Actions,
			opts.Colors,
		)(c)
	}
}

// CenterAlignedTopAppBar displays information and actions at the top of a screen.
// The title is centered horizontally.
func CenterAlignedTopAppBar(
	title Composable,
	options ...TopAppBarOption,
) Composable {
	return func(c Composer) Composer {
		opts := DefaultTopAppBarOptions()
		for _, option := range options {
			if option == nil {
				continue
			}
			option(&opts)
		}

		return surface.Surface(
			box.Box(
				c.Sequence(
					// Layer 1: Navigation Icon and Actions
					row.Row(
						c.Sequence(
							// Navigation Icon
							c.When(opts.NavigationIcon != nil,
								box.Box(
									surface.Surface(
										opts.NavigationIcon,
										surface.WithContentColor(opts.Colors.NavigationIconContentColor),
										surface.WithColor(color.NRGBA{}),
									),
									box.WithAlignment(layout.W),
									box.WithModifier(padding.Padding(4, 0, 0, 0)),
								),
							),
							// Spacer to push Actions to end
							spacer.SpacerWeight(1),
							// Actions
							c.When(
								len(opts.Actions) > 0,
								row.Row(
									c.Sequence(opts.Actions...),
									row.WithAlignment(row.Middle),
									row.WithModifier(padding.Padding(0, 0, 4, 0)),
								),
							),
						),
						row.WithModifier(size.FillMaxWidth()),
						row.WithModifier(size.Height(64)),
						row.WithAlignment(row.Middle),
					),
					// Layer 2: Centered Title
					box.Box(
						c.When(
							title != nil,
							surface.Surface(
								title,
								surface.WithContentColor(opts.Colors.TitleContentColor),
								surface.WithColor(color.NRGBA{}),
							),
						),
						box.WithAlignment(layout.Center),
						box.WithModifier(size.FillMax()), // Consume space to allow centering
					),
				),
				box.WithModifier(size.FillMaxWidth()),
				box.WithModifier(size.Height(64)),
			),
			surface.WithModifier(opts.Modifier),
			surface.WithColor(opts.Colors.ContainerColor),
		)(c)
	}
}

// TwoRowsTopAppBar is an internal component to layout Medium and Large TopAppBars.
func TwoRowsTopAppBar(
	modifier Modifier,
	title Composable,
	titleBottomPadding int,
	maxHeight int,
	navigationIcon Composable,
	actions []Composable,
	colors TopAppBarColors,
) Composable {
	return func(c Composer) Composer {
		return surface.Surface(
			column.Column(
				c.Sequence(
					// Top Row: Nav Icon + Actions (No Title)
					SingleRowTopAppBar(
						EmptyModifier,
						nil, // No title in top row for expanded state
						navigationIcon,
						actions,
						TopAppBarColors{
							// Use transparent container for nested app bar to avoid layering issues?
							// Or use same colors. SingleRowTopAppBar sets container color.
							ContainerColor:             color.NRGBA{}, // Transparent, let parent surface color show
							NavigationIconContentColor: colors.NavigationIconContentColor,
							ActionIconContentColor:     colors.ActionIconContentColor,
						},
					),
					// Bottom Row: Title
					box.Box(
						c.When(
							title != nil,
							surface.Surface(
								// Default style for Headline?
								// For now, just render title.
								// Ideally we should apply Material3 typography h5 or h6.
								// But we don't have typography passed in easily yet.
								// Rely on user passing Text with style.

								// Using ProvideTextStyle? Not yet implemented in foundation?
								// Just render.
								title,
								surface.WithContentColor(colors.TitleContentColor),
								surface.WithColor(color.NRGBA{}),
							),
						),
						box.WithModifier(size.FillMaxWidth()),
						box.WithModifier(weight.Weight(1)), // Fill remaining height
						box.WithAlignment(layout.SW),       // Start, Bottom
						box.WithModifier(padding.Padding(16, 0, 16, titleBottomPadding)),
					),
				),
				column.WithModifier(size.FillMaxWidth()),
				column.WithModifier(size.Height(maxHeight)),
			),
			surface.WithModifier(modifier),
			surface.WithColor(colors.ContainerColor),
		)(c)
	}
}

// MediumTopAppBar displays information and actions at the top of a screen.
// It has a larger height than SmallTopAppBar.
func MediumTopAppBar(
	title Composable,
	options ...TopAppBarOption,
) Composable {
	return func(c Composer) Composer {
		opts := DefaultTopAppBarOptions()
		for _, option := range options {
			if option == nil {
				continue
			}
			option(&opts)
		}

		return TwoRowsTopAppBar(
			opts.Modifier,
			title,
			16,  // Bottom padding - Reduced from 24 to avoid cutoff
			112, // Max height (64 + 48)
			opts.NavigationIcon,
			opts.Actions,
			opts.Colors,
		)(c)
	}
}

// LargeTopAppBar displays information and actions at the top of a screen.
// It has the largest height.
func LargeTopAppBar(
	title Composable,
	options ...TopAppBarOption,
) Composable {
	return func(c Composer) Composer {
		opts := DefaultTopAppBarOptions()
		for _, option := range options {
			if option == nil {
				continue
			}
			option(&opts)
		}

		return TwoRowsTopAppBar(
			opts.Modifier,
			title,
			28,  // Bottom padding
			152, // Max height (64 + 88)
			opts.NavigationIcon,
			opts.Actions,
			opts.Colors,
		)(c)
	}
}
