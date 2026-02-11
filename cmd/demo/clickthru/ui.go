package main

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/foundation/stopclickthru"
	"github.com/zodimo/go-compose/compose/material3/iconbutton"
	"github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/modifiers/clickable"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"
	mwIcons "golang.org/x/exp/shiny/materialdesign/icons"
)

func UI() api.Composable {
	return func(c api.Composer) api.Composer {
		return column.Column(
			c.Sequence(
				row.Row(
					c.Sequence(
						text.TitleSmall("Clickable row"),
						spacer.Weight(1),
						stopclickthru.StopClickThru(text.TitleSmall("Not clickable text")),
						spacer.Weight(1),
						iconbutton.Filled(func() {
							fmt.Println("Icon clicked")
						}, mwIcons.AVAVTimer, ""),
					),
					row.WithModifier(
						clickable.OnClick(func() {
							fmt.Println("Row clicked")
						}).Then(size.FillMaxWidth()),
					),
					row.WithAlignment(row.Middle),
				),
				column.Column(
					c.Sequence(
						text.TitleSmall("Clickable Column"),
						spacer.Weight(1),
						stopclickthru.StopClickThru(text.TitleSmall("Not clickable text")),
						spacer.Weight(1),
						iconbutton.Filled(func() {
							fmt.Println("Icon clicked")
						}, mwIcons.AVAVTimer, ""),
					),
					column.WithModifier(
						clickable.OnClick(func() {
							fmt.Println("Column clicked")
						}).Then(size.FillMaxHeight()),
					),
					column.WithAlignment(column.Middle),
				),
			),
		)(c)
	}
}
