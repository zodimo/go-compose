package main

import (
	"fmt"
	"image/color"

	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	ftext "github.com/zodimo/go-compose/compose/foundation/text"
	"github.com/zodimo/go-compose/compose/material3"
	"github.com/zodimo/go-compose/compose/material3/icon"
	"github.com/zodimo/go-compose/compose/material3/navigationdrawer"
	"github.com/zodimo/go-compose/compose/material3/navigationrail"
	"github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	uiText "github.com/zodimo/go-compose/compose/ui/text"
	"github.com/zodimo/go-compose/compose/ui/text/style"
	"github.com/zodimo/go-compose/modifiers/clickable"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/theme"

	"github.com/zodimo/go-compose/pkg/api"

	"gioui.org/widget"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type ColorDescriptor = theme.ColorDescriptor

var colorHelper = theme.ColorHelper

func UI() api.Composable {
	return func(c api.Composer) api.Composer {
		// State
		key := c.GenerateID()
		path := c.GetPath()
		// Manual unique key usage
		selKey := key.String() + "/" + path.String() + "/selected"
		drawerKey := key.String() + "/" + path.String() + "/drawer"

		selectedValue := c.State(selKey, func() any { return 0 })
		drawerOpen := c.State(drawerKey, func() any { return false })

		selectedIndex := selectedValue.Get().(int)
		isDrawerOpen := drawerOpen.Get().(bool)

		return navigationdrawer.ModalNavigationDrawer(
			// Drawer Content
			func(c api.Composer) api.Composer {
				return column.Column(
					func(c api.Composer) api.Composer {
						// Drawer Header / Headline
						text.TextWithStyle("GoCompose", text.TypestyleHeadlineSmall,
							ftext.WithModifier(padding.Padding(28, 24, 16, 24)),
						)(c)

						items := []struct {
							Label string
							Icon  []byte
						}{
							{"Home", icons.ActionHome},
							{"Search", icons.ActionSearch},
							{"Settings", icons.ActionSettings},
						}

						for i, item := range items {
							isSelected := selectedIndex == i
							navigationdrawer.NavigationDrawerItem( // Changed from navigationrail.NavigationRailItem
								isSelected,
								//onClick
								func() {
									selectedValue.Set(i)
									drawerOpen.Set(false) // Close drawer
								},

								//Icon
								c.If(
									isSelected,
									icon.Icon(
										item.Icon,
										icon.WithColor(colorHelper.ColorSelector().SecondaryRoles.OnContainer),
										icon.WithModifier(size.Size(24, 24)),
									),
									icon.Icon(
										item.Icon,
										icon.WithColor(colorHelper.ColorSelector().SurfaceRoles.OnVariant),
										icon.WithModifier(size.Size(24, 24)),
									),
								),

								func(c api.Composer) api.Composer {
									var textColor graphics.Color
									if isSelected {
										textColor = material3.Theme(c).ColorScheme().SecondaryContainer.OnColor
									} else {
										textColor = material3.Theme(c).ColorScheme().SurfaceVariant.OnColor
									}

									return text.LabelMedium(
										item.Label,
										ftext.WithTextStyleOptions(
											uiText.WithColor(textColor),
										),
									)(c)
								},
								ui.EmptyModifier,
							)(c)
							spacer.Uniform(4)(c) // Changed spacer size
						}
						return c
					},
					column.WithModifier(padding.Padding(12, 12, 12, 12)), // Added modifier
				)(c)
			},
			// Main Content
			func(c api.Composer) api.Composer {
				return row.Row(
					func(c api.Composer) api.Composer {
						// Navigation Rail (Collapsed state)
						navigationrail.NavigationRail(
							ui.EmptyModifier,
							func(c api.Composer) api.Composer {
								// Header (Menu Icon to toggle drawer)
								return box.Box(
									func(c api.Composer) api.Composer {
										return icon.Icon(
											icons.NavigationMenu,
											icon.WithColor(colorHelper.SpecificColor(graphics.FromNRGBA(color.NRGBA{A: 255}))),
											icon.WithModifier(size.Size(24, 24)),
										)(c)
									},
									box.WithAlignment(box.Center),
									box.WithModifier(
										size.FillMaxWidth().
											Then(size.Height(48)). // Minimum touch target height
											Then(clickable.OnClick(func() {
												drawerOpen.Set(true)
											}, clickable.WithClickable(c.State("menu_click", func() any { return &widget.Clickable{} }).Get().(*widget.Clickable)))),
									),
								)(c)
							},
							func(c api.Composer) api.Composer {
								items := []struct {
									Label string
									Icon  []byte
								}{
									{"Home", icons.ActionHome},
									{"Search", icons.ActionSearch},
									{"Settings", icons.ActionSettings},
								}

								for i, item := range items {
									isSelected := selectedIndex == i
									navigationrail.NavigationRailItem(
										isSelected,
										func() {
											selectedValue.Set(i)
										},
										//Icon
										c.If(
											isSelected,
											icon.Icon(
												item.Icon,
												icon.WithColor(colorHelper.SpecificColor(
													material3.Theme(c).ColorScheme().SecondaryContainer.OnColor,
												)),
												icon.WithModifier(size.Size(24, 24)),
											),
											icon.Icon(
												item.Icon,
												icon.WithColor(colorHelper.SpecificColor(
													material3.Theme(c).ColorScheme().SurfaceVariant.OnColor,
												)),
												icon.WithModifier(size.Size(24, 24)),
											),
										),
										//label
										func(c api.Composer) api.Composer {
											var textColor graphics.Color
											if isSelected {
												textColor = material3.Theme(c).ColorScheme().SecondaryContainer.OnColor // theme.ColorHelper.ColorSelector().SecondaryRoles.OnContainer
											} else {
												textColor = material3.Theme(c).ColorScheme().SurfaceVariant.OnColor // theme.ColorHelper.ColorSelector().SurfaceRoles.OnVariant
											}

											return text.TextWithStyle(
												item.Label,
												text.TypestyleLabelMedium,
												ftext.WithTextStyleOptions(
													uiText.WithColor(textColor),
												),
												ftext.WithAlignment(style.TextAlignMiddle),
											)(c)
										},
										ui.EmptyModifier,
									)(c)
									spacer.Uniform(12)(c)
								}
								return c
							},
						)(c)

						// Main Content Body
						column.Column(
							func(c api.Composer) api.Composer {
								text.TextWithStyle(fmt.Sprintf("Selected Page: %d", selectedIndex), text.TypestyleHeadlineMedium)(c)
								return c
							},
							column.WithModifier(size.FillMax()),
							column.WithAlignment(column.Middle), // Center alignment
						)(c)

						return c
					},
					row.WithModifier(size.FillMax()),
				)(c)
			},
			navigationdrawer.WithIsOpen(isDrawerOpen),
			navigationdrawer.WithOnClose(func() {
				drawerOpen.Set(false)
			}),
		)(c)
	}
}
