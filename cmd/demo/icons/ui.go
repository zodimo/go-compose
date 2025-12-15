package main

import (
	"image/color"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/foundation/icon"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/material3/iconbutton"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/scale"
	"github.com/zodimo/go-compose/theme"

	"github.com/zodimo/go-compose/pkg/api"
)

func UI(c api.Composer) api.LayoutNode {
	var rows []api.Composable
	chunkSize := 40

	themeManager := theme.GetThemeManager()

	// Demo IconButtons
	var demoIconButtons []api.Composable
	if len(List) >= 4 {
		// Using first 4 icons for the buttons
		demoIconButtons = append(demoIconButtons,
			iconbutton.Standard(func() {}, List[0].Data, "Standard"),
			iconbutton.Filled(func() {}, List[1].Data, "Filled"),
			iconbutton.FilledTonal(func() {}, List[2].Data, "FilledTonal"),
			iconbutton.Outlined(func() {}, List[3].Data, "Outlined"),
		)
		rows = append(rows, row.Row(compose.Sequence(demoIconButtons...)))
	}

	colors := []color.NRGBA{
		{R: 213, G: 0, B: 0, A: 255},    // Red
		{R: 197, G: 17, B: 98, A: 255},  // Pink
		{R: 170, G: 0, B: 255, A: 255},  // Purple
		{R: 98, G: 0, B: 234, A: 255},   // Deep Purple
		{R: 48, G: 79, B: 254, A: 255},  // Indigo
		{R: 41, G: 98, B: 255, A: 255},  // Blue
		{R: 0, G: 145, B: 234, A: 255},  // Light Blue
		{R: 0, G: 184, B: 212, A: 255},  // Cyan
		{R: 0, G: 191, B: 165, A: 255},  // Teal
		{R: 0, G: 200, B: 83, A: 255},   // Green
		{R: 100, G: 221, B: 23, A: 255}, // Light Green
		{R: 174, G: 234, B: 0, A: 255},  // Lime
		{R: 255, G: 214, B: 0, A: 255},  // Yellow
		{R: 255, G: 171, B: 0, A: 255},  // Amber
		{R: 255, G: 109, B: 0, A: 255},  // Orange
		{R: 221, G: 44, B: 0, A: 255},   // Deep Orange
	}

	for i := 0; i < len(List); i += chunkSize {
		end := i + chunkSize
		if end > len(List) {
			end = len(List)
		}

		var rowItems []api.Composable
		for j, def := range List[i:end] {
			// Range over colors
			col := colors[(i+j)%len(colors)]

			// Create a scaling effect (oscillating size)
			// Scale between 0.8 and 1.5
			scaleFactor := 0.8 + 0.7*float32((i+j)%10)/10.0

			rowItems = append(rowItems, icon.Icon(def.Data,
				icon.WithColor(themeManager.ColorDescriptor(col)),
				icon.WithModifier(scale.Scale(scaleFactor)),
				icon.WithModifier(padding.All(3)),
			))
		}
		rows = append(rows, row.Row(compose.Sequence(rowItems...)))
	}

	c = column.Column(
		compose.Sequence(rows...),
	)(c)

	return c.Build()
}
