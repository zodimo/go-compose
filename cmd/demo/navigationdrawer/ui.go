package main

import (
	"fmt"
	"image/color"

	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/foundation/material3/button"
	"github.com/zodimo/go-compose/compose/foundation/material3/navigationdrawer"
	"github.com/zodimo/go-compose/compose/foundation/material3/text"
	ftext "github.com/zodimo/go-compose/compose/foundation/text"
	"github.com/zodimo/go-compose/modifiers/background"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"
)

func UI() api.Composable {
	return func(c api.Composer) api.Composer {
		// State for drawer type selection
		drawerType := c.State("drawerType", func() any { return "Modal" }).Get().(string)
		setDrawerType := func(t string) {
			c.State("drawerType", nil).Set(t)
		}

		// State for drawer open/closed
		isOpen := c.State("drawerOpen", func() any { return false }).Get().(bool)
		setIsOpen := func(o bool) {
			c.State("drawerOpen", nil).Set(o)
		}

		// Selected item state
		selectedItem := c.State("selectedItem", func() any { return "Inbox" }).Get().(string)
		setSelectedItem := func(i string) {
			c.State("selectedItem", nil).Set(i)
		}

		// Drawer Content
		drawerContent := func(c api.Composer) api.Composer {
			return column.Column(
				func(c api.Composer) api.Composer {
					text.Text("GoCompose Mail", text.TypestyleTitleMedium, ftext.WithModifier(padding.Padding(24, 24, 24, 12)))(c)

					items := []string{"Inbox", "Outbox", "Favorites", "Trash"}
					for _, item := range items {
						navigationdrawer.NavigationDrawerItem(
							selectedItem == item,
							func() { setSelectedItem(item); setIsOpen(false) }, // Close drawer on selection for Modal
							nil, // Icon
							func(c api.Composer) api.Composer {
								text.Text(item, text.TypestyleBodyMedium)(c)
								return c
							},
							api.EmptyModifier,
						)(c)
					}
					return c
				},
				column.WithModifier(padding.All(12)),
			)(c)
		}

		// Main Content
		mainContent := func(c api.Composer) api.Composer {
			return column.Column(
				func(c api.Composer) api.Composer {
					// Header with Toggle
					row.Row(
						func(c api.Composer) api.Composer {
							text.Text("Navigation Drawer Demo", text.TypestyleHeadlineSmall)(c)
							spacer.SpacerWidth(16)(c)
							// Toggle Button (only meaningful for Modal/Dismissible)
							button.Filled(
								func() { setIsOpen(!isOpen) },
								fmt.Sprintf("%v Drawer", isOpen),
							)(c)
							return c
						},
						row.WithModifier(padding.All(16)),
					)(c)

					spacer.SpacerHeight(24)(c)

					// Type Switcher
					row.Row(
						func(c api.Composer) api.Composer {
							types := []string{"Modal", "Dismissible", "Permanent"}
							for _, t := range types {
								isSelected := drawerType == t
								label := t
								if isSelected {
									label = "> " + t + " <"
								}
								// Simple way to show selection state visually via label for now, or use Tonal button
								if isSelected {
									button.Filled(
										func() { setDrawerType(t); setIsOpen(false) },
										label,
									)(c)
								} else {
									button.Text( // Using Text button for unselected
										func() { setDrawerType(t); setIsOpen(false) },
										label,
									)(c)
								}
								spacer.SpacerWidth(8)(c)
							}
							return c
						},
						row.WithModifier(padding.All(16)),
					)(c)

					spacer.SpacerHeight(24)(c)

					text.Text("Current Selection: "+selectedItem, text.TypestyleBodyLarge)(c)

					return c
				},
				column.WithModifier(size.FillMax()),
			)(c)
		}

		// Render based on type
		return box.Box(
			func(c api.Composer) api.Composer {
				switch drawerType {
				case "Modal":
					navigationdrawer.ModalNavigationDrawer(
						drawerContent,
						mainContent,
						navigationdrawer.WithIsOpen(isOpen),
						navigationdrawer.WithOnClose(func() { setIsOpen(false) }),
					)(c)
				case "Dismissible":
					navigationdrawer.DismissibleNavigationDrawer(
						drawerContent,
						mainContent,
						navigationdrawer.WithIsOpen(isOpen),
						navigationdrawer.WithOnClose(func() { setIsOpen(false) }),
					)(c)
				case "Permanent":
					navigationdrawer.PermanentNavigationDrawer(
						drawerContent,
						mainContent,
						api.EmptyModifier,
					)(c)
				}
				return c
			},
			box.WithModifier(
				api.EmptyModifier.
					Then(size.FillMax()).
					Then(background.Background(color.NRGBA{255, 255, 255, 255})),
			),
		)(c)
	}
}
