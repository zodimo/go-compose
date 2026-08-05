# Tasks: Runtime App

Reference: design.md (how), specs/runtime-app/spec.md (contracts).

## 1. Phase 1: Seam Cleanup (DrawCommand leak)

- [x] 1.1 Make `internal/render.DrawCommand.Apply(ops *op.Ops)` package-internal or replace with the unexported `applyTo(b Backend)` path (spec: "Frame application requires no engine types in application code" — no exported method accepts `*op.Ops`)
- [x] 1.2 Update the only consumer of the public `Apply` (currently the demo mains via `cmd.Apply(gtx.Ops)`); keep a package-private gio replay helper in `internal/render` (e.g. `ApplyToGio(cmd DrawCommand, ops *op.Ops)`) for runtime/ to call
- [x] 1.3 `make test` green and `make check-api` green after the seam change

## 2. Phase 2: runtime.App Core

- [x] 2.1 Add `runtime/app.go`: `Application` struct, `App(...)` convenience constructor, `NewApp(...)`, `AddWindow(WindowOptions)`, `Run() error` (public surface gioui-free per api-purity rules)
- [x] 2.2 Add `runtime/app_options.go`: `AppOption`, `AppConfig`, `Title`, `Size(unit.DpSize)`, `Content(api.Composable)`, `WithRootContext(context.Context)`, `QuitWhenLastWindowCloses(bool)` (default true); `WindowOptions{Title, Size, Content}`
- [x] 2.3 Implement the per-window driver `runWindow` (design D3): create gio `app.Window` with Title/Size (via `internal/unitconvert.DpToGioUnitUnsafe`), per-window `store.NewPersistentState(WithRootContext(a.ctx))`, `store.Subscribe → win.Invalidate()` (with `Unsubscribe` on exit), theme init via `theme.GetThemeManager()`, composer per frame, `runtime.Run(...)`, apply, `e.Frame(gtx.Ops)`
- [x] 2.4 Implement platform-aware `Run()` exit semantics (design D2): exit monitor select on `{first window error, ctx.Done, allWindowsClosed}`; macOS path (`app.Main()` + monitor `os.Exit` with design D9 exit-code mapping); Linux/Windows/JS path returning `<-exitCh`
- [x] 2.5 Add validation: zero windows registered → error before `app.Main()` (spec: "No windows registered"); >1 window on iOS/Android/WASM → error before opening any window (spec: "Multi-window rejected on unsupported platform")
- [x] 2.6 Auto-provision `platform.LocalWindow` in the frame driver (spec: "Window is provided to composition"); wrap UI with `CompositionLocalProvider1(platform.LocalWindow, platform.NewWindow(win), ui)` per window
- [x] 2.7 En locale default in frame driver (design D7): set `gtx.Locale = system.Locale{Language: "en", Direction: system.LTR}` before theme init
- [x] 2.8 Unit tests for `Run()` semantics where feasible with a fake window driver seam (exit-on-last-close nil, ctx cancel → context.Canceled, first-error-wins, validation errors, quit-disabled behavior)
- [x] 2.9 `make test` green, `make check-api` green (spec: "Public API purity check passes")

## 3. Phase 3: Demo Migration

- [x] 3.1 Migrate the canonical demo (`cmd/demo/main.go`) to `runtime.App(runtime.Title(...), runtime.Size(unit.NewDpSize(...)), runtime.Content(UI())).Run()`; delete the gio setup block; verify `go run ./cmd/demo/main.go`-equivalent renders
- [x] 3.2 Migrate the remaining 47 demo `main.go` files to the same pattern (mechanical; zero gioui imports remaining in any demo main)
- [x] 3.3 Update `cmd/demo/kitchen/main.go`: replace signal/context/lifecycle boilerplate with `WithRootContext`; keep `QuitWhenLastWindowCloses` semantics matching today's clean-exit-on-signal behavior
- [x] 3.4 Update `cmd/demo/fileexplorer/main.go`: remove the manual `CompositionLocalProvider1(platform.LocalWindow, ...)` wrapper (App now auto-provides it; spec: "Window is provided to composition")
- [x] 3.5 Handle the 4 non-standard demos (`segmentedbutton` raw-int size, `effect`/`stateflow` continuous `window.Invalidate()`, `switch` no-locale) — verify behavior preserved or document intentional change
- [x] 3.6 Verify the 3 flow demos (`flow_shared`, `flow_state`, `flow_operators`) are untouched (console-only, no gio)
- [x] 3.7 Confirm `go build ./...` passes and `go vet ./...` is clean

## 4. Phase 4: Verification + Docs

- [x] 4.1 `make test` full green — golden baselines unchanged (proves no rendering behavior change)
- [x] 4.2 `make check-api` green — public API purity preserved (spec: "Application setup requires no engine imports" scenario "Public API purity check passes")
- [x] 4.3 Smoke test `go run ./cmd/demo/kitchen/` (desktop) and `go run ./cmd/go-compose serve ./cmd/demo/kitchen/` (WASM)
- [x] 4.4 Write `docs/runtime_app_migration.md`: before/after demo main, `Run() error` contract (nil / context.Canceled / DestroyEvent.Err), macOS caveat (Run never returns; monitor os.Exit), `QuitWhenLastWindowCloses` semantics, multi-window example
- [x] 4.5 Commit in phase-sized commits (Phase 1 seam, Phase 2 core, Phase 3 demos) with migration notes
