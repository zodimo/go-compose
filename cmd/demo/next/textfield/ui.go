package main

import (
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	foundationText "github.com/zodimo/go-compose/compose/foundation/next/text"
	"github.com/zodimo/go-compose/compose/foundation/next/text/input"
	"github.com/zodimo/go-compose/compose/material3/text"
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

		return column.Column(
			c.Sequence(
				// Title
				text.Text("BasicTextField Demo (Next)", text.TypestyleHeadlineMedium),
				spacer.Height(24),

				// Section: Basic Text Field
				text.Text("Basic TextField", text.TypestyleTitleMedium),
				spacer.Height(8),
				foundationText.BasicTextField(basicState),
				spacer.Height(16),

				// Section: Single Line
				text.Text("Single Line TextField", text.TypestyleTitleMedium),
				spacer.Height(8),
				foundationText.BasicTextField(
					singleLineState,
					foundationText.WithLineLimits(input.TextFieldLineLimitsSingleLine),
				),
				spacer.Height(16),

				// Section: Max Length (10 chars)
				text.Text("Max Length (10 chars)", text.TypestyleTitleMedium),
				spacer.Height(8),
				foundationText.BasicTextField(
					maxLengthState,
					foundationText.WithInputTransformation(input.MaxLengthTransformation(10)),
				),
				spacer.Height(16),

				// Section: Digits Only
				text.Text("Digits Only", text.TypestyleTitleMedium),
				spacer.Height(8),
				foundationText.BasicTextField(
					digitsOnlyState,
					foundationText.WithInputTransformation(input.DigitsOnlyTransformation()),
				),
				spacer.Height(24),

				// Footer
				text.Text("✓ Using TextFieldState + TextLayoutController", text.TypestyleBodySmall),
			),
		)(c)
	}
}
