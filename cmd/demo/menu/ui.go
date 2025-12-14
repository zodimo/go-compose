package main

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/foundation/icon"
	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/material3/button"
	"github.com/zodimo/go-compose/compose/foundation/material3/menu"
	"github.com/zodimo/go-compose/compose/foundation/material3/scaffold"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/theme"

	"gioui.org/unit"
	mdicons "golang.org/x/exp/shiny/materialdesign/icons"
)

func UI() api.Composable {
	return func(c api.Composer) api.Composer {
		// State for menus
		expanded1 := c.State("menu1", func() any { return false })
		expanded2 := c.State("menu2", func() any { return false })

		return scaffold.Scaffold(
			func(c api.Composer) api.Composer {
				return column.Column(
					func(c api.Composer) api.Composer {

						// Demo 1: Simple Dropdown
						box.Box(
							func(c api.Composer) api.Composer {
								button.Filled(
									func() {
										expanded1.Set(true)
									},
									"Open Menu",
								)(c)

								menu.DropdownMenu(
									expanded1.Get().(bool),
									func() { expanded1.Set(false) },
									func(c api.Composer) api.Composer {
										menu.DropdownMenuItem(
											"Item 1 (Selected)",
											func() {
												fmt.Println("Item 1 Clicked")
												expanded1.Set(false)
											},
											menu.WithLeadingIcon(func(c api.Composer) api.Composer {
												return icon.Icon(
													mdicons.ActionDone,
													icon.WithColorDescriptor(theme.GetThemeManager().ColorRoleDescriptors().PrimaryRoles.Primary),
												)(c)
											}),
										)(c)
										menu.DropdownMenuItem(
											"Item 2 (Disabled)",
											func() {},
											menu.WithEnabled(false),
										)(c)
										menu.DropdownMenuItem(
											"Item 3",
											func() {
												fmt.Println("Item 3 Clicked")
												expanded1.Set(false)
											},
										)(c)
										return c
									},
									// Options
									menu.WithModifier(api.EmptyModifier),
									menu.WithOffset(unit.Dp(0), unit.Dp(0)),
								)(c)
								return c
							},
						)(c)

						// Spacer
						box.Box(func(c api.Composer) api.Composer { return c }, box.WithModifier(padding.Vertical(32, 32)))(c)

						// Demo 2: Leading/Trailing Icons
						box.Box(
							func(c api.Composer) api.Composer {
								button.FilledTonal(
									func() {
										expanded2.Set(true)
									},
									"Menu with Icons",
								)(c)

								menu.DropdownMenu(
									expanded2.Get().(bool),
									func() { expanded2.Set(false) },
									func(c api.Composer) api.Composer {
										menu.DropdownMenuItem(
											"Edit",
											func() { expanded2.Set(false) },
											menu.WithLeadingIcon(func(c api.Composer) api.Composer {
												return icon.Icon(
													mdicons.ContentCreate,
													icon.WithColorDescriptor(theme.GetThemeManager().ColorRoleDescriptors().SurfaceRoles.OnSurface),
												)(c)
											}),
										)(c)
										menu.DropdownMenuItem(
											"Settings",
											func() { expanded2.Set(false) },
											menu.WithLeadingIcon(func(c api.Composer) api.Composer {
												return icon.Icon(
													mdicons.ActionSettings,
													icon.WithColorDescriptor(theme.GetThemeManager().ColorRoleDescriptors().SurfaceRoles.OnSurface),
												)(c)
											}),
										)(c)
										menu.DropdownMenuItem(
											"Share",
											func() { expanded2.Set(false) },
											menu.WithLeadingIcon(func(c api.Composer) api.Composer {
												return icon.Icon(
													mdicons.SocialShare,
													icon.WithColorDescriptor(theme.GetThemeManager().ColorRoleDescriptors().SurfaceRoles.OnSurface),
												)(c)
											}),
										)(c)
										return c
									},
									// Options
									menu.WithOffset(unit.Dp(10), unit.Dp(10)), // Offset example
								)(c)
								return c
							},
						)(c)

						return c
					},
					column.WithModifier(padding.All(32)),
				)(c)
			},
			// Scaffold Options
		)(c)
	}
}
