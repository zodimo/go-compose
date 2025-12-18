package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/runtime"
	"github.com/zodimo/go-compose/state"
	"github.com/zodimo/go-compose/store"
	"github.com/zodimo/go-compose/theme"
)

func main() {
	go func() {
		w := new(app.Window)
		w.Option(
			app.Title("StateFlow Demo"),
			app.Size(unit.Dp(800), unit.Dp(600)),
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
	store := store.NewPersistentState(map[string]state.MutableValue{})
	store.SetOnStateChange(func() {
		window.Invalidate()
	})

	rt := runtime.NewRuntime()
	themeManager := theme.GetThemeManager()

	for {
		switch frameEvent := window.Event().(type) {
		case app.DestroyEvent:
			return frameEvent.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, frameEvent)

			// Initialize Theme (M3)
			gtx = themeManager.Material3ThemeInit(gtx)

			composer := UI(compose.NewComposer(store))

			callOp := rt.Run(gtx, composer.Build())
			callOp.Add(gtx.Ops)
			frameEvent.Frame(gtx.Ops)
			window.Invalidate()
		}
	}
}
