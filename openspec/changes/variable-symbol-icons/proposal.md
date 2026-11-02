# Proposal

## Why

`compose/material3/icon` renders Material Symbols by drawing the icon's codepoint as **text**
through Gio's shaper. Gio's `font.Font` carries only `Typeface`, `Style` and `Weight`, and its
shaper never calls typesetting's `SetVariations` — so the font's variable axes are unreachable
from a `TextStyle`. The bundled `MaterialSymbolsOutlined.ttf` is a variable font with four axes
(`FILL 0..1`, `GRAD -50..200`, `opsz 20..48`, `wght 100..700`), yet every icon today renders at
the fixed default instance. Filled/weighted/graded symbols — the entire point of Material
Symbols — cannot be expressed.

Rendering through a text shaper also caps icons at what a font glyph can be: a single
solid-colored shape, positioned by text layout. Gradients, per-icon strokes, and clipping are
impossible, and each icon pays for a full text layout pass.

`symbols.go` itself is not the problem — it is 4,102 correct `name → codepoint` constants. The
limitation is entirely in the render path.

## What Changes

- Add a **vector render path for `SymbolName` icons**: parse the bundled variable TTF once with
  `github.com/go-text/typesetting`, apply `Face.SetVariations`, extract the `GlyphOutline` for
  the codepoint, and fill it as a `gioui.org/op/clip.Path`. Verified: the same glyph under
  `FILL=0` and `FILL=1` produces distinct outlines, and both fill headlessly.
- Expose the four axes as `IconOption`s — `WithFill`, `WithGrade`, `WithOpticalSize`,
  `WithWeight` — each translating to a design-space `font.Variation`.
- **Size-adaptive `opsz`**: when `opsz` is not set explicitly, derive it from the resolved icon
  size (clamped to the axis range) so optical sizing matches Material's behavior.
- Add **brush and stroke fill** for symbol icons: `WithBrush(Brush)` for gradient fills via
  `paint.LinearGradientOp`-style shaders, and `WithStroke(...)` for outlined icons. This is the
  capability a font glyph cannot provide.
- Support **Outlined, Rounded and Sharp** styles via `WithSymbolStyle(SymbolStyle)`, backed by
  embedded font assets (see Impact for the size trade-off).
- Add an outline cache keyed by `(codepoint, style, axes, upem-scaled size)` so repeated layouts
  reuse converted paths instead of re-extracting outlines every frame.
- Keep the public surface compatible: `Icon(source, options...)`, `Symbol(...)`, the `SymbolName`
  constants, `WithSize`, `WithColor`, `WithModifier` and `WithSymbolSize` all keep working.
  `IconBytes` icons keep their existing `widget.NewIcon`/IconVG path unchanged.
- Update the `cmd/demo` icon demos to exercise axes, styles, brush/stroke and optical sizing.

## Capabilities

### New Capabilities

- `symbol-icons`: how `SymbolName` icons are rendered — variable-axis selection, style
  selection, size-adaptive optical sizing, brush/stroke fill, and the geometry cache that keeps
  the vector path affordable.

### Modified Capabilities

None. No existing capability's requirements change.

## Impact

- **Code**: `compose/material3/icon/` gains the vector renderer, axis/style options, outline
  cache and their tests; `icon.go`'s `iconFromSymbol` is replaced. `assets/fonts/fonts.go` gains
  Rounded and Sharp faces. `cmd/demo/iconsymbols` and `cmd/demo/icon-unified` grow the new usages.
- **Dependencies**: adds a direct `github.com/go-text/typesetting` dependency (already present
  transitively through Gio at v0.3.4) and `gioui.org/op/clip` / `gioui.org/op/paint` usage.
- **Binary size**: this is the significant cost. Each style's variable TTF is ~10 MB
  (`MaterialSymbolsOutlined.ttf` is 10.2 MB), so embedding all three adds roughly 20 MB to any
  binary importing `assets/fonts`. The design phase must settle this: embed all three, embed
  Outlined only and accept a font path for the rest, or gate styles behind a build tag.
- **Behavior**: default icons (no axis options) must render identically to today's glyph
  rendering at the default instance — this is the key regression risk and needs a golden test.
- **Public API**: additive only inside `compose/material3/icon`; no signatures change.
