package main

import (
	"fmt"
	"go-compose-dev/compose"
	"go-compose-dev/compose/foundation/layout/column"
	"go-compose-dev/compose/foundation/layout/row"
	"go-compose-dev/compose/foundation/material/button"
	m3Button "go-compose-dev/compose/foundation/material3/button"
	m3Divider "go-compose-dev/compose/foundation/material3/divider"
	"go-compose-dev/compose/ui/graphics/shape"

	"go-compose-dev/compose/foundation/text"
	"go-compose-dev/internal/modifiers/background"
	"go-compose-dev/internal/modifiers/clickable"
	"go-compose-dev/internal/modifiers/clip"
	"go-compose-dev/internal/modifiers/padding"
	"go-compose-dev/internal/modifiers/size"
	"go-compose-dev/internal/modifiers/weight"
	"go-compose-dev/pkg/api"
	"image/color"
)

func UI(c api.Composer) api.LayoutNode {

	counterCell := c.State("counter", func() any { return 0 })

	c = column.Column(
		compose.Sequence(
			row.Row(compose.Sequence(
				column.Column(
					compose.Sequence(
						text.Text(fmt.Sprintf("state=%v", counterCell.Get()),
							text.WithTextStyleOptions(
								text.StyleWithColor(color.NRGBA{R: 255, G: 255, B: 255, A: 255}),
							),
						),
					),
					// column.WithModifier(size.Size(200, 200, size.SizeRequired())),
					column.WithModifier(clickable.OnClick(func() {
						fmt.Println("First Column clicked!!")
					})),
					column.WithModifier(weight.Weight(1)),
					column.WithModifier(background.Background(color.NRGBA{R: 0, G: 0, B: 200, A: 200})),
					column.WithModifier(padding.All(20)),
					column.WithModifier(background.Background(color.NRGBA{R: 0, G: 100, B: 200, A: 200})),
				),
				column.Column(
					compose.Sequence(),
					column.WithModifier(weight.Weight(1)),
					column.WithModifier(background.Background(color.NRGBA{R: 150, G: 0, B: 0, A: 200})),
				),
				column.Column(
					compose.Sequence(
						button.Button(func() {
							counterCell.Set(counterCell.Get().(int) + 1)
							fmt.Println("Button clicked!")
						}, "click me",
							button.WithModifier(size.FillMax())),
					),
					column.WithModifier(size.Size(50, 50)),
					column.WithModifier(background.Background(color.NRGBA{R: 100, G: 0, B: 0, A: 200})),
					column.WithModifier(clickable.OnClick(func() {
						fmt.Println("Last Column clicked!!")
					})),
					column.WithModifier(clip.Clip(shape.ShapeCircle)),
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
		),
		column.WithModifier(size.FillMax()),
		column.WithModifier(background.Background(color.NRGBA{R: 200, G: 0, B: 0, A: 50})),

		column.WithAlignment(column.Middle),
	)(c)

	return c.Build()

}
