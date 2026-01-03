package main

import (
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	fText "github.com/zodimo/go-compose/compose/foundation/text"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	uiText "github.com/zodimo/go-compose/compose/ui/text"
	"github.com/zodimo/go-compose/pkg/api"
)

func UI(c api.Composer) api.LayoutNode {
	root := column.Column(
		c.Sequence(
			fText.Text(
				"Hello World",
			),
			fText.Text("🚦 Traffic Light"),
			fText.Text(
				"Hello World",
				fText.WithTextStyleOptions(
					uiText.WithColor(graphics.ColorRed),
				),
			),
		),
	)

	return root(c).Build()
}
