package main

import (
	"go-compose-dev/compose/foundation/layout/column"
	"go-compose-dev/compose/foundation/layout/row"
	"go-compose-dev/compose/foundation/layout/spacer"
	"go-compose-dev/compose/foundation/material3/badge"
	"go-compose-dev/compose/foundation/material3/text"
	"go-compose-dev/pkg/api"
	"image/color"
)

func UI() api.Composable {
	return func(c api.Composer) api.Composer {
		return column.Column(
			func(c api.Composer) api.Composer {
				text.Text("Badge Components", text.TypestyleHeadlineMedium)(c)
				spacer.SpacerHeight(16)(c)

				// Small badge (dot)
				text.Text("Small Badge (Dot)", text.TypestyleTitleSmall)(c)
				spacer.SpacerHeight(8)(c)
				badge.Badge()(c)

				spacer.SpacerHeight(16)(c)

				// Large badge with number
				text.Text("Large Badge (Number)", text.TypestyleTitleSmall)(c)
				spacer.SpacerHeight(8)(c)
				badge.Badge(badge.WithText("999+"))(c)

				spacer.SpacerHeight(16)(c)

				// BadgedBox examples
				text.Text("BadgedBox Examples", text.TypestyleTitleSmall)(c)
				spacer.SpacerHeight(8)(c)

				row.Row(func(c api.Composer) api.Composer {
					// Dot badge on text
					badge.BadgedBox(
						badge.Badge(),
						text.Text("🔔", text.TypestyleHeadlineMedium),
					)(c)

					spacer.SpacerWidth(24)(c)

					// Number badge on icon
					badge.BadgedBox(
						badge.Badge(badge.WithText("5")),
						text.Text("✉️", text.TypestyleHeadlineMedium),
					)(c)

					spacer.SpacerWidth(24)(c)

					// Custom color badge
					badge.BadgedBox(
						badge.Badge(
							badge.WithText("!"),
							badge.WithContainerColor(color.NRGBA{R: 0, G: 200, B: 0, A: 255}),
						),
						text.Text("📦", text.TypestyleHeadlineMedium),
					)(c)

					return c
				})(c)

				return c
			},
		)(c)
	}
}
