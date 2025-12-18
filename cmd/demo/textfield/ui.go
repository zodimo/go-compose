package main

import (
	"fmt"

	"gioui.org/unit"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	m3text "github.com/zodimo/go-compose/compose/foundation/material3/text"
	"github.com/zodimo/go-compose/compose/foundation/material3/textfield"
	"github.com/zodimo/go-compose/pkg/api"
)

func UI(c api.Composer) api.LayoutNode {
	filledText := c.State("filled_text", func() any { return "" })
	outlinedText := c.State("outlined_text", func() any { return "" })

	root := column.Column(
		func(c api.Composer) api.Composer {

			// Filled
			textfield.Filled(
				filledText.Get().(string),
				func(s string) { filledText.Set(s) },
				"Filled Text Field",
				textfield.WithSingleLine(true),
			)(c)

			spacer.Height(int(unit.Dp(16)))(c)

			// Outlined
			textfield.Outlined(
				outlinedText.Get().(string),
				func(s string) { outlinedText.Set(s) },
				"Outlined Text Field",
				textfield.WithSingleLine(true),
			)(c)

			spacer.Height(int(unit.Dp(16)))(c)

			// Display values
			m3text.Text(fmt.Sprintf("Filled value: %s", filledText.Get().(string)), m3text.TypestyleBodyLarge)(c)
			m3text.Text(fmt.Sprintf("Outlined value: %s", outlinedText.Get().(string)), m3text.TypestyleBodyLarge)(c)

			return c
		},
	)

	return root(c).Build()
}
