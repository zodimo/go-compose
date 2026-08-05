package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/zodimo/go-compose/compose/effect"
	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/text"
	"github.com/zodimo/go-compose/compose/material3/button"
	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/runtime"
)

func main() {
	if err := runtime.App(
		runtime.Title("LaunchedEffect Demo"),
		runtime.Size(unit.NewDpSize(800, 600)),
		runtime.Content(UI()),
	).Run(); err != nil {
		log.Fatal(err)
	}
}

func UI() api.Composable {
	return func(c api.Composer) api.Composer {
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
						button.Filled(func() {
							counter.Set(counter.Get().(int) + 1)
						}, "Increment Counter (Restarts Effect)",
							button.WithModifier(padding.All(20)),
						),
					),
				),
			),
		)(c)

		return c
	}
}
