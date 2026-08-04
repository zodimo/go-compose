# Render Backend Migration

This document describes the breaking API changes introduced by the **abstract-render-backend** change (branch `abtract-render-backend`, OpenSpec change `openspec/changes/abstract-render-backend/`). The framework's public API surface (`compose/`, `modifiers/`, `theme/`, `runtime/`, `pkg/`) no longer exposes any `gioui.org` type in an exported signature. All engine conversions moved to seam packages under `internal/` and are not importable by downstream consumers.

Last Updated: 2026-08-04

## Breaking Changes Overview

| Change | Type | Affected Packages | Migration Action |
|--------|------|-------------------|------------------|
| `runtime.Run` returns `render.DrawCommand` | signature | `runtime/` | `cmd.Apply(gtx.Ops)` instead of `cmd.Add(gtx.Ops)` |
| `TextFieldWidget`/`FilledTextFieldWidget` reshape | type-identity | `compose/material3/textfield` | Editor via wrapper methods |
| `LazyListState`/`LazyGridState` field change | type-identity | `compose/foundation/lazy` | `s.List.L.Position` |
| `SnackbarData` clickables | type-identity | `compose/material3/snackbar` | `*clickable.GioClickable` |
| `platform.LocalWindow` | signature | `compose/ui/platform` | `platform.NewWindow(window)` |
| `theme.BasicTheme`/`ThemeInterface` | type-identity | `theme/`, `compose/material` | field access unchanged; engine via `GioThemeForEngine` |
| `box.Direction`/`Stack`/`StackChild`, `column.Spacing`/`Alignment` | type-identity | `compose/foundation/layout/{box,column,row}` | `box.NW`, `column.SpaceEvenly` |
| `shape.Shape`/`shape.Outline` | signature | `compose/ui/graphics/shape` | custom shapes wrap `clipconvert` |
| `intl.Locale`, `padding.RTL` | type-identity | `compose/ui/text/intl`, `modifiers/padding` | `intl.TextDirectionRTL` |
| `unit` conversions | removed-helper | `compose/ui/unit` | use `compose.LocalDensity` |
| font conversions | removed-helper | `compose/ui/text{,/font}` | `internal/textconvert` |
| `card.GioImage` | removed | `compose/material3/card` | image composables in `CardContents` |
| `clickable.GioClickable` | type-identity | `modifiers/clickable` | `clickable.NewGioClickable()` |
| `style.TextAlign`/`style.LineBreak` | enum | `compose/ui/text/style` | same constant names/values |

## Type Identity Changes

The following exported types were previously transparent aliases of gioui types. They are now framework-owned defined types. Constant names and values are preserved; only assignability to the gioui equivalents is gone.

| Component | Old (gioui alias) | New (go-compose type) | Migration |
|-----------|-------------------|----------------------|-----------|
| `theme.BasicTheme` | `= material.Theme` | `struct { Bg, Fg, ContrastBg, ContrastFg color.NRGBA }` | Field access (`t.Fg`, `t.ContrastBg`) unchanged |
| `theme.Theme` | `= token.Theme` | `struct { Scheme *token.Scheme }` | Access via `.Scheme` |
| `box.Direction` | `= layout.Direction` | `uint8` enum, consts `NW, N, NE, E, SE, S, SW, W, Center` | `layout.NW` → `box.NW` |
| `box.Stack` | `= layout.Stack` | `struct { Alignment Direction }` | `box.Stack{Alignment: box.Center}` |
| `box.StackChild` | `= layout.StackChild` | `struct { Expanded bool; Widget func(layoutnode.LayoutContext) layoutnode.LayoutDimensions }` | Layout callback types changed |
| `column.Spacing` | `= layout.Spacing` | `uint8` enum (`SpaceEnd..SpaceEvenly`) | `layout.SpaceAround` → `column.SpaceAround` |
| `column.Alignment` / `row.Alignment` | `= layout.Alignment` | `uint8` enum (`Start, End, Middle, Baseline`) | `layout.Middle` → `column.Middle` |
| `intl.Locale` | `= system.Locale` | `struct { Language string; Direction TextDirection }` + `TextDirectionLTR`/`TextDirectionRTL` | `locale.Direction == system.RTL` → `locale.Direction == intl.TextDirectionRTL` |
| `padding.RTL` | `= system.RTL` (type `system.Direction`) | `const RTL TextDirection` (framework `layoutnode.LayoutDirection`) | Compare against `padding.RTL` |
| `clickable.GioClickable` | `= widget.Clickable` | `struct { impl *internalclickable.Impl }` | `clickable.NewGioClickable()`, `.Clicked(gtx)` |

## Removed Conversion Helpers

All gioui-conversion helpers were removed from public packages and relocated to seam packages under `internal/` (not importable by external consumers). Do not convert — pass framework types and let the frame boundary convert.

| Helper | Old location | New location |
|--------|-------------|--------------|
| `TextUnit.AsGioSp()` | `compose/ui/unit` | removed — density via `compose.LocalDensity` |
| `DpToGioUnit`, `DpToGioUnitUnsafe` | `compose/ui/unit` | `internal/unitconvert` |
| `TextUnitToGioSp`, `TextUnitToGioSpUnsafe`, `TextUnitToGioDp` | `compose/ui/unit` | `internal/unitconvert` |
| `DensityFromLayoutContext` | `compose/ui/unit` | `internal/unitconvert` |
| `ToGioFont`, `ToGioWeight`, `ToGioStyle` | `compose/ui/text/font` | `internal/textconvert` |
| `FromGioFont`, `FromGioTypeface` | `compose/ui/text/font` | `internal/textconvert` |
| `TextAlignToGioTextAlignment`, `FromGioTextAlign` | `compose/ui/text/style` | `internal/textconvert` |
| `LineBreakToGioWrapPolicy`, `GioWrapPolicyToLineBreak` | `compose/ui/text/style` | `internal/textconvert` |

> [!NOTE]
> `TextStyleFromGioFont` remains public but its parameter changed from `giofont.Font` to the framework-owned `font.Font`. Build fonts from framework types (`font.FontFamily`/`FontWeight`/`FontStyle`) — never import `gioui.org/font` for conversions.

## Signature Changes

### `runtime.Runtime.Run` — returns `render.DrawCommand`

**Before** (SHA `b134820`):

```go
type Runtime interface {
    Run(LayoutContext, api.Composer, api.Composable) op.CallOp
}
// caller:
callOp := runtime.Run(gtx, composer, UI())
callOp.Add(gtx.Ops)
```

**After** (current):

```go
type Runtime interface {
    Run(ctx any, composer api.Composer, ui api.Composable) render.DrawCommand
}
// caller:
cmd := runtime.Run(gtx, composer, UI())
cmd.Apply(gtx.Ops)
```

The gio `layout.Context` is still passed as `any`; the seam conversion happens inside `Run`. `DrawCommand` retains `Apply(*op.Ops)` for the gio path, so app-shell code changes by one line.

### `theme.ThemeManager.Material3ThemeInit` — returns `any`

**Before:** `Material3ThemeInit(gtx layout.Context) layout.Context`
**After:** `Material3ThemeInit(gtx any) any` — callers must assert back:

```go
gtx = themeManager.Material3ThemeInit(gtx).(layout.Context)
```

### `platform.LocalWindow` — opaque window

**Before:** provided `*app.Window`

```go
platform.LocalWindow, window, // raw *app.Window
```

**After:** provides `*platform.Window` (opaque); construct with `NewWindow`, read back seam-only:

```go
platform.LocalWindow, platform.NewWindow(window), // framework-owned wrapper
```

Engine-bound code reaches the gio window via `w.PlatformWindow().(*app.Window)`.

### `shape.Outline` — `clipconvert` wrappers

`CreateOutline` now takes the framework `Metric{PxPerDp, PxPerSp}` and `Outline` methods operate on `internal/clipconvert` wrappers instead of raw `*op.Ops`/`clip.*` types. Custom shape implementations must wrap engine ops (`clipconvert.NewStack(...)`, `clipconvert.NewClipOp(...)`, `clipconvert.NewPathSpec(...)`). Built-in shapes (`RoundedCornerShape`, `CircleShape`, etc.) are unchanged in construction.

### `ThemeInterface.GioMaterialTheme()`

**Before:** `*gioMaterial.Theme`
**After:** `*BasicTheme` (framework-owned). `LocalGioMaterialTheme` was removed → use `LocalBasicTheme`; Door-2 engine-bound code uses `GioThemeForEngine(c) any` and asserts the `*gioMaterial.Theme`.

## Enum Redefinitions

`compose/ui/text/style` enums are now framework-owned `int` types with identical constant names and values:

| Type | Constants |
|------|-----------|
| `style.TextAlign` | `TextAlignUnspecified=99, TextAlignStart=0, TextAlignEnd=1, TextAlignMiddle=2` |
| `style.LineBreak` | `LineBreakParagraph=0, LineBreakHeading=1, LineBreakSimple=2, LineBreakUnspecified=99` |

Users set `style.TextAlignStart` / `style.LineBreakSimple` exactly as before — only assignability to gioui `text.Alignment`/`text.WrapPolicy` is gone (engine conversion happens in `internal/textconvert`). Foundation re-exports keep their names: `foundation/text` `Alignment = style.TextAlign`, `WrapPolicy = style.LineBreak`.

> [!NOTE]
> `compose/ui/next/text/style` defines a separate, expanded pipeline: `TextAlign` has `Left/Right/Center/Justify/Start/End` and `LineBreak` is a packed int32 built via `LineBreakOf(strategy, strictness, wordBreak)`. Consumers of `next/text` use these, not `compose/ui/text/style`.

## Consumer Migration Examples

### App shell (all demos, e.g. `cmd/demo/snackbar/main.go`)

**Before** (SHA `7beb0ef^`):

```go
gtx := app.NewContext(&ops, frameEvent)
gtx.Locale = enLocale
gtx = themeManager.Material3ThemeInit(gtx) // returned layout.Context, no cast

composer := compose.NewComposer(api.ComposerWithStore(store))
cmd := runtime.Run(gtx, composer, UI())
cmd.Apply(gtx.Ops)
frameEvent.Frame(gtx.Ops)
```

**After** (current):

```go
gtx := app.NewContext(&ops, frameEvent)
gtx.Locale = enLocale
gtx = themeManager.Material3ThemeInit(gtx).(layout.Context) // opaque any → cast back

composer := compose.NewComposer(api.ComposerWithStore(store))
cmd := runtime.Run(gtx, composer, UI())
cmd.Apply(gtx.Ops)
frameEvent.Frame(gtx.Ops)
```

### `pkg/x/fileexplorer` — `RememberExplorer`

**Before** (SHA `7beb0ef`): callback took `*explorer.Explorer` (gioui).

```go
onOpenFile, _ := fileexplorer.RememberExplorer(c, func(expl *explorer.Explorer) {
    file, err := expl.ChooseFile("png", "jpeg", "jpg")
    ...
})
```

**After** (current): callback takes the framework wrapper `*explorerwrap.Explorer`; cross the seam with `.ToGio()`.

```go
onOpenFile, _ := fileexplorer.RememberExplorer(c, func(expl *explorerwrap.Explorer) {
    file, err := expl.ToGio().ChooseFile("png", "jpeg", "jpg")
    ...
})
```

External modules never name the wrapper type — the callback parameter is inferred.

### Layout constants

**Before:** `column.WithSpacing(layout.SpaceAround)`, `column.WithAlignment(layout.Middle)`, `box.WithAlignment(layout.Center)`
**After:** `column.WithSpacing(column.SpaceAround)`, `column.WithAlignment(column.Middle)`, `box.WithAlignment(box.Center)`

### Clickables

**Before:** `return &widget.Clickable{}` and `.Get().(*widget.Clickable)`
**After:** `clickable.NewGioClickable()` and `.Get().(*clickable.GioClickable)`; `.Clicked(gtx)` instead of `.Click()`.

### Card images

**Before:** `graphics.ImageResource{ ImageOp: paint.NewImageOp(img) }`
**After:** `graphics.NewImageResource(img)`; place images as composables inside `card.CardContents(card.ContentCover(fImage.Image(...)), ...)`.

### Widget structs (internal shape change, call sites unchanged)

- `TextFieldWidget`/`FilledTextFieldWidget`: `widget.Editor` → `*textinput.Editor`; `Prefix`/`Suffix` → `layoutnode.GioLayoutWidget`. The high-level `textfield.Filled(...)`/`Outlined(...)` API is unchanged.
- `LazyListState`/`LazyGridState`: `List widget.List` → `List *widgetstate.List`; scroll access `s.List.Position.First` → `s.List.L.Position.First`, `s.List.L.ScrollToEnd()`.
- `SnackbarData`: `ActionClickable`/`DismissClickable` → `*clickable.GioClickable` (`.Clicked(gtx)`).

### CLI tooling

`cmd/go-compose` required **no migration** — it shells out to `go build` and never calls the changed APIs (its only `runtime` import is stdlib, for `GOOS`/`GOARCH`).

## Migration Checklist

- [ ] Update `runtime.Run` call sites: `cmd := runtime.Run(gtx, composer, UI()); cmd.Apply(gtx.Ops)`
- [ ] Cast `Material3ThemeInit` result: `.(layout.Context)`
- [ ] Replace `layout.NW`/`layout.Center` with `box.NW`/`box.Center`
- [ ] Replace `layout.SpaceAround`/`layout.Middle` with `column.SpaceAround`/`column.Middle`
- [ ] Wrap the window: `platform.NewWindow(window)`
- [ ] Replace `&widget.Clickable{}` with `clickable.NewGioClickable()`
- [ ] Construct images via `graphics.NewImageResource(img)`
- [ ] Replace `system.RTL` comparisons with `intl.TextDirectionRTL` / `padding.RTL`
- [ ] Remove `AsGioSp`/`DpToGioUnit`/`TextUnitToGio*`/`ToGioFont` usage — pass framework types instead
- [ ] Update `LazyListState` scroll access to `.List.L.Position`
- [ ] Run `make check-api` — the analyzer now gates exported-signature purity (exit non-zero on `gioui.org` leakage)

## Version Note

This is a breaking change to the public API. Version it accordingly: `make version` shows the current tag; `make tag-patch` bumps patch (`v0.1.X` → `v0.1.X+1`), `make tag-minor` bumps minor. Tags are lightweight and act as the project's changelog. The `make check-api` target (added by this change) enforces purity going forward.
