package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/zodimo/go-compose/compose"
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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	rootContext, cancelFunc := context.WithCancel(ctx)

	go func() {
		defer cancelFunc()
		w := new(app.Window)
		w.Option(
			app.Title("Component Showcase"),
			app.Size(unit.Dp(1024), unit.Dp(768)),
		)

		if err := Run(w, rootContext); err != nil {
			log.Fatal(err)
		}

	}()

	go func() {
		<-rootContext.Done()
		fmt.Println("Root Context done")
		os.Exit(0)
	}()

	app.Main()

}

func Run(window *app.Window, rootContext context.Context) error {

	enLocale := system.Locale{Language: "en", Direction: system.LTR}

	var ops op.Ops

	store := store.NewPersistentState(
		map[string]state.MutableValue{},
		store.WithRootContext(rootContext),
	)
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

			// M3 Widget Requirement
			gtx = themeManager.Material3ThemeInit(gtx)

			composer := compose.NewComposer(store)
			layoutNode := UI(composer)

			callOp := runtime.Run(gtx, layoutNode)
			callOp.Add(gtx.Ops)
			frameEvent.Frame(gtx.Ops)

		}
	}

}
