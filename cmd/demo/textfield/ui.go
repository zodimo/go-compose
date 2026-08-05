package main

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/foundation/next/text/input"
	foundationTextField "github.com/zodimo/go-compose/compose/foundation/next/textfield"
	"github.com/zodimo/go-compose/compose/material3/text"
	m3text "github.com/zodimo/go-compose/compose/material3/text"
	m3TextField "github.com/zodimo/go-compose/compose/material3/textfield"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/modifiers/background"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"
)

func UI() api.Composable {
	return func(c api.Composer) api.Composer {
		filledText := c.State("filled_text", func() any { return "" })
		outlinedText := c.State("outlined_text", func() any { return "" })
		multilineStateVal := c.State("multiline_state", func() any {
			return input.NewTextFieldState("Line 1\nLine 2\nLine 3")
		})
		multilineState := multilineStateVal.Get().(*input.TextFieldState)

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

				spacer.Height(16),
				m3TextField.Filled(
					filledText.Get().(string),
					func(s string) { filledText.Set(s) },
					m3TextField.WithSingleLine(true),
				),
				spacer.Height(16),
				m3TextField.Filled(
					filledText.Get().(string),
					nil,
					m3TextField.WithSingleLine(true),
					m3TextField.WithPlaceholder("no label, not state"),
				),
				spacer.Height(16),
				m3TextField.Filled(
					filledText.Get().(string),
					func(_ string) {},
					m3TextField.WithSingleLine(true),
					m3TextField.WithLabel("Filled Text Field with noop onchange"),
				),

				spacer.Height(16),
				// Outlined
				m3TextField.Outlined(
					outlinedText.Get().(string),
					func(s string) { outlinedText.Set(s) },
					m3TextField.WithLabel("Outlined Text Field"),
					m3TextField.WithSingleLine(true),
				),
				text.BodySmall(fmt.Sprintf("count: %d", len(outlinedText.Get().(string)))),

				spacer.Height(16),
				m3TextField.Outlined(
					outlinedText.Get().(string),
					func(s string) { outlinedText.Set(s) },
					m3TextField.WithSingleLine(true),
				),
				spacer.Height(16),
				m3TextField.Outlined(
					outlinedText.Get().(string),
					nil,
					m3TextField.WithSingleLine(true),
					m3TextField.WithPlaceholder("no label, not state"),
				),
				spacer.Height(16),
				m3TextField.Outlined(
					outlinedText.Get().(string),
					func(_ string) {},
					m3TextField.WithSingleLine(true),
					m3TextField.WithLabel("Outlined Text Field with noop onchange"),
				),

				spacer.Height(16),
				// Display values
				m3text.BodyLarge(fmt.Sprintf("Filled value: %s", filledText.Get().(string))),
				m3text.BodyLarge(fmt.Sprintf("Outlined value: %s", outlinedText.Get().(string))),

				spacer.Height(24),
				// foundation/next
				text.TitleMedium("foundation/next"),
				foundationTextField.BasicTextField(
					multilineState,
					func(value string) {
						multilineState.SetTextAndPlaceCursorAtEnd(value)
					},
					foundationTextField.WithLineLimits(input.NewMultiLineWithLimits(3, 6)),
					foundationTextField.WithModifier(
						background.Background(graphics.ColorLightGray).
							Then(size.FillMaxWidth()).
							Then(padding.All(16)),
					),
				),
				text.BodySmall("Multiline: grows from 3 lines, scrolls after 6"),
			),
		)

		return root(c)
	}
}
