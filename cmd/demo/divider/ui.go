package main

import (
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/material3/divider"
	"github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/modifiers/weight"
	"github.com/zodimo/go-compose/pkg/api"
)

func UI() api.Composable {
	return func(c api.Composer) api.Composer {

		return column.Column(
			c.Sequence(
				horizontalDividerDemo(),
				verticalDividerDemo(),
			),
		)(c)
	}
}

func horizontalDividerDemo() api.Composable {
	return func(c api.Composer) api.Composer {
		return column.Column(
			c.Sequence(
				text.HeadlineMedium("Horizontal Divider Demo"),
				spacer.Height(16),
				divider.HorizontalDivider(),
				spacer.Height(16),
				divider.HorizontalDivider(divider.WithThickness(2)),
				spacer.Height(16),
				divider.HorizontalDivider(divider.WithThickness(4), divider.WithColor(graphics.ColorRed)),
			),
			column.WithSpacing(column.SpaceSides),
			column.WithAlignment(column.Middle),
			column.WithModifier(weight.Weight(1).Then(size.FillMaxWidth())),
		)(c)
	}
}

func verticalDividerDemo() api.Composable {
	return func(c api.Composer) api.Composer {
		return column.Column(
			c.Sequence(
				text.HeadlineMedium("Vertical Divider Demo"),
				spacer.Height(16),
				row.Row(
					c.Sequence(
						spacer.Width(16),
						divider.VerticalDivider(),
						spacer.Width(16),
						divider.VerticalDivider(divider.WithThickness(2)),
						spacer.Width(16),
						divider.VerticalDivider(divider.WithThickness(4), divider.WithColor(graphics.ColorRed)),
					),
					row.WithSpacing(row.SpaceSides),
					row.WithAlignment(row.Middle),
					row.WithModifier(size.FillMaxHeight()),
				),
			),
			column.WithSpacing(column.SpaceSides),
			column.WithAlignment(column.Middle),
			column.WithModifier(weight.Weight(1).Then(size.FillMaxWidth())),
		)(c)
	}
}
