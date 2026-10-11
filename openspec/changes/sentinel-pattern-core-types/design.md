## Context

`go-sentinel-helper` v0.2.0 models "unspecified" as a reserved bit-pattern per type plus a
Companion — one value bundling ten canonical operations — so callers never write `== nil` or
compare against a raw sentinel. The skill's package contract is unambiguous: a sentinel-carrying
package exports **exactly one value, its Companion**, and the sentinel stays package-private,
reachable only through `XCompanion.Unspecified()`.

The current state of the compose core surface is inconsistent:

- **Conforming (3):** `compose/ui/text/text_style.go`, `span_style.go`, `paragraph_style.go`
  register `companions.NewStructCompanion` with package-private sentinels
  (`textStyleUnspecified`, …). These are the model for the rest.
- **Non-conforming (~22 types):** every other core sentinel type exports a raw
  `var XUnspecified` and, at best, exposes operations as exported package-level funcs or
  receiver methods. `compose/ui/geometry/offset.go:14` is representative: an exported
  `OffsetUnspecified` var, plus `IsOffset()`, `IsUnspecified()`, `TakeOrElse()` methods that
  hard-code `== OffsetUnspecified`. There is no companion, so there is no interface to test,
  merge, or default through.
- **Divergent fork (1 package):** `pkg/sentinel` reimplements the library — its own
  `Companion` interfaces, its own generic `NewScalarCompanion` taking **nine** arguments (it
  omits `isUnspecified`) registered with nine args (`strings.go:74-84`), and its own helpers.
  It is consumed by **17 files** across `compose/foundation/form/components/*`,
  `compose/material3/textfield/*`, `compose/material3/next/textfield/*` and
  `compose/ui/graphics/shadow.go`.

Measured call-site impact within scope: **314 references to exported `*Unspecified` sentinels
across `compose/ui/**`**, spread over ~24 sentinel names (top: `TextUnitUnspecified` 40,
`TextForegroundStyleUnspecified` 31, `ColorUnspecified` 29, `SizeUnspecified` 27,
`FontWeightUnspecified` 22, `ShadowUnspecified` 21, `OffsetUnspecified` 20).

User decisions that fix the shape of this change: (1) **hard-remove** the exported sentinels —
no deprecated aliases; (2) scope is **`compose/ui/**` primitives only**; (3) `pkg/sentinel` **is
in scope** and must be aligned with the v0.2.0 contract.

## Goals / Non-Goals

**Goals:**

- Every in-scope sentinel type exports exactly one value: `XCompanion`.
- No exported `*Unspecified` symbol remains in `compose/ui/**` or `pkg/sentinel`.
- All 314 in-scope call sites migrate to the companion API.
- `pkg/sentinel` is deleted and all 17 importers use `go-sentinel-helper` directly, so only one
  sentinel implementation exists in the module.
- Per-type conformance tests covering the semantics that survive review.

**Non-Goals:**

- `compose/foundation/**`, `compose/material3/**` (incl. `material3/next/button/*`), `theme/`,
  `pkg/adapters/**` — these keep their exported sentinels for a follow-up change.
- Changing sentinel *representations* (e.g. `NaN` for `Dp`, packed `TextUnit`). The sentinel
  value for each type is already correct; only its visibility and the companion wiring change.
- Introducing `Coalesce`/option-`Copy` for value types that cannot be nil.
- **`pkg/sentinel`'s non-importing helpers** (`Sentinel`, `TakeOrElseValue`,
  `TakeOrElsePointer`) — verified to have no callers outside the package, so they are deleted
  with it rather than ported.
- Any change to rendering, layout, or text-shaping behavior. This is an API-surface and
  correctness refactor; runtime output must be byte-identical.

## Decisions

### D1 — Companion per type, package-private sentinel and ops

For each in-scope type: rename `var XUnspecified` → `xUnspecified` (package-private), add
package-private `isSpecifiedX` / `isUnspecifiedX` / `takeOrElseX` / `mergeX` / `stringX` /
`sameX` / `semanticEqualX` / `equalX` / `copyX`, and register one exported `XCompanion`.

*Scalar vs struct:* value types (`Offset`, `Size`, `CornerRadius`, `Dp`, `DpOffset`, `DpSize`,
`TextUnit`, `Color`, `FontWeight`, `FontStyle`, `BaselineShift`, `Shape` implementations) use
`NewScalarCompanion` — ten args, `copy` is `func(T) T`. Pointer/interface types (`Shadow`,
`Brush`, `LocaleList`, `TextIndent`, `LineHeightStyle`, `TextForegroundStyle`,
`TextGeometricTransform`, `PlatformTextStyle`, `PlatformSpanStyle`, `PlatformParagraphStyle`,
and the `next/text` `TextStyle`/`SpanStyle`/`ParagraphStyle` family) use `NewStructCompanion` —
the same ten plus `coalesceX` before `copyX`, with `copy` as `func(T, ...func(T)) T`.

*Why:* the three existing conforming companions already establish the exact shape, and the
library's `NewScalarCompanion`/`NewStructCompanion` fix the argument order. Reusing it keeps one
implementation of NaN-aware equality (`floatutils`) rather than per-type reimplementations.

*Alternative rejected:* a generic `Sentinel[T]` helper to shrink the boilerplate. Rejected — it
would either hide the sentinel behind a closure (per the skill's performance caveat) or require
generics on types the library already parameterises, and it diverges from the three companions
already in the tree.

### D2 — Predicate logic in the companion, not on the type

Receiver methods that exist today only to express sentinel logic — `Offset.IsUnspecified()`,
`Offset.IsOffset()`, `Offset.TakeOrElse()` — are **removed** (they are part of the same exported
sentinel family: they encode "compare against the raw sentinel"). The same expression becomes
`OffsetCompanion.IsUnspecified(o)` / `IsSpecified(o)` / `TakeOrElse(o, def)`.

Genuine domain methods stay: `Offset.X()`, `IsValid()`, `GetDistance()`, `String()`, arithmetic.
`String()` keeps its human-readable form but is no longer the canonical stringifier — the
companion's `stringX` is.

*Why:* the contract's whole point is "one public source of truth for what is unspecified." A
method that asks the same question as the companion is exactly the duplication the contract
forbids. *Risk:* receivers are convenient; removing them touches call sites that use
`o.IsUnspecified()` rather than the sentinel name — these are counted in the 314 and handled in
the same pass.

### D3 — pkg/sentinel is deleted, not aligned

Delete `pkg/sentinel` outright and migrate its **17 importers** to
`github.com/zodimo/go-sentinel-helper` directly. Inspection of the package shows it is not a
thin re-export but a **parallel fork of the same library**: its own `Companion`,
`ScalarCompanion` and `StructCompoanion` interfaces (`companion.go`), its own 9-argument
`NewScalarCompanion` missing `isUnspecified` (`companions.go:15`), and its own
`Sentinel`/`TakeOrElseValue`/`TakeOrElsePointer` helpers (`helper.go`). Keeping it aligned would
preserve a second, redundant implementation of the exact abstraction the library exists to
provide — so it goes.

*The mapping is total and mechanical.* Every use site falls into one of four patterns:

| `pkg/sentinel` usage | sites | replacement |
|---|---|---|
| `sentinel.TakeOrElseString(a, "")` | 18 | `stringutils.StringsCompanion.TakeOrElse(a, "")` |
| `sentinel.StringValueUnspecified` | 11 | `stringutils.StringsCompanion.Unspecified()` |
| `sentinel.StringUnspecified` (deprecated alias) | 6 | `stringutils.StringsCompanion.Unspecified()` |
| `sentinel.TakeOrElse(a, b, unspec)` (`shadow.go` only) | 3 | per-field companion `TakeOrElse` |

*Why:* the user's instruction is to replace usage, and the contract's value comes from there
being one implementation. A fork cannot satisfy the contract no matter how it is aligned.

*On the four helper exports with no library equivalent:* `Sentinel`, `TakeOrElseValue` and
`TakeOrElsePointer` have **no callers outside the package**, and the local
`NewScalarCompanion` is superseded by `sentinel/companions.NewScalarCompanion`. All four are
dropped rather than ported. This also resolves the previously-open question about
`pkg/sentinel`'s exported per-op funcs — they disappear with the package.

*Caveat the change must respect:* `go-sentinel-helper`'s own `stringutils` and `intutils`
**do** export `StringValueUnspecified` / `IntValueUnspecified` consts. The "exactly one export"
rule is therefore enforceable on *our* packages but not on the vendored library; migrated code
calls `StringsCompanion.Unspecified()` rather than the const so the companion stays the single
door on our side. (Scenario "The library's own exported sentinel is not used as a comparison".)

### D4 — Hard removal with a compile-checked migration

Delete exported sentinels outright (per the user's decision). The compiler is the migration
oracle: after removal, `go build ./...` enumerates every broken call site. This is chosen over
deprecated aliases precisely because aliases would let the contract violation persist
invisibly — the point of the change is that the raw sentinel stops being reachable.

*Ordering to keep the tree buildable:* introduce companions and rename sentinels
**per package**, updating that package's call sites in the same commit, then move outward. The
module will not compile mid-migration; that is accepted and the change lands as one unit.

## Risks / Trade-offs

- **~314 call sites across ~24 types → large, mechanical, merge-conflict-prone diff.** →
  Migrate package-by-package in the order (geometry, unit, text/unit, graphics, graphics/shape,
  text/font, text/intl, next/text) and keep each package's edit self-contained. Verify each
  package with `go build ./compose/ui/<pkg>/...` before moving on.
- **Transposed `isSpecified`/`isUnspecified` compiles cleanly and inverts silently.** →
  Every companion gets a test asserting `IsSpecified(Unspecified())==false` and
  `IsUnspecified(Unspecified())==true`; the skill calls this out as the canonical guard.
- **`Merge` direction inverted silently.** → Test both `Merge(base, Unspecified())==base` and
  `Merge(base, override)==override` per type.
- **`NaN`-based sentinels make `==` comparisons wrong.** → `sameX` and `semanticEqualX` delegate
  to `floatutils`; `semanticEqualX(Unspecified(), Unspecified())` must be true, with a test.
- **Removing receiver methods may break non-`compose/ui` callers** (demos, `cmd/`, adapters). →
  These are out of the *type* scope but not out of *compile* scope; the migration sweep is
  `go build ./...`, and demos in `cmd/demo/**` that reference the sentinels are updated in the
  same change as ordinary call sites.
- **`pkg/sentinel`'s exported helper funcs** (`IsSpecifiedString`, `MergeString`, …) are a
  weaker contract violation than the sentinels. → Resolved by D3: the whole package is deleted,
  so the funcs go with it. Verified that no non-`pkg/sentinel` call site depends on them
  (the only names used externally are `TakeOrElseString`, `StringValueUnspecified`,
  `StringUnspecified`, `TakeOrElse`).
- **Deleting `pkg/sentinel` touches 17 files outside the `compose/ui/**` type scope**, in
  `foundation` and `material3`. → Accepted and unavoidable: the user's instruction is to
  replace the package's usage. These are *import* migrations (call-site renames), not sentinel
  conversions — the affected types keep their own sentinels.
- **`Shape` is an interface with implementations** (`ShapeUnspecifiedImpl`, `CutCornerShape`,
  `RoundedCornerShape`). → Treat `Shape` as a struct/interface companion (pointer semantics) and
  keep each concrete impl's sentinel private; confirm the interface can carry a companion before
  wiring.

## Migration Plan

1. `go.mod`: promote `go-sentinel-helper` from `// indirect` to a direct require.
2. Per package, in order: add package-private ops + `XCompanion`; rename `XUnspecified` →
   `xUnspecified`; delete sentinel-only receiver methods; update that package's call sites.
3. Sweep the wider tree: `go build ./...` → fix each remaining reference (demos, `cmd/`,
   adapters outside scope).
4. Migrate the 17 `pkg/sentinel` importers to `go-sentinel-helper` per the D3 mapping
   (`TakeOrElseString` → `StringsCompanion.TakeOrElse`, `StringValueUnspecified` /
   `StringUnspecified` → `StringsCompanion.Unspecified()`, `shadow.go`'s three-arg
   `TakeOrElse` → per-field companion calls); add the `stringutils` import and drop the
   `pkg/sentinel` import in each.
5. Delete `pkg/sentinel/` (all 7 files) once no importer references it; confirm
   `grep -r "pkg/sentinel" --include=*.go .` returns nothing.
6. Add conformance tests; run `go test ./...`, then `go doc -all <pkg> | grep -i unspecified`
   for every in-scope package (expect no exported sentinel).
7. Update `docs/sentinel_pattern.md` and `AGENTS.md` with the core-surface contract and the
   out-of-scope list, and note that `pkg/sentinel` is gone in favour of the library.

**Rollback:** the change is one revertable unit; no data migration, no persisted state.

## Open Questions

- Should `compose/ui/next/**` and `compose/ui/text/**`'s parallel `TextStyle`/`SpanStyle`/
  `ParagraphStyle` (the `next/` mirror) be unified onto the same companion as their non-`next`
  counterparts, or keep separate companions for now? (Scope says "include", but unification vs.
  parallel registration is undecided.)
- Do `compose/foundation/**` and `compose/material3/**` follow in a second change for their *own*
  sentinels, and if so, does `fontWeight`/`Color` re-export in `material3/alias.go` change the
  ordering? (Note: their *imports* of `pkg/sentinel` are migrated now, regardless.)
