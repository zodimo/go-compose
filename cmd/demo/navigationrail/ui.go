package main

import (
	"fmt"
	"image/color"

	"github.com/zodimo/go-compose/compose/foundation/icon"
	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/foundation/material3/navigationdrawer"
	"github.com/zodimo/go-compose/compose/foundation/material3/navigationrail"
	"github.com/zodimo/go-compose/compose/foundation/material3/text"
	ftext "github.com/zodimo/go-compose/compose/foundation/text"
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
						text.Text("GoCompose", text.TypestyleHeadlineSmall,
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
									tm := theme.GetThemeManager()
									m3 := tm.GetMaterial3Theme()
									var textColor color.NRGBA
									if isSelected {
										textColor = m3.Scheme.SecondaryContainer.OnColor.AsNRGBA()
									} else {
										textColor = m3.Scheme.SurfaceVariant.OnColor.AsNRGBA()
									}

									return text.Text(
										item.Label,
										text.TypestyleLabelMedium,
										ftext.WithTextStyleOptions(
											ftext.StyleWithColor(textColor),
										),
										// ftext.WithAlignment(ftext.Middle), // align to the middle horizontally not vertically
										// ftext.WithModifier(background.Background(theme.ColorHelper.SpecificColor(color.NRGBA{R: 0, G: 0, B: 200, A: 200})))),

									)(c)
								},
								api.EmptyModifier,
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
							api.EmptyModifier,
							func(c api.Composer) api.Composer {
								// Header (Menu Icon to toggle drawer)
								return box.Box(
									func(c api.Composer) api.Composer {
										return icon.Icon(
											icons.NavigationMenu,
											icon.WithColor(colorHelper.SpecificColor(color.NRGBA{A: 255})),
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
												icon.WithColor(colorHelper.ColorSelector().SecondaryRoles.OnContainer),
												icon.WithModifier(size.Size(24, 24)),
											),
											icon.Icon(
												item.Icon,
												icon.WithColor(colorHelper.ColorSelector().SurfaceRoles.OnVariant),
												icon.WithModifier(size.Size(24, 24)),
											),
										),
										//label
										func(c api.Composer) api.Composer {
											tm := theme.GetThemeManager()
											m3 := tm.GetMaterial3Theme()
											var textColor color.NRGBA
											if isSelected {
												textColor = m3.Scheme.SecondaryContainer.OnColor.AsNRGBA()
											} else {
												textColor = m3.Scheme.SurfaceVariant.OnColor.AsNRGBA()
											}

											return text.Text(
												item.Label,
												text.TypestyleLabelMedium,
												ftext.WithTextStyleOptions(
													ftext.StyleWithColor(textColor),
												),
												ftext.WithAlignment(ftext.Middle),
											)(c)
										},
										api.EmptyModifier,
									)(c)
									spacer.Uniform(12)(c)
								}
								return c
							},
						)(c)

						// Main Content Body
						column.Column(
							func(c api.Composer) api.Composer {
								text.Text(fmt.Sprintf("Selected Page: %d", selectedIndex), text.TypestyleHeadlineMedium)(c)
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
