package main

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/material3/scaffold"
	"github.com/zodimo/go-compose/compose/foundation/material3/tab"
	"github.com/zodimo/go-compose/compose/foundation/material3/text"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/modifiers/weight"
	"github.com/zodimo/go-compose/pkg/api"
)

type Composable = api.Composable
type Composer = api.Composer

func UI() Composable {
	return func(c Composer) Composer {
		// State for selected tab
		selectedTabIndex := c.State("selectedTabIndex", func() any { return 0 }) // State logic to be added

		return scaffold.Scaffold(
			column.Column(
				c.Sequence(
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
