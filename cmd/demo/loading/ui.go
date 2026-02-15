package main

import (
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/material3/appbar"
	"github.com/zodimo/go-compose/compose/material3/progress"
	"github.com/zodimo/go-compose/compose/material3/scaffold"
	"github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"
)

func UI() api.Composable {
	return func(c api.Composer) api.Composer {
		return scaffold.Scaffold(
			func(c api.Composer) api.Composer {
				return column.Column(
					c.Sequence(
						text.TextWithStyle("Default Loading Indicator", text.TypestyleBodyLarge),
						progress.LoadingIndicator(),
						spacer.Height(16),
						text.TextWithStyle("Loading Indicator with size 200x200", text.TypestyleBodyLarge),
						progress.LoadingIndicator(
							progress.WithModifier(size.Size(200, 200)),
						),
					),
					column.WithModifier(size.FillMax()),
					column.WithSpacing(column.SpaceEvenly),
					column.WithAlignment(column.Middle),
				)(c)
			},
			scaffold.WithTopBar(
				appbar.TopAppBar(
					text.TextWithStyle("Loading Demo", text.TypestyleTitleMedium),
				),
			),
		)(c)
	}
}
