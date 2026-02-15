package main

import (
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/material3/icon"
	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/pkg/api"
)

func UI() api.Composable {
	return func(c api.Composer) api.Composer {
		c = column.Column(
			c.Sequence(
				// Search icon
				icon.Symbol(icon.SymbolSearch, icon.WithSymbolSize(unit.Sp(48))),

				// Home icon
				icon.Symbol(icon.SymbolHome, icon.WithSymbolSize(unit.Sp(48))),

				// Settings icon
				icon.Symbol(icon.SymbolSettings, icon.WithSymbolSize(unit.Sp(100))),
			),
			column.WithSpacing(column.SpaceEvenly),
			column.WithAlignment(column.Middle),
		)(c)

		return c
	}
}
