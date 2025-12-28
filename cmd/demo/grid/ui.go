package main

import (
	"fmt"
	"image/color"

	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/foundation/lazy"
	"github.com/zodimo/go-compose/compose/foundation/material3/text"
	ftext "github.com/zodimo/go-compose/compose/foundation/text"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/modifiers/background"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/theme"
)

func UI() api.Composable {
	return func(c api.Composer) api.Composer {
		state := lazy.RememberLazyGridState(c)

		return column.Column(
			c.Sequence(
				// Title for Fixed Grid
				text.Text("Fixed Grid (3 columns)", text.TypestyleHeadlineMedium,
					ftext.WithModifier(padding.All(16)),
				),

				// LazyVerticalGrid with Fixed(3) columns
				lazy.LazyVerticalGrid(
					lazy.Fixed(3),
					func(scope lazy.LazyGridScope) {
						scope.Items(12, nil, func(index int) api.Composable {
							return GridItem(index)
						})
					},
					lazy.WithGridModifier(size.Height(300)),
					lazy.WithGridState(state),
				),

				spacer.Height(16),

				// Title for Adaptive Grid
				text.Text("Adaptive Grid (min 100dp)", text.TypestyleHeadlineMedium,
					ftext.WithModifier(padding.All(16)),
				),

				// LazyVerticalGrid with Adaptive sizing
				lazy.LazyVerticalGrid(
					lazy.Adaptive(100),
					func(scope lazy.LazyGridScope) {
						scope.Items(15, nil, func(index int) api.Composable {
							return GridItem(index + 100) // Offset to distinguish
						})
					},
					lazy.WithGridModifier(size.FillMax()),
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
		text.Text(fmt.Sprintf("%d", index), text.TypestyleTitleLarge),
		box.WithModifier(
			size.Height(80).
				Then(size.FillMaxWidth()).
				Then(background.Background(theme.ColorHelper.SpecificColor(graphics.FromNRGBA(bgColor)))).
				Then(padding.All(8)),
		),
		box.WithAlignment(box.Center),
	)
}
