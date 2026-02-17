package main

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/material3/chip"
	"github.com/zodimo/go-compose/compose/material3/icon"
	"github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/pkg/api"
)

func UI() api.Composable {
	return func(c api.Composer) api.Composer {
		// State for Filter chips
		selectedState := c.State("filter_selected", func() any { return false })
		selected := selectedState.Get().(bool)

		return column.Column(
			c.Sequence(
				// Title
				text.HeadlineMedium(
					"Chip Components",
				),
				spacer.Height(16),
				// Assist Chip Row
				row.Row(
					c.Sequence(
						chip.AssistChip(func() { fmt.Println("Assist Chip Clicked") }, "Assist Chip"),
						spacer.Width(16),
						chip.AssistChip(func() { fmt.Println("Assist with Icon") }, "With Icon",
							chip.WithLeadingIcon(icon.Icon(icon.SymbolStar)),
						),
					),
				),
				spacer.Height(16),
				// Filter Chip Row
				row.Row(
					c.Sequence(
						c.If(
							selected,
							chip.FilterChip(
								func() { selectedState.Set(!selected) },
								"Selected",
								chip.WithSelected(selected),
								chip.WithLeadingIcon(icon.Icon(icon.SymbolCheck)),
							),
							chip.FilterChip(
								func() { selectedState.Set(!selected) },
								"Filter Chip",
								chip.WithSelected(selected),
							),
						),
					),
				),
				spacer.Height(16),
				// Input Chip
				chip.InputChip(
					func() { fmt.Println("Input Clicked") },
					"Input Chip",
					chip.WithTrailingIcon(icon.Icon(icon.SymbolClose)),
				),
				spacer.Height(16),
				// Suggestion Chip
				row.Row(
					c.Sequence(
						chip.SuggestionChip(func() { fmt.Println("Suggestion 1") }, "Suggestion 1"),
						spacer.Width(8),
						chip.SuggestionChip(func() { fmt.Println("Suggestion 2") }, "Suggestion 2"),
					),
				),
			),
		)(c)
	}
}
