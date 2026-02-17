package main

import (
	"fmt"
	"log"
	"os"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/material3/appbar"
	"github.com/zodimo/go-compose/compose/material3/scaffold"
	mswitch "github.com/zodimo/go-compose/compose/material3/switch"
	"github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/runtime"
	"github.com/zodimo/go-compose/store"
	"github.com/zodimo/go-compose/theme"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/unit"
)

func main() {
	go func() {
		w := new(app.Window)
		w.Option(
			app.Title("Switch Demo"),
			app.Size(unit.Dp(400), unit.Dp(700)),
		)

		if err := Run(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	app.Main()
}

func Run(window *app.Window) error {
	var ops op.Ops

	store := store.NewPersistentState()
	store.Subscribe(func() {
		window.Invalidate()
	})

	runtime := runtime.NewRuntime()

	themeManager := theme.GetThemeManager()

	for {
		switch frameEvent := window.Event().(type) {
		case app.DestroyEvent:
			return frameEvent.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, frameEvent)

			// Init Theme
			gtx = themeManager.Material3ThemeInit(gtx)

			composer := compose.NewComposer(api.ComposerWithStore(store))
			callOp := runtime.Run(gtx, composer, UI())
			callOp.Add(gtx.Ops)
			frameEvent.Frame(gtx.Ops)
		}
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
