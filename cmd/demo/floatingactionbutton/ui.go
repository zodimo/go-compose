package main

import (
	"fmt"
	"image/color"

	"go-compose-dev/compose/foundation/icon"
	"go-compose-dev/compose/foundation/layout/box"
	"go-compose-dev/compose/foundation/material3/floatingactionbutton"
	"go-compose-dev/compose/foundation/material3/scaffold"
	"go-compose-dev/compose/foundation/material3/text"
	"go-compose-dev/internal/modifiers/size"
	"go-compose-dev/pkg/api"

	"gioui.org/layout"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

func UI() api.Composable {
	return func(c api.Composer) api.Composer {
		// State for click count
		count := c.State("click_count", func() any { return 0 })

		return scaffold.Scaffold(
			// Content
			func(c api.Composer) api.Composer {
				return box.Box(
					text.Text(
						fmt.Sprintf("FAB Clicked: %d", count.Get().(int)),
						text.TypestyleBodyLarge, // Added Typestyle
					),
					box.WithAlignment(layout.Center),
					box.WithModifier(size.FillMax()),
				)(c)
			},
			// FAB
			scaffold.WithFloatingActionButton(
				floatingactionbutton.FloatingActionButton(
					func() {
						current := count.Get().(int)
						count.Set(current + 1)
						fmt.Println("FAB Clicked!")
					},
					icon.Icon(icons.ContentAdd),
					floatingactionbutton.WithContainerColor(color.NRGBA{R: 0, G: 100, B: 200, A: 255}), // Explicit color for now
					floatingactionbutton.WithContentColor(color.NRGBA{R: 255, G: 255, B: 255, A: 255}),
				),
			),
		)(c)
	}
}
