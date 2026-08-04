# Design: Abstract Render Backend

## Context

go-compose is a Go reimplementation of Jetpack Compose's Material 3 API built on gioui.org. The public API is `Composable func(Composer) Composer`; composition is driven by a zipper-style `Composer` over a `LayoutNode` tree; layout/draw is delegated to gioui's `layout.Context` and `op.Ops`; the app entry point is gioui's `app.Window` event loop.

The engine is currently welded into the framework at every layer. 189 files import `gioui.org` (549 import lines). The critical mechanical facts:

1. **The alias wall** (`internal/layoutnode/alias.go`) — the project *already has* a chokepoint, but it is transparent:
   ```go
   type LayoutContext   = layout.Context     // gioui type alias
   type LayoutDimensions = layout.Dimensions // gioui type alias
   type LayoutConstraints= layout.Constraints
   type GioLayoutWidget = layout.Widget      // func(layout.Context) layout.Dimensions
   type DrawOp          = op.CallOp
   ```
   Go type aliases (`=`) preserve identity: `LayoutContext` **is** `layout.Context`. Every consumer sees gioui types. Every modifier, every component, the runtime, and the demos flow through this wall.

2. **`runtime.Run` returns `op.CallOp`** (`runtime/runtime.go`) — the runtime interface is gioui-typed in its return value.

3. **Exported aliases and conversions leak gioui into the public API**:
   - `theme.BasicTheme = material.Theme`, `theme.ThemeManager` methods use `*material.Theme`
   - `compose/material/theme.go`: `ThemeInterface.GioMaterialTheme() *gioMaterial.Theme`, `LocalGioMaterialTheme`
   - `box.Direction/Stack/StackChild`, `column.Spacing/Alignment` alias `layout.*`
   - `intl.Locale = system.Locale`, `padding.RTL = system.RTL`
   - `card.GioImage = widget.Image`, `clickable.GioClickable = widget.Clickable`
   - `compose/ui/unit/gio_helpers.go`: exported `DpToGioUnit*`, `TextUnitToGio*` returning `unit.Dp/unit.Sp`
   - `compose/ui/text/{,font,style}`: `ToGioFont`, `FromGioFont`, `AsGioSp`, `TextAlignToGio*`, `LineBreakToGio*`, and `style.TextAlign`/`style.LineBreak` typed on gioui enums
   - `shape.Outline` returns `clip.Stack`/`clip.Op`/`clip.PathSpec`
   - `platform.LocalWindow = *app.Window`
   - `TextFieldWidget` embeds `widget.Editor` + `gesture.Click`; `LazyListState.List` embeds `widget.List`; `SnackbarData` holds `widget.Clickable` fields

4. **The frame pipeline today** (7 stages, from `cmd/demo/*/main.go` + `runtime/models.go`):
   ```
   demo main loop:
     window.Event() → app.FrameEvent
       gtx := app.NewContext(&ops, frameEvent)     // gioui layout.Context, per-frame
       composer := compose.NewComposer(store)
       callOp := runtime.Run(gtx, composer, UI())  // returns op.CallOp
       callOp.Add(gtx.Ops)                          // replay recorded draw macro
       frameEvent.Frame(gtx.Ops)                    // submit to renderer
   runtime.Run:
     composer.StartFrame(); defer composer.EndFrame()
     node := ui(composer).Build()                   // zipper builds LayoutNode tree
     coord := NewNodeCoordinator(node)              // wrap in modifier chain
     macro := op.Record(gtx.Ops)                    // record layout
     coord.Layout(gtx)                              // expand modifiers, measure, place
     coord.PointerPhase(gtx)                        // input pass (records, discards)
     _ = macro.Stop()
     return coord.Draw(gtx)                         // record draw ops → op.CallOp
   ```

**The composition zipper — the engine-agnostic heart.** The `Composer` is a zipper over the `LayoutNode` tree: it holds a `focus` pointer (the node currently being composed) and a `path` stack (the route back to the root). As user composables run, they navigate the tree:

- `c.StartBlock(key)` pushes the current focus onto the path stack and descends into a child node;
- modifier calls (`c.Modifier(...)`) and widget-constructor calls (`c.SetWidgetConstructor(...)`) attach data to the focus node;
- `c.EndBlock()` pops back up to the parent, and the function returns `c` so the caller can continue composing siblings.

`Build()` replays the path stack to the root and returns the complete `LayoutNode` tree. Critically, **the entire tree is rebuilt from scratch on every frame**: composition is not retained between frames — only the `PersistentState` slots (keyed by `Identifier`, managed by the frame-lifecycle-aware store) survive. This is the same model as Jetpack Compose's recomposition, with one difference: Compose skips recomposition via keying diffs and `remember` invalidation, whereas this implementation re-runs the composables each frame and relies on the store to preserve state across rebuilds.

**Why the zipper matters for this design:** it means framework state that must live *across* frames lives in the store — never in the composer, the tree, or anything attached to them. The tree, and therefore every context and coordinator derived from it, is per-frame *by construction*. This is the invariant EC1's frame-token guard extends; see EC1.

**The modifier chain and nodeCoordinator — the Phase-2 ripple surface.** When `Build()` returns the tree, `NewNodeCoordinator(node)` (`internal/layoutnode/constructor.go`) recursively wraps every `LayoutNode` in a `nodeCoordinator`. Each coordinator holds a `layoutCallChain` — initially the node's own widget constructor (`GioLayoutWidget`, i.e. `func(layout.Context) layout.Dimensions`). Then `Expand()` walks the node's modifier chain and, for each modifier, dispatches on its phase:

- `AttachLayoutModifier` — the modifier wraps the current chain with its layout behavior (e.g. `size` records the child, measures it, emits an `op.Offset`, and returns adjusted dimensions);
- `AttachDrawModifier` — the modifier wraps the chain with draw behavior (e.g. `background` records the child's draw, fills behind it, replays the child);
- `AttachPointerInputModifier` — the modifier adds a pointer-filter wrapper for input hit-testing.

The result is a nested sandwich: each modifier contributes an outer layer that calls the inner `layoutCallChain.Layout(gtx)` and post-processes the returned `LayoutDimensions`. `Layout()` then invokes the innermost widget, and the `LayoutDimensions` thread outward through every wrapper. This is why a single `LayoutContext`/`LayoutDimensions` type change in Phase 2 propagates to *every* modifier node and *every* component — they all participate in the sandwich, and the compiler walks the full closure. This is also where op-recording happens: `Draw()` re-invokes the chain inside `op.Record(...)/macro.Stop()` and returns the `op.CallOp` — precisely the call site Phase 3 replaces with the backend's `Record()`.

The state layer (`state/`, `pkg/flow/`) is 100% engine-free (0 gioui imports in 31 files). ~half of Material 3 components (card, dialog, snackbar, tab, `next/button`) are already pure composition with zero gioui imports. The `next/` directory is the project's existing de-gioification pattern: rewrite a component as pure composition over primitives, leaving gioui at the primitives.

## Goals / Non-Goals

**Goals:**

- **Signature-level API purity** (user decision Q3): no `gioui.org` type appears in any exported symbol's signature (params, returns, fields, embedded types, type aliases, underlying types) in `compose/`, `modifiers/`, `theme/`, `runtime/`, `pkg/`.
- **A real internal seam**: `internal/layoutnode` + `internal/render` become the only packages that import engine types (`gioui.org/layout`, `op`, `widget`, `text`, `font`); everything above consumes framework-owned types.
- **A second backend that proves the seam**: a software/recording backend producing deterministic draw-call lists → golden regression tests (user decision Q1-c; this is the concrete "second consumer" that justifies the abstraction).
- **Enforcement tooling**: `make check-api` fails on any leaked signature, with a whitelist for seam packages.
- **Delegated text** (user decision Q2): `compose/foundation/next/text` stays engine-side behind the seam; only its public types are signature-cleaned.
- **Engine swap becomes a new-backend exercise** (user decision Q1-b, "after"): the design must not paint us into a corner, but building a non-gio backend is NOT in this change.

**Non-Goals:**

- Owning the layout math (flex/stack/measure) — `layout.Flex`/`layout.Stack` remain gioui's for now (deferred Phase 4; the seam keeps this replaceable later).
- Full text shaper/editor abstraction (3+ months of work, rejected per Q2).
- A facade with per-field accessors over `layout.Context` (rejected, see D2).
- Import-level purity — internal implementation files *should* use gioui; "implementation detail" means exactly that.
- Replacing gioui with another engine in this change.
- `runtime.App` event-loop encapsulation beyond the `runtime.Run` signature change.

## Decisions

### D1 — The seam is a 2D draw IR, not a widget-level abstraction

Every serious framework places the engine boundary at the *draw-call* level, below layout, above composition:

| Framework | Seam | IR before backend |
|---|---|---|
| Compose → Skiko/Canvas | `Canvas` interface (~30 methods) | GraphicsLayer display list |
| Flutter → Impeller/Skia | `dart:ui` (Canvas, SceneBuilder) | Layer tree → Scene |
| Gio → OpenGL/Metal/Vulkan | `gpu.GPU` + `driver.Device` | `op.Ops` |

Rationale: composition/state are framework concerns; layout is framework concern; **drawing is the engine's**. Abstracting at the widget level would double the API surface without solving layout or text — the actual couplings. Gio itself already proves the pattern *and* ships `gpu/headless`, the same idea as our software backend.

Alternatives considered: (a) widget-level wrapper layer — rejected (surface bloat, doesn't decouple layout/text); (b) full scene-graph IR (`op.Ops` clone) — rejected (over-engineering until a second *rendering* engine exists; the test backend only needs draw-call capture, not a replayable IR).

### D2 — Defined types + `ToGio()` conversion at the seam, not a facade (user decision: option b)

gioui's `layout.Context` is a fat struct: `Ops`, `Metric`, `Constraints`, `Locale`, `Queue`, `Fonts`, `InputState`, `now`, plus methods. A facade with per-field accessors (`GetOps()`, `GetMetric()`, …) was the alternative (option a): ~15 accessor methods, each a gioui-shaped hole in the API — a *leaky* abstraction that costs more than it hides, and re-hides the exact types we're trying to remove.

Option b (chosen): `LayoutContext` becomes a **defined type** whose shape is framework-owned. Engine access is a single, auditable conversion at the seam:

```go
package layoutnode

type LayoutContext struct {
    ops         *render.Ops   // Phase 3: draw IR handle (Phase 1–2: opaque handle)
    constraints LayoutConstraints // framework-owned
    density     Density           // compose/ui/unit.Density
    locale      Locale            // compose/ui/text/intl.Locale (redefined, Phase 1)
    frameToken  uint64            // guards stale rehydration (see EC1)
}

// ToGio is the ONLY sanctioned crossing from framework types to engine types.
// It must be called during the frame that created c; frameToken asserts this.
func (c LayoutContext) ToGio() *gioLayout.Context { ... }
```

Why this is honest rather than leaky: the conversion is confined to one function in one package. Nothing above the seam imports gioui. Components that genuinely need engine services take one of two doors (EC4) — a framework-typed seam accessor built on demand, or explicit engine-bound classification.

### D3 — Purity scope: exported signatures only (user decision Q3)

The rule is: **no gioui type in an exported symbol's signature of the public package trees.** Internal bodies may (and should) keep using gioui. This is what makes the change tractable: of 189 files, only ~25 need signature surgery; the rest are already fine (their gioui use is internal implementation).

Enforcement (D5) makes this a *contract*, not a hope. Note the asymmetry: internal implementation may use gioui, but internal *types that flow into public signatures* (e.g. `TextFieldWidget`'s embedded `widget.Editor`) cannot — those get re-shaped in Phase 1.

### D4 — Phased migration, boundary-first

Order matters because each phase is independently valuable and revertible:

1. **Phase 1 — Boundary cleanup** (public API purity): the ~25 leaking files, mechanical, no behavior change. Delivers the visible win first.
2. **Phase 2 — Alias wall → defined types**: `internal/layoutnode/alias.go` conversions. Compiler-driven (every `Layout`/`Draw` signature changes; the compiler finds every touchpoint). Makes the *internal* boundary real.
3. **Phase 3 — Draw IR + software backend**: the actual seam + golden tests. Only now does behavior-capture exist to prove phases 1–2 were behavior-preserving.
4. **Deferred**: own layout math (only when a second engine demands it); `runtime.App`; text internals.

Rationale for this order: Phase 1 is low-risk and delivers the product win; Phase 2 uses the compiler as a checklist; Phase 3 needs phases 1–2's stable framework types to define the backend surface against.

### D5 — Enforcement: a `go/ast`-based API-purity checker wired into `make check-api`

The project has no CI and no linter (AGENTS.md). Enforcement must be a single, runnable thing. Design: a checker (test or small command using `golang.org/x/tools` — already a dependency, v0.46.0) that:

1. Loads packages `compose/...`, `modifiers/...`, `theme/...`, `runtime/`, `pkg/...` (via `go/packages`).
2. Walks exported declarations (`go/ast`), inspects signature types: params, returns, struct fields, embedded types, type aliases, and **underlying types** of defined types (catches `type TextAlign gioText.Alignment`).
3. Fails if any referenced type's package path starts with `gioui.org`, **unless** the referencing package is in the seam whitelist (`internal/layoutnode`, `internal/render`, and their subpackages).
4. Reports the offending package + symbol, exits non-zero.

**The 25 known leaks (inventory from the exploration) become the initial test fixtures** — the checker must flag all of them before Phase 1 and none after. This turns a "grep-based hope" into a regression suite.

Alternative considered: a full `go/analysis` pass — more powerful (can check call sites, not just signatures) but heavier; signature-level purity doesn't need it. Keep the analyzer dumb; the rule is dumb by design.

### D6 — Backend registration via blank imports (Gio's own pattern)

```go
// internal/render/render.go
var backends = map[string]BackendFactory{}

func Register(name string, f BackendFactory) { backends[name] = f }

// internal/render/gio/gio.go (blank-imported by the app shell)
import _ "github.com/zodimo/go-compose/internal/render/gio"
```

Mirrors `gioui.org/gpu`'s own `init()` registration (`_ "gioui.org/gpu/internal/opengl"`). The test suite blank-imports `internal/render/software` instead. Compile-time backend selection, no build tags.

### D7 — Text is engine-delegated (user decision Q2)

Compose's model: the framework does layout and exposes `TextLayoutResult`; the platform does shaping. We adopt the same: `compose/foundation/next/text` (shaper, editor, selection, IME, clipboard — 124 gioui references) remains engine-side. The public text *types* (`TextStyle`, `TextUnit`, `TextAlign`, `LineBreak`, font descriptors) become signature-clean; their engine conversion (`ToGioFont` etc.) moves to `internal/`. Consequences: text behavior is whatever gioui's shaper does — accepted; we do not re-implement shaping.

### D8 — The `next/` pattern is the component migration vehicle

Components already migrated to `next/`-style pure composition (card, dialog, snackbar, tab, `next/button`) need only signature-level cleanup of the primitives they use. gioui-bound components (button-old, slider, textfield, radiobutton, progress, badge) get classified (EC4) — for this change, **Door 2 (engine-bound behind the seam)** for all of them, keeping scope controlled. When Phase 2 forces their `Layout` signatures to change anyway, prefer rewriting as pure composition (the `next/` way) where the rewrite is cheaper than the adapter.

## Data Structures (Phase-3 target shapes)

```go
// ============ internal/render/backend.go — THE SEAM ============
// A 2D draw IR modeled on Compose's Canvas. Derived from the actual
// gioui op vocabulary used by modifier nodes (audited in Phase 3):
//   op.Record/Stop, op.Offset, op.Affine, paint.ColorOp, paint.Fill,
//   paint.FillShape, paint.PushOpacity, clip.Rect, clip.Path, clip.Stroke,
//   layout.Background, layout.Inset, widget.Clickable, material.Clickable
package render

type Backend interface {
    // Ops lifecycle (per frame)
    BeginFrame(w, h float32, density float32) *Ops
    EndFrame()

    // State stack
    Save() int
    Restore(count int)

    // Transforms (used by scale/, shadow/ modifiers)
    Translate(dx, dy float32)
    Scale(sx, sy float32)

    // Clipping
    ClipRect(r Rect, radius CornerRadius)
    ClipPath(ops []PathOp)

    // Paint
    FillRect(r Rect, c Color)
    FillShape(sh Shape, c Color)
    PushOpacity(a float32)

    // Text + Images
    DrawTextLayout(layout TextLayout, at Point, c Color)
    DrawImage(img Image, src, dst Rect)

    // Recording (macro semantics — see EC2)
    Record() CallOp        // begin recording draw calls
    Add(ops *Ops)          // replay a CallOp into this backend
}

// CallOp is opaque: consumers hold it, only the backend applies it.
type CallOp struct{ backend Backend; id uint64 }

// DrawCommand: what runtime.Run returns (replaces op.CallOp).
type DrawCommand interface{ applyTo(b Backend) }

// ============ internal/render/software/canvas.go — the test backend ============
package software

type DrawCall struct {
    Op    string  // "FillRect", "DrawText", ...
    Rect  Rect
    Color Color
    // ... per-op payload
}

type Canvas struct {
    calls []DrawCall  // deterministic order — no maps
    img   *image.RGBA // optional rasterization for pixel-level asserts
}

func (c *Canvas) Calls() []DrawCall { return c.calls } // golden-test input

// ============ internal/layoutnode/context.go — framework-owned context ============
package layoutnode

type LayoutContext struct {
    ops         *render.Ops
    constraints LayoutConstraints
    density     Density
    locale      Locale
    frameToken  uint64
}

func (c LayoutContext) ToGio() *gioLayout.Context // seam-only, EC1-guarded

// ============ runtime/runtime.go — after Phase 1 ============
package runtime

type Runtime interface {
    Run(LayoutContext, api.Composer, api.Composable) render.DrawCommand
}

// ============ internal/api-check/ — the enforcement tool ============
package apicheck

// Check(pkgs []string, whitelist []string) (violations []Violation)
// Violation{Pkg, Symbol, GioType}
// Invoked by `make check-api` and as a Go test in the repo.
```

## Architecture (after all phases)

```
┌──────────────────────────────────────────────────────────────┐
│ PUBLIC API — 0 gioui types (enforced by make check-api)       │
│ Composable func(Composer) Composer                            │
│ modifiers/, theme/, compose/ui/unit, pkg/                     │
│   framework-owned: Dp, Sp, TextUnit, Density, Locale,         │
│   TextStyle, TextAlign, LineBreak, Color, Rect, Offset,       │
│   Shape, ImageResource (opaque), Window (opaque)              │
├──────────────────────────────────────────────────────────────┤
│ FRAMEWORK LAYER — engine-agnostic                             │
│ api.Composer (zipper), LayoutNode tree, modifier chains       │
│ state/, pkg/flow/ (already 100% clean)                        │
│ material3 pure-composition components                         │
├──────────────────────────────────────────────────────────────┤
│ THE SEAM — internal/layoutnode + internal/render              │
│ LayoutContext (defined type, ToGio at seam)                   │
│ render.Backend interface + render.Ops/CallOp                  │
│ ⚠ only packages allowed to import gioui.org/layout,op,widget, │
│   text,font                                                   │
├──────────────────────────────────────────────────────────────┤
│ BACKENDS (blank-import registration)                          │
│ gio/       → op.Ops → gpu.GPU (production)                    │
│ software/  → DrawCall list + *image.RGBA (golden tests)       │
└──────────────────────────────────────────────────────────────┘
```

Frame flow after Phase 3:

```
runtime.Run:
  backend := render.ActiveBackend()            // gio or software
  ops := backend.BeginFrame(...)               // framework owns frame start
  ctx := layoutnode.NewContext(ops, density, locale)
  node := ui(composer).Build()
  coord := NewNodeCoordinator(node)
  coord.Layout(ctx)                            // framework types only
  coord.PointerPhase(ctx)
  return coord.Draw(ctx)                       // render.DrawCommand
app shell (demo main / runtime.App):
  cmd := runtime.Run(...)
  backend.Apply(cmd)                           // gio: callOp.Add(gtx.Ops); software: append DrawCalls
  frameEvent.Frame(gtx.Ops)                    // gio only; software backends skip
```

## Edge Cases

### EC1 — `ToGio()` validity across the frame lifecycle

Today `layout.Context` is created per-frame by `app.NewContext(&ops, frameEvent)`. After the change, the framework builds `LayoutContext` per frame inside `runtime.Run`. Danger: a component stashing `LayoutContext` across frames would `ToGio()` into a stale context (wrong `Ops` pointer, wrong `Metric`). Mitigation: `frameToken` — an incrementing uint64 assigned at frame start; `ToGio()` panics if the stored token ≠ the backend's current token. This turns a silent cross-frame bug into a loud one, mirroring gioui's own contract (a `layout.Context` is single-frame; the current code relies on the same discipline without a guard).

This is not a new discipline the design imposes — it is the *existing* invariant made explicit. Because the zipper rebuilds the `LayoutNode` tree from scratch every frame (see Context, "The composition zipper"), no node, coordinator, or context legitimately survives a frame boundary; only store-resident `PersistentState` does. The frame token therefore asserts what the architecture already guarantees, catching the two ways code can violate it: (a) a component caching the context in a struct field, and (b) an engine-bound component hoisting `ToGio()` outside the frame it was created in (e.g. into a goroutine, or into a `defer` that runs after `EndFrame`).

### EC2 — Macro replay semantics (`op.Record`/`op.CallOp`)

The coordinator records draw ops into a macro, returns `op.CallOp`; the demo replays with `callOp.Add(gtx.Ops)`. Gioui macros capture state at record time and replay at `Add` time. Two subtleties:

1. **Discard path**: the pointer phase does `macro := op.Record(gtx.Ops); …; _ = macro.Stop()` — records then **drops** (input pass produces no visuals). The backend `Record()` must preserve drop semantics: `software.Canvas` must *not* append draw calls for a dropped recording. Design: `Record()` returns `CallOp` whose `applyTo` is the only append point; dropping = never calling `applyTo`.
2. **Opaqueness**: `render.CallOp` hides `op.CallOp`; the gioui backend keeps real macros inside. Golden tests assert the software backend's `DrawCall` list, so a record/playback mismatch (e.g. a modifier that records but forgets to `Add`) shows up as missing calls — this is exactly why Phase 3 must precede any future engine swap.

### EC3 — Type-identity break: aliases → owned types

Today `box.Direction` **is** `layout.Direction`, so users pass `layout.NW`. After redefinition it's a distinct type; `layout.NW` no longer assigns. This is the core breaking change and it is *inherent to the goal* (you cannot hide a type and keep its identity). Mitigation: go-compose exports its own constants (`box.NW`, `box.NE`, …) mirroring the gioui set; the seam converts. All demos migrate in the same commit. The analyzer flags residual gioui types in *our* signatures; it cannot police user code — that's the nature of a type-identity change, and pre-1.0 (project tags are v0.0.x) is the right time for it.

### EC4 — The two doors for engine-hungry components

Textfield's IME handling needs `gtx.Queue`; `next/text`'s editor needs `gtx.Fonts`, clipboard, semantic. Two sanctioned exits:

- **Door 1 (seam accessor)**: framework-typed method added to `LayoutContext` on demand, e.g. `func (c LayoutContext) Input() InputSource`. Built only when a component actually needs it — the facade grows by demand, not preemptively (this is what keeps D2 honest).
- **Door 2 (engine-bound)**: the component's public API is clean, but its internals call `ToGio()` and use gioui directly, classified in code comments as engine-bound.

Policy for this change: `next/text` → Door 2 (consistent with the text-delegation decision); slider/textfield/radiobutton/progress/badge → Door 2. Door 1 opens only when Phase 3's backend interface needs a framework-typed equivalent (e.g. pointer events for the software backend's golden tests).

### EC5 — Analyzer whitelist is per-package, not per-symbol

Per-symbol whitelisting is a smell (it becomes "whitelist the thing I just wrote"). The whitelist is `internal/layoutnode`, `internal/render`, and subpackages. Every other package — including other `internal/` packages like `internal/composer` — is judged by the same rule: exported signatures must not contain gioui types. (Internal packages' *bodies* may import gioui freely; the rule only inspects exported symbols, matching the "signature-level" contract.)

### EC6 — `ToGio()` must use framework values as source of truth

`gtx.Metric` is derived from Density; `gtx.Locale` from Locale. If the framework's Density/Locale and the gioui context disagree (e.g. a mid-frame density change), `ToGio()` must rehydrate from **framework values** — the framework is the source of truth; the engine context is a projection. This prevents the engine from silently overriding framework state (the current code has no such ordering because there is only one source: gioui's).

### EC7 — Determinism required for golden tests

`software.Canvas` must produce identical `DrawCall` lists for identical input. Constraints: slice-based ordering only (no map iteration), no wall-clock time (animation modifiers already use injectable `timeNowFunc`), no pointer addresses in payloads, stable string formatting. Golden baselines are committed; a flaky test is treated as a backend bug, not a test problem.

## Risks / Trade-offs

- **[Breaking public API]** Type-identity changes and removed helpers (Phase 1) → Pre-1.0 project; document migration in the change's final commit; migrate all demos in-repo; keep phases as separate commits so the break is reviewable.
- **[Phase 2 churn is large]** Every `LayoutNode.Layout/Draw` signature and every modifier changes → Compiler-driven; do it package-by-package keeping the tree green; LSP rename for mechanical parts; golden tests from Phase 3 prove behavior preservation retroactively.
- **[Layout math remains gioui-owned]** If a second engine cannot reuse `layout.Flex`/`layout.Stack`, Phase 4 becomes mandatory → The seam isolates the coupling to `internal/layoutnode`; the defined `LayoutContext` means only the seam package wires layout, so Phase 4 is a contained swap, not a rewrite.
- **[Text behaviors are engine-specific]** Selection, IME, clipboard semantics are whatever gioui provides → Accepted by the delegation decision; documented as engine-bound; public API remains stable.
- **[Analyzer false positives/negatives]** The checker might miss indirect leaks (e.g. a type from a non-gio package that wraps gio) or flag benign ones → Seed with the 25 known leaks as fixtures; extend the fixture list as discovered; keep the rule dumb (signature-level only).
- **[Scope creep into full engine swap]** The seam makes swapping *possible*, which may invite it → Non-goal for this change; the software backend is the only second consumer and it exists for tests.

## Migration Plan

1. **Phase 1** (boundary cleanup, ~1–2 weeks): ship checker + fixtures first (red), then fix the ~25 files (green). Public packages stop exporting gioui types; conversion helpers move to `internal/`; `runtime.Run` returns `render.DrawCommand`. Demos updated. Commit.
2. **Phase 2** (alias wall, ~2–3 weeks): convert `internal/layoutnode` aliases to defined types; `ToGio()` introduced; compiler drives the sweep across modifiers/components. Commit per package group.
3. **Phase 3** (draw IR + software backend, ~4–6 weeks): audit modifier op vocabulary → define `render.Backend`; implement `gio` + `software` backends; wire `make test` golden runs; classify engine-bound components. Commit.
4. **Rollback**: each phase is independently revertible (git revert of that phase's commits); no behavior change until Phase 3 (golden tests then *prove* phases 1–2 changed nothing observable). If a phase stalls, the tree remains green at the previous phase's boundary.

## Open Questions

1. **`ToGio()` home**: `internal/layoutnode` (owns `LayoutContext`) vs `internal/render` (owns backend). Leaning `layoutnode` — context-specific conversion; resolved during Phase 2 when the first real caller appears.
2. **Custom renderers in this change**: all Door 2 (recommended — scope control), or migrate slider/radiobutton/progress to backend primitives since they're pure paint ops? Recommendation stands: Door 2; backend primitives only for what golden tests need.
3. **`runtime.App` (event-loop encapsulation)**: in this change or a follow-up? Recommendation: only the `runtime.Run` signature change here; full `runtime.App` is Phase 5, deferred.
4. **Checker form**: standalone `go test` in `internal/api-check` vs a `cmd/` binary? Recommendation: Go test (runs under `make test` for free) with a `make check-api` alias.
