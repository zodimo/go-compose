## Why

gioui.org types leak into go-compose's public API. `runtime.Run` returns `op.CallOp`, `theme.BasicTheme` is an alias for `material.Theme`, `box.Direction` aliases `layout.Direction`, conversion helpers like `DpToGioUnit` and `ToGioFont` are exported from `compose/ui/unit` and `compose/ui/text`, and structs like `TextFieldWidget` embed `widget.Editor` directly. The project's identity is a declarative Jetpack Compose Material 3 API in Go — gioui is the engine that powers it, not the product. Today a consumer cannot adopt the framework without learning gioui's type system, and the engine cannot be swapped or tested headlessly without breaking the API.

The fix is to make the render engine an implementation detail: the public API surface must be gioui-free (signature-level purity), and the engine boundary must be a real seam (`internal/`) that a second backend — starting with a software/recording backend for golden tests — can implement.

## What Changes

- **Hides gioui behind the public API.** Exported aliases to gioui types become go-compose-owned types: `theme.BasicTheme`, `box.Direction`, `column.Spacing`/`column.Alignment`, `intl.Locale`, `padding.RTL`, `card.GioImage`, `clickable.GioClickable`. **BREAKING** — these types change identity; code depending on their gioui behavior must migrate.
- **Removes gioui from exported signatures.** `runtime.Run` returns an opaque draw command instead of `op.CallOp`; `shape.Outline` no longer returns `clip.Stack`/`clip.Op`/`clip.PathSpec`; `LocalWindow` no longer exposes `*app.Window`; `ThemeInterface.GioMaterialTheme()` is replaced. **BREAKING** — function signatures change.
- **Moves engine conversion helpers out of public packages.** `DpToGioUnit*`, `TextUnitToGio*`, `AsGioSp`, `ToGioFont`/`FromGio*`, `TextAlignToGio*`, `LineBreakToGio*` relocate to `internal/` (or are deleted). **BREAKING** — these exported functions are removed from `compose/ui/unit`, `compose/ui/text`, `compose/ui/text/style`, `compose/ui/text/font`.
- **Redefines enums as go-compose-owned types.** `style.TextAlign` and `style.LineBreak` stop using gioui types (`text.Alignment`, `text.WrapPolicy`) as their underlying type.
- **Converts the internal alias wall into defined types.** `internal/layoutnode`: `LayoutContext`, `LayoutDimensions`, `LayoutConstraints`, `DrawOp`, `GioLayoutWidget` become defined types with `ToGio()`-style conversions confined to the seam. No public API change here, but the internal boundary becomes real.
- **Introduces an engine registration seam.** A `internal/render` backend interface (2D draw IR, modeled on Compose's `Canvas` / Gio's `gpu.GPU`) with a gioui implementation and a software/recording implementation for golden tests.
- **Adds tooling that enforces purity.** An API-purity analyzer (a test or `go/ast`-based checker run via `make check-api`) fails when any `gioui.org` type appears in an exported signature of the public packages, with a whitelist for the seam packages.
- **Delegates text shaping to the engine.** `compose/foundation/next/text` stays engine-side behind the seam; only its public types get signature-cleaned. No attempt to re-abstract the text shaper/editor.

## Capabilities

### New Capabilities

- `public-api-purity`: Contract that no exported symbol in `compose/`, `modifiers/`, `theme/`, `runtime/`, or `pkg/` references a `gioui.org` type in its signature (parameters, returns, fields, embedded types, type aliases/underlying types). Behavioral, verifiable contract for the whole public surface.
- `render-backend-seam`: Contract describing the internal seam: `internal/layoutnode` + `internal/render` are the only packages that may import engine types (`gioui.org/layout`, `op`, `widget`, `text`, `font`); framework-owned types flow above the seam; conversion happens at the seam only; a second backend (software/recording) can render the same composable tree without framework changes.
- `api-purity-analyzer`: Contract for the enforcement tooling: how it is invoked (`make check-api`), what it reports, exit behavior, and the whitelist mechanism for seam packages.
- `golden-testing`: Contract that a software/recording backend captures draw-call lists from any composable, enabling golden/snapshot regression tests, and that existing components produce identical output pre- and post-refactor.

### Modified Capabilities

<!-- No existing specs; this is the first OpenSpec change in the repository. -->

## Impact

- **Public API (breaking):** `compose/`, `modifiers/`, `theme/`, `runtime/`, `pkg/` — type identity changes (aliases → owned types), removed conversion helpers, changed signatures (`runtime.Run`, `shape.Outline`, `LocalWindow`, `ThemeInterface`).
- **Internal seam:** `internal/layoutnode` (alias wall → defined types), new `internal/render` (backend interface + gioui + software implementations).
- **Components:** All Material 3 / foundation components transitively affected through `LayoutContext` type change (compiler-driven migration). Custom renderers (slider, textfield, radiobutton, progress, badge) keep gioui internally behind the seam in this change; only their public signatures change.
- **Text:** `compose/foundation/next/text` internals untouched (engine-delegated); public text types (`TextStyle`, `TextUnit`, `TextAlign`, `LineBreak`, font types) cleaned.
- **Tooling/build:** new `make check-api` target; `Makefile` gains an API-purity check; no CI currently exists (AGENTS.md), so enforcement is local + optional CI hook.
- **Entry points:** demo `main.go` event loops change where they consume `runtime.Run`'s new return type; gioui `app` imports in demos remain acceptable (app glue, not compose API).
- **Dependencies:** no new external dependencies; `golang.org/x/tools` (already a dependency) used for the analyzer.
