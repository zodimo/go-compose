package main

import (
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/text"
	m3icon "github.com/zodimo/go-compose/compose/material3/icon"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/pkg/api"

	"golang.org/x/exp/shiny/materialdesign/icons"
)

// UI demonstrates unified color and size for both IconBytes and SymbolName sources.
// Each row shows the same size and color applied to both icon types side-by-side.
func UI() api.Composable {
	return func(c api.Composer) api.Composer {
		// Define test configurations: size and color pairs
		type iconConfig struct {
			label string
			size  unit.Dp
			color graphics.Color
		}

		configs := []iconConfig{
			{"Small (16dp) - Red", 16, graphics.NewColorSrgb(213, 0, 0, 255)},
			{"Medium (24dp) - Blue", 24, graphics.NewColorSrgb(41, 98, 255, 255)},
			{"Large (32dp) - Green", 32, graphics.NewColorSrgb(0, 200, 83, 255)},
			{"XLarge (48dp) - Purple", 48, graphics.NewColorSrgb(170, 0, 255, 255)},
			{"XXLarge (64dp) - Orange", 64, graphics.NewColorSrgb(255, 109, 0, 255)},
		}

		var rows []api.Composable

		// Header row
		rows = append(rows,
			row.Row(
				c.Sequence(
					text.Text("IconBytes (SVG)", text.WithModifier(padding.All(8))),
					text.Text("SymbolName (Font)", text.WithModifier(padding.All(8))),
					text.Text("Size/Color", text.WithModifier(padding.All(8))),
				),
				row.WithSpacing(row.SpaceEvenly),
			),
		)

		// For each config, show IconBytes and SymbolName side by side
		for _, cfg := range configs {
			// IconBytes: Home icon from materialdesign/icons
			iconBytes := m3icon.Icon(
				m3icon.IconBytes(icons.ActionHome),
				m3icon.WithSize(cfg.size),
				m3icon.WithColor(cfg.color),
			)

			// SymbolName: Home symbol from Material Symbols font
			iconSymbol := m3icon.Icon(
				m3icon.SymbolHome,
				m3icon.WithSize(cfg.size),
				m3icon.WithColor(cfg.color),
			)

			// Label
			label := text.Text(cfg.label, text.WithModifier(padding.All(8)))

			rows = append(rows,
				row.Row(
					c.Sequence(
						wrapWithPadding(iconBytes),
						wrapWithPadding(iconSymbol),
						label,
					),
					row.WithSpacing(row.SpaceEvenly),
					row.WithAlignment(row.Middle),
				),
			)
		}

		c = column.Column(
			c.Sequence(rows...),
			column.WithModifier(padding.All(16)),
		)(c)

		return c
	}
}

func wrapWithPadding(composable api.Composable) api.Composable {
	return func(c api.Composer) api.Composer {
		c = composable(c)
		return c
	}
}
