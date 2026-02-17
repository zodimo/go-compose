package main

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/lazy"
	ftext "github.com/zodimo/go-compose/compose/foundation/text"
	"github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/modifiers/background"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"
)

func UI() api.Composable {
	return func(c api.Composer) api.Composer {
		state := lazy.RememberLazyListState(c)

		return column.Column(
			func(c api.Composer) api.Composer {
				lazy.LazyColumn(
					func(scope lazy.LazyListScope) {
						for i := 0; i < 5; i++ {
							// Sticky Header
							headerIndex := i
							scope.StickyHeader(fmt.Sprintf("header-%d", i), func(c api.Composer) api.Composer {
								text.HeadlineSmall(
									fmt.Sprintf("Sticky Header %d", headerIndex),
									ftext.WithModifier(
										background.Background(graphics.NewColorSrgb(200, 200, 200, 255)).
											Then(size.FillMaxWidth()),
									),
								)(c)
								return c
							})

							// Items
							scope.Items(15, nil, func(index int) api.Composable {
								return func(c api.Composer) api.Composer {
									text.BodyMedium(
										fmt.Sprintf("Item %d - %d", headerIndex, index),
										ftext.WithModifier(padding.All(8)),
									)(c)
									return c
								}
							})
						}
					},
					lazy.WithModifier(size.FillMax()),
					lazy.WithState(state),
				)(c)
				return c
			},
			column.WithModifier(size.FillMax()),
		)(c)
	}
}
