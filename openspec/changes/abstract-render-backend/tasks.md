# Tasks: Abstract Render Backend

Reference: design.md (how), specs/{public-api-purity,render-backend-seam,api-purity-analyzer,golden-testing}/spec.md (contracts).

## 1. API-Purity Analyzer (fixtures first — red state)

- [x] 1.1 Scaffold `internal/api-check` package: load public packages (`compose/...`, `modifiers/...`, `theme/...`, `runtime/`, `pkg/...`) via `go/packages` (uses `golang.org/x/tools`, already a dependency)
- [x] 1.2 Implement exported-signature inspector: params, returns, struct fields, embedded types, type aliases, and underlying types of defined types; flag any type whose package path starts with `gioui.org`
- [x] 1.3 Implement seam whitelist: `internal/layoutnode`, `internal/render`, and subpackages are exempt from the import rule but still subject to the exported-signature rule
- [x] 1.4 Seed regression fixtures from the known leak inventory: every pre-cleanup leak site (unit gio_helpers, text/font converters, theme aliases, box/column aliases, intl.Locale, padding.RTL, card.GioClickable, clickable.GioClickable, shape.Outline, LocalWindow, TextFieldWidget, LazyListState/GridState, SnackbarData, runtime.Run) must be flagged
- [x] 1.5 Add `make check-api` target invoking the analyzer; document exit behavior (non-zero on violation)
- [x] 1.6 Verify analyzer is in RED state: flags every fixture, exits non-zero (spec: api-purity-analyzer — "Known leaks are regression fixtures")

## 2. Phase 1: Public API Boundary Cleanup

- [x] 2.1 Create `internal/render` stub package: opaque `DrawCommand` type + `CallOp` type (no behavior yet); add to analyzer whitelist
- [x] 2.2 Redefine `style.TextAlign` and `style.LineBreak` as go-compose-owned enum types with own constants; move `TextAlignToGio*`/`LineBreakToGio*`/`FromGio*` conversions to `internal/` (spec: public-api-purity — "Text enum types are distinct", "Engine conversion helpers are not public")
- [x] 2.3 Redefine `intl.Locale` as own type (drop `system.Locale` alias) and `padding.RTL` as own constant; update `compose/ui/text/intl` + `compose/ui/next/text/intl`
- [x] 2.4 Redefine `theme.BasicTheme` as framework-owned theme type; refactor `theme.ThemeManager` methods to not use `*material.Theme`
- [x] 2.5 Replace `compose/material` `ThemeInterface.GioMaterialTheme()` and `LocalGioMaterialTheme` with a framework-owned theme interface/provider (no `*gioMaterial.Theme` in signatures)
- [x] 2.6 Replace `box.Direction`/`box.Stack`/`box.StackChild` and `column.Spacing`/`column.Alignment` with go-compose-owned types + mirrored constants (`box.NW`, `box.NE`, …); move gioui conversions to the seam
- [x] 2.7 Replace `card.GioImage` and `clickable.GioClickable` aliases with go-compose-owned types
- [x] 2.8 Move `compose/ui/unit/gio_helpers.go` conversions (`DpToGioUnit*`, `TextUnitToGio*`) into `internal/`; remove exported `AsGioSp` from `TextUnit`
- [x] 2.9 Move text/font conversions (`ToGioFont`, `FromGioFont`, `ToGioWeight`, `ToGioStyle`, `FromGioTypeface`, `TextStyleFromGioFont`, `ToGioFont(ts *TextStyle)`) out of `compose/ui/text{,/font,style}` into `internal/`
- [x] 2.10 Redefine `shape.Shape`/`shape.Outline` so no method returns `clip.Stack`/`clip.Op`/`clip.PathSpec`; introduce framework-owned outline types with seam conversion
- [x] 2.11 Replace `platform.LocalWindow` (`*app.Window`) with an opaque framework-owned window type
- [x] 2.12 Reshape `TextFieldWidget`/`FilledTextFieldWidget`, `LazyListState`/`LazyGridState`, and `SnackbarData` to drop embedded/field gioui types (`widget.Editor`, `gesture.Click`, `widget.List`, `widget.Clickable`); expose framework-owned state types
- [x] 2.13 Change `runtime.Runtime.Run` to return `render.DrawCommand` (opaque) instead of `op.CallOp`; implement the gioui application path internally
- [x] 2.14 Migrate all in-repo consumers: `cmd/demo/*/main.go`, `cmd/go-compose`, `pkg/x/fileexplorer` (`RememberExplorer` no longer takes `*explorer.Explorer`), tests
- [x] 2.15 Run analyzer to GREEN: all former fixtures pass, exit zero (spec: api-purity-analyzer — "Fixture coverage after cleanup", public-api-purity)

## 3. Phase 2: Alias Wall → Defined Types

- [x] 3.1 Define `LayoutContext` as a framework-owned defined type (no gioui fields); introduce per-frame construction inside the runtime
- [x] 3.2 Define `LayoutDimensions`, `LayoutConstraints`, and `GioLayoutWidget` as defined types (not aliases)
- [x] 3.3 Introduce `ToGio()` conversions in `internal/layoutnode` (seam-only, frame-token guarded per design EC1)
- [x] 3.4 Convert `internal/layoutnode` internals (`models.go`, `constructor.go`, `widget_models.go`, `layout_node.go`) to the defined types
- [x] 3.5 Convert `modifiers/*` node.go files package-by-package (compiler-driven; keep tree green between packages)
- [x] 3.6 Convert `compose/foundation` (layout/box, column, row, overlay, text, lazy, image, icon) to the defined types
- [x] 3.7 Convert `compose/material3` components; where a gioui-bound component is cheaper to rewrite as pure composition, do so following the `next/` pattern
- [x] 3.8 Classify remaining engine-hungry components (slider, textfield, radiobutton, progress, badge, `next/text` editor) as engine-bound (Door 2) with explicit code comments; their public APIs stay clean
- [~] 3.9 Ensure non-seam packages no longer directly import `gioui.org/layout`, `op`, `widget`, `text`, `font` (analyzer import rule green; spec: render-backend-seam — "Seam packages are the only engine-type importers")
  - INTERPRETATION: the import rule (Rule B) is implemented in the analyzer (`-all` mode, seam whitelist) but NOT enforced by `make check-api`, which gates on Rule A (public signature purity). Engine-bound components (slider, textfield, radiobutton, progress, badge, next/text editor) are classified Door 2 and legitimately import engine packages internally (design D3: "internal implementation files should use gioui"). Enforcing Rule B repo-wide would contradict Door 2; `make check-api -all` remains available as a stricter audit.
- [x] 3.10 Full `make test` green; `make check-api` green (spec: render-backend-seam — "Framework context is engine-agnostic", "Conversion to engine types is confined to the seam")

## 4. Phase 3: Draw IR + Software Backend

- [ ] 4.1 Audit the gioui op vocabulary used by modifier nodes (op.Record/Stop, op.Offset, op.Affine, paint.ColorOp/Fill/FillShape/PushOpacity, clip.Rect/Path/Stroke, layout.Background/Inset) and derive the `render.Backend` interface
- [ ] 4.2 Implement `internal/render`: `Ops`, `CallOp`, `DrawCommand`, backend registration (blank-import pattern, design D6)
- [ ] 4.3 Implement the gioui backend: maps `render.Backend` calls to `op.Ops`/`op/clip`/`op/paint`; preserves macro record/replay semantics
- [ ] 4.4 Implement the software backend: ordered `DrawCall` list (deterministic — no map iteration, no wall-clock, injectable `timeNowFunc`) plus optional `image.RGBA` rasterization
- [ ] 4.5 Wire `runtime.Run` and the node coordinator through the backend; preserve discard semantics for the pointer pass (recorded-but-dropped recordings emit no draw calls — spec: golden-testing — "Discarded recordings emit no draw calls")
- [ ] 4.6 Add golden-test harness: render covered components through the software backend, compare against committed baselines (spec: golden-testing — "Golden baselines are committed", "Refactor preserves rendering behavior")
- [ ] 4.7 Generate initial golden baselines for covered components and commit them
- [ ] 4.8 Full `make test` green (including golden runs) and `make check-api` green (spec: golden-testing — "Software backend captures draw calls", "Software backend output is deterministic", render-backend-seam — "Backends are registrable", "Runtime interface is engine-agnostic")

## 5. Final Verification

- [ ] 5.1 Verify no `gioui.org` type appears in any exported signature of `compose/`, `modifiers/`, `theme/`, `runtime/`, `pkg/` (analyzer + manual spot-check of the former leak inventory)
- [ ] 5.2 Verify all demos still build and run (`go run ./cmd/demo/kitchen/`, `go build ./...`)
- [ ] 5.3 Commit in phase-sized commits (analyzer + Phase 1, Phase 2, Phase 3) with migration notes for the breaking API changes
