package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/runtime"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err := runtime.App(
		runtime.Title("Component Showcase"),
		runtime.Size(unit.NewDpSize(1024, 768)),
		runtime.Content(UI()),
		runtime.WithRootContext(ctx),
	).Run()

	if err != nil && !errors.Is(err, context.Canceled) {
		log.Fatal(err)
	}
}
