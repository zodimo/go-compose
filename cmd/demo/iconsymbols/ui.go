package main

import (
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/text"
	"github.com/zodimo/go-compose/compose/material3/icon"
	"github.com/zodimo/go-compose/compose/ui/geometry"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/pkg/api"
)

// UI shows the capabilities a font glyph could not express: variable-font axes,
// optical sizing, the three symbol styles, gradient fills and strokes.
func UI() api.Composable {
	return func(c api.Composer) api.Composer {
		label := func(s string) api.Composable {
			return text.Text(s, text.WithModifier(padding.All(8)))
		}
		iconRow := func(items ...api.Composable) api.Composable {
			return row.Row(c.Sequence(items...),
				row.WithSpacing(row.SpaceEvenly),
				row.WithAlignment(row.Middle),
				row.WithModifier(padding.All(8)),
			)
		}

		// The fill axis: outlined through filled.
		fillRow := iconRow(
			icon.Icon(icon.SymbolFavorite, icon.WithFill(0), icon.WithSize(unit.Dp(40))),
			icon.Icon(icon.SymbolFavorite, icon.WithFill(0.5), icon.WithSize(unit.Dp(40))),
			icon.Icon(icon.SymbolFavorite, icon.WithFill(1), icon.WithSize(unit.Dp(40))),
		)

		// The weight axis at a fixed size.
		weightRow := iconRow(
			icon.Icon(icon.SymbolStar, icon.WithWeight(100), icon.WithSize(unit.Dp(40))),
			icon.Icon(icon.SymbolStar, icon.WithWeight(400), icon.WithSize(unit.Dp(40))),
			icon.Icon(icon.SymbolStar, icon.WithWeight(700), icon.WithSize(unit.Dp(40))),
		)

		// Optical size follows the rendered size unless set explicitly.
		opticalRow := iconRow(
			icon.Icon(icon.SymbolSearch, icon.WithSize(unit.Dp(16))),
			icon.Icon(icon.SymbolSearch, icon.WithSize(unit.Dp(32))),
			icon.Icon(icon.SymbolSearch, icon.WithSize(unit.Dp(48))),
		)

		// Styles: outlined, rounded, sharp. Rounded and sharp need the
		// symbols_extra_fonts build tag, otherwise they fall back to outlined.
		styleRow := iconRow(
			icon.Icon(icon.SymbolHome, icon.WithSymbolStyle(icon.SymbolOutlined), icon.WithSize(unit.Dp(40))),
			icon.Icon(icon.SymbolHome, icon.WithSymbolStyle(icon.SymbolRounded), icon.WithSize(unit.Dp(40))),
			icon.Icon(icon.SymbolHome, icon.WithSymbolStyle(icon.SymbolSharp), icon.WithSize(unit.Dp(40))),
		)

		// A gradient fill, which no font glyph can do.
		gradient := graphics.LinearGradientBrush(
			[]graphics.Color{
				graphics.NewColorSrgb(255, 64, 129, 255),
				graphics.NewColorSrgb(41, 98, 255, 255),
			},
			geometry.NewOffset(0, 0),
			geometry.NewOffset(0, 64),
			graphics.TileModeClamp,
		)
		brushRow := iconRow(
			icon.Icon(icon.SymbolFavorite, icon.WithBrush(gradient), icon.WithSize(unit.Dp(64))),
			icon.Icon(icon.SymbolStar, icon.WithBrush(gradient), icon.WithSize(unit.Dp(64))),
		)

		// Stroked icons, drawn as outlines at a chosen width.
		strokeRow := iconRow(
			icon.Icon(icon.SymbolFavorite, icon.WithStroke(unit.Dp(1)), icon.WithSize(unit.Dp(48))),
			icon.Icon(icon.SymbolStar, icon.WithStroke(unit.Dp(2)), icon.WithSize(unit.Dp(48))),
		)

		return column.Column(
			c.Sequence(
				label("Fill axis: 0, 0.5, 1"),
				fillRow,
				label("Weight axis: 100, 400, 700"),
				weightRow,
				label("Optical size follows size: 16dp, 32dp, 48dp"),
				opticalRow,
				label("Styles: outlined, rounded, sharp"),
				styleRow,
				label("Gradient brush"),
				brushRow,
				label("Stroke: 1dp, 2dp"),
				strokeRow,
			),
			column.WithModifier(padding.All(16)),
		)(c)
	}
}
