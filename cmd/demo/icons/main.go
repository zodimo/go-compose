package main

import (
	"log"
	"os"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/lifecycle"
	"github.com/zodimo/go-compose/runtime"
	"github.com/zodimo/go-compose/state"
	"github.com/zodimo/go-compose/store"
	"github.com/zodimo/go-compose/theme"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/unit"
)

func main() {

	go func() {
		w := new(app.Window)
		w.Option(app.Title("Compose Icons"))
		w.Option(app.Size(unit.Dp(1250), unit.Dp(800)))

		if err := Run(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	app.Main()

}

func Run(window *app.Window) error {

	enLocale := system.Locale{Language: "en", Direction: system.LTR}

	var ops op.Ops

	store := store.NewPersistentState(map[string]state.MutableValue{})
	store.Subscribe(func() {
		window.Invalidate()
	})

	runtimeOptions := []runtime.RuntimeOption{}
	if lifecycleAwareStore, ok := store.(lifecycle.FrameLifecycleAwarePersistentState); ok {
		runtimeOptions = append(runtimeOptions,
			runtime.WithOnStartFrame(lifecycleAwareStore.StartFrame),
			runtime.WithOnEndFrame(lifecycleAwareStore.EndFrame),
		)
	}

	runtime := runtime.NewRuntime(runtimeOptions...)

	themeManager := theme.GetThemeManager()

	for {
		switch frameEvent := window.Event().(type) {
		case app.DestroyEvent:
			return frameEvent.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, frameEvent)

			gtx.Locale = enLocale

			// M3 Widget Requirement
			gtx = themeManager.Material3ThemeInit(gtx)

			composer := compose.NewComposer(store)

			callOp := runtime.Run(gtx, composer, UI())
			callOp.Add(gtx.Ops)
			frameEvent.Frame(gtx.Ops)

		}
	}

}
