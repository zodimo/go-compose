package main

import (
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/material3/appbar"
	"github.com/zodimo/go-compose/compose/foundation/material3/progress"
	"github.com/zodimo/go-compose/compose/foundation/material3/scaffold"
	"github.com/zodimo/go-compose/compose/foundation/material3/text"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"
)

func UI(c api.Composer) api.LayoutNode {
	return scaffold.Scaffold(
		func(c api.Composer) api.Composer {
			return column.Column(
				func(c api.Composer) api.Composer {
					text.Text("Default Loading Indicator", text.TypestyleBodyLarge)(c)
					progress.LoadingIndicator()(c)
					return c
				},
				column.WithModifier(size.FillMax()),
				column.WithSpacing(column.SpaceEvenly),
				column.WithAlignment(column.Middle),
			)(c)
		},
		scaffold.WithTopBar(
			appbar.TopAppBar(
				text.Text("Loading Demo", text.TypestyleTitleMedium),
			),
		),
	)(c).Build()
}
