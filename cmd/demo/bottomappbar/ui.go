package main

import (
	"image/color"

	"gioui.org/layout"

	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/material3/bottomappbar"
	"github.com/zodimo/go-compose/compose/foundation/material3/iconbutton"
	m3text "github.com/zodimo/go-compose/compose/foundation/material3/text"
	"github.com/zodimo/go-compose/compose/foundation/text"
	"github.com/zodimo/go-compose/pkg/api" // Import api for Composable/Composer

	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/modifiers/weight"

	mdicons "golang.org/x/exp/shiny/materialdesign/icons"
)

func UI() api.Composable {
	return func(c api.Composer) api.Composer {
		return column.Column(
			c.Sequence(
				// Content Area
				box.Box(
					func(c api.Composer) api.Composer {
						return m3text.Text(
							"Content Area",
							m3text.TypestyleBodyLarge, // valid style
							text.WithTextStyleOptions(
								text.StyleWithColor(color.NRGBA{A: 255}),
							),
						)(c)
					},
					box.WithModifier(weight.Weight(1)),
					box.WithAlignment(layout.Center),
				),
				// Bottom App Bar with FAB
				bottomappbar.BottomAppBar(
					row.Row(
						c.Sequence(
							IconButton(mdicons.NavigationMenu),
							IconButton(mdicons.ActionSearch),
							IconButton(mdicons.ContentSave),
							IconButton(mdicons.SocialShare),
						),
					),
					bottomappbar.WithFloatingActionButton(
						func(c api.Composer) api.Composer {
							// Simulate FAB
							return iconbutton.Filled(
								func() { /* click */ },
								mdicons.ContentAdd,
								"Add",
							)(c)
						},
					),
				),
			),
			column.WithModifier(size.FillMax()),
		)(c)
	}
}

func IconButton(iconData []byte) api.Composable {
	return func(c api.Composer) api.Composer {
		return iconbutton.Standard(
			func() {
				// OnClick
			},
			iconData,
			"Icon",
		)(c)
	}
}
