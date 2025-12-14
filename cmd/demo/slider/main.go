package main

import (
	"log"
	"os"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/runtime"
	"github.com/zodimo/go-compose/internal/state"
	"github.com/zodimo/go-compose/internal/store"
	"github.com/zodimo/go-compose/theme"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/unit"
)

func main() {
	go func() {
		w := new(app.Window)
		w.Option(app.Title("Slider Demo"))
		w.Option(app.Size(unit.Dp(400), unit.Dp(800)))

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
	// runtime := runtime.NewRuntime() // Not used directly in this loop style often

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
			// UI returns Composable, we need to execute it or wrapped it
			// Checking other examples: UI(composer) returns LayoutNode if it matches signature
			// But here UI() returns Composable.
			// So:
			rootComposable := UI()
			// We need to build the tree.
			// compose.NewComposer returns *composer.
			// Composable is func(Composer) Composer.

			// Execute composable
			rootComposable(composer)

			// Build to get LayoutNode
			layoutNode := composer.Build()

			rt := runtime.NewRuntime()
			callOp := rt.Run(gtx, layoutNode)
			callOp.Add(gtx.Ops)
			frameEvent.Frame(gtx.Ops)

			// Invalidate for animations/interactions if needed
			// window.Invalidate() // Using aggressive invalidation for demo
		}
	}
}
