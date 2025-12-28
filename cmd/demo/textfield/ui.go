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
		c.Sequence(
			// Filled
			textfield.Filled(
				filledText.Get().(string),
				func(s string) { filledText.Set(s) },
				textfield.WithLabel("Filled Text Field"),
				textfield.WithSingleLine(true),
			),
			spacer.Height(int(unit.Dp(16))),
			// Outlined
			textfield.Outlined(
				outlinedText.Get().(string),
				func(s string) { outlinedText.Set(s) },
				textfield.WithLabel("Outlined Text Field"),
				textfield.WithSingleLine(true),
			),
			spacer.Height(int(unit.Dp(16))),
			// Display values
			m3text.Text(fmt.Sprintf("Filled value: %s", filledText.Get().(string)), m3text.TypestyleBodyLarge),
			m3text.Text(fmt.Sprintf("Outlined value: %s", outlinedText.Get().(string)), m3text.TypestyleBodyLarge),
		),
	)

	return root(c).Build()
}
