package main

import (
	"log"

	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/material3/button"
	"github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/compose/material3/tooltip"
	"github.com/zodimo/go-compose/pkg/api"
)

func UI() api.Composable {
	return func(c api.Composer) api.Composer {
		return column.Column(
			func(c api.Composer) api.Composer {
				text.HeadlineMedium("Tooltip Demo")(c)
				spacer.Height(24)(c)

				// Example 1: Tooltip on a button
				text.TitleSmall("Button with Tooltip")(c)
				spacer.Height(8)(c)
				tooltip.Tooltip(
					"Click to submit",
					button.Filled(func() { log.Println("Clicked") }, "Submit"),
				)(c)

				spacer.Height(24)(c)

				// Example 2: Tooltip on an icon
				text.TitleSmall("Icon with Tooltip")(c)
				spacer.Height(8)(c)
				tooltip.Tooltip(
					"Notifications",
					text.HeadlineMedium("🔔"),
				)(c)

				spacer.Height(24)(c)

				// Example 3: Multiple tooltips in a row
				text.TitleSmall("Multiple Tooltips")(c)
				spacer.Height(8)(c)
				row.Row(func(c api.Composer) api.Composer {
					tooltip.Tooltip(
						"Home",
						text.HeadlineMedium("🏠"),
					)(c)
					spacer.Width(16)(c)

					tooltip.Tooltip(
						"Settings",
						text.HeadlineMedium("⚙️"),
					)(c)
					spacer.Width(16)(c)

					tooltip.Tooltip(
						"Profile",
						text.HeadlineMedium("👤"),
					)(c)

					return c
				})(c)

				return c
			},
		)(c)
	}
}
