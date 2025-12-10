package main

import (
	"go-compose-dev/compose"
	"go-compose-dev/compose/foundation/layout/column"
	"go-compose-dev/compose/foundation/layout/row"
	"go-compose-dev/internal/modifiers/background"
	"go-compose-dev/internal/modifiers/padding"
	"go-compose-dev/internal/modifiers/size"
	"go-compose-dev/internal/modifiers/weight"
	"go-compose-dev/pkg/api"
	"image/color"
)

func UI(c api.Composer) api.LayoutNode {

	c = column.Column(
		compose.Sequence(
			row.Row(compose.Sequence(
				column.Column(
					compose.Sequence(),
					// column.WithModifier(size.Size(200, 200, size.SizeRequired())),
					column.WithModifier(weight.Weight(1)),
					column.WithModifier(background.Background(color.NRGBA{R: 0, G: 0, B: 200, A: 200})),
					column.WithModifier(padding.PaddingAll(20)),
					column.WithModifier(background.Background(color.NRGBA{R: 0, G: 100, B: 200, A: 200})),
				),
				column.Column(
					compose.Sequence(),
					column.WithModifier(weight.Weight(1)),
					column.WithModifier(background.Background(color.NRGBA{R: 150, G: 0, B: 0, A: 200})),
				),
				column.Column(
					compose.Sequence(),
					column.WithModifier(size.Size(50, 50)),
					column.WithModifier(background.Background(color.NRGBA{R: 100, G: 0, B: 0, A: 200})),
				),
			),
				row.WithModifier(size.Size(500, 300)),
				row.WithModifier(background.Background(color.NRGBA{R: 0, G: 200, B: 0, A: 200})),
			),
		),
		column.WithModifier(size.FillMax()),
		column.WithModifier(background.Background(color.NRGBA{R: 200, G: 0, B: 0, A: 50})),
	)(c)

	return c.Build()

}
