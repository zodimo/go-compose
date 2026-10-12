## Why

The `go-sentinel-helper` contract requires that a sentinel-carrying package export **exactly one
value — its Companion** — and reach the reserved "unspecified" value only through
`XCompanion.Unspecified()`. Today only three core compose types conform
(`TextStyleCompanion`, `SpanStyleCompanion`, `ParagraphStyleCompanion` in `compose/ui/text/`);
every other core type exports a raw `var XUnspecified`, giving two public sources of truth for
"unspecified" and no companion to test, merge, default or compare through. Because the sentinels
are exported, misuse (`== FooUnspecified` instead of `Companion.IsSpecified`) is invisible to the
compiler and un-inspectable by `go doc`.

## What Changes

- **BREAKING**: Remove every exported `XUnspecified` sentinel from the in-scope `compose/ui/**`
  primitive types. The sentinel becomes package-private; the only public way to get it is
  `XCompanion.Unspecified()`.
- **BREAKING**: Delete `pkg/sentinel` entirely (7 files, ~1,100 lines) and replace all 17
  importers with `github.com/zodimo/go-sentinel-helper` — `sentinel.TakeOrElseString(a, "")`
  becomes `stringutils.StringsCompanion.TakeOrElse(a, "")`, and
  `sentinel.StringValueUnspecified` / `sentinel.StringUnspecified` become
  `stringutils.StringsCompanion.Unspecified()`. `pkg/sentinel` is a parallel fork of the same
  library (its own `Companion` interfaces, its own 9-argument `NewScalarCompanion`, a
  `StructCompoanion` typo), so the two implementations must not coexist.
- Add a package-private ten-operation set (eleven/twelve for pointer types) and register one
  exported `XCompanion` per sentinel type, following `go-sentinel-helper` v0.2.0
  `NewScalarCompanion` / `NewStructCompanion`.
- Migrate all ~314 call sites in `compose/ui/**` off the exported sentinels, routing through
  the new companions.
- Delete `pkg/sentinel`'s local generic `NewScalarCompanion` — along with the package itself —
  since it omits the `isUnspecified` argument and diverges from the library contract used by the
  three conforming companions.
- Add conformance tests per type asserting the semantics that survive review: sentinel ≠ zero
  value, merge direction, unspecified-equals-unspecified, and (struct) coalesce/copy.

## Capabilities

### New Capabilities

- `sentinel-companions`: the contract that every sentinel-carrying type in the compose core
  surface exports exactly one Companion, keeps its sentinel package-private, wires the ten
  canonical operations with the specified semantics, and validates the result against the
  `go-sentinel-helper` package contract.

### Modified Capabilities

<!-- None: openspec/specs/ contains only form-binding and form-engine, neither of which changes. -->

## Impact

- **Code (in scope)**: `compose/ui/geometry/*` (Offset, Size, CornerRadius),
  `compose/ui/unit/dp.go` (Dp, DpOffset, DpSize), `compose/ui/text/unit/text_unit.go` (TextUnit),
  `compose/ui/graphics/*` (Color, Brush, Shadow), `compose/ui/graphics/shape/*`,
  `compose/ui/text/font/*` (FontWeight, FontStyle), `compose/ui/text/intl/locale_list.go`,
  `compose/ui/next/text/**` (TextStyle, SpanStyle, ParagraphStyle, platform styles, style/*).
- **Code (deleted)**: `pkg/sentinel/` — `companion.go`, `companions.go`, `helper.go`, `ints.go`,
  `strings.go` and their tests. `pkg/sentinel`'s four exported helpers with no library
  equivalent (`Sentinel`, `TakeOrElseValue`, `TakeOrElsePointer`, and the local
  `NewScalarCompanion`) are dropped rather than ported; the first three have no callers outside
  the package and the constructor is superseded by `sentinel/companions.NewScalarCompanion`.
- **Importers to migrate (17)**: `compose/foundation/form/components/*` (11 files),
  `compose/material3/textfield/*` (3), `compose/material3/next/textfield/*` (2), and
  `compose/ui/graphics/shadow.go`. No test file outside `pkg/sentinel` imports it.
- **Out of scope** (follow-up change): `compose/foundation/**`, `compose/material3/**`
  (including `material3/next/button/*`), `theme/`, `pkg/adapters/**`. These still export
  sentinels after this change.
- **API**: breaking — callers must replace `FooUnspecified` with `FooCompanion.Unspecified()`
  and predicate checks with `FooCompanion.IsSpecified/IsUnspecified`.
- **Dependencies**: `github.com/zodimo/go-sentinel-helper` moves from indirect to a direct
  require (currently `// indirect` at v0.2.0).
- **Docs**: `docs/sentinel_pattern.md` and `AGENTS.md` gain a note on the core-surface contract.
