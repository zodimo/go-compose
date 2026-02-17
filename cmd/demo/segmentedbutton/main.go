package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/op"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/runtime"
	"github.com/zodimo/go-compose/store"
	"github.com/zodimo/go-compose/theme"
)

func main() {
	go func() {
		w := new(app.Window)
		w.Option(app.Title("Segmented Button Demo"))
		w.Option(app.Size(600, 400))

		err := run(w)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(window *app.Window) error {
	var ops op.Ops
	themeManager := theme.GetThemeManager()

	store := store.NewPersistentState()
	store.Subscribe(func() {
		window.Invalidate()
	})

	runtime := runtime.NewRuntime()

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			gtx = themeManager.Material3ThemeInit(gtx)

			composer := compose.NewComposer(api.ComposerWithStore(store))
			callOp := runtime.Run(gtx, composer, UI())
			callOp.Add(gtx.Ops)
			e.Frame(gtx.Ops)
		}
	}
}
