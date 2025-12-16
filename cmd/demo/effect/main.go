package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/effect"
	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/material/button"
	"github.com/zodimo/go-compose/compose/foundation/text"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/runtime"
	"github.com/zodimo/go-compose/state"
	"github.com/zodimo/go-compose/store"
	"github.com/zodimo/go-compose/theme"
)

func main() {
	go func() {
		w := new(app.Window)
		w.Option(
			app.Title("LaunchedEffect Demo"),
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

			composer := compose.NewComposer(store)
			layoutNode := UI(composer)

			callOp := rt.Run(gtx, layoutNode)
			callOp.Add(gtx.Ops)
			frameEvent.Frame(gtx.Ops)
			window.Invalidate()
		}
	}
}

func UI(c api.Composer) api.LayoutNode {
	counter := c.State("counter", func() any { return 0 })
	effectStatus := c.State("effect_status", func() any { return "Waiting..." })

	// Wrap everything in a Box so LaunchedEffect is just a sibling, not the Root
	c = box.Box(
		c.Sequence(
			// LaunchedEffect that reacts to counter
			effect.LaunchedEffect(func(ctx context.Context) {
				currentCount := counter.Get().(int)

				effectStatus.Set(fmt.Sprintf("Effect STARTED for %d", currentCount))
				fmt.Printf("Effect STARTED for %d\n", currentCount)

				select {
				case <-time.After(2 * time.Second):
					if ctx.Err() == nil {
						effectStatus.Set(fmt.Sprintf("Effect FINISHED for %d", currentCount))
						fmt.Printf("Effect FINISHED for %d\n", currentCount)
					}
				case <-ctx.Done():
					// This might effectively be overwritten by the next effect starting immediately
					// but we should see it in logs at least.
					// Note: If we set state here, it might be racey with the new effect setting "STARTED".
					// But cancel() is called synchronously before next effect starts?
					// No, cancel() is sync, but the goroutine might take a microsecond to wake up and print.
					// The NEXT effect starts immediately after cancel().
					// So "STARTED for N+1" might overwrite "CANCELLED for N".
					// We'll log to console to verify cancellation.
					fmt.Printf("Effect CANCELLED for %d\n", currentCount)
				}
			}, counter.Get()),

			column.Column(
				c.Sequence(
					text.Text(fmt.Sprintf("Counter: %d", counter.Get()), text.WithModifier(padding.All(10))),
					text.Text(fmt.Sprintf("Status: %v", effectStatus.Get()), text.WithModifier(padding.All(10))),
					button.Button(func() {
						counter.Set(counter.Get().(int) + 1)
					}, "Increment Counter (Restarts Effect)",
						button.WithModifier(padding.All(20)),
					),
				),
			),
		),
	)(c)

	return c.Build()
}
