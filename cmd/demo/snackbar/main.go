package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/runtime"
	"github.com/zodimo/go-compose/store"
	"github.com/zodimo/go-compose/theme"
)

func main() {
	go func() {
		w := new(app.Window)
		w.Option(app.Title("Snackbar Demo"))
		w.Option(app.Size(unit.Dp(400), unit.Dp(600)))

		if err := run(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(window *app.Window) error {
	enLocale := system.Locale{Language: "en", Direction: system.LTR}
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
			gtx.Locale = enLocale
			gtx = themeManager.Material3ThemeInit(gtx)

			composer := compose.NewComposer(api.ComposerWithStore(store))
			callOp := runtime.Run(gtx, composer, UI())
			callOp.Add(gtx.Ops)
			frameEvent.Frame(gtx.Ops)
		}
	}
}
