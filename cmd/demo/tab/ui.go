package main

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/material3/scaffold"
	"github.com/zodimo/go-compose/compose/material3/tab"
	"github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/modifiers/weight"
	"github.com/zodimo/go-compose/pkg/api"
)

type Composable = api.Composable
type Composer = api.Composer

func UI() Composable {
	return func(c Composer) Composer {
		// State for selected tab
		selectedTabIndex := c.State("selectedTabIndex", func() any { return 0 })
		selectedTabIndex2 := c.State("selectedTabIndex2", func() any { return 0 })

		return scaffold.Scaffold(
			column.Column(
				c.Sequence(
					// Tab Row 1: Default theme colors
					text.TextWithStyle("Default Theme Colors", text.TypestyleLabelLarge),
					tab.TabRow(
						selectedTabIndex.Get().(int),
						func(c Composer) Composer {
							titles := []string{"Tab 1", "Tab 2", "Tab 3"}
							for i, title := range titles {
								index := i // Capture loop variable
								tab.Tab(
									selectedTabIndex.Get().(int) == index,
									func() {
										fmt.Printf("Clicked Tab %d\n", index)
										selectedTabIndex.Set(index)
									},
									text.TextWithStyle(title, text.TypestyleLabelMedium),
									tab.WithModifier(weight.Weight(1)),
								)(c)
							}
							return c
						},
					),

					// Spacing
					spacer.Height(24),

					// Tab Row 2: Custom theme colors
					text.TextWithStyle("Custom Theme Colors (Secondary)", text.TypestyleLabelLarge),
					tab.TabRow(
						selectedTabIndex2.Get().(int),
						func(c Composer) Composer {
							titles := []string{"Home", "Profile", "Settings"}
							for i, title := range titles {
								index := i // Capture loop variable
								tab.Tab(
									selectedTabIndex2.Get().(int) == index,
									func() {
										fmt.Printf("Clicked Tab %d\n", index)
										selectedTabIndex2.Set(index)
									},
									text.TextWithStyle(title, text.TypestyleLabelMedium),
									tab.WithModifier(weight.Weight(1)),
								)(c)
							}
							return c
						},
						tab.WithTabRowModifier(size.FillMaxWidth()),
					),

					// Spacing
					spacer.Height(24),

					// Explanatory text
					text.TextWithStyle(
						"Tab colors now use theme-aware styling. "+
							"They automatically adapt to light/dark themes.",
						text.TypestyleBodyMedium,
					),
				),
				column.WithModifier(
					size.FillMaxWidth().
						Then(size.FillMaxHeight()).
						Then(padding.All(16)),
				),
			),
		)(c)
	}
}
