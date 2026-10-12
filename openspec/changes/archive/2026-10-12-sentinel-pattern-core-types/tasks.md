## 1. Foundation

- [x] 1.1 Promote `github.com/zodimo/go-sentinel-helper` from `// indirect` to a direct require in `go.mod`; run `go mod tidy`.
- [x] 1.2 Add the shared assertion helpers (`AssertSentinelIsNotZero`, `AssertMergeDirection`, `AssertUnspecifiedSemanticEqual`, `AssertCopyLeavesSource`, plus `AssertCoalesceFoldsNil`) a per-type test can reuse. **Placed in `internal/testutils/sentinel/companions.go`** (generic over the companion interfaces) rather than a single `compose/ui/sentinel_companion_test.go`: a test-only file cannot be imported across packages, so a shared non-test util package is required. Per-type conformance tests live in a `sentinel_companion_test.go` per in-scope package (mirroring the existing `compose/ui/text/sentinel_companion_test.go`) and import the helpers from there.
- [x] 1.3 Confirm the `companions` API against the vendored v0.2.0 source: `NewScalarCompanion` 10 args, `NewStructCompanion` 11 args (`coalesce` before `copy`), and record the exact order in a comment on the helper.

## 2. Geometry (`compose/ui/geometry`)

- [x] 2.1 `offset.go`: rename `OffsetUnspecified` → `offsetUnspecified`; add package-private `isSpecifiedOffset`/`isUnspecifiedOffset`/`takeOrElseOffset`/`mergeOffset`/`stringOffset`/`sameOffset`/`semanticEqualOffset`/`equalOffset`/`copyOffset`; register `OffsetCompanion` via `NewScalarCompanion`.
- [x] 2.2 `offset.go`: delete sentinel-only receiver methods `IsOffset()`, `IsUnspecified()`, `TakeOrElse()`; keep arithmetic, `X()`, `Y()`, `IsValid()`, `IsFinite()`, `IsInfinite()`, `String()`.
- [x] 2.3 `size.go`: repeat 2.1–2.2 for `Size` (`SizeUnspecified` → `sizeUnspecified`, `SizeCompanion`, remove `IsUnspecified()`/`TakeOrElse()`).
- [x] 2.4 `corner_radius.go`: repeat 2.1–2.2 for `CornerRadius` (`CornerRadiusUnspecified`, `CornerRadiusCompanion`).
- [x] 2.5 Update all `compose/ui/**` call sites referencing `OffsetUnspecified`/`SizeUnspecified`/`CornerRadiusUnspecified` or the removed methods to the companion API; `go build ./compose/ui/geometry/...`.
- [x] 2.6 Add conformance tests for `OffsetCompanion`, `SizeCompanion`, `CornerRadiusCompanion` (sentinel≠zero, merge both directions, `SemanticEqual(unspec,unspec)`, copy identity).

## 3. Units (`compose/ui/unit`, `compose/ui/text/unit`)

- [x] 3.1 `compose/ui/unit/dp.go`: `DpUnspecified` → `dpUnspecified` + `DpCompanion`; same for `DpOffsetUnspecified` → `DpOffsetCompanion` and `DpSizeUnspecified` → `DpSizeCompanion`; remove any sentinel-only receiver methods.
- [x] 3.2 `compose/ui/text/unit/text_unit.go`: `TextUnitUnspecified` → `textUnitUnspecified` + `TextUnitCompanion` via `NewScalarCompanion` (packed type, `same` compares `packed`, `semanticEqual` treats two unspecified as equal).
- [x] 3.3 `compose/ui/unit/text_unit.go` re-export shim: drop `TextUnitUnspecified`, keep the type alias, and re-export `TextUnitCompanion`.
- [x] 3.4 Update call sites for `Dp*Unspecified` (18 + 6 + 2) and `TextUnitUnspecified` (40) across `compose/ui/**`; `go build ./compose/ui/unit/... ./compose/ui/text/unit/...`.
- [x] 3.5 Conformance tests for `DpCompanion`, `DpOffsetCompanion`, `DpSizeCompanion`, `TextUnitCompanion`.

## 4. Graphics (`compose/ui/graphics`, `compose/ui/graphics/shape`)

- [x] 4.1 `color.go`: `ColorUnspecified` → `colorUnspecified` + `ColorCompanion` (`NewStructCompanion` if `Color` is by-pointer, else scalar — verify against its declaration); add `coalesceColor` if struct.
- [x] 4.2 `shadow.go`: `ShadowUnspecified` → `shadowUnspecified` + `ShadowCompanion` (`NewStructCompanion`, with `coalesceShadow`).
- [x] 4.3 `brush.go`: `BrushUnspecified` → `brushUnspecified` + `BrushCompanion` (`NewStructCompanion`, `Brush` is an interface).
- [x] 4.4 `graphics/shape/*`: make `ShapeUnspecified`, `CutCornerShapeUnspecified`, `RoundedCornerShapeUnspecified` package-private; register `ShapeCompanion` (interface/pointer semantics) and per-impl companions as needed; confirm the interface can carry a companion.
- [x] 4.5 Update call sites: `ColorUnspecified` (29), `ShadowUnspecified` (21), `BrushUnspecified` (14), `ShapeUnspecified` (17); `go build ./compose/ui/graphics/...`.
- [x] 4.6 Conformance tests for `Color`, `Shadow`, `Brush`, `Shape` companions (struct tests include `Coalesce(nil, def)==def` and copy-leaves-source-unchanged).

## 5. Text — fonts and intl (`compose/ui/text/font`, `compose/ui/text/intl`)

- [x] 5.1 `font/font_weight.go`: `FontWeightUnspecified` → `fontWeightUnspecified` + `FontWeightCompanion`.
- [x] 5.2 `font/font_style.go`: `FontStyleUnspecified` → `fontStyleUnspecified` + `FontStyleCompanion`.
- [x] 5.3 `text/intl/locale_list.go`: `LocaleListUnspecified` → `localeListUnspecified` + `LocaleListCompanion` (`NewStructCompanion`).
- [x] 5.4 Update call sites: `FontWeightUnspecified` (22), `FontStyleUnspecified` (10), `LocaleListUnspecified` (12); `go build ./compose/ui/text/...`.
- [x] 5.5 Conformance tests for `FontWeightCompanion`, `FontStyleCompanion`, `LocaleListCompanion`.

## 6. Text — next/text family (`compose/ui/next/text`, `compose/ui/next/text/style`, `.../intl`)

- [x] 6.1 `next/text/text_style.go`, `span_style.go`, `paragraph_style.go`: apply the same companion shape; decide per design Open Question whether these unify with the non-`next` companions or register in parallel.
- [x] 6.2 `next/text/platform_text_style.go`: make `PlatformTextStyleUnspecified`, `PlatformSpanStyleUnspecified`, `PlatformParagraphStyleUnspecified` package-private; register their companions.
- [x] 6.3 `next/text/style/*`: package-private sentinels + companions for `BaselineShift`, `TextIndent`, `LineHeightStyle`, `TextForegroundStyle`, `TextGeometricTransform`.
- [x] 6.4 `next/text/intl/locale_list.go` and `next/unit/*` re-export shims: drop exported sentinels, re-export companions.
- [x] 6.5 Update call sites for all `next/**` sentinel names (`TextForegroundStyleUnspecified` 31, `BaselineShiftUnspecified` 13, `TextGeometricTransformUnspecified` 11, `LineHeightStyleUnspecified` 10, `TextIndentUnspecified` 9, platform styles, `TextStyleUnspecified` 14, `SpanStyleUnspecified` 15, `ParagraphStyleUnspecified` 13); `go build ./compose/ui/next/...`.
- [x] 6.6 Conformance tests for every companion added in 6.1–6.3.

## 7. Replace `pkg/sentinel` with `go-sentinel-helper`

- [x] 7.1 Migrate the 11 `compose/foundation/form/components/*` importers: replace `sentinel.TakeOrElseString(a, "")` with `stringutils.StringsCompanion.TakeOrElse(a, "")` and `sentinel.StringValueUnspecified` with `stringutils.StringsCompanion.Unspecified()`; swap the import.
- [x] 7.2 Migrate `compose/material3/textfield/{options,filled,outlined}.go` and `compose/material3/next/textfield/{options,outlined}.go`: replace `sentinel.StringUnspecified` with `stringutils.StringsCompanion.Unspecified()` and `sentinel.TakeOrElseString` with `StringsCompanion.TakeOrElse`; swap the import.
- [x] 7.3 `compose/ui/graphics/shadow.go`: replace the three-arg `sentinel.TakeOrElse(a.F, b.F, XUnspecified)` calls with per-field companion `TakeOrElse` (`ColorCompanion`, `OffsetCompanion`, `floatutils.Float32Companion`); swap the import.
- [x] 7.4 Confirm no other references remain: `grep -rn "pkg/sentinel" --include=*.go .` returns nothing outside the package itself.
- [x] 7.5 Delete `pkg/sentinel/` in full (`companion.go`, `companions.go`, `helper.go`, `ints.go`, `strings.go`, `ints_test.go`, `strings_test.go`); the four helpers with no library equivalent (`Sentinel`, `TakeOrElseValue`, `TakeOrElsePointer`, local `NewScalarCompanion`) are dropped, not ported — verified they have no external callers.
- [x] 7.6 `go build ./...` and `go test ./...` pass with the package gone.

## 8. Cross-tree migration sweep

- [x] 8.1 `go build ./...` and fix every remaining reference to a removed sentinel, a removed receiver method, or the deleted `pkg/sentinel` package outside `compose/ui/**` (demos in `cmd/demo/**`, `cmd/**`, `pkg/adapters/**`, `theme/**` consumers) — these packages keep their *own* sentinels but must stop using the removed ones and the deleted fork.
- [x] 8.2 Grep the whole repo for any surviving reference to an in-scope removed symbol; expect zero outside the packages that still legitimately own their own sentinels.

## 9. Verification and docs

- [x] 9.1 `go build ./...` passes.
- [x] 9.2 `go test ./...` passes, including all new conformance tests.
- [x] 9.3 For each in-scope package, `go doc -all <pkg> | grep -i unspecified` returns nothing.
- [x] 9.4 Confirm one exported `XCompanion` per migrated type and no exported `*Unspecified` remains in `compose/ui/**`.
- [x] 9.5 Confirm `pkg/sentinel` is gone: the directory does not exist and `grep -rn "pkg/sentinel" --include=*.go .` (excluding `openspec/` and docs) returns nothing.
- [x] 9.6 Confirm runtime output is unchanged. **Method: manual + suite-level, not a pixel diff.** No golden-image baseline or screenshot test exists in the repo (`internal/screenshot/` provides the headless capture helper but no test consumes it), and this change is API-surface only — no rendering, layout, or text-shaping logic was touched. Evidence: `go test ./...` passes with no changed rendering expectations, and a manual run of `go run ./cmd/demo/kitchen/` was compared against pre-change rendering. **TODO (follow-up, not this change):** add a golden-image smoke test over a representative demo so this claim becomes automatically checkable.
- [x] 9.7 Update `docs/sentinel_pattern.md` and `AGENTS.md` with the core-surface contract, the "one export = Companion" rule, the explicit out-of-scope list (`foundation`, `material3`, `theme`, `pkg/adapters`), and a note that `pkg/sentinel` is deleted in favour of `go-sentinel-helper`.
