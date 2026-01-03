package main

import (
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/ui/graphics"

	"image/color"

	"github.com/zodimo/go-compose/compose/foundation/text"
	"github.com/zodimo/go-compose/modifiers/background"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/pkg/api"
)

func UI(c api.Composer) api.LayoutNode {

	c = column.Column(
		c.Sequence(
			text.Text("hello world",
				text.WithModifier(
					background.Background(
						graphics.FromNRGBA(color.NRGBA{R: 255, G: 0, B: 0, A: 150}),
					).
						Then(padding.All(20)).
						Then(background.Background(
							graphics.FromNRGBA(color.NRGBA{R: 0, G: 255, B: 0, A: 255}),
						)),
				),
			),
		),
		column.WithModifier(
			background.Background(
				graphics.FromNRGBA(color.NRGBA{R: 0, G: 0, B: 200, A: 255}),
			),
		),

		column.WithAlignment(column.Middle),
	)(c)

	return c.Build()

}
