package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/op"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/runtime"
	"github.com/zodimo/go-compose/state"
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

func run(w *app.Window) error {
	var ops op.Ops
	themeManager := theme.GetThemeManager()

	persistentStore := store.NewPersistentState(map[string]state.MutableValue{})
	rt := runtime.NewRuntime()

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			gtx = themeManager.Material3ThemeInit(gtx)

			composer := compose.NewComposer(persistentStore)
			rootComposer := UI()(composer)
			layoutNode := rootComposer.Build()

			_ = rt.Run(gtx, layoutNode)
			e.Frame(gtx.Ops)
		}
	}
}
