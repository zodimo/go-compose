package main

import (
	"fmt"

	"gioui.org/unit"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/material3/text"
	m3text "github.com/zodimo/go-compose/compose/material3/text"
	m3TextField "github.com/zodimo/go-compose/compose/material3/textfield"
	"github.com/zodimo/go-compose/pkg/api"
)

func UI() api.Composable {
	return func(c api.Composer) api.Composer {
		filledText := c.State("filled_text", func() any { return "" })
		outlinedText := c.State("outlined_text", func() any { return "" })

		root := column.Column(
			c.Sequence(

				text.TitleLarge("Material3 Text Fields"),

				// Filled
				m3TextField.Filled(
					filledText.Get().(string),
					func(s string) { filledText.Set(s) },
					m3TextField.WithLabel("Filled Text Field"),
					m3TextField.WithSingleLine(true),
				),
				text.BodySmall(fmt.Sprintf("count: %d", len(filledText.Get().(string)))),

				spacer.Height(int(unit.Dp(16))),
				m3TextField.Filled(
					filledText.Get().(string),
					func(s string) { filledText.Set(s) },
					m3TextField.WithSingleLine(true),
				),
				spacer.Height(int(unit.Dp(16))),
				m3TextField.Filled(
					filledText.Get().(string),
					nil,
					m3TextField.WithSingleLine(true),
					m3TextField.WithPlaceholder("no label, not state"),
				),
				spacer.Height(int(unit.Dp(16))),
				m3TextField.Filled(
					filledText.Get().(string),
					func(_ string) {},
					m3TextField.WithSingleLine(true),
					m3TextField.WithLabel("Filled Text Field with noop onchange"),
				),

				spacer.Height(int(unit.Dp(16))),
				// Outlined
				m3TextField.Outlined(
					outlinedText.Get().(string),
					func(s string) { outlinedText.Set(s) },
					m3TextField.WithLabel("Outlined Text Field"),
					m3TextField.WithSingleLine(true),
				),
				text.BodySmall(fmt.Sprintf("count: %d", len(outlinedText.Get().(string)))),

				spacer.Height(int(unit.Dp(16))),
				m3TextField.Outlined(
					outlinedText.Get().(string),
					func(s string) { outlinedText.Set(s) },
					m3TextField.WithSingleLine(true),
				),
				spacer.Height(int(unit.Dp(16))),
				m3TextField.Outlined(
					outlinedText.Get().(string),
					nil,
					m3TextField.WithSingleLine(true),
					m3TextField.WithPlaceholder("no label, not state"),
				),
				spacer.Height(int(unit.Dp(16))),
				m3TextField.Outlined(
					outlinedText.Get().(string),
					func(_ string) {},
					m3TextField.WithSingleLine(true),
					m3TextField.WithLabel("Outlined Text Field with noop onchange"),
				),

				spacer.Height(int(unit.Dp(16))),
				// Display values
				m3text.TextWithStyle(fmt.Sprintf("Filled value: %s", filledText.Get().(string)), m3text.TypestyleBodyLarge),
				m3text.TextWithStyle(fmt.Sprintf("Outlined value: %s", outlinedText.Get().(string)), m3text.TypestyleBodyLarge),
			),
		)

		return root(c)
	}
}
