package main

import (
	"fmt"
	"go-compose-dev/compose"
	"go-compose-dev/compose/foundation/layout/column"
	"go-compose-dev/compose/foundation/material3/button"
	"go-compose-dev/compose/foundation/material3/menu"
	m3Text "go-compose-dev/compose/foundation/material3/text"
	"go-compose-dev/internal/modifiers/size"
	"go-compose-dev/pkg/api"
)

func UI(c api.Composer) api.LayoutNode {
	expanded := c.State("menu_expanded", func() any { return false })

	c = column.Column(
		compose.Sequence(
			button.Filled(func() {
				expanded.Set(!expanded.Get().(bool))
			}, "Toggle Menu"),

			menu.DropdownMenu(
				expanded.Get().(bool),
				func() { expanded.Set(false) },
				compose.Sequence(
					menu.DropdownMenuItem(
						func() {
							fmt.Println("Item 1 clicked")
							expanded.Set(false)
						},
						m3Text.Text("Item 1", m3Text.TypestyleBodyMedium),
					),
					menu.DropdownMenuItem(
						func() {
							fmt.Println("Item 2 clicked")
							expanded.Set(false)
						},
						m3Text.Text("Item 2", m3Text.TypestyleBodyMedium),
					),
				),
			),
		),
		column.WithModifier(size.FillMax()),
		column.WithAlignment(column.Middle),
	)(c)

	return c.Build()
}
