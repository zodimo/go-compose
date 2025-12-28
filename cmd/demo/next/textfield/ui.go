package main

import (
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	foundationText "github.com/zodimo/go-compose/compose/foundation/next/text"
	"github.com/zodimo/go-compose/compose/foundation/next/text/input"
	"github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/modifiers/background"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/theme"
)

// Create text field states at package level for persistence across frames
var (
	basicState      = input.NewTextFieldState("Hello, World!")
	singleLineState = input.NewTextFieldState("Type here...")
	maxLengthState  = input.NewTextFieldState("")
	digitsOnlyState = input.NewTextFieldState("")
)

func UI() api.Composable {
	return func(c api.Composer) api.Composer {

		modifier := background.Background(theme.ColorHelper.SpecificColor(graphics.ColorLightGray)).
			Then(size.FillMaxWidth()).
			Then(padding.All(16))

		return column.Column(
			c.Sequence(
				// Title
				text.Text("BasicTextField Demo (Next)", text.TypestyleHeadlineMedium),
				spacer.Height(24),

				// Section: Basic Text Field
				text.Text("Basic TextField", text.TypestyleTitleMedium),
				spacer.Height(8),
				foundationText.BasicTextField(
					basicState,
					func(value string) {
						basicState.SetTextAndPlaceCursorAtEnd(value)
					},
					foundationText.WithTextFieldModifier(modifier),
				),
				spacer.Height(16),

				// Section: Single Line
				text.Text("Single Line TextField", text.TypestyleTitleMedium),
				spacer.Height(8),
				foundationText.BasicTextField(
					singleLineState,
					func(value string) {
						singleLineState.SetTextAndPlaceCursorAtEnd(value)
					},
					foundationText.WithLineLimits(input.TextFieldLineLimitsSingleLine),
					foundationText.WithTextFieldModifier(modifier),
				),
				spacer.Height(16),

				// Section: Max Length (10 chars)
				text.Text("Max Length (10 chars)", text.TypestyleTitleMedium),
				spacer.Height(8),
				foundationText.BasicTextField(
					maxLengthState,
					func(value string) {
						maxLengthState.SetTextAndPlaceCursorAtEnd(value)
					},
					foundationText.WithInputTransformation(input.MaxLengthTransformation(10)),
					foundationText.WithTextFieldModifier(modifier),
				),
				spacer.Height(16),

				// Section: Digits Only
				text.Text("Digits Only", text.TypestyleTitleMedium),
				spacer.Height(8),
				foundationText.BasicTextField(
					digitsOnlyState,
					func(value string) {
						digitsOnlyState.SetTextAndPlaceCursorAtEnd(value)
					},
					foundationText.WithInputTransformation(input.DigitsOnlyTransformation()),
					foundationText.WithTextFieldModifier(modifier),
				),
				spacer.Height(24),

				// Footer
				text.Text("✓ Using TextFieldState + EditableTextLayoutController", text.TypestyleBodySmall),
			),
		)(c)
	}
}
