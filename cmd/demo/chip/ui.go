package main

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/material3/chip"
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
					func(c api.Composer) api.Composer {
						chip.AssistChip(func() { fmt.Println("Assist Chip Clicked") }, "Assist Chip")(c)
						spacer.Width(16)(c)
						chip.AssistChip(func() { fmt.Println("Assist with Icon") }, "With Icon",
							chip.WithLeadingIcon(func(c api.Composer) api.Composer {
								return text.TextWithStyle("★", text.TypestyleBodyMedium)(c)
							}),
						)(c)
						return c
					},
				),
				spacer.Height(16),
				// Filter Chip Row
				row.Row(
					func(c api.Composer) api.Composer {
						label := "Filter Chip"
						if selected {
							label = "Selected"
						}

						var leadingIcon api.Composable
						if selected {
							leadingIcon = func(c api.Composer) api.Composer { return text.TextWithStyle("✓", text.TypestyleBodyMedium)(c) }
						}

						chip.FilterChip(
							func() { selectedState.Set(!selected) },
							label,
							chip.WithSelected(selected),
							chip.WithLeadingIcon(leadingIcon),
						)(c)
						return c
					},
				),
				spacer.Height(16),
				// Input Chip
				row.Row(
					func(c api.Composer) api.Composer {
						chip.InputChip(func() { fmt.Println("Input Clicked") }, "Input Chip",
							chip.WithTrailingIcon(func(c api.Composer) api.Composer {
								return text.TextWithStyle("×", text.TypestyleBodyMedium)(c)
							}),
						)(c)
						return c
					},
				),
				spacer.Height(16),
				// Suggestion Chip
				row.Row(
					func(c api.Composer) api.Composer {
						chip.SuggestionChip(func() { fmt.Println("Suggestion 1") }, "Suggestion 1")(c)
						spacer.Width(8)(c)
						chip.SuggestionChip(func() { fmt.Println("Suggestion 2") }, "Suggestion 2")(c)
						return c
					},
				),
			),
		)(c)
	}
}
