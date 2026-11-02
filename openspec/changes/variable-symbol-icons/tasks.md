# Tasks

## 1. Land the vector skeleton before building on it

- [ ] 1.1 Add a `symbolFont` registry in `compose/material3/icon` that parses embedded TTF bytes into a `*typesetting/font.Font` lazily, once per style, keeping the parsed value for the process lifetime. Verify: unit test parses each available style twice and asserts the same pointer is returned, and that `Upem()` is 960 for the bundled font.
- [ ] 1.2 Implement `outlineToPath(segments []font.Segment, upem uint16, sizePx int) clip.PathSpec`, mapping `SegmentOpMoveTo/LineTo/QuadTo/CubeTo` onto `clip.Path` and closing each contour. Verify: golden test on `SymbolHome` asserts a non-empty path and stable segment count against the 30-segment baseline recorded from the bundled font.
- [ ] 1.3 Replace `iconFromSymbol` with a vector renderer that extracts `GlyphOutline` for `NominalGlyph(rune)` and fills it with the resolved color, returning `layout.Dimensions` matching the requested pixel size. Verify: headless test renders `SymbolSearch` into an `op.Ops` and asserts ops were recorded and dimensions equal the requested size.
- [ ] 1.4 Delete the text-based symbol path (`text.Text` usage in `iconFromSymbol`, `MaterialSymbolsFontFamily`, the `Dp→Sp` conversion) and confirm `go build ./...` still compiles for all demos. Verify: `go build ./...` and `go vet ./...` pass; `grep -r MaterialSymbolsFontFamily` returns no hits.
- [ ] 1.5 Add a default-render golden test: render a representative icon set with no options and compare against a baseline captured from the pre-change text path, documenting in the test any accepted sub-pixel difference. Verify: `go test ./compose/material3/icon/...` passes and the diff rationale is recorded in the test file comment.

## 2. Expose the variable axes

- [ ] 2.1 Add `WithFill`, `WithGrade`, `WithOpticalSize` and `WithWeight` options, with range constants for each axis and a `hasX bool` alongside each value so an explicit `0` stays distinguishable from unset. Verify: table-driven test asserts each option marks its axis as provided, and that unprovided axes report unset.
- [ ] 2.2 Clamp each axis to its range in one resolution step (fill `0..1`, grade `-50..200`, opsz `20..48`, weight `100..700`), and assert clamped values feed both the renderer and the cache key. Verify: tests for grade `500` → `200` and fill `-1` → `0`, matching the spec scenarios.
- [ ] 2.3 Translate resolved design-space axis values into `[]font.Variation` and apply them via `Face.SetVariations` during extraction. Verify: test asserts `FILL=1` and `FILL=0` produce different outlines while both are non-empty, and that `wght=700` yields strictly larger outline bounds than `wght=100`.
- [ ] 2.4 Add `SymbolStyle` with `SymbolOutlined` (default), `SymbolRounded` and `SymbolSharp`, plus `WithSymbolStyle`. Verify: test asserts the three style constants resolve and that an unknown style falls back to outlined rather than rendering blank.

## 3. Size-adaptive optical sizing

- [ ] 3.1 Derive `opsz` from the resolved icon size when it was not provided, mapping dp onto the axis range and clamping, treating explicit `opsz` as authoritative. Verify: tests cover 20dp → axis minimum, 48dp and above → axis maximum, and 48dp with explicit `24` → `24`, matching the spec scenarios.
- [ ] 3.2 Make the resolved pixel size the single source of truth for both geometry scale and optical size, removing any remaining font-size derivation. Verify: test asserts a 32dp icon reports roughly 32px dimensions and selects an `opsz` strictly between the axis minimum and maximum.

## 4. Brush and stroke painting

- [ ] 4.1 Add `WithBrush(graphics.Brush)` with precedence brush > color > ambient `LocalContentColor`, filling the icon geometry with the brush shader. Verify: test renders with a linear-gradient brush and asserts a gradient paint op was recorded instead of a solid fill.
- [ ] 4.2 Add `WithStroke(width unit.Dp, ...)` painting via `clip.Stroke`, returning the same dimensions regardless of stroke width. Verify: test asserts a stroke op is recorded at the requested width and that reported size is unchanged versus the unstroked icon.
- [ ] 4.3 Document weight versus stroke in the package doc comments so the two knobs' independence is explicit. Verify: `go doc ./compose/material3/icon` shows both options with the distinction stated.

## 5. Caching

- [ ] 5.1 Add a bounded LRU cache keyed by `(symbol, style, clamped axes, sizePx)` holding converted geometry, with a mutex-guarded map plus `container/list` mirroring the existing `GlobalIconCache` shape. Verify: test asserts repeated identical layouts convert once (counting extractions), and that two different weights yield two entries.
- [ ] 5.2 Add an ops-level cache keyed additionally by resolved paint and stroke width, recording into a per-entry `op.Ops` and replaying an `op.CallOp`, so paint changes do not force geometry reconversion. Verify: test toggles color between frames and asserts no new extraction occurred while the rendered ops differ.
- [ ] 5.3 Evict least-recently-used entries once the cache is full and assert correctness after eviction. Verify: test renders more distinct configurations than the cache limit, then re-renders an early configuration and asserts correct geometry is produced.

## 6. Styles, assets and demos

- [ ] 6.1 Embed the Rounded and Sharp variable TTFs in `assets/fonts` behind an opt-in `SymbolCollection(styles ...SymbolStyle)` accessor, leaving `Collection()` returning outlined only so default binaries do not grow. Verify: `go test ./assets/fonts/...` asserts `Collection()` excludes rounded/sharp faces and `SymbolCollection(SymbolRounded, SymbolSharp)` includes them.
- [ ] 6.2 Add a test asserting a missing style asset surfaces as a diagnosable error at composition time and falls back to outlined, never a blank glyph. Verify: test stubs the registry to omit a style and asserts the fallback plus a non-blank render.
- [ ] 6.3 Update `cmd/demo/iconsymbols` and `cmd/demo/icon-unified` to exercise axes, styles, optical sizing, brush and stroke. Verify: both demos build and `go run ./cmd/demo/iconsymbols/` launches and renders the new variants.
- [ ] 6.4 Add a benchmark covering cold extraction versus cached replay for one icon configuration. Verify: `go test ./compose/material3/icon/ -bench . -run '^$'` reports both paths and the cached path is measurably faster.

## 7. Integration checks

- [ ] 7.1 Run the full suite and confirm no regressions outside the icon package, particularly in `compose/material3` consumers of `icon.Icon`. Verify: `go test ./...` passes.
- [ ] 7.2 Confirm the public surface is additive only: `Icon`, `Symbol`, `WithSize`, `WithColor`, `WithModifier`, `WithSymbolSize`, `IconBytes` and every `SymbolName` constant behave as before. Verify: `go doc ./compose/material3/icon` diffed against the pre-change surface shows only additions.

## Workflow follow-up

- Run `openspec validate variable-symbol-icons --strict` after implementation.
- Archive the change once the demos and suite are green.
