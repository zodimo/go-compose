## Why

The render-backend seam made the framework's *rendering* surface gioui-free, but every demo app (`cmd/demo/*/main.go`) still hand-rolls ~30 lines of gio setup: window creation, the `Event()` loop, theme bootstrap, store wiring, and `cmd.Apply(gtx.Ops)` — the one place a gioui type (`*op.Ops`) still leaks through the seam's own return path. New users' first encounter with the framework is gio's window/event model, and 48 files × 5 gioui imports means every gio breaking change costs 48 edits.

## What Changes

- **NEW** `runtime.App` entry point that owns the gio window + event loop, so application code no longer imports gioui during setup.
- **NEW** `App.Run() error` — blocks until close; `nil` on last-window-close, `context.Canceled` on root-context cancel, propagates `DestroyEvent.Err` failures.
- **NEW** Multi-window support via `NewApp()` + `AddWindow()` (desktop; single-window enforced on iOS/Android/WASM per gio limits).
- **NEW** `runtime.App` options: `Title(string)`, `Size(unit.DpSize)`, `Content(api.Composable)`, `WithRootContext(context.Context)`.
- **NEW** Auto-provisioning of `platform.LocalWindow` in the App frame driver (currently only the fileexplorer demo does this manually).
- **REMOVE** the `DrawCommand.Apply(ops *op.Ops)` gioui leak — replay moves inside App; the seam method dies.
- **CHANGE** 48 demo `main.go` files: gio imports removed, replaced by the `runtime.App` one-liner.
- **NON-GOAL (documented)**: input IR / text shaping / backend-carried frame context (Option A scope). The software/headless path and golden harness keep their own entry.
- **NON-GOAL (documented)**: per-window theme isolation — `theme.GetThemeManager()` singleton is kept as-is.

## Capabilities

### New Capabilities
- `runtime-app`: The `runtime.App` application entry — window lifecycle ownership, event-loop encapsulation, multi-window support, exit/error semantics, and `platform.LocalWindow` provisioning.

### Modified Capabilities
<!-- No existing spec requirements change: render-backend-seam's DrawCommand.Apply is
     internal/render implementation, not a spec-level requirement. -->

## Impact

- **`runtime/`** — new `app.go` (public API), `window.go` (per-window driver), `options.go` (expanded); Door 2 pattern: public signatures gioui-free, impl files may import gio (same as `theme/manager.go`).
- **`internal/render/command.go`** — `DrawCommand.Apply(ops *op.Ops)` removed or unexported; `applyTo(b Backend)` becomes the only replay path.
- **`compose/ui/platform/local_window.go`** — used (unchanged); `LocalWindow` now provided by App.
- **`store/`** — `WithRootContext` + per-window `PersistentState` created by App; `store.Subscribe → window.Invalidate()` moves inside App. Optionally uses the currently-unused `StartFrameTriggerReceiver`/`EndFrameTriggerReceiver` seam.
- **`cmd/demo/*/main.go`** (48 files) — setup block replaced with `runtime.App(...).Run()`; gioui imports removed.
- **`cmd/demo/kitchen/main.go`** — context/lifecycle pattern moves into `WithRootContext`; `fileexplorer` keeps `LocalWindow` usage via App auto-provision.
- **`pkg/api`** — no change; `api.Composable` already the content type.
- **gioui**: unchanged dependency; `app.Main()`, `app.Window`, `unit.Dp` consumed only inside `runtime/`.
- **API purity**: `make check-api` must stay green (new public signatures contain no gioui types).
