package main

import (
	"fmt"

	"go-compose-dev/compose"
	"go-compose-dev/compose/foundation/layout/column"
	"go-compose-dev/compose/foundation/material3/scaffold"
	"go-compose-dev/compose/foundation/material3/tab"
	"go-compose-dev/compose/foundation/material3/text"
	"go-compose-dev/internal/modifiers/size"
	"go-compose-dev/internal/modifiers/weight"
	"go-compose-dev/pkg/api"
)

type Composable = api.Composable
type Composer = api.Composer

func UI() Composable {
	return func(c Composer) Composer {
		// State for selected tab
		selectedTabIndex := c.State("selectedTabIndex", func() any { return 0 }) // State logic to be added

		return scaffold.Scaffold(
			column.Column(
				compose.Sequence(
					// Tab Row 1: Text only
					text.Text("Primary TabRow", text.TypestyleLabelLarge),
					tab.TabRow(
						selectedTabIndex.Get().(int),
						func(c Composer) Composer {
							titles := []string{"Tab 1", "Tab 2", "Tab 3"}
							for i, title := range titles {
								index := i // Capture loop variable
								tab.Tab(
									selectedTabIndex.Get().(int) == index,
									func() {
										// State update placeholder
										fmt.Printf("Clicked Tab %d\n", index)
										selectedTabIndex.Set(index)
									},
									text.Text(title, text.TypestyleLabelMedium),
									tab.WithModifier(weight.Weight(1)),
								)(c)
							}
							return c
						},
					),
				),
				column.WithModifier(size.FillMaxWidth().Then(size.FillMaxHeight())),
			),
		)(c)
	}
}
