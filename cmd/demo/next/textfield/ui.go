package main

import (
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/foundation/next/text/input"
	foundationTextField "github.com/zodimo/go-compose/compose/foundation/next/textfield"
	"github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/modifiers/background"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"
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

		modifier := background.Background(graphics.ColorLightGray).
			Then(size.FillMaxWidth()).
			Then(padding.All(16))

		return column.Column(
			c.Sequence(
				// Title
				text.HeadlineMedium("BasicTextField Demo (Next)"),
				spacer.Height(24),

				// Section: Basic Text Field
				text.TitleMedium("Basic TextField"),
				spacer.Height(8),
				foundationTextField.BasicTextField(
					basicState,
					func(value string) {
						basicState.SetTextAndPlaceCursorAtEnd(value)
					},
					foundationTextField.WithModifier(modifier),
				),
				spacer.Height(16),

				// Section: Single Line
				text.TitleMedium("Single Line TextField"),
				spacer.Height(8),
				foundationTextField.BasicTextField(
					singleLineState,
					func(value string) {
						singleLineState.SetTextAndPlaceCursorAtEnd(value)
					},
					foundationTextField.WithLineLimits(input.TextFieldLineLimitsSingleLine),
					foundationTextField.WithModifier(modifier),
				),
				spacer.Height(16),

				// Section: Max Length (10 chars)
				text.TitleMedium("Max Length (10 chars)"),
				spacer.Height(8),
				foundationTextField.BasicTextField(
					maxLengthState,
					func(value string) {
						maxLengthState.SetTextAndPlaceCursorAtEnd(value)
					},
					foundationTextField.WithInputTransformation(input.MaxLengthTransformation(10)),
					foundationTextField.WithModifier(modifier),
				),
				spacer.Height(16),

				// Section: Digits Only
				text.TitleMedium("Digits Only"),
				spacer.Height(8),
				foundationTextField.BasicTextField(
					digitsOnlyState,
					func(value string) {
						digitsOnlyState.SetTextAndPlaceCursorAtEnd(value)
					},
					foundationTextField.WithInputTransformation(input.DigitsOnlyTransformation()),
					foundationTextField.WithModifier(modifier),
				),
				spacer.Height(24),

				// Footer
				text.BodySmall("✓ Using TextFieldState + EditableTextLayoutController"),
			),
		)(c)
	}
}
