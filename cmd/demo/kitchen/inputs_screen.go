package main

import (
	"gioui.org/layout"

	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/material3/card"
	m3text "github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/compose/material3/textfield"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/state"
)

// InputsScreen shows text fields and typography
func InputsScreen(c api.Composer) api.Composable {
	filledText := c.State("inp_text", func() any { return "" })
	outlinedText := c.State("outlined_text", func() any { return "" })

	_ = state.MustRemember(c, "input-viewmodel", func() *InputViewModel {
		return NewInputViewModel()
	})

	return func(c api.Composer) api.Composer {
		return column.Column(
			c.Sequence(
				SectionTitle("Text Field Filled"),
				spacer.Height(8),

				textfield.Filled(
					filledText.Get().(string),
					func(s string) { filledText.Set(s) },
					textfield.WithSingleLine(true),
				),

				spacer.Height(24),

				textfield.Filled(
					filledText.Get().(string),
					func(s string) { filledText.Set(s) },
					textfield.WithLabel("Filled Text Field"),
					textfield.WithSingleLine(true),
				),

				spacer.Height(24),

				textfield.Filled(
					filledText.Get().(string),
					func(s string) { filledText.Set(s) },
					textfield.WithLabel("Filled Text disabled"),
					textfield.WithSingleLine(true),
					textfield.WithEnabled(false),
				),

				spacer.Height(24),

				SectionTitle("Text Field Outlined"),
				spacer.Height(8),
				textfield.Outlined(
					outlinedText.Get().(string),
					func(s string) { outlinedText.Set(s) },
					textfield.WithSingleLine(true),
				),
				spacer.Height(24),
				textfield.Outlined(
					outlinedText.Get().(string),
					func(s string) { outlinedText.Set(s) },
					textfield.WithLabel("Outlined Text Field"),
					textfield.WithSingleLine(true),
				),
				spacer.Height(24),
				textfield.Outlined(
					outlinedText.Get().(string),
					func(s string) { outlinedText.Set(s) },
					textfield.WithLabel("Outlined Text disabled"),
					textfield.WithSingleLine(true),
					textfield.WithEnabled(false),
				),
				spacer.Height(24),
				SectionTitle("Card"),
				spacer.Height(8),
				card.Filled(
					card.CardContents(
						card.Content(
							box.Box(
								func(c api.Composer) api.Composer {
									column.Column(c.Sequence(
										m3text.TextWithStyle("Card Title", m3text.TypestyleTitleMedium),
										spacer.Height(8),
										m3text.TextWithStyle("This is card content demonstrating the Card surface component.", m3text.TypestyleBodyMedium),
									))(c)
									return c
								},
								box.WithModifier(padding.All(16)),
								box.WithAlignment(layout.NW),
							),
						),
					),
					card.WithModifier(size.Width(300)),
				),
			),
			column.WithModifier(padding.All(16)),
		)(c)
	}
}
