# Runtime App Migration Guide

This guide documents the migration to `runtime.App` introduced in the `runtime-app` change.

## Overview

Previously, Go Compose applications and demo mains required extensive Gio UI boilerplate in `main.go`:
- Direct instantiation of `app.Window`
- Spawning a goroutine with `app.Main()` blocking the main thread
- Manual loop over `window.Event()` handling `app.FrameEvent` and `app.DestroyEvent`
- Manual management of `op.Ops`, `themeManager.Material3ThemeInit`, `system.Locale`, and `store.PersistentState` subscriptions
- Direct application of `DrawCommand` onto `gtx.Ops`

With `runtime.App`, the entire lifecycle, per-window event loop, theming, state invalidation, window context provisioning, and platform-specific main thread semantics are handled automatically.

## Before and After

### Before

```go
package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
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
		w.Option(app.Title("My App"))
		w.Option(app.Size(unit.Dp(800), unit.Dp(600)))

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
			gtx = themeManager.Material3ThemeInit(gtx).(layout.Context)

			composer := compose.NewComposer(api.ComposerWithStore(store))
			cmd := runtime.Run(gtx, composer, UI())
			cmd.Apply(gtx.Ops)
			frameEvent.Frame(gtx.Ops)
		}
	}
}
```

### After (Single-Window Convenience)

```go
package main

import (
	"log"

	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/runtime"
)

func main() {
	if err := runtime.App(
		runtime.Title("My App"),
		runtime.Size(unit.NewDpSize(800, 600)),
		runtime.Content(UI()),
	).Run(); err != nil {
		log.Fatal(err)
	}
}
```

## Multi-Window Support

For desktop platforms supporting multiple windows, use `runtime.NewApp(...)`:

```go
package main

import (
	"log"

	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/runtime"
)

func main() {
	app := runtime.NewApp()
	app.AddWindow(runtime.WindowOptions{
		Title:   "Primary Window",
		Size:    unit.NewDpSize(800, 600),
		Content: PrimaryUI(),
	})
	app.AddWindow(runtime.WindowOptions{
		Title:   "Secondary Tool Window",
		Size:    unit.NewDpSize(400, 300),
		Content: ToolUI(),
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
```

> **Note on Mobile & Web**: On single-window platforms (iOS, Android, WASM/JS), attempting to register more than 1 window returns a validation error before any window is created.

## `Run() error` Semantics

1. **Clean Exit**: When all windows are closed by the user (with `QuitWhenLastWindowCloses(true)` default), `Run()` returns `nil`.
2. **Cancellation**: If a root context is provided via `runtime.WithRootContext(ctx)` and cancelled, `Run()` returns `context.Canceled`.
3. **Window Failure**: If any window encounters an unrecoverable error during its lifecycle, `Run()` returns the first error encountered.
4. **macOS Lifecycle**: On macOS, `app.Main()` must run on the main OS thread and does not return upon normal termination. In this environment, the exit monitor calls `os.Exit(0)` on clean exit or `os.Exit(1)` on error.

## Automatic Composition Locals

Each window's frame loop automatically binds `platform.LocalWindow` with an abstraction over the native window, allowing child composables to request window invalidation or access window-level state without importing Gio UI packages.
