## Context

The render-backend seam (abstract-render-backend, archived) made the framework's public API gioui-free: `make check-api` is GREEN, and only `internal/layoutnode` + `internal/render` import engine types. But the *application shell* was left as boilerplate. Every one of the 48 demo `main.go` files contains the same ~30-line gio setup block:

```go
func main() {
	go func() {
		w := new(app.Window)
		w.Option(app.Title("Demo"), app.Size(unit.Dp(800), unit.Dp(600)))
		if err := Run(w); err != nil { log.Fatal(err) }
		os.Exit(0)
	}()
	app.Main()
}

func Run(window *app.Window) error {
	enLocale := system.Locale{Language: "en", Direction: system.LTR}
	var ops op.Ops
	store := store.NewPersistentState()
	store.Subscribe(func() { window.Invalidate() })
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

The design doc for the seam explicitly deferred this to "Phase 5: `runtime.App` event-loop encapsulation" (design.md decisions Q3, non-goals). This change implements that deferred phase — Option A scope: wrap the gio window loop only; do not build an input IR or move the frame-context through `Backend.BeginFrame`.

**Key constraint discovered during design**: the Backend interface's `BeginFrame(w, h, density)` cannot carry the gio frame context (Metric, Queue, Fonts, InputState) — the demo's `app.NewContext(&ops, frameEvent)` is the only source of those. Therefore `runtime.Run(ctx any, ...)` keeps its opaque-ctx signature, and App passes the gio `layout.Context` to it exactly as demos do today. This is the accepted Option A trade-off.

## Goals / Non-Goals

**Goals:**
- Application setup (window, event loop, store wiring, theme init, frame apply) is owned by `runtime.App`; demo `main.go` files contain zero gioui imports.
- `App.Run() error` blocks until close and reports the close cause: `nil` (last window closed), `context.Canceled` (root context cancelled), or the window's `DestroyEvent.Err`.
- Multi-window on desktop (`NewApp()` + `AddWindow()`), single-window degradation on iOS/Android/WASM.
- The `DrawCommand.Apply(ops *op.Ops)` gioui leak in the seam dies; replay moves inside App.
- API purity preserved: `make check-api` stays GREEN; new public signatures contain no gioui types.
- Existing behavior preserved: demos render identically after migration (golden baselines unchanged).

**Non-Goals:**
- Input IR / pointer-key-IME abstraction through the seam (Door 1 opening).
- Text shaping abstraction (design D7 of the seam).
- `Backend.BeginFrame` carrying the gio frame context (the Phase 4 wiring gap).
- Unifying the software/headless golden-harness entry with App (Option C) — golden harness keeps its own `RenderComponent` entry.
- Per-window theme isolation — `theme.GetThemeManager()` singleton remains; documented limitation.
- Replacing gioui or supporting a second *interactive* engine.

## Decisions

### D1 — Public API shape: `App()` convenience + `NewApp()` multi-window

```go
// runtime/app.go — public surface (no gioui types in signatures)

// App is the framework-owned application entry. It owns the window event
// loop(s), the per-window store, and the frame render pipeline.
type App struct {
	ctx        context.Context
	cancel     context.CancelFunc
	quitOnLast bool
	windows    []*windowConfig
}

// App creates a single-window application. Equivalent to NewApp() with one
// AddWindow(WindowOptions{...}).
func App(opts ...AppOption) *App

// NewApp creates a (potentially) multi-window application.
// Register windows with AddWindow, then call Run.
func NewApp(opts ...AppOption) *App

// AddWindow registers a window to be opened when Run is called.
func (a *App) AddWindow(opts WindowOptions) *App

// Run blocks until the application closes. Returns nil when the last window
// closes normally, context.Canceled when the root context is cancelled, or
// the DestroyEvent.Err of a window that failed. On macOS, Run never returns
// on desktop (Cocoa run loop); the monitor goroutine handles process exit.
// Run must be called from the main goroutine.
func (a *App) Run() error

// --- options ---
type AppOption func(*AppConfig)

type AppConfig struct {
	Title      string
	Size       unit.DpSize
	Content    api.Composable
	RootCtx    context.Context
	QuitOnLast bool
}

func Title(t string) AppOption
func Size(size unit.DpSize) AppOption                    // framework units
func Content(ui api.Composable) AppOption
func WithRootContext(ctx context.Context) AppOption
func QuitWhenLastWindowCloses(b bool) AppOption          // default true

// WindowOptions carries per-window settings for AddWindow.
type WindowOptions struct {
	Title   string
	Size    unit.DpSize
	Content api.Composable
}
```

Rationale: the `App(...)` variadic form keeps the 49 single-window demos to one line; `NewApp`+`AddWindow` recovers gio's desktop multi-window capability. Both funnel to the same internal machinery. `unit.DpSize` is the framework's own float32-based size type (`compose/ui/unit/dp.go`), converted to `gioUnit.Dp` via the existing `internal/unitconvert.DpToGioUnitUnsafe`.

Alternative considered: `App(opts...)` returning immediately and requiring an explicit `app.Main()` call — rejected because that leaks gio's `app.Main` into user code, defeating the goal.

### D2 — `Run() error` exit semantics are platform-aware

gio's `app.Main()` behavior is platform-dependent (verified in fork source):

| Platform | `osMain()` | Returns? | Window creation gated on Main? |
|---|---|---|---|
| Linux/Unix | `select {}` | never | no — driver runs own loop |
| Windows | `select {}` | never | no |
| JS/WASM | `select {}` | never | no |
| macOS | `C.gio_main()` → `[NSApp run]` | never | **yes** — `newWindow` blocks on `<-launched`, closed by `applicationDidFinishLaunching` (os_macos.go:997,1001) |
| Android/iOS | returns immediately (linkname main) | yes | n/a — platform creates the window |

Because `app.Main()` never returns on desktop, a uniform "call app.Main() then return an error" is impossible. Design:

```go
func (a *App) Run() error {
	// guard: must run on main goroutine
	exitCh := make(chan error, 1)

	// 1. Register the exit monitor: nil on last-window-close, err otherwise.
	go func() {
		select {
		case err := <-a.windowExit:   // first fatal window error wins
			exitCh <- err
		case <-a.ctx.Done():          // root context cancelled
			exitCh <- a.ctx.Err()     // context.Canceled / DeadlineExceeded
		case <-a.allWindowsClosed:    // WaitGroup done, no errors
			if a.quitOnLast {
				exitCh <- nil
			}
			// if !quitOnLast: stay alive; caller manages lifecycle
		}
	}()

	// 2. macOS: the Cocoa run loop must run on the main goroutine and never
	//    returns. The monitor handles exit; we never reach <-exitCh.
	if runtime.GOOS == "darwin" {
		for _, w := range a.windows {
			go a.runWindow(w)
		}
		go func() { <-exitCh; os.Exit(exitCode(err)) }()
		app.Main()            // never returns on desktop
		return nil            // unreachable
	}

	// 3. Linux/Windows/JS: windows are not gated on app.Main; block and return.
	for _, w := range a.windows {
		go a.runWindow(w)
	}
	return <-exitCh
}
```

The `exitCh` select picks the *first* terminal event: a fatal window error wins over "last window closed"; root-context cancel wins over both only if it fires first. This mirrors signal-handling semantics (kitchen demo today: SIGINT → clean exit).

macOS divergence is documented, not papered over: `Run()` blocks forever on desktop macOS because the Cocoa run loop is mandatory and non-returning; the monitor goroutine calls `os.Exit(0)`/`os.Exit(1)` (the pattern gio's own multiwindow example uses). Android/iOS return from `app.Main()` immediately and the platform drives the lifecycle; App treats them as single-window (see D5).

### D3 — Per-window `PersistentState` + invalidation loop moves inside App

```go
// runtime/window.go — per-window driver (Door 2: may import gio)

func (a *App) runWindow(cfg *windowConfig) {
	defer a.wg.Done()

	win := new(app.Window)
	if cfg.Title != "" { win.Option(app.Title(cfg.Title)) }
	if cfg.Size != unit.DpSizeZero {
		win.Option(app.Size(
			unitconvert.DpToGioUnitUnsafe(cfg.Size.Width),
			unitconvert.DpToGioUnitUnsafe(cfg.Size.Height),
		))
	}

	st := store.NewPersistentState(store.WithRootContext(a.ctx))
	sub := st.Subscribe(func() { win.Invalidate() })
	defer sub.Unsubscribe()

	themeManager := theme.GetThemeManager()   // singleton, per D6

	var ops op.Ops
	for {
		switch e := win.Event().(type) {
		case app.DestroyEvent:
			if e.Err != nil {
				select {
				case a.windowExit <- e.Err:
				default: // first error wins
				}
			}
			return
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			gtx.Locale = localeFor(gtx)                  // en LTR default; see D7
			gtx = themeManager.Material3ThemeInit(gtx).(layout.Context)
			composer := compose.NewComposer(api.ComposerWithStore(st))
			ui := cfg.Content
			// LocalWindow auto-provision (D4)
			if cfg.useLocalWindow {
				ui = compose.CompositionLocalProvider1(
					platform.LocalWindow,
					platform.NewWindow(win),
					ui,
				)
			}
			cmd := runtime.NewRuntime().Run(gtx, composer, ui)
			cmd.Apply(gtx.Ops)          // ← moved inside App; the seam leak dies
			e.Frame(gtx.Ops)
		}
	}
}
```

State model: **one store per window** (matches today's per-demo isolation). The FrameLifecycleHandler GC sweep runs per window-frame; a window's "stop" = App stops calling `StartFrame` for that window (un-remembering everything) — the lifecycle hooks (`OnCleared`, `OnForgotten`, scoped-context cancels) all fire through existing machinery with zero new code.

`store.Subscribe → window.Invalidate()` is the reactive redraw loop; `Subscription.Unsubscribe()` is now called (no demo did before — the App owns the cleanup).

### D4 — `platform.LocalWindow` auto-provisioning

Today only `cmd/demo/fileexplorer` wraps its UI in `CompositionLocalProvider1(platform.LocalWindow, platform.NewWindow(window), UI())`. App provides the window to composition automatically for every window (gated by an option so the rare component that must not see a window can opt out). `compose/ui/platform.Window` is already the opaque seam type (`win any` field); `NewWindow(win)` wraps the `*app.Window`. This is exactly parallel to how `runtime.Run` already auto-provides `LocalDensity`.

### D5 — Multi-window capability gating

gio: "More than one Window is not supported on iOS, Android, WebAssembly" (window.go doc). Design:
- `NewApp()` with >1 `AddWindow` on a non-desktop platform → `Run()` returns a validation error before opening any window (fail fast, no partial state).
- Desktop (Linux/Windows/macOS/BSD): N windows, one goroutine each, one `PersistentState` each — gio's supported model.
- Cross-window state sharing = shared `MutableValue`s / shared Go objects passed to both UIs (gio multiwindow example's pattern; `sync.WaitGroup` broadcast via root context for shutdown).

The store decision intentionally does NOT create a shared-store registry: per-window stores keep GC semantics per window, and cross-window reactivity is achieved by sharing the *values*, not the store.

### D6 — Theme singleton kept as-is (documented limitation)

`theme.GetThemeManager()` is a deprecated global singleton mutating `gtx.Values` per frame. Multi-window means two goroutines may call `Material3ThemeInit` concurrently — it takes `tm.mu.Lock()` internally (theme/manager.go:91), so it is goroutine-safe today; the limitation is *shared token theme* (no per-window theme isolation), not a race. Decision: keep the singleton, document per-window themes as future work. Do NOT bundle theme isolation into this change — it would expand scope into the token system.

### D7 — Locale handling

Demos set `gtx.Locale = system.Locale{Language: "en", Direction: system.LTR}` before theme init. App does the same by default (deterministic en/LTR, matching 44/48 demos); a future `WithLocale` option can extend. Rationale: keep today's behavior identical; locale configurability is a separate concern.

### D8 — `DrawCommand` replay path cleanup

`internal/render/command.go` currently exposes `DrawCommand.Apply(ops *op.Ops)` (gioui leak) alongside unexported `applyTo(b Backend)`. With App owning the loop:
- `Apply(ops *op.Ops)` becomes package-internal (or is replaced by `applyTo` invoked from `runtime`'s gio-backed Apply), so no gioui type appears in any exported method.
- `applyTo(b Backend)` remains the backend-agnostic path used by the software/golden path.
- `runtime.Run` still returns `render.DrawCommand`; only the *consumer* changes (App instead of demo main).

### D9 — Exit code mapping

```go
func exitCode(err error) int {
	switch {
	case err == nil || errors.Is(err, context.Canceled):
		return 0     // clean exit (last window closed, or signal-driven cancel)
	default:
		return 1     // window failure (DestroyEvent.Err)
	}
}
```
Used only on macOS (the `os.Exit` path) and by callers who want to mirror it. Linux/Windows/JS callers do `if err := app.Run(); err != nil { log.Fatal(err) }` and let Go's default exit code flow.

## Architecture Diagram

```
┌──────────────────────────────────────────────────────────────────────────┐
│  main()  [demo main.go — ZERO gio imports]                              │
│                                                                          │
│  runtime.App(                                                            │
│      runtime.Title("Button Demo"),                                       │
│      runtime.Size(unit.NewDpSize(800, 600)),                             │
│      runtime.Content(UI()),                                              │
│  ).Run()  ──► error ──► log.Fatal / errors.Is(err, context.Canceled)     │
└──────────────┬───────────────────────────────────────────────────────────┘
               │
               ▼
┌──────────────────────────────────────────────────────────────────────────┐
│  runtime.App  (runtime/app.go — public, gioui-free signatures)           │
│                                                                          │
│  ┌──────────────┐   ┌──────────────┐   ┌──────────────┐                 │
│  │ window #0    │   │ window #1    │   │ window #N    │  (desktop)      │
│  │ goroutine    │   │ goroutine    │   │ goroutine    │                 │
│  │              │   │              │   │              │                 │
│  │ gio app.Win  │   │ gio app.Win  │   │ gio app.Win  │                 │
│  │ own store    │   │ own store    │   │ own store    │                 │
│  │ theme init   │   │ theme init   │   │ theme init   │                 │
│  │ composer/    │   │ composer/    │   │ composer/    │                 │
│  │ runtime.Run  │   │ runtime.Run  │   │ runtime.Run  │                 │
│  │ LocalWindow  │   │ LocalWindow  │   │ LocalWindow  │                 │
│  │ cmd.Apply    │   │ cmd.Apply    │   │ cmd.Apply    │  ◄─ leak dies   │
│  └──────┬───────┘   └──────┬───────┘   └──────┬───────┘                 │
│         │                  │                  │                          │
│         ▼                  ▼                  ▼                          │
│  ┌─────────────────────────────────────────────────────────────┐        │
│  │  exit monitor:  select {                                   │        │
│  │    first window error ──► exitCh <- err                    │        │
│  │    a.ctx.Done()       ──► exitCh <- ctx.Err()              │        │
│  │    allWindowsClosed   ──► exitCh <- nil (if quitOnLast)    │        │
│  │  }                                                         │        │
│  └─────────────────────────────────────────────────────────────┘        │
│                                                                          │
│  macOS: app.Main() on main goroutine (never returns) + monitor os.Exit   │
│  others: Run() returns <-exitCh                                          │
└──────────────────────────────────────────────────────────────────────────┘
               │
               ▼
   internal/render (seam) — DrawCommand.Apply now package-internal;
   runtime.Run(ctx any, composer, ui) unchanged signature
```

## Edge Cases

### EC1 — First fatal error wins
Two windows fail simultaneously (e.g., GPU init on both): `a.windowExit` is buffered (cap 1) with a `default`-guard non-blocking send; the first error is delivered, later ones dropped. The monitor's select picks whichever terminal event arrives first. Deterministic enough: no lost errors (first one is always delivered), no panic (non-blocking send).

### EC2 — Root context cancelled before Run
`WithRootContext(ctx)` with an already-cancelled ctx: `Run()` still opens windows (frames render), but the monitor's `<-a.ctx.Done()` fires immediately → `Run()` returns `context.Canceled` on the first monitor evaluation. Behavior: windows flash then the app closes. Documented; callers should check `ctx.Err()` before `Run()` if they want to avoid it.

### EC3 — Zero windows registered
`NewApp()` + `Run()` with no `AddWindow`: validation error returned synchronously ("runtime.App: no windows registered") — fail fast rather than block forever on `app.Main()`.

### EC4 — Multi-window on unsupported platform
`NewApp()` + 2 `AddWindow` on JS: validation error before any window opens. Single window on JS works as today (one canvas). Platform check: `runtime.GOOS == "ios" || "android"` or `runtime.GOARCH == "wasm"`.

### EC5 — `QuitWhenLastWindowCloses(false)`
All windows close, no error: monitor does NOT send nil; `Run()` blocks (on Linux/Windows/JS) or stays in the Cocoa loop (macOS). The app must be closed via root-context cancel or an external signal. This is the "document app" model (like gio's multiwindow example staying until Cmd-Q).

### EC6 — Theme init under concurrency
Two window goroutines call `Material3ThemeInit` concurrently: safe today (mutex in themeManager), shared `tokenTheme` means both windows see the same theme. Documented limitation, not a regression — today each demo is one window/one goroutine.

### EC7 — The frame-context gap (Option A boundary)
`runtime.Run(ctx any, ...)` still type-asserts to `layout.Context` (runtime/models.go:25-28 panics otherwise). App always passes a real gio `layout.Context` from `app.NewContext`, so the panic path is unreachable from App. The `any` remains a fig leaf — documented as the accepted Option A trade-off; Phase-B work (input IR) would remove it.

## Risks / Trade-offs

- **[48-file migration churn]** → Mechanical, compiler-driven; each demo's `main.go` shrinks to the App one-liner; golden baselines prove no behavior change; migrate in one commit.
- **[os.Exit inside a library (macOS path)]** → gio's own multiwindow example does exactly this; the Cocoa run loop leaves no alternative on desktop macOS. Confined to the `darwin` build path; documented.
- **[Concurrent theme singleton]** → Goroutine-safe today (mutex); only isolation is missing. Documented non-goal; revisit when per-window themes matter.
- **[`QuitWhenLastWindowCloses(false)` can look like a hang]** → Intentional: document-app model. Root-context cancel remains the universal shutdown path.
- **[macOS divergence in Run() error contract]** → `Run()` never returns on desktop macOS. Callers wanting uniform error handling must use the root-context path (which DOES propagate `context.Canceled` even on macOS via... the monitor's os.Exit code 0 — acceptable; document).
- **[Scope creep into Phase B (input IR)]** → Explicit non-goal; the Option A boundary (frame-context gap) is stated in Context and Non-Goals so implementers don't wander.

## Migration Plan

1. **Phase 1 — seam cleanup**: make `DrawCommand.Apply` package-internal; verify `make test` + `make check-api` green.
2. **Phase 2 — `runtime.App` core**: `app.go` (public types + options), `window.go` (per-window driver), `Run()` with platform-aware exit; unit tests with an injected fake window-driver where feasible.
3. **Phase 3 — demo migration**: rewrite all 48 `cmd/demo/*/main.go` to `runtime.App(...).Run()`; delete the gio setup block; update `cmd/demo/kitchen` (root context) and `cmd/demo/fileexplorer` (LocalWindow now auto-provided — remove manual wrapper).
4. **Phase 4 — verification**: `make test` green (golden baselines unchanged), `make check-api` green, `go build ./...`, `go run ./cmd/demo/kitchen/` smoke test.
5. **Rollback**: revert the demo-migration commit (tree returns to per-demo setup); runtime.App commits are independently revertible. No behavior change to rendering (golden baselines prove it).
6. **Migration notes**: `docs/runtime_app_migration.md` — before/after for a typical demo main, the new `Run() error` contract, macOS caveat, `QuitWhenLastWindowCloses` semantics.

## Open Questions

1. **`useLocalWindow` opt-out**: should LocalWindow provision be default-on for all windows, or opt-in (current fileexplorer behavior)? Default-on matches `LocalDensity` precedent; opt-in is safer for exotic components. Lean: default-on, with a `WithoutLocalWindow` per-window flag if ever needed (YAGNI — add when requested).
2. **`WithLocale` option**: default en/LTR today; add `WithLocale(unit.Locale)` now or later? Lean: later (out of scope; 4/48 demos already skip locale).
3. **`App.Run` on Android/iOS**: platform creates the single window and `app.Main()` returns immediately — does App need an Android-specific loop entry (linkname pattern) in this change, or is desktop+WASM coverage sufficient? Lean: desktop+WASM now; mobile entry documented as follow-up.
