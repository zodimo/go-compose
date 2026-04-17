package main

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/lazy"
	ftext "github.com/zodimo/go-compose/compose/foundation/text"
	"github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/modifiers/weight"
	"github.com/zodimo/go-compose/pkg/api"
)

func UI() api.Composable {
	return func(c api.Composer) api.Composer {
		return column.Column(
			c.Sequence(
				text.HeadlineMedium("With Scrollbars:", ftext.WithModifier(padding.All(16))),
				box.Box(
					func(c api.Composer) api.Composer {
						columnState := lazy.RememberLazyListState(c)
						rowState := lazy.RememberLazyListState(c)

						return column.Column(
							c.Sequence(
								text.HeadlineMedium("Horizontal Row:", ftext.WithModifier(padding.All(16))),
								lazy.LazyRow(
									func(scope lazy.LazyListScope) {
										scope.Items(50, nil, func(index int) api.Composable {
											return row.Row(
												func(c api.Composer) api.Composer {
													text.BodyLarge(fmt.Sprintf("H%d", index), ftext.WithModifier(padding.All(16)))(c)
													return c
												},
											)
										})
									},
									lazy.WithModifier(size.Height(80)),
									lazy.WithState(rowState),
								),
								text.HeadlineMedium("Lazy List Demo", ftext.WithModifier(padding.All(16))),
								lazy.LazyColumn(
									func(scope lazy.LazyListScope) {
										// 100 items
										scope.Items(100, nil, func(index int) api.Composable {
											return func(c api.Composer) api.Composer {
												text.BodyLarge(fmt.Sprintf("Item %d", index), ftext.WithModifier(padding.All(8)))(c)
												return c
											}
										})
									},
									lazy.WithModifier(size.FillMax()),
									lazy.WithState(columnState),
								),
							),
							column.WithModifier(size.FillMax()),
						)(c)
					},
					box.WithModifier(weight.Weight(1)),
				),
				text.HeadlineMedium("Without Scrollbars:", ftext.WithModifier(padding.All(16))),
				box.Box(
					func(c api.Composer) api.Composer {
						columnState := lazy.RememberLazyListState(c)
						rowState := lazy.RememberLazyListState(c)

						return column.Column(
							c.Sequence(
								text.HeadlineMedium("Horizontal Row:", ftext.WithModifier(padding.All(16))),
								lazy.LazyRow(
									func(scope lazy.LazyListScope) {
										scope.Items(50, nil, func(index int) api.Composable {
											return row.Row(
												func(c api.Composer) api.Composer {
													text.BodyLarge(fmt.Sprintf("H%d", index), ftext.WithModifier(padding.All(16)))(c)
													return c
												},
											)
										})
									},
									lazy.WithModifier(size.Height(80)),
									lazy.WithState(rowState),
									lazy.WithScrollbar(false),
								),
								text.HeadlineMedium("Lazy List Demo", ftext.WithModifier(padding.All(16))),
								lazy.LazyColumn(
									func(scope lazy.LazyListScope) {
										// 100 items
										scope.Items(100, nil, func(index int) api.Composable {
											return func(c api.Composer) api.Composer {
												text.BodyLarge(fmt.Sprintf("Item %d", index), ftext.WithModifier(padding.All(8)))(c)
												return c
											}
										})
									},
									lazy.WithModifier(size.FillMax()),
									lazy.WithState(columnState),
									lazy.WithScrollbar(false),
								),
							),
							column.WithModifier(size.FillMax()),
						)(c)
					},
					box.WithModifier(weight.Weight(1)),
				),
			),
			column.WithModifier(size.FillMax()),
		)(c)
	}
}
