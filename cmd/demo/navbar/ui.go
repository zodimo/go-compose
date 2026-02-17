package main

import (
	"fmt"

	"gioui.org/layout"

	"github.com/zodimo/go-compose/compose/foundation/icon"
	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/material3/navigationbar"
	m3text "github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/modifiers/weight"
	"github.com/zodimo/go-compose/pkg/api"

	mdicons "golang.org/x/exp/shiny/materialdesign/icons"
)

func UI() api.Composable {
	return func(c api.Composer) api.Composer {
		// State for selected item
		// We use a pointer to int for the state value
		selectedIdxVal := c.State("nav_selected_index", func() any { return 0 })
		selectedIdx := selectedIdxVal.Get().(int)

		items := []struct {
			Label string
			Icon  []byte
		}{
			{"Home", mdicons.ActionHome},
			{"Search", mdicons.ActionSearch},
			{"Settings", mdicons.ActionSettings},
		}

		return column.Column(
			func(c api.Composer) api.Composer {
				// Content Area
				box.Box(
					func(c api.Composer) api.Composer {
						return m3text.DisplayMedium(fmt.Sprintf("Selected: %s", items[selectedIdx].Label))(c)
					},
					box.WithModifier(weight.Weight(1)),
					box.WithAlignment(layout.Center),
				)(c)

				// Navigation Bar
				navigationbar.NavigationBar(
					func(c api.Composer) api.Composer {
						for i, item := range items {
							idx := i // capture loop variable
							navigationbar.NavigationBarItem(
								selectedIdx == idx,
								func() {
									selectedIdxVal.Set(idx)
								},
								func(c api.Composer) api.Composer {
									return icon.Icon(item.Icon)(c)
								},
								func(c api.Composer) api.Composer {
									return m3text.LabelMedium(item.Label)(c)
								},
							)(c)
						}
						return c
					},
				)(c)
				return c
			},
			column.WithModifier(size.FillMax()),
		)(c)
	}
}
