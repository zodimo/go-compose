package main

import (
	"log"

	"go-compose-dev/compose/foundation/layout/column"
	"go-compose-dev/compose/foundation/layout/row"
	"go-compose-dev/compose/foundation/layout/spacer"
	"go-compose-dev/compose/foundation/material3/button"
	"go-compose-dev/compose/foundation/material3/text"
	"go-compose-dev/compose/foundation/material3/tooltip"
	"go-compose-dev/pkg/api"
)

func UI() api.Composable {
	return func(c api.Composer) api.Composer {
		return column.Column(
			func(c api.Composer) api.Composer {
				text.Text("Tooltip Demo", text.TypestyleHeadlineMedium)(c)
				spacer.SpacerHeight(24)(c)

				// Example 1: Tooltip on a button
				text.Text("Button with Tooltip", text.TypestyleTitleSmall)(c)
				spacer.SpacerHeight(8)(c)
				tooltip.Tooltip(
					"Click to submit",
					button.Filled(func() { log.Println("Clicked") }, "Submit"),
				)(c)

				spacer.SpacerHeight(24)(c)

				// Example 2: Tooltip on an icon
				text.Text("Icon with Tooltip", text.TypestyleTitleSmall)(c)
				spacer.SpacerHeight(8)(c)
				tooltip.Tooltip(
					"Notifications",
					text.Text("🔔", text.TypestyleHeadlineMedium),
				)(c)

				spacer.SpacerHeight(24)(c)

				// Example 3: Multiple tooltips in a row
				text.Text("Multiple Tooltips", text.TypestyleTitleSmall)(c)
				spacer.SpacerHeight(8)(c)
				row.Row(func(c api.Composer) api.Composer {
					tooltip.Tooltip(
						"Home",
						text.Text("🏠", text.TypestyleHeadlineMedium),
					)(c)
					spacer.SpacerWidth(16)(c)

					tooltip.Tooltip(
						"Settings",
						text.Text("⚙️", text.TypestyleHeadlineMedium),
					)(c)
					spacer.SpacerWidth(16)(c)

					tooltip.Tooltip(
						"Profile",
						text.Text("👤", text.TypestyleHeadlineMedium),
					)(c)

					return c
				})(c)

				return c
			},
		)(c)
	}
}
