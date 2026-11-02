# Design

## Context

See `proposal.md` — Why for motivation and `specs/symbol-icons/spec.md` for the behavior contract.

Current state that shapes the approach:

- `iconFromSymbol` (`compose/material3/icon/icon.go:98`) builds a `text.Text` with
  `MaterialSymbolsFontFamily` and the codepoint as its string. Icon size is expressed as a font
  size; there is no path to the font's axes.
- Gio's shaper is a dead end for this. `gioui.org/font.Font` carries only `Typeface`, `Style`,
  `Weight`; `text.Parameters` has no variation field, and `text/shaper.go` resolves faces purely
  by that metadata. `go-text/typesetting`'s `Face.SetVariations` is reached only if you own the
  `*font.Face`, which Gio's shaper does.
- The bundled `assets/fonts/MaterialSymbolsOutlined.ttf` (10.2 MB) is variable:
  `FILL 0..1`, `GRAD -50..200`, `opsz 20..48`, `wght 100..700`, `upem = 960`, 7 named instances.
- `compose/ui/graphics.Path` exists as an *interface* with no concrete implementation, so it
  cannot be reused for this. Gio's `op/clip.Path` is the usable primitive.
- The existing `IconBytes` path is unrelated to font limitations — `widget.NewIcon` decodes IconVG
  to an image. It stays as-is.

A probe confirmed the approach end to end: parsing the TTF with `font.ParseTTF`, calling
`Face.SetVariations`, taking `GlyphData(gid)` as `font.GlyphOutline`, converting its `Segment`s
to `clip.Path` and filling it produces valid geometry, headlessly, and `FILL=0` vs `FILL=1`
yields distinct outlines.

## Goals / Non-Goals

**Goals:**

- Own the `*font.Face` directly so every axis is settable, and make axis choice a first-class,
  typed part of the icon API.
- Keep the vector path cheap enough for real layouts: convert once per configuration, replay ops
  per frame.
- Preserve exact rendering for icons that set no new options (see Risks — this is the main one).
- Keep the change additive: `Icon(source, options...)`, `Symbol(...)` and the 4,102 generated
  constants keep working unchanged.

**Non-Goals:**

- Replacing the `IconBytes` / IconVG path, or converting it to vector.
- Replacing `symbols.go` or changing the generator's output shape. The constants are correct and
  are not the bottleneck.
- Animating axes across frames. The design must not *preclude* it (geometry is keyed on effective
  axes, so a frame with a new axis value simply misses the cache), but no animation API ships here.
- Text glyph variation support generally — this is scoped to icons.
- Removing `WithSymbolSize` or any other deprecated shim.

## Decisions

### D1. Own the font face and extract outlines; do not go through the shaper

Verified alternative — pre-instantiating static faces per axis combination and registering them
via `text.WithCollection`: rejected. Axes become a discrete menu rather than continuous values,
`FILL` × `GRAD` × `opsz` × `wght` × 3 styles is a combinatorial explosion of faces, and it does
not unlock brush or stroke fill.

Approach: parse each style's TTF once into a `*typesetting/font.Font`, create a
`font.Face` per layout with `font.NewFace`, call `SetVariations` with design-space values, look up
the glyph via `NominalGlyph(rune)`, take `GlyphData`, and convert `font.GlyphOutline.Segments`
into a `clip.Path` scaled by `iconSizePx / upem`.

`SetVariations` takes `[]font.Variation{Tag, Value}` in **design units** (not normalized), which
matches how the axis ranges are documented, so no coordinate normalization is needed at the call
site.

### D2. Vector path is the only symbol path; keep a flat cutover

The proposal settled on vector outlines for symbol icons. Rather than keeping text rendering as
a fallback, `iconFromSymbol` is replaced outright, with the text path deleted. A dual path would
mean two rendering behaviors to keep in sync, two golden-test baselines, and the text path would
silently diverge on axes anyway. The regression risk this creates is accepted and covered by
golden tests (see Risks).

### D3. Axis options are typed, clamped values, not raw variation settings

Four explicit options — `WithFill(float32)`, `WithGrade(float32)`, `WithOpticalSize(float32)`,
`WithWeight(float32)` — each carrying a range constant, rather than exposing
`WithVariations([]font.Variation)`. Rationale:

- The four axes are fixed by the font format; there is no fifth axis to be forward-compatible
  with.
- Typed options let the renderer clamp at a single boundary (D4) and let the cache key be
  computed from a canonical struct rather than an arbitrary slice whose order affects equality.
- A raw escape hatch would leak `go-text/typesetting` types into the public API of a Material
  component package.

Considered: a `SymbolVariations` struct passed in one option. Rejected as equivalent but less
discoverable through `go doc`.

Note that the project's sentinel convention (`XCompanion.Unspecified()`) applies to `compose/ui/**`
types; the icon package is in `compose/material3/**`, which the AGENTS.md sentinel rule explicitly
excludes. Axis "unset" is therefore modeled as "not provided", using an explicit
`hasX bool` alongside each `float32` in the options struct — not as a sentinel value — so that an
explicit `0` for fill or grade remains distinguishable from "unset".

### D4. Clamping happens once, at option resolution

Each axis is clamped to its range when `IconOptions` are resolved (before the composable runs),
and the *effective* axis values — including the auto-derived `opsz` — are what the renderer and
the cache key see. This means:

- Out-of-range values can never produce an out-of-range variation call.
- The cache key is canonical: two calls that differ only in a clamped-away out-of-range value hit
  the same entry.
- The size-adaptive rule ("explicit opsz wins") is expressed as: if `opsz` was not provided,
  derive from size; otherwise use the provided value, clamped.

### D5. Size-adaptive optical size derives from resolved dp, not sp

The renderer resolves the icon's pixel size first (today: `DpRoundToPx` of the requested size, or
the inherited font size when unset), then maps that dp value onto the `opsz` range by clamping
`dp × 96/72` — i.e. treating the icon size as a CSS pixel size, which is what `opsz` means.
The existing `Dp → Sp` conversion used for font sizing is removed with the text path.

### D6. Paint is a brush-or-color, and stroke is a separate geometry pass

Extend the options with `WithBrush(graphics.Brush)` and `WithStroke(width unit.Dp, cap, join)`.
Precedence: brush > color > ambient `LocalContentColor`. Stroke, when non-zero, paints the
outline with `paint.FillShape` on the *stroke* version of the path instead of the filled one;
Gio has no stroke-to-fill conversion for arbitrary paths, so the design uses
`clip.Stroke{Path: ..., Width: ...}.Op()` — verified to exist in this Gio version — and reports
the same `layout.Dimensions` regardless of stroke width, so stroke does not perturb layout.

### D7. Two-level cache: converted path per (symbol, style, axes, size), ops per (path, paint)

The probe showed conversion is the expensive part (parsing, variation setup, outline extraction,
segment→path). Recording ops is cheap. So:

- **Level 1** — `pathKey{symbol, style, fill, grade, opsz, wght, sizePx}` → `clip.Path`-derived
  `clip.Op` factory state. Bounded LRU, mirroring the existing `GlobalIconCache` shape
  (map + container/list, mutex-guarded, limit ~256).
- **Level 2** — reuse the existing pattern of recording into a per-entry `op.Ops` and replaying a
  `op.CallOp`, keyed additionally by the resolved paint and stroke width. This is exactly what
  `iconWidgetConstructor` already does for `IconBytes`; the symbol path adopts the same shape so
  the two paths stay symmetrical.

The two-level split matters because paint changes frequently (hover, theme) while axes rarely do.

### D8. Style is a font asset, selected up front

`SymbolStyle` (Outlined default, Rounded, Sharp) maps to a lazily-parsed `*font.Font` held in a
package-level registry, parsed once on first use from embedded bytes. Verified: Rounded and Sharp
variable TTFs are published at the same upstream path as Outlined.

The 20 MB cost is real, so the registry exposes the three assets through `assets/fonts` with the
new faces behind an explicit opt-in rather than unconditional appends to `Collection()`, so a
binary that only ever renders outlined icons does not pay for the other two. `Collection()`
keeps returning Outlined only; a new `SymbolCollection(styles ...SymbolStyle)` returns the
requested faces. This keeps today's default binary size unchanged.

### D9. Errors surface at composition time, not mid-draw

A missing style asset or an unmapped codepoint is knowable during composition. The design
validates style availability when options resolve and falls back to Outlined with a diagnostic —
satisfying the spec's "fails visibly" scenario without risking a panic inside `Layout`, where
there is no channel to report.

## Risks / Trade-offs

- **Default-render visual drift** → The vector path rasterizes differently from the text shaper
  even at identical axes (hinting, advance-width handling, baseline alignment). Mitigation: a
  golden test capturing a representative icon set rendered with no options, compared before and
  after; icon box alignment must be derived so the glyph's visual bounds sit where the shaper
  placed them. If a pixel-exact match proves impossible, the acceptable outcome is "visually
  identical at 24dp", documented in the test rather than silently loosened.

- **Binary size from three variable fonts (~30 MB total)** → Mitigation: D8's opt-in
  `SymbolCollection`, Outlined remains the only default face; the added cost is paid only by
  apps that select Rounded or Sharp.

- **Up-front font parse cost** → Parsing a 10 MB TTF on first icon render would stall the first
  frame. Mitigation: parse lazily on first use per style, and keep the parsed `*font.Font` for
  the process lifetime; the existing `NewGlobalIconCache`-in-composer-state pattern shows the
  project already accepts first-use costs, but this one is measurable and belongs in a benchmark.

- **`op/clip` API surface is narrower than it looks** → The probe hit several dead ends
  (`clip.Outline` has no `Push`, `ops.Bytes()` does not exist, `font.NewFace` is a package
  function). Mitigation: the task list carries an explicit early step to land a compiling
  end-to-end renderer skeleton with its tests before any option plumbing, so API mismatches
  surface before the option surface is built on top.

- **Stroke geometry differs from font-rendered thickness** → Weight and stroke are related but
  independent knobs; a user may expect `WithWeight(700)` and `WithStroke` to combine
  monotonically. Mitigation: document that weight changes the glyph design and stroke adds an
  independent outline; the spec deliberately keeps stroke width from affecting icon size.

- **Cache key correctness** → Two configurations that should be identical must produce identical
  keys, and configurations that differ must not collide. Mitigation: keys are built from the
  canonical post-clamp axis struct plus integer pixel size, never from float Dp values.

## Migration Plan

Purely additive for callers; no call-site migration required. Rollback is reverting the icon
package, since no other package's API changes. The `IconBytes` path and `symbols.go` are
untouched, so a revert restores exactly today's behavior.

## Open Questions

None — the branch sizes, style set, feature depth and compatibility posture were settled before
this design. Remaining unknowns (exact LRU limit, exact stroke default cap/join) are tunable
implementation details that cannot change the specs or the task breakdown.
