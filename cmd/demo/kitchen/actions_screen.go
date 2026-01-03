package main

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/material3/button"
	"github.com/zodimo/go-compose/compose/material3/iconbutton"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/pkg/api"

	mdicons "golang.org/x/exp/shiny/materialdesign/icons"
)

// ActionsScreen shows buttons, FAB, and icon buttons
func ActionsScreen(c api.Composer) api.Composable {
	return func(c api.Composer) api.Composer {
		return column.Column(
			c.Sequence(
				SectionTitle("Buttons"),
				spacer.Height(8),
				row.Row(c.Sequence(
					button.Filled(func() { fmt.Println("Filled clicked") }, "Filled"),
					spacer.Width(8),
					button.Filled(func() { fmt.Println("Filled clicked") }, "Filled", button.WithEnabled(false)),
					spacer.Width(8),
					button.Outlined(func() { fmt.Println("Outlined clicked") }, "Outlined"),
					spacer.Width(8),
					button.Outlined(func() { fmt.Println("Outlined clicked") }, "Outlined", button.WithEnabled(false)),
				)),
				spacer.Height(8),
				row.Row(c.Sequence(
					button.Text(func() { fmt.Println("Text clicked") }, "Text"),
					spacer.Width(8),
					button.Elevated(func() { fmt.Println("Elevated clicked") }, "Elevated"),
				)),

				spacer.Height(24),
				// SectionTitle("Floating Action Button"),
				// FABs temporarily disabled due to runtime panic in border modifier
				/*
					spacer.Height(8),
					row.Row(c.Sequence(
						floatingactionbutton.FloatingActionButton(
							func() { fmt.Println("FAB clicked") },
							icon.Icon(mdicons.ContentAdd),
						),
						spacer.Width(16),
						floatingactionbutton.FloatingActionButton(
							func() { fmt.Println("Small FAB clicked") },
							icon.Icon(mdicons.ContentAdd),
							floatingactionbutton.WithSize(floatingactionbutton.FabSizeSmall),
						),
					), row.WithAlignment(row.Middle)),
				*/

				// spacer.Height(24),
				SectionTitle("Icon Buttons"),
				spacer.Height(8),
				row.Row(c.Sequence(
					iconbutton.Standard(func() {}, mdicons.ActionFavorite, "Favorite"),
					spacer.Width(8),
					iconbutton.Standard(func() {}, mdicons.ActionSearch, "Search"),
					spacer.Width(8),
					iconbutton.Standard(func() {}, mdicons.ActionSettings, "Settings"),
				)),
			),
			column.WithModifier(padding.All(16)),
		)(c)
	}
}
