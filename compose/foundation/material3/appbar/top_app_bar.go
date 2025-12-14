package appbar

import (
	"image/color"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/material3/surface"
	padding_modifier "github.com/zodimo/go-compose/internal/modifiers/padding"
	"github.com/zodimo/go-compose/internal/modifiers/size"
	"github.com/zodimo/go-compose/internal/modifiers/weight"

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
	return surface.Surface(
		func(c Composer) Composer {

			return row.Row(
				compose.Sequence(
					func(c Composer) Composer {
						// Navigation Icon
						if navigationIcon != nil {
							box.Box(
								surface.Surface(
									navigationIcon,
									surface.WithContentColor(colors.NavigationIconContentColor),
									surface.WithColor(color.NRGBA{}), // Transparent background
								),
								box.WithAlignment(layout.W),
								box.WithModifier(padding_modifier.Padding(4, 0, 0, 0)), // Start(4)
							)(c)
						} else {
							// Spacer if no navigation icon but we want alignment consistency?
							// Material 3 doesn't strictly require a spacer if missing.
						}
						return c
					},

					// Title
					box.Box(
						func(c Composer) Composer {
							if title != nil {
								return surface.Surface(
									title,
									surface.WithContentColor(colors.TitleContentColor),
									surface.WithColor(color.NRGBA{}), // Transparent
								)(c)
							}
							return c
						},
						box.WithModifier(weight.Weight(1)),                    // Occupy remaining space
						box.WithAlignment(layout.W),                           // Align text to start
						box.WithModifier(padding_modifier.Horizontal(16, 16)), // Horizontal(16, 16)
					),
					box.Box(compose.Id(), box.WithModifier(weight.Weight(1))),
					func(c Composer) Composer {
						// Actions
						if len(actions) > 0 {
							row.Row(
								compose.Sequence(actions...),
								row.WithAlignment(row.Middle),
								row.WithModifier(padding_modifier.Padding(0, 0, 4, 0)), // End(4)
							)(c)
						}
						return c
					},
				),
				row.WithModifier(size.FillMaxWidth()),
				row.WithModifier(size.Height(64)), // Standard Height
				row.WithAlignment(row.Middle),     // Vertical Alignment
			)(c)
		},
		surface.WithModifier(modifier),
		surface.WithColor(colors.ContainerColor),
	)
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
			func(c Composer) Composer {
				return box.Box(
					func(c Composer) Composer {
						// Layer 1: Navigation Icon and Actions
						row.Row(
							func(c Composer) Composer {
								// Navigation Icon
								if opts.NavigationIcon != nil {
									box.Box(
										surface.Surface(
											opts.NavigationIcon,
											surface.WithContentColor(opts.Colors.NavigationIconContentColor),
											surface.WithColor(color.NRGBA{}),
										),
										box.WithAlignment(layout.W),
										box.WithModifier(padding_modifier.Padding(4, 0, 0, 0)),
									)(c)
								}

								// Spacer to push Actions to end
								box.Box(compose.Id(), box.WithModifier(weight.Weight(1)))(c)

								// Actions
								if len(opts.Actions) > 0 {
									// Wrap actions in a Row with Middle alignment
									row.Row(
										compose.Sequence(opts.Actions...),
										row.WithAlignment(row.Middle),
										row.WithModifier(padding_modifier.Padding(0, 0, 4, 0)),
									)(c)
								}
								return c
							},
							row.WithModifier(size.FillMaxWidth()),
							row.WithModifier(size.Height(64)),
							row.WithAlignment(row.Middle),
						)(c)

						// Layer 2: Centered Title
						box.Box(
							func(c Composer) Composer {
								if title != nil {
									return surface.Surface(
										title,
										surface.WithContentColor(opts.Colors.TitleContentColor),
										surface.WithColor(color.NRGBA{}),
									)(c)
								}
								return c
							},
							box.WithAlignment(layout.Center),
							box.WithModifier(size.FillMax()), // Consume space to allow centering
						)(c)
						return c
					},
					box.WithModifier(size.FillMaxWidth()),
					box.WithModifier(size.Height(64)),
				)(c)
			},
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
	return surface.Surface(
		func(c Composer) Composer {
			return column.Column(
				func(c Composer) Composer {
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
					)(c)

					// Bottom Row: Title
					box.Box(
						func(c Composer) Composer {
							if title != nil {
								// Apply specific typography style here if needed?
								// Taking title as is for now.
								return surface.Surface(
									func(c Composer) Composer {
										// Default style for Headline?
										// For now, just render title.
										// Ideally we should apply Material3 typography h5 or h6.
										// But we don't have typography passed in easily yet.
										// Rely on user passing Text with style.

										// Using ProvideTextStyle? Not yet implemented in foundation?
										// Just render.
										return title(c)
									},
									surface.WithContentColor(colors.TitleContentColor),
									surface.WithColor(color.NRGBA{}),
								)(c)
							}
							return c
						},
						box.WithModifier(size.FillMaxWidth()),
						box.WithModifier(weight.Weight(1)), // Fill remaining height
						box.WithAlignment(layout.SW),       // Start, Bottom
						box.WithModifier(padding_modifier.Padding(16, 0, 16, titleBottomPadding)),
					)(c)
					return c
				},
				column.WithModifier(size.FillMaxWidth()),
				column.WithModifier(size.Height(maxHeight)),
			)(c)
		},
		surface.WithModifier(modifier),
		surface.WithColor(colors.ContainerColor),
	)
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
