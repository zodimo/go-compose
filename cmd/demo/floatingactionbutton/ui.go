package main

import (
	"fmt"
	"image/color"

	"github.com/zodimo/go-compose/compose/foundation/icon"
	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/material3/floatingactionbutton"
	"github.com/zodimo/go-compose/compose/foundation/material3/scaffold"
	"github.com/zodimo/go-compose/compose/foundation/material3/text"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/theme"

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
					floatingactionbutton.WithContainerColor(theme.ColorHelper.SpecificColor(graphics.FromNRGBA(color.NRGBA{R: 0, G: 100, B: 200, A: 255}))), // Explicit color for now
					floatingactionbutton.WithContentColor(theme.ColorHelper.SpecificColor(graphics.FromNRGBA(color.NRGBA{R: 255, G: 255, B: 255, A: 255}))),
				),
			),
		)(c)
	}
}
