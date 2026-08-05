# Runtime App Specification

## Purpose

Define the behavioral contract for `runtime.App` — the framework-owned application entry that owns the platform window and event loop so application code (demo mains, user apps) does not import engine types during setup.

## Requirements

### Requirement: Application setup requires no engine imports

The system SHALL provide a `runtime.App` entry point such that application setup — window creation, event loop, per-window state store, theme initialization, and frame application — requires no direct engine (gioui) imports in application code. All public `runtime.App` signatures SHALL be free of engine types.

#### Scenario: Single-window demo main without engine imports
- **WHEN** a demo `main.go` constructs `runtime.App(runtime.Title("Demo"), runtime.Size(unit.NewDpSize(800, 600)), runtime.Content(UI()))` and calls `Run()`
- **THEN** the file compiles with zero `gioui.org` imports and the window opens with the given title and size

#### Scenario: Public API purity check passes
- **WHEN** `make check-api` runs against the new `runtime.App` public surface
- **THEN** no exported signature of `runtime.App`, its options, or `WindowOptions` references a `gioui.org` type

### Requirement: Run blocks until close and reports the close cause

The system SHALL provide `App.Run() error` that blocks until the application closes and reports why it closed: `nil` when the last window closes normally, `context.Canceled` when the root context is cancelled, or the `DestroyEvent.Err` of a window that failed. The system SHALL return a validation error when no windows are registered.

#### Scenario: Last window closes normally
- **WHEN** the application has one window and the user closes it
- **THEN** `Run()` returns `nil`

#### Scenario: Root context cancelled
- **WHEN** the application's root context (from `WithRootContext`) is cancelled
- **THEN** `Run()` returns `context.Canceled`

#### Scenario: Window failure propagates
- **WHEN** a window is destroyed prematurely with an error (e.g., GPU/driver failure)
- **THEN** `Run()` returns that window's `DestroyEvent.Err`

#### Scenario: No windows registered
- **WHEN** `NewApp()` is constructed with no `AddWindow` calls and `Run()` is invoked
- **THEN** `Run()` returns a validation error without opening any window

### Requirement: Multi-window support on desktop

The system SHALL support multiple windows on desktop platforms (Linux, Windows, macOS), each with its own event loop and its own per-window state store. The system SHALL reject multi-window configuration on platforms where the engine does not support it (iOS, Android, WebAssembly) with a validation error before opening any window.

#### Scenario: Two windows open on desktop
- **WHEN** `NewApp()` registers two windows via `AddWindow` and `Run()` is invoked on a desktop platform
- **THEN** both windows open, each renders its own content, and `Run()` returns `nil` when both are closed

#### Scenario: Multi-window rejected on unsupported platform
- **WHEN** `NewApp()` registers more than one window via `AddWindow` and `Run()` is invoked on WebAssembly
- **THEN** `Run()` returns a validation error and no window opens

#### Scenario: Per-window state isolation
- **WHEN** two windows each remember state under the same key
- **THEN** each window observes its own value; changes in one window do not affect the other

### Requirement: Store and invalidation loop are framework-owned

The system SHALL create the per-window persistent state store inside `runtime.App` and SHALL wire state changes to window invalidation (redraw) without application code involvement.

#### Scenario: State change triggers redraw
- **WHEN** application code mutates a remembered `MutableValue` in a window's store
- **THEN** the owning window is invalidated and a new frame is rendered without any application-side subscription code

#### Scenario: Frame lifecycle is bracketed
- **WHEN** a frame renders through `runtime.App`
- **THEN** the window store's `StartFrame`/`EndFrame` bracket every frame, so values left un-remembered are cleaned up via the existing lifecycle machinery

### Requirement: Window is provided to composition

The system SHALL provide the platform window to the composition tree via the framework-owned `platform.LocalWindow` composition local for every App-created window, so components can access window services without engine imports.

#### Scenario: Component reads LocalWindow
- **WHEN** a component reads `platform.LocalWindow.Current(c)` inside an App-created window
- **THEN** it receives a non-nil framework `*platform.Window` wrapping the live platform window

### Requirement: Exit behavior is configurable

The system SHALL provide a `QuitWhenLastWindowCloses` option (default true). When false, the application SHALL NOT exit when the last window closes and SHALL remain alive until the root context is cancelled or an external signal terminates it.

#### Scenario: Quit disabled keeps app alive
- **WHEN** `QuitWhenLastWindowCloses(false)` is set, the last window closes, and no error occurred
- **THEN** `Run()` does not return and the process remains alive

#### Scenario: Quit disabled still responds to root context
- **WHEN** `QuitWhenLastWindowCloses(false)` is set and the root context is cancelled
- **THEN** `Run()` returns `context.Canceled`

### Requirement: Frame application requires no engine types in application code

The system SHALL apply the frame's draw command inside `runtime.App`; application code SHALL NOT call any engine-typed apply method. The seam's `DrawCommand` SHALL NOT expose an engine-typed apply method in its public interface.

#### Scenario: DrawCommand has no engine-typed apply
- **WHEN** inspecting the exported surface of `internal/render.DrawCommand`
- **THEN** no exported method accepts a `*op.Ops` (or other engine type) parameter
