package main

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/text"
	m3Button "github.com/zodimo/go-compose/compose/material3/button"
	m3Divider "github.com/zodimo/go-compose/compose/material3/divider"
	m3Text "github.com/zodimo/go-compose/compose/material3/text"

	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"
)

func UI() api.Composable {
	return func(c api.Composer) api.Composer {
		// State for toggle
		showDetails := c.State("show_details", func() any { return false })

		// State for counter (If/Else condition)
		count := c.State("count", func() any { return 0 })

		c = column.Column(
			c.Sequence(
				m3Text.TextWithStyle("Control Flow Demo", m3Text.TypestyleHeadlineMedium),
				m3Divider.Divider(m3Divider.WithModifier(padding.Vertical(16, 16))),

				// Test 'If'
				m3Text.TextWithStyle("1. If/Else (Click to toggle)", m3Text.TypestyleTitleMedium),
				m3Button.Filled(func() {
					showDetails.Set(!showDetails.Get().(bool))
				}, "Toggle Details"),

				c.If(showDetails.Get().(bool),
					m3Text.TextWithStyle("Details are SHOWN! This block is visible because condition is true.", m3Text.TypestyleBodyMedium),
					m3Text.TextWithStyle("Details are HIDDEN. This block is visible because condition is false.", m3Text.TypestyleBodyMedium),
				),

				m3Divider.Divider(m3Divider.WithModifier(padding.Vertical(16, 16))),

				// Test 'When'
				m3Text.TextWithStyle("2. When (Visible only when count > 5)", m3Text.TypestyleTitleMedium),
				row.Row(c.Sequence(
					m3Button.Outlined(func() {
						count.Set(count.Get().(int) - 1)
					}, "-"),
					m3Text.TextWithStyle(fmt.Sprintf("Count: %d", count.Get().(int)), m3Text.TypestyleDefault, text.WithModifier(padding.Horizontal(16, 16))),
					m3Button.Outlined(func() {
						count.Set(count.Get().(int) + 1)
					}, "+"),
				), row.WithAlignment(row.Middle)),

				c.When(count.Get().(int) > 5,
					m3Text.TextWithStyle("Count is greater than 5! (This text appears via 'When')", m3Text.TypestyleBodyMedium),
				),

				m3Divider.Divider(m3Divider.WithModifier(padding.Vertical(16, 16))),

				// Test 'Range'
				m3Text.TextWithStyle(fmt.Sprintf("3. Range (Loop %d times)", count.Get().(int)), m3Text.TypestyleTitleMedium),
				c.Range(count.Get().(int), func(i int) api.Composable {
					return m3Text.TextWithStyle(fmt.Sprintf("Item #%d", i), m3Text.TypestyleBodyMedium)
				}),

				m3Divider.Divider(m3Divider.WithModifier(padding.Vertical(16, 16))),

				// Test 'Key'
				m3Text.TextWithStyle("4. Key (Stable Identity)", m3Text.TypestyleTitleMedium),
				c.Key("my-stable-block",
					m3Text.TextWithStyle("This block has a stable key 'my-stable-block'", m3Text.TypestyleBodyMedium),
				),
			),
			column.WithModifier(size.FillMax().Then(padding.All(24))),
		)(c)

		return c
	}
}
