package main

import (
	"fmt"
	"image/color"

	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	ftext "github.com/zodimo/go-compose/compose/foundation/text"
	"github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/modifiers/background"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/pkg/x/m3table"
)

func UI() api.Composable {
	return func(c api.Composer) api.Composer {

		return column.Column(
			c.Sequence(
				// Title for Fixed Grid
				text.HeadlineMedium("Table",
					ftext.WithModifier(padding.All(16)),
				),

				m3table.Table(
					[]m3table.Column{
						{
							Header: text.BodyLarge("Header 1"),
							Width:  100,
						},
						{
							Header: text.BodyLarge("Header 2"),
							Width:  200,
						},
						{
							Header: text.BodyLarge("Header 3"),
							Width:  300,
						},
					},
					10,
					func(row, col int) api.Composable {
						return text.BodyLarge(fmt.Sprintf("Cell %d,%d", row, col))
					},
				),
			),
			column.WithModifier(size.FillMax()),
		)(c)
	}
}

// GridItem creates a single grid item with colored background
func GridItem(index int) api.Composable {
	// Simple alternating colors for visual distinction
	colors := []color.NRGBA{
		{R: 234, G: 221, B: 255, A: 255}, // Primary container
		{R: 232, G: 222, B: 248, A: 255}, // Secondary container
		{R: 255, G: 216, B: 228, A: 255}, // Tertiary container
	}
	bgColor := colors[index%len(colors)]

	return box.Box(
		text.TitleLarge(fmt.Sprintf("%d", index)),
		box.WithModifier(
			size.Height(80).
				Then(size.FillMaxWidth()).
				Then(background.Background(graphics.FromNRGBA(bgColor))).
				Then(padding.All(8)),
		),
		box.WithAlignment(box.Center),
	)
}
