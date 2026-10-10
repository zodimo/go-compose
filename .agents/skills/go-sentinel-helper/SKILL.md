---
name: go-sentinel-helper
description: Work with github.com/zodimo/go-sentinel-helper — the Go library for representing "unspecified"/optional values without pointers or wrapper types. Use when (1) calling an existing companion (StringsCompanion, IntValueCompanion, Float32Companion, BooleanValueCompanion, protobufwrapper.*Companion) to test, merge, default, compare, or copy a value, or (2) defining a new sentinel type and registering its own companion with NewScalarCompanion/NewStructCompanion — including choosing a sentinel value, wiring the ten canonical operations, and validating the result against the package contract. Also use when reviewing code that imports sentinel/*, or when a codebase has values that must distinguish "explicitly set to zero" from "never set".
---

# Go Sentinel Helper

A value that can be "unspecified" is normally modelled with a pointer (`*int`), a wrapper
(`sql.NullString`), or a `hasX bool` field. All three cost either a heap allocation, a second
field, or an `== nil` check scattered through business code. This library replaces them with a
single reserved bit-pattern per type plus a **Companion** — one generic value bundling the
operations you need — so callers never write `== nil` and never allocate.

Two jobs bring you here. **Using** a companion that already exists, and **defining** one for a
domain type of your own. The rules differ enough to keep separate.

- Using an existing companion → [Use a companion](#use-a-companion)
- Defining a new sentinel type → [Define your own sentinel](#define-your-own-sentinel)
- Reading the full pattern catalog, anti-patterns, or the machine-readable package contract →
  [`references/patterns.md`](references/patterns.md)

---

## The model in one paragraph

Everything is built around ten canonical operations, and **the only public value a helper
package exports is its Companion**:

```go
c := stringutils.StringsCompanion     // a companions.ScalarCompanion[string]

c.Unspecified()          // the reserved sentinel value for this type
c.IsSpecified(v)         // v is a real value (not the sentinel)
c.IsUnspecified(v)       // v is the sentinel
c.String(v)              // human-readable, e.g. `StringValue{"hi"}`
c.TakeOrElse(v, def)     // v if specified, else def
c.Merge(a, b)            // b if b is specified, else a  (composition)
c.Same(a, b)             // identity
c.SemanticEqual(a, b)    // semantic equality
c.Equal(a, b)            // Same fast-path, else SemanticEqual
c.Copy(v)                // scalar: identity. struct: Copy(v, opts...)
```

The distinction that makes the library worth its weight: the sentinel is chosen so it is
**semantically impossible in the domain**. `""` stays a meaningful string because
`stringutils` uses `"\x00unspecified"` as its sentinel; `0` stays a meaningful integer because
`intutils` uses `math.MinInt`. So `IsSpecified(v)` answers "did anyone actually set this?"

`ScalarCompanion[T]` is for value types (`string`, `int`, `float32`, a struct passed by value).
`StructCompanion[T]` is for pointer types (`*TextStyle`, `*wrapperspb.StringValue`) and adds
`Coalesce(ptr, def)` plus option-based `Copy`. See
[`references/patterns.md`](references/patterns.md) for which pattern to reach for.

---

## Use a companion

Pick the package matching the type, call the companion, and never look at the sentinel directly.

| Type you have | Package | Companion |
|---|---|---|
| `string` | `sentinel/stringutils` | `StringsCompanion` |
| `int` | `sentinel/intutils` | `IntValueCompanion` |
| `float32` / `float64` | `sentinel/floatutils` | `Float32Companion` / `Float64Companion` |
| tri-state bool | `sentinel/boolutils` | `BooleanValueCompanion` |
| `*wrapperspb.*Value` | `sentinel/protobufwrapper` | `StringValueCompanion`, `Int32ValueCompanion`, … |

```go
import "github.com/zodimo/go-sentinel-helper/sentinel/stringutils"

c := stringutils.StringsCompanion

c.IsSpecified("hello")          // true
c.IsSpecified("")               // true  — "" is a real value, not the sentinel
c.IsUnspecified("")             // false
c.Merge("a", "b")               // "b"   — b wins when it is specified
c.TakeOrElse("", "default")     // ""    — "" is specified, so the default is ignored
c.TakeOrElse("", "")            // "" — use c.Unspecified() to actually express "unset"
c.String("hi")                  // `StringValue{"hi"}`
```

### Defaulting and layering

`TakeOrElse` applies a default; `Merge` layers an override onto a base, keeping the base
wherever the override is unspecified. This is the whole point of the pattern — the "unset"
case flows through without a single `if field != nil`.

```go
import "github.com/zodimo/go-sentinel-helper/sentinel/intutils"

// override wins only where it was actually set
func Resolve(base, override int) int {
    return intutils.IntValueCompanion.Merge(base, override)
}
```

### Tri-state booleans

`true` and `false` are both valid, so a bit-pattern sentinel is impossible. `boolutils`
promotes `bool` to a three-state `BooleanValue` struct. Canonical operations go through the
companion; a few domain methods live on the type itself.

```go
import "github.com/zodimo/go-sentinel-helper/sentinel/boolutils"

var v = boolutils.BooleanValueFrom(true)   // or boolutils.BooleanValueCompanion.Unspecified()

v.IsTrue()                  // true
v.IsFalse()                 // false (it is explicitly true)
v.IsSpecified()             // true
v.Bool()                    // false when unspecified — the safe read
v.BoolOrElse(false)         // explicit fallback

boolutils.BooleanValueCompanion.Merge(a, b)   // canonical merge
```

Reach for `BoolOrElse` (or the companion's `TakeOrElse`) rather than `Bool()` when "unset"
must be distinguishable from "set to false" — `Bool()` collapses both to `false`.

### Protobuf wrappers

The nine `wrapperspb.*Value` companions are nil-safe struct companions: `nil` is treated as
unspecified and coalesced to a package-private empty singleton. `Copy` deep-copies and applies
your mutators to the *copy*.

```go
import (
    "github.com/zodimo/go-sentinel-helper/sentinel/protobufwrapper"
    "google.golang.org/protobuf/types/known/wrapperspb"
)

c := protobufwrapper.StringValueCompanion

c.IsSpecified(nil)                                       // false
c.Coalesce(nil, c.Unspecified())                         // the singleton, not nil
c.Merge(nil, &wrapperspb.StringValue{Value: "x"}).Value  // "x"
c.Copy(&wrapperspb.StringValue{Value: "x"}, func(v *wrapperspb.StringValue) {
    v.Value = "y"
}).Value                                                 // "y" — original untouched
```

### The three rules that follow from the design

**Define the sentinel package-private — never export it.** Whether the type is one of the
library's own (`stringutils`, `intutils`, `floatutils`, `boolutils`, `protobufwrapper`) or a
custom type of your own, the sentinel is an unexported package-level symbol
(`stringValueUnspecified`, `intValueUnspecified`, `float32Unspecified`, `opacityUnspecified`,
...) or a literal written inline in the package-private funcs. Nothing else. There is no
"export it so callers can compare" exception, no exported `MyTypeUnspecified`, no exported
sentinel const. Every package the pattern touches follows the same shape.

**Obtain the sentinel only through `c.Unspecified()`.** `c.Unspecified()` is the only public
way to get the reserved value back — for the library's companions and for yours. One public
source of truth for "what is unspecified" is how a domain type ends up comparing against the
wrong value, so the companion is deliberately the single door. When your type is built on
another type's sentinel, call that type's companion — e.g.
`floatutils.Float32Companion.Unspecified()` to seed your own unexported sentinel — rather than a
constant, and keep your own sentinel unexported either way:

```go
// right — unexported, seeded from the wrapped type's companion
var opacityUnspecified = Opacity(floatutils.Float32Companion.Unspecified())

// wrong — a second public source of truth for "unspecified"
const OpacityUnspecified = Opacity(0)          // exported sentinel: forbidden
var OpacityUnspecified = Opacity(kNaN)         // exported sentinel: forbidden
```

A custom package exports **exactly one value: its Companion**. If a caller needs the sentinel,
it calls `YourCompanion.Unspecified()`.

**Call the companion, and where the hot path matters, call the underlying func.** The
companion is a package-level interface value, so its methods are indirect calls that the
compiler cannot inline. That is fine for correctness and fine for almost all code — but if you
are in a per-frame render loop, delegate to the concrete helper (e.g. `floatutils.IsInfinite`)
or your own package-private func rather than routing through the interface. See
[Performance](#performance).

---

## Define your own sentinel

The reference implementations are `sentinel/stringutils` (scalar), `sentinel/protobufwrapper`
(struct), and `sentinel/boolutils` (value struct). Read one before starting — this workflow
mirrors them.

### Step 1 — Choose the sentinel representation

Pick exactly one per type; never mix. The choice is driven by whether a bit-pattern exists that
the domain can never legitimately produce.

| Your type | Sentinel | Pattern name | Companion kind |
|---|---|---|---|
| Domain primitive where `0` is common (`type Dp float32`) | `NaN` | Primitive | Scalar |
| Integer | `math.MinInt` | Primitive | Scalar |
| String where `""` is meaningful | `"\x00unspecified"` | Primitive | Scalar |
| Bool (both values valid) | tri-state enum in a struct | Value struct | Scalar |
| Struct of optional fields | pointer + singleton | Complex | Struct |
| Packed value that must block implicit conversion | `struct{ packed int64 }` | Type-safe wrapper | Scalar |

`0` is a fine sentinel when `0` is rare — but the moment `0` becomes a real value, you have
silently broken the pattern. When in doubt, pick a value the domain cannot produce.
[`references/patterns.md`](references/patterns.md) covers all four patterns with worked code.

### Step 2 — Write the ten operations, package-private

Lowercase every one of them. The funcs are the implementation; the companion is the interface.

```go
package opacity

import (
    "fmt"

    "github.com/zodimo/go-sentinel-helper/sentinel/companions"
    "github.com/zodimo/go-sentinel-helper/sentinel/floatutils"
)

// The domain type.
type Opacity float32

// 1. Sentinel — a package-level var (or const for primitives).
var opacityUnspecified = Opacity(floatutils.Float32Companion.Unspecified())

func isSpecifiedOpacity(o Opacity) bool {
    return floatutils.Float32Companion.IsSpecified(float32(o))
}

func isUnspecifiedOpacity(o Opacity) bool {
    return floatutils.Float32Companion.IsUnspecified(float32(o))
}

func takeOrElseOpacity(o, def Opacity) Opacity {
    if !isSpecifiedOpacity(o) {
        return def
    }
    return o
}

// Merge prefers the incoming specified value.
func mergeOpacity(a, b Opacity) Opacity {
    if isSpecifiedOpacity(b) {
        return b
    }
    return a
}

func stringOpacity(o Opacity) string {
    if isUnspecifiedOpacity(o) {
        return "Opacity{Unspecified}"
    }
    return fmt.Sprintf("Opacity{%v}", float32(o))
}

func sameOpacity(a, b Opacity) bool { return a == b }

func semanticEqualOpacity(a, b Opacity) bool {
    return floatutils.Float32Companion.SemanticEqual(float32(a), float32(b))
}

// Equal: Same fast-path, then semantic.
func equalOpacity(a, b Opacity) bool {
    if sameOpacity(a, b) {
        return true
    }
    return semanticEqualOpacity(a, b)
}

// Copy is identity for immutable value types.
func copyOpacity(o Opacity) Opacity { return o }
```

Two things carry real weight here. `Merge` takes `b`'s specified values over `a`'s — get the
direction wrong and every override silently inverts. And `semanticEqual` must treat two
unspecified values as **equal**: both are "unset", and comparing them as NaN ≠ NaN makes every
round-trip equality test fail. For floats this is why you delegate to the float utility rather
than using `==`.

### Step 3 — Register the companion and export nothing else

Argument order is fixed by `NewScalarCompanion`:

```go
// The ONLY export of the package:
var OpacityCompanion = companions.NewScalarCompanion(
    func() Opacity { return opacityUnspecified }, // 1. Unspecified
    stringOpacity,                                // 2. String
    isSpecifiedOpacity,                           // 3. IsSpecified
    isUnspecifiedOpacity,                         // 4. IsUnspecified
    takeOrElseOpacity,                            // 5. TakeOrElse
    mergeOpacity,                                 // 6. Merge
    sameOpacity,                                  // 7. Same
    semanticEqualOpacity,                         // 8. SemanticEqual
    equalOpacity,                                 // 9. Equal
    copyOpacity,                                  // 10. Copy (scalar: T -> T)
)
```

For a pointer type, use `NewStructCompanion`, which takes the same ten in the same positions
plus `coalesceX` before `copyX`, and whose `copy` has the signature
`func(T, ...func(T)) T`:

```go
var StyleCompanion = companions.NewStructCompanion(
    unspecifiedStyle, stringStyle, isSpecifiedStyle, isUnspecifiedStyle,
    takeOrElseStyle, mergeStyle, sameStyle, semanticEqualStyle, equalStyle,
    coalesceStyle,   // 11. Coalesce(ptr, def T) T — nil -> def
    copyStyle,       // 12. Copy(s T, options ...func(T)) T
)
```

Check the exact parameter order against `companions.go` before you rely on memory — a
transposed `isSpecified`/`isUnspecified` compiles cleanly and inverts silently.

### Step 4 — Enforce the package contract

A helper package exports **exactly one value: its Companion** — that holds for your custom
package exactly as it does for the library's. Do not export a sentinel value, a sentinel const,
or a `MyTypeUnspecified` var: `IntValueUnspecified`, `StringValueUnspecified`,
`Float32Unspecified` and any equivalent for a type of your own all belong in the unexported
world, reachable only through `Companion.Unspecified()`. The only extra exports the library
keeps are genuine helpers over non-sentinel concerns: `floatutils` exports
`Float32Infinite`/`FloatInfinite`, `IsInfinite`, `Float32Equals`, `Float64Equals`, and
`boolutils` exports the `BooleanValue` constructors and its tri-state methods. What must never
happen is a domain type exporting two ways to ask "is this specified" — one of which is the
raw sentinel.

Before you call it done, confirm with `go doc -all` that the only new public symbol is the
companion (plus anything you deliberately documented as an extra). Nothing named
`*Unspecified` should appear:

```bash
go doc -all ./sentinel/opacity
go doc -all ./sentinel/opacity | grep -i unspecified   # expect no exported sentinel symbol
```

The full machine-readable contract and the anti-pattern list live in
[`references/patterns.md`](references/patterns.md). The anti-patterns are short and worth
reading once — several (methods on nil receivers, `func()` parameters in hot paths, a
`isSpecified bool` field) look reasonable and quietly cost an allocation or a panic.

### Step 5 — Test the semantics, not the plumbing

The bugs that survive review are semantic: merge direction, unspecified-equals-unspecified, and
copy not actually copying. `sentinel/stringutils/strings_test.go` and
`sentinel/protobufwrapper/*_test.go` are the models. Cover, at minimum:

- `IsSpecified(sentinel) == false` and `IsSpecified(zeroValue) == true` — the case that proves
  your sentinel is not the zero value.
- `Merge(base, sentinel) == base` and `Merge(base, override) == override`.
- `SemanticEqual(unspecified, unspecified) == true`.
- For struct companions: `Coalesce(nil, def) == def`, and `Copy` leaves the source unmodified.

---

## Performance

The pattern's promise is zero heap allocation at the call site:

- **Scalar sentinel** — the constant lives in `.rodata` / a register. ~0 B, inlined.
- **Struct `nil`** — a pointer on the stack.
- **Struct custom** — one allocation, paid only when a caller actually builds a new value.

The one caveat is the companion interface itself: its methods are indirect calls and do not
inline. `references/patterns.md` § Performance has the escape-analysis commands; the short
version is:

```bash
go test -bench=. -benchmem
go test -gcflags="-m" ./... 2>&1 | grep -E "does not escape|can inline"
```

If a benchmark shows an allocation you did not expect, the usual cause is a `func()` parameter
reaching a hot path, or a sentinel that escaped to the heap because it was built inside a
closure instead of at package level.

---

## Reference

- [`references/patterns.md`](references/patterns.md) — the four sentinel patterns with full
  worked code (primitive, packed, complex/pointer, type-safe wrapper), the `Copy` functional-options
  pattern, semantic levels of "unspecified", anti-patterns, and the machine-readable
  `COMPANION_PACKAGE_CONTRACT`.
- The library's own [`docs/sentinel_pattern.md`](../../docs/sentinel_pattern.md) — the
  narrative reference this skill condenses.
