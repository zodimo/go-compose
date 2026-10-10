# Sentinel Patterns — Deep Reference

Disclosed reference for the `go-sentinel-helper` skill. Read this when you are choosing a
sentinel representation, writing a struct companion, implementing option-based `Copy`, or
checking a package against the contract.

## Contents

- [1. Semantic levels of "unspecified"](#1-semantic-levels-of-unspecified)
- [2. Pattern 2-A — Primitive / value type](#2-pattern-2-a--primitive--value-type)
- [3. Pattern 2-B — Packed primitive with Copy options](#3-pattern-2-b--packed-primitive-with-copy-options)
- [4. Pattern 2-C — Complex objects (pointer + singleton)](#4-pattern-2-c--complex-objects-pointer--singleton)
- [5. Pattern 2-D — Type-safe struct wrappers](#5-pattern-2-d--type-safe-struct-wrappers)
- [6. Sentinel choice guide](#6-sentinel-choice-guide)
- [7. Anti-patterns](#7-anti-patterns)
- [8. Performance and verification](#8-performance-and-verification)
- [9. COMPANION_PACKAGE_CONTRACT](#9-companion_package_contract)

---

## 1. Semantic levels of "unspecified"

Three distinct states get conflated unless you name them. Keeping them separate is what lets a
theme system, a config layer, or an API patch express "leave this alone" without a pointer
chase at every level.

| Level | Value | Meaning | Typical use |
|---|---|---|---|
| Object absent | `nil` | "Nothing provided" | caller omitted the parameter |
| Object present, all fields deferred | `xUnspecified` singleton | "Use ambient for everything" | whole-object default |
| Field absent | `xFieldUnspecified` | "Use ambient for this field" | field-level override |

The sentinel at every level is **unexported**; the only public way to reach it is
`XCompanion.Unspecified()`. Uppercase `XUnspecified` names below are shorthand for the concept,
not exported identifiers.

```go
Text("hi")                                       // nil      -> theme
Text("hi", WithTextStyle(textStyleUnspecified))  // singleton-> theme
Text("hi", WithTextStyle(&TextStyle{size: 20}))  // partial  -> theme colour + 20
```

The structure companion collapses `nil` into the singleton via `Coalesce`, so the two top rows
behave identically in the merge path while staying distinguishable at the API boundary.

---

## 2. Pattern 2-A — Primitive / value type

The value lives in a primitive; the sentinel is a bit-pattern impossible in the domain.

```go
type Dp float32

var f32 = floatutils.Float32Companion

var dpUnspecified = Dp(f32.Unspecified())

func isSpecifiedDp(d Dp) bool   { return f32.IsSpecified(float32(d)) }
func isUnspecifiedDp(d Dp) bool { return f32.IsUnspecified(float32(d)) }

func takeOrElseDp(d, def Dp) Dp {
    if !isSpecifiedDp(d) {
        return def
    }
    return d
}

func mergeDp(a, b Dp) Dp {
    if isSpecifiedDp(b) {
        return b
    }
    return a
}

func stringDp(d Dp) string { return fmt.Sprintf("Dp{%v}", float32(d)) }
func sameDp(a, b Dp) bool  { return float32(a) == float32(b) }
func semanticEqualDp(a, b Dp) bool {
    return f32.SemanticEqual(float32(a), float32(b))
}
func equalDp(a, b Dp) bool { return semanticEqualDp(a, b) }
func copyDp(d Dp) Dp       { return d }

var DpCompanion = companions.NewScalarCompanion(
    func() Dp { return dpUnspecified },
    stringDp, isSpecifiedDp, isUnspecifiedDp,
    takeOrElseDp, mergeDp, sameDp, semanticEqualDp, equalDp, copyDp,
)
```

Delegating to `Float32Companion` for the predicate and semantic-equality operations is
deliberate: NaN handling is subtle enough that one implementation should own it. Note the
sentinel is a package-level `var` here rather than a `const` — `math.NaN()` is not a constant
expression, so it cannot be `const`. It still does not allocate, because it is initialised once
at package init and copied by value.

**Guarantee**: the sentinel lives in `.data`/`.rodata`, is copied into a register, and never
touches the heap.

---

## 3. Pattern 2-B — Packed primitive with Copy options

Same as 2-A, but the type packs several components, so a plain `Copy` could not distinguish
"set this component to zero" from "leave it alone". Functional options backed by sentinel
fields resolve it: each option starts at the sentinel, and you walk backward over the
constructors.

```go
type ColorCopyOptions struct{ Alpha, Red, Green, Blue float32 }
type ColorCopyOption func(*ColorCopyOptions) ColorCopyOptions

func CopyWithAlpha(a float32) ColorCopyOption {
    return func(o *ColorCopyOptions) ColorCopyOptions {
        o.Alpha = a
        return *o
    }
}

func (c Color) Copy(opts ...ColorCopyOption) Color {
    f32 := floatutils.Float32Companion
    o := ColorCopyOptions{
        Alpha: f32.Unspecified(),
        Red:   f32.Unspecified(),
        Green: f32.Unspecified(),
        Blue:  f32.Unspecified(),
    }
    for _, opt := range opts {
        opt(&o)
    }
    return NewColor(
        f32.TakeOrElse(o.Alpha, c.Alpha()),   // option wins if specified
        f32.TakeOrElse(o.Red, c.Red()),
        f32.TakeOrElse(o.Green, c.Green()),
        f32.TakeOrElse(o.Blue, c.Blue()),
        colorspace.Get(c.ColorSpaceId()),
    )
}
```

Because the type is still a value type, it stays a **ScalarCompanion** — `Copy` is the only
operation whose shape changes. Option slices are fine here because `Copy` is not a hot path;
the same construct in `Merge` would allocate.

---

## 4. Pattern 2-C — Complex objects (pointer + singleton)

Structs of optional fields are pointers. The companion wraps a **singleton** sentinel, and
every helper is nil-safe, so business code never writes `== nil`.

```go
type TextStyle struct {
    color      Color
    fontSize   Dp
    fontWeight FontWeight
}

var textStyleUnspecified = &TextStyle{
    color:      colorUnspecified,
    fontSize:   dpUnspecified,
    fontWeight: fontWeightUnspecified,
}

func isSpecifiedTextStyle(s *TextStyle) bool { return s != nil && s != textStyleUnspecified }
func isUnspecifiedTextStyle(s *TextStyle) bool { return !isSpecifiedTextStyle(s) }

func takeOrElseTextStyle(s, def *TextStyle) *TextStyle {
    if !isSpecifiedTextStyle(s) {
        return def
    }
    return s
}

// Merge: nil-safe and field-by-field.
func mergeTextStyle(a, b *TextStyle) *TextStyle {
    a = coalesceTextStyle(a, textStyleUnspecified)
    b = coalesceTextStyle(b, textStyleUnspecified)
    if a == textStyleUnspecified {
        return b
    }
    if b == textStyleUnspecified {
        return a
    }
    return &TextStyle{
        color:      mergeColor(b.color, a.color),
        fontSize:   mergeDp(b.fontSize, a.fontSize),
        fontWeight: mergeFontWeight(b.fontWeight, a.fontWeight),
    }
}

func stringTextStyle(s *TextStyle) string {
    if !isSpecifiedTextStyle(s) {
        return "TextStyle{Unspecified}"
    }
    return fmt.Sprintf("TextStyle{color: %s, fontSize: %s, fontWeight: %s}",
        s.color, s.fontSize, s.fontWeight)
}

// Same: identity, with nil folded onto the singleton.
func sameTextStyle(a, b *TextStyle) bool {
    if a == nil && b == nil {
        return true
    }
    if a == nil {
        return b == textStyleUnspecified
    }
    if b == nil {
        return a == textStyleUnspecified
    }
    return a == b
}

func semanticEqualTextStyle(a, b *TextStyle) bool {
    a = coalesceTextStyle(a, textStyleUnspecified)
    b = coalesceTextStyle(b, textStyleUnspecified)
    return semanticEqualColor(a.color, b.color) &&
        semanticEqualDp(a.fontSize, b.fontSize) &&
        semanticEqualFontWeight(a.fontWeight, b.fontWeight)
}

func equalTextStyle(a, b *TextStyle) bool {
    if sameTextStyle(a, b) {
        return true
    }
    return semanticEqualTextStyle(a, b)
}

func coalesceTextStyle(ptr, def *TextStyle) *TextStyle {
    if ptr == nil {
        return def
    }
    return ptr
}

// Copy: start from the sentinel so unspecified option fields defer to the source.
func copyTextStyle(s *TextStyle, options ...func(*TextStyle)) *TextStyle {
    if !isSpecifiedTextStyle(s) {
        return textStyleUnspecified
    }
    opt := *textStyleUnspecified
    for _, option := range options {
        option(&opt)
    }
    return &TextStyle{
        color:      takeOrElseColor(opt.color, s.color),
        fontSize:   takeOrElseDp(opt.fontSize, s.fontSize),
        fontWeight: takeOrElseFontWeight(opt.fontWeight, s.fontWeight),
    }
}

var TextStyleCompanion = companions.NewStructCompanion(
    func() *TextStyle { return textStyleUnspecified },
    stringTextStyle, isSpecifiedTextStyle, isUnspecifiedTextStyle,
    takeOrElseTextStyle, mergeTextStyle,
    sameTextStyle, semanticEqualTextStyle, equalTextStyle,
    coalesceTextStyle, copyTextStyle,
)
```

Call `Copy` with mutators:
`TextStyleCompanion.Copy(s, func(o *TextStyle) { o.fontSize = 0 })` — the option sets one
field on a scratch value initialised from the sentinel, and every field it did not touch falls
back to `s`. This is what makes "set to zero" distinguishable from "no change".

**Memory**: an 8-byte pointer on the stack to a singleton in `.data`. Zero heap bytes at the
call site; the single allocation in `copyTextStyle` is paid only when a caller builds a new
style.

---

## 5. Pattern 2-D — Type-safe struct wrappers

A single-field struct around a packed primitive, used when the underlying primitive would allow
unwanted implicit conversions. `TextUnit(24)` becomes a compile error, forcing `Sp(24)`.

```go
// Packs unit type (Sp/Em) and value into 64 bits.
type TextUnit struct{ packed int64 }

var textUnitUnspecified = TextUnit{
    packed: packRaw(int64(TextUnitTypeUnspecified), float32(math.NaN())),
}

func packRaw(unitType int64, v float32) int64 {
    valBits := int64(math.Float32bits(v)) & 0xFFFFFFFF
    return unitType | valBits
}

func Sp(value float32) TextUnit { return TextUnit{packed: packRaw(int64(TextUnitTypeSp), value)} }
func Em(value float32) TextUnit { return TextUnit{packed: packRaw(int64(TextUnitTypeEm), value)} }

func isSpecifiedTextUnit(tu TextUnit) bool {
    return tu.rawType() != int64(TextUnitTypeUnspecified)
}

func semanticEqualTextUnit(a, b TextUnit) bool {
    if a.Type() != b.Type() {
        return false
    }
    if !isSpecifiedTextUnit(a) {
        return !isSpecifiedTextUnit(b)      // two unspecified values are equal
    }
    return floatutils.Float32Equals(a.Value(), b.Value(), floatutils.Float32EqualityThreshold)
}

func sameTextUnit(a, b TextUnit) bool { return a.packed == b.packed }

func equalTextUnit(a, b TextUnit) bool {
    if sameTextUnit(a, b) {
        return true
    }
    return semanticEqualTextUnit(a, b)
}

func copyTextUnit(tu TextUnit) TextUnit { return tu }
```

Because the type is passed by value with no nil possible, it registers a **ScalarCompanion** —
`Coalesce` is not needed.

**Merge and Copy are implementation-dependent** for packed types. Choose by semantics:

- **Whole-value replacement** (atomic types like `TextUnit`): `if isSpecified(b) { return b }; return a`.
- **Component overlay** (struct option types like `Color`): overlay only the specified fields,
  as in §3.

Reference: `sentinel/boolutils` uses this pattern for its tri-state `BooleanValue`, then keeps
domain methods (`IsTrue`, `IsFalse`, `Bool`, `BoolOrElse`, `Format`) on the type because those
are the reason the wrapper exists. Canonical operations still go through
`BooleanValueCompanion`.

---

## 6. Sentinel choice guide

| Situation | Sentinel | Why |
|---|---|---|
| Float, `0` is common | `NaN` | `0` stays meaningful |
| Int, `0` is common | `math.MinInt` | `0` stays meaningful; no NaN available |
| String, `""` is meaningful | `"\x00unspecified"` | NUL-prefixed, impossible in user data |
| String, `""` never meaningful | `""` | cheapest; the zero value is the sentinel |
| Bool | tri-state enum | both `true` and `false` are valid |
| Enum | `iota` member | reserve one ordinal |
| Packed primitive | a reserved type tag | `TextUnitTypeUnspecified` |

The test to apply: **can the domain ever legitimately produce this value?** If yes, it is not a
sentinel. `0` as a sentinel is the most common mistake, because it works in testing (where
nobody sets zero) and fails in production (where somebody does).

---

## 7. Anti-patterns

Reject these on sight:

```go
type Color interface{ IsSpecified() bool }        // interface forces an allocation
var unspecifiedColor *Color = nil                 // nil interface -> panic on method call
func (ts *TextStyle) IsSpecified() bool           // method on a possibly-nil receiver
func (ts *TextStyle) TakeOrElse(def *TextStyle) *TextStyle  // same problem
func TakeOrElseColor(block func() Color) Color    // lambda escapes to the heap
type T struct { isSpecified bool }                // flag field costs a second field
var TextStyleUnspecified = &TextStyle{...}        // exported sentinel: a 2nd source of truth
const ColorUnspecified = 0                        // exported sentinel: same problem
// exporting the ten funcs                        // only the companion may be exported
```

Each looks defensible in isolation. The method-on-receiver ones are the trap: they read nicer
than a free function, and they are exactly what you cannot do when `nil` is a valid input. A
`func()` parameter in a signature that runs per frame is the other silent cost.

---

## 8. Performance and verification

The contract is zero heap bytes at the call site:

| Case | Cost |
|---|---|
| Scalar sentinel | 0 B, ~2 ns, inlined |
| Struct `nil` | 0 B, ~5 ns |
| Struct custom | 24 B struct alloc, ~20 ns, only when constructing |

The companion itself is a package-level interface value; its methods are indirect calls that do
not inline. Keep the underlying funcs (not the interface) in hot paths when inlining matters.

```bash
# Benchmark allocation counts
go test -bench=. -benchmem

# Escape analysis — look for "does not escape" and "can inline"
go test -gcflags="-m" ./... 2>&1 | grep -E "does not escape|can inline|escapes to heap"

# Whole-module escape check
go build -gcflags="-m" ./... 2>&1 | grep "sentinel/"
```

An unexpected allocation almost always traces to one of two causes: a `func()` parameter on a
hot path, or a sentinel built inside a closure instead of at package level.

---

## 9. COMPANION_PACKAGE_CONTRACT

```go
// COMPANION_PACKAGE_CONTRACT
// Every sentinel helper package exports exactly ONE public value: its Companion.
//
//   var XCompanion companions.ScalarCompanion[T]  // value types
//   var XCompanion companions.StructCompanion[T]  // pointer types
//
// The companion MUST be constructed from these ten package-private funcs:
//   func xUnspecified() T                     // sentinel / singleton
//   func stringX(T) string                    // stringification
//   func isSpecifiedX(T) bool                 // predicate
//   func isUnspecifiedX(T) bool               // inverse predicate
//   func takeOrElseX(src, def T) T            // fallback
//   func mergeX(a, b T) T                     // composition
//   func sameX(a, b T) bool                   // identity
//   func semanticEqualX(a, b T) bool          // semantic equality
//   func equalX(a, b T) bool                  // composed equality
//   func copyX(...) ...                       // scalar: Copy(T) T
//                                             // struct: Copy(T, ...func(T)) T
// Plus, for StructCompanion only:
//   func coalesceX(ptr, def T) T              // nil coalescing
//
// Additional symbols are allowed but MUST be documented and purposeful.
// A sentinel value MUST NEVER be exported — not as a var, const, or func.
// Sentinels MUST be package-level constants or variables, package-private, reachable
// only through Companion.Unspecified().
// No function may accept func() parameters in hot paths (forces heap escape).
// No method may have a pointer receiver that checks for nil.
// END_CONTRACT
```

Argument order for the two constructors, verbatim from `sentinel/companions/companions.go`:

```go
func NewScalarCompanion[T any](
    unspecified func() T,
    stringer func(T) string,
    isSpecified func(T) bool,
    isUnspecified func(T) bool,
    takeOrElse func(src, def T) T,
    merge func(a, b T) T,
    same func(a, b T) bool,
    semanticEqual func(a, b T) bool,
    equal func(a, b T) bool,
    copy func(s T) T,
) ScalarCompanion[T]

func NewStructCompanion[T any](
    unspecified func() T,
    stringer func(T) string,
    isSpecified func(T) bool,
    isUnspecified func(T) bool,
    takeOrElse func(src, def T) T,
    merge func(a, b T) T,
    same func(a, b T) bool,
    semanticEqual func(a, b T) bool,
    equal func(a, b T) bool,
    coalesce func(ptr, def T) T,
    copy func(s T, options ...func(T)) T,
) StructCompanion[T]
```
