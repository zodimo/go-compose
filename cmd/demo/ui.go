package main

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/material/button"
	m3Button "github.com/zodimo/go-compose/compose/foundation/material3/button"
	m3Card "github.com/zodimo/go-compose/compose/foundation/material3/card"
	m3Divider "github.com/zodimo/go-compose/compose/foundation/material3/divider"
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"

	"image/color"

	"github.com/zodimo/go-compose/compose/foundation/text"
	"github.com/zodimo/go-compose/modifiers/background"
	"github.com/zodimo/go-compose/modifiers/clickable"
	"github.com/zodimo/go-compose/modifiers/clip"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/modifiers/weight"
	"github.com/zodimo/go-compose/pkg/api"
)

func UI(c api.Composer) api.LayoutNode {

	counterCell := c.State("counter", func() any { return 0 })

	c = column.Column(
		c.Sequence(
			row.Row(c.Sequence(
				column.Column(
					c.Sequence(
						text.Text(fmt.Sprintf("state=%v", counterCell.Get()),
							text.WithTextStyleOptions(
								text.StyleWithColor(color.NRGBA{R: 255, G: 255, B: 255, A: 255}),
							),
						),
					),
					// column.WithModifier(size.Size(200, 200, size.SizeRequired())),
					column.WithModifier(clickable.OnClick(func() {
						fmt.Println("First Column clicked!!")
						counterCell.Set(counterCell.Get().(int) + 1)
					})),
					column.WithModifier(weight.Weight(1)),
					column.WithModifier(background.Background(color.NRGBA{R: 0, G: 0, B: 200, A: 200})),
					column.WithModifier(padding.All(20)),
					column.WithModifier(background.Background(color.NRGBA{R: 0, G: 100, B: 200, A: 200})),
				),
				column.Column(
					c.Sequence(
						button.Button(func() {
							counterCell.Set(counterCell.Get().(int) + 1)
							fmt.Println("Button clicked!")
						}, "click me",
						// button.WithModifier(size.FillMax())
						),
					),
					column.WithModifier(weight.Weight(1)),
					column.WithModifier(background.Background(color.NRGBA{R: 150, G: 0, B: 0, A: 200})),
				),
				column.Column(
					c.Sequence(),
					column.WithModifier(clip.Clip(shape.ShapeCircle)),
					column.WithModifier(size.Size(100, 50)),
					column.WithModifier(background.Background(color.NRGBA{R: 100, G: 0, B: 0, A: 200})),
					column.WithModifier(clickable.OnClick(func() {
						fmt.Println("Last Column clicked!!")
					})),
				),
			),
				row.WithModifier(size.Size(500, 300)),
				row.WithModifier(background.Background(color.NRGBA{R: 0, G: 200, B: 0, A: 200})),
			),
			text.Text("hello world",
				text.Selectable(),
				text.WithAlignment(text.Middle),
				text.WithTextStyleOptions(
					text.StyleWithColor(color.NRGBA{R: 255, G: 255, B: 255, A: 255}),
				),
				text.WithModifier(background.Background(color.NRGBA{R: 100, G: 0, B: 0, A: 150})),
				text.WithModifier(padding.All(20)),
				text.WithModifier(background.Background(color.NRGBA{R: 200, G: 0, B: 50, A: 50})),
			),

			m3Button.Text(func() {
				fmt.Println("M3 Text Button clicked!")
			}, "Hello M3 Text Button"),
			m3Button.Outlined(func() {
				fmt.Println("M3 Outlined Button clicked!")
			}, "Hello M3 Outlined Button"),
			m3Button.Filled(func() {
				fmt.Println("M3 Filled Button clicked!")
			}, "Hello M3 Filled Button",
				m3Button.WithModifier(padding.All(20)),
			),
			m3Divider.Divider(),
			m3Button.FilledTonal(func() {
				fmt.Println("M3 Filled Tonal Button clicked!")
			}, "Hello M3 FilledTonal Button",
				m3Button.WithModifier(padding.All(20)),
			),
			m3Divider.Divider(),
			m3Card.Filled(m3Card.CardContents(
				m3Card.Content(text.Text("Filled")),
			)),
		),
		column.WithModifier(size.FillMax()),
		column.WithModifier(background.Background(color.NRGBA{R: 200, G: 0, B: 0, A: 50})),

		column.WithAlignment(column.Middle),
	)(c)

	return c.Build()

}
