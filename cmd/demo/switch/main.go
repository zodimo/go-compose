package main

import (
	"fmt"
	"log"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/material3/appbar"
	"github.com/zodimo/go-compose/compose/material3/scaffold"
	mswitch "github.com/zodimo/go-compose/compose/material3/switch"
	"github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/runtime"
)

func main() {
	if err := runtime.App(
		runtime.Title("Switch Demo"),
		runtime.Size(unit.NewDpSize(400, 700)),
		runtime.Content(UI()),
	).Run(); err != nil {
		log.Fatal(err)
	}
}

func UI() api.Composable {
	return func(c compose.Composer) api.Composer {
		checked1 := c.State("switch_state_1", func() any { return false })
		checked2 := c.State("switch_state_2", func() any { return false })
		c = scaffold.Scaffold(
			column.Column(
				c.Sequence(
					// Case 1: Switch in a row, text then switch
					row.Row(
						c.Sequence(
							text.BodyMedium("Enable Feature"),
							spacer.Width(20),
							mswitch.Switch(
								checked1.Get().(bool),
								func(b bool) {
									checked1.Set(b)
									fmt.Printf("Switch toggled 1: %v\n", b)
								},
							),
						),
						row.WithSpacing(row.SpaceStart),
						row.WithAlignment(row.Middle),
					),
					text.BodyMedium("Switch should be to the right to text above"),
					// Case 2: Switch in a row, switch then text
					row.Row(
						c.Sequence(
							mswitch.Switch(
								checked2.Get().(bool),
								func(b bool) {
									checked2.Set(b)
									fmt.Printf("Switch toggled 2: %v\n", b)
								},
							),
							spacer.Width(20),
							text.BodyMedium("Enable Feature"),
						),
						row.WithSpacing(row.SpaceStart),
						row.WithAlignment(row.Middle),
					),
					// Case 2: Just confirmation text
					text.BodyMedium("Switch should be to the right to text above"),
				),
				column.WithSpacing(column.SpaceAround),
				column.WithAlignment(column.Middle),
			),
			scaffold.WithTopBar(
				appbar.TopAppBar(
					func(c compose.Composer) compose.Composer {
						return text.TitleLarge("Switch Demo")(c)
					},
				),
			),
		)(c)
		return c
	}
}
