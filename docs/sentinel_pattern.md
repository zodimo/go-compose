# The Unspecified Pattern: Complete Go Reference  
*Declarative UI Composition with Zero-Allocation Semantics*

---

> **AGENT CONTRACT (revised for `go-sentinel-helper`)**: Every exported type `T` that carries a
> sentinel must expose **exactly one public value — its Companion**:
>
> 1. `var TCompanion = companions.NewScalarCompanion(...)` (value types) or
>    `companions.NewStructCompanion(...)` (pointer/interface types) — the only exported
>    sentinel-family symbol.
>
> Everything else is **package-private**: the reserved "unspecified" value is
> `tUnspecified` (lower-case) and the ten canonical operations are lower-case functions
> (`isSpecifiedT`, `isUnspecifiedT`, `takeOrElseT`, `mergeT`, `stringT`, `sameT`,
> `semanticEqualT`, `equalT`, `copyT`, plus `coalesceT` for struct types). Callers reach them
> only through `TCompanion.Unspecified()`, `TCompanion.IsSpecified(v)`,
> `TCompanion.TakeOrElse(src, def)`, etc.
>
> There is **no exported `TUnspecified`** in the compose core surface. Two public sources of
> truth for "unspecified" are exactly what the contract forbids.

> The library argument order (verified against `go-sentinel-helper` v0.2.0) is fixed:
>
> ```go
> companions.NewScalarCompanion(unspecified, stringer, isSpecified, isUnspecified,
>     takeOrElse, merge, same, semanticEqual, equal, copy)                       // 10 args
> companions.NewStructCompanion(unspecified, stringer, isSpecified, isUnspecified,
>     takeOrElse, merge, same, semanticEqual, equal, coalesce, copy)             // 11 args
> ```
>
> `isUnspecified` is the **fourth** argument. Transposing `isSpecified`/`isUnspecified`
> compiles cleanly and inverts the semantics silently — every companion has a conformance test
> asserting `IsSpecified(Unspecified()) == false` and `IsUnspecified(Unspecified()) == true`.

### Scope of the contract

**In scope (this change, conforming):** `compose/ui/geometry`, `compose/ui/unit`,
`compose/ui/text/unit`, `compose/ui/text/font`, `compose/ui/text/intl`,
`compose/ui/text/style`, `compose/ui/text`, `compose/ui/graphics`,
`compose/ui/graphics/shape`, and the whole `compose/ui/next/**` mirror.

**Out of scope (follow-up change; still export raw sentinels):** `compose/foundation/**`,
`compose/material3/**`, `theme/`, `pkg/adapters/**`. These own *their own* sentinels and are
untouched here — but they no longer import the deleted `pkg/sentinel` fork.

> **Naming carry-over to clean up in the follow-up.** Two out-of-scope packages still export a
> symbol *named* `ColorUnspecified`: `compose/material3/alias.go` and
> `pkg/adapters/material3/color_set.go`. Both are already aliases of the core companion
> (`var ColorUnspecified = graphics.ColorCompanion.Unspecified()`), so they are not a contract
> violation — but the name shadows the core `graphics.ColorUnspecified` this change removed, and
> two different things now answer to `ColorUnspecified` depending on import path. The follow-up
> change that converts `foundation`/`material3` to companions should retire these aliases in
> favour of `graphics.ColorCompanion.Unspecified()` at the call sites.

> **Receiver methods on in-scope types are gone.** Every method that merely re-expressed an
> existing companion (`Color.IsSpecified()`, `Dp.IsUnspecified()`, `TextUnit.IsSpecified()`,
> …) was removed; callers use `XCompanion.IsSpecified(v)` / `IsUnspecified(v)` instead, so the
> companion is the single door even for predicates. Receiver methods that remain belong to
> types with **no companion** (`Linearity`, `SubpixelTextPositioning`, `LineHeightStyleTrim`,
> `LineHeightStyleMode`) or outside the contract (`layout.Axis`).

> **`pkg/sentinel` is deleted.** It was a parallel fork of `go-sentinel-helper` (its own
> `Companion` interfaces, a 9-argument `NewScalarCompanion` missing `isUnspecified`, a
> `StructCompoanion` typo). All importers now use `github.com/zodimo/go-sentinel-helper`
> directly: `sentinel.TakeOrElseString(a, "")` → `stringutils.StringsCompanion.TakeOrElse(a, "")`,
> and `sentinel.StringValueUnspecified` / `sentinel.StringUnspecified` →
> `stringutils.StringsCompanion.Unspecified()`.

---

## Comments on MergeT and CopyT
- They have the same logic, but different levelt of granularity
- MergeT is used for composition, CopyT is used for mutation
- Prefer incoming Specified values over current values

## 0. Quick Glance Cheat-Sheet

| Type Family | Go Declaration | Sentinel | Is-Specified Check | Alloc/Cost | Use-Case |
|-------------|----------------|----------|--------------------|------------|----------|
| **1-A Primitive** | `type Dp float32` | `var DpUnspecified = Dp(NaN)` | `d.IsSpecified()` | **0 B / 2 ns** | dp, simple units |
| **1-B Packed Primitive** | `type Color uint64` | `const ColorUnspecified = 0` | `c != ColorUnspecified` | **0 B / 2 ns** | color with copy options |
| **1-C Complex** | `type TextStyle struct{ ... }` | `var TextStyleUnspecified = &TextStyle{ ... }` | `IsSpecifiedTextStyle(s)` | **0 B if nil / 24 B if custom** | style, modifier |
| **1-D Type-Safe Wrapper** | `type TextUnit struct{ packed int64 }` | `var TextUnitUnspecified = TextUnit{...}` | `tu.IsSpecified()` | **0 B / 2 ns** | packed unit with type safety |

---

## 1. The Four Sentinel Patterns

Choose **once per type**—never mix.

### 1-A Primitive / Value Types (zero-allocation sentinel)

```go
// unit.Dp

import "github.com/zodimo/go-compose/pkg/floatutils"


// The value is stored as a float32.
type Dp float32

// 1. `TUnspecified` – sentinel 
var DpUnspecified = Dp(floatutils.Float32Unspecified)

// 2. `IsSpecified` – predicate
func (d Dp) IsUnspecified() bool {
	return floatutils.IsSpecified(d)
}

// 3. `TakeOrElse` – 2-param fallback
func (d Dp) TakeOrElse(def Dp) Dp {
	if d.IsSpecified() {
		return d
	}
	return def
}

// 4. `Merge` – composition merge
func MergeDp(a, b Dp) Dp {
	if b.IsSpecified() {
		return b
	}
	return a
}

// 5. `String` – stringification
func (d Dp) String() string {
	if d.IsUnspecified() {
		return "Dp{Unspecified}"
	}
	return fmt.Sprintf("Dp{%.1f}", d)
}

// 6. `Coalesce` – nil coalescing
// Not applicable

// 7. `Same` – identity
// Not applicable

// 8. `SemanticEqual` – semantic equality
// Not applicable

// 9. `Equal` – semantic equality
func (d Dp) Equal(other Dp) bool {
	return d == other
}

// 10. `Copy` – copy
func (d Dp) Copy() Dp {
	return d
}
```

**Sentinel choice guide**  
- `0` → OK when 0 is rare   
- `floatutils.Float32Unspecified, floatutils.Float64Unspecified` for floats → use when 0 is common
- `math.MinInt,math.MaxInt,` for integers →   
- `iota + sentinel` → for enums

**Guarantees**: lives in `.rodata`, copied into register/stack → **zero heap bytes**.

---

### 1-B Packed Primitive / Packed Value Types (zero-allocation sentinel)

```go
//Same as 1-A, but with a different Copy method

type ColorCopyOptions struct {
	Alpha, Red, Green, Blue float32
}
type ColorCopyOption func(*ColorCopyOptions) ColorCopyOptions

func CopyWithAlpha(alpha float32) ColorCopyOption {
	return func(o *ColorCopyOptions) ColorCopyOptions {
		o.Alpha = alpha
		return *o
	}
}

func CopyWithRed(red float32) ColorCopyOption {
	return func(o *ColorCopyOptions) ColorCopyOptions {
		o.Red = red
		return *o
	}
}

func CopyWithGreen(green float32) ColorCopyOption {
	return func(o *ColorCopyOptions) ColorCopyOptions {
		o.Green = green
		return *o
	}
}

func CopyWithBlue(blue float32) ColorCopyOption {
	return func(o *ColorCopyOptions) ColorCopyOptions {
		o.Blue = blue
		return *o
	}
}

// Copy creates a new color with modified components.
func (c Color) Copy(opts ...ColorCopyOption) Color {
	var o ColorCopyOptions = ColorCopyOptions{
		Alpha: floatutils.Float32Unspecified,
		Red:   floatutils.Float32Unspecified,
		Green: floatutils.Float32Unspecified,
		Blue:  floatutils.Float32Unspecified,
	}
	for _, opt := range opts {
		opt(&o)
	}

	id := c.ColorSpaceId()
	space := colorspace.Get(id)
	return NewColor(
		floatutils.TakeOrElse(o.Alpha, c.Alpha()),
		floatutils.TakeOrElse(o.Red, c.Red()),
		floatutils.TakeOrElse(o.Green, c.Green()),
		floatutils.TakeOrElse(o.Blue, c.Blue()),
		space,
	)
}

```



### 1-C Complex Object Types (nullable pointer + singleton)

See section below for full example.

---

### 1-D Type-Safe Struct Wrappers (zero-allocation, compile-time safety)

Use this pattern when you need to **prevent implicit type conversions** while maintaining **zero-allocation value semantics**. This is for single-field struct wrappers around packed primitives.

**When to use**:
- The underlying type is a primitive (`int64`, `uint64`) but allows unwanted implicit conversions
- You want `Foo(24)` to be a compile error, forcing users to use `NewFoo(24)`
- The type is semantically atomic (not a collection of optional fields)
- Zero-allocation is critical

**Example: TextUnit**

The contract shape for a packed value type: a package-private sentinel, package-private operations,
and one exported `XCompanion`. There is no exported `TextUnitUnspecified` and no sentinel-only
receiver method — callers reach everything through the companion.

```go
// TextUnit packs unit type (Sp/Em) and float value into 64 bits.
// Struct wrapper prevents `TextUnit(24)` - must use `Sp(24)` or `Em(1.5)`.
type TextUnit struct {
	packed int64
}

// 1. Sentinel – package-private value.
// Kotlin uses a packed value with NaN for Unspecified.
var textUnitUnspecified = TextUnit{packed: packRaw(int64(TextUnitTypeUnspecified), float32(math.NaN()))}

// Internal packing function
func packRaw(unitType int64, v float32) int64 {
	valBits := int64(math.Float32bits(v)) & 0xFFFFFFFF
	return unitType | valBits
}

// Constructors (the ONLY way to create valid TextUnits)
func Sp(value float32) TextUnit {
	return TextUnit{packed: packRaw(int64(TextUnitTypeSp), value)}
}

func Em(value float32) TextUnit {
	return TextUnit{packed: packRaw(int64(TextUnitTypeEm), value)}
}

// 2. Operations – package-private functions, in the go-sentinel-helper order.
func isSpecifiedTextUnit(tu TextUnit) bool {
	return tu.rawType() != int64(TextUnitTypeUnspecified)
}

func isUnspecifiedTextUnit(tu TextUnit) bool {
	return tu.rawType() == int64(TextUnitTypeUnspecified)
}

func takeOrElseTextUnit(tu, def TextUnit) TextUnit {
	if isSpecifiedTextUnit(tu) {
		return tu
	}
	return def
}

// Merge: whole-value replacement for atomic packed types.
func mergeTextUnit(a, b TextUnit) TextUnit {
	if isSpecifiedTextUnit(b) {
		return b
	}
	return a
}

func stringTextUnit(tu TextUnit) string {
	if isUnspecifiedTextUnit(tu) {
		return "TextUnit{Unspecified}"
	}
	return fmt.Sprintf("TextUnit{%v %s}", tu.Value(), tu.Type())
}

// same/semanticEqual delegate to floatutils; two unspecified values are equal.
func sameTextUnit(a, b TextUnit) bool          { return a == b }
func semanticEqualTextUnit(a, b TextUnit) bool { return a == b } // NaN-aware via packRaw identity
func equalTextUnit(a, b TextUnit) bool         { return a == b }

// Copy: identity for immutable value types.
func copyTextUnit(tu TextUnit) TextUnit { return tu }

// 3. ONE exported sentinel-family symbol: the companion.
// isUnspecified is the FOURTH argument (after isSpecified) — see the library's
// NewScalarCompanion: unspecified, string, isSpecified, isUnspecified,
// takeOrElse, merge, same, semanticEqual, equal, copy.
var TextUnitCompanion = companions.NewScalarCompanion(
	func() TextUnit { return textUnitUnspecified },
	stringTextUnit,
	isSpecifiedTextUnit,
	isUnspecifiedTextUnit,
	takeOrElseTextUnit,
	mergeTextUnit,
	sameTextUnit,
	semanticEqualTextUnit,
	equalTextUnit,
	copyTextUnit,
)
```

Callers never see `textUnitUnspecified`. They write:

```go
tu := TextUnitCompanion.Unspecified()               // the reserved value
if TextUnitCompanion.IsSpecified(tu) { ... }        // predicate
tu = TextUnitCompanion.TakeOrElse(tu, Sp(14))       // default
merged := TextUnitCompanion.Merge(a, b)             // overlay
```

> [!NOTE]
> **Merge and Copy are implementation-dependent.** The packed value may contain multiple logical components. Choose the pattern that matches your type's semantics:
>
> **Option A: Whole-value replacement** (for atomic types like `TextUnit`)
> ```go
> func mergeTextUnit(a, b TextUnit) TextUnit {
>     if isSpecifiedTextUnit(b) { return b }
>     return a
> }
> // Copy: identity function `func(tu TextUnit) TextUnit { return tu }`
> ```
>
> **Option B: Struct options pattern** (for component types like `Color`)
> ```go
> func copyColor(c Color, opts ...ColorCopyOption) Color {
>     o := ColorCopyOptions{Alpha: floatutils.Float32Unspecified, Red: floatutils.Float32Unspecified, ...}
>     for _, opt := range opts { opt(&o) }
>     return NewColor(
>         floatutils.TakeOrElse(o.Alpha, c.Alpha()),
>         floatutils.TakeOrElse(o.Red, c.Red()),
>         ...
>     )
> }
> // Merge: use copyColor to overlay components
> ```

**Key differences from Pattern 1-A/1-B**:
- Uses `struct { packed T }` instead of `type T primitive`
- Constructors are **mandatory** - no implicit conversion from primitives
- Sentinel is a package-private struct value, not a primitive constant

**Key differences from Pattern 1-C**:
- Passed by **value**, not pointer
- Sentinel is a **value**, not a pointer singleton
- Registers via `NewScalarCompanion` (not `NewStructCompanion`)
- `Coalesce` is not needed (no nil possible)

**Guarantees**: Same as Pattern 1-A - lives in registers/stack → **zero heap bytes**.

**Re-export note**: some types (e.g. `TextUnit`) are defined in an inner package and
re-exported by a sibling. The companion is re-exported as an alias —
`compose/ui/unit/text_unit.go` has `var TextUnitCompanion = textunit.TextUnitCompanion` — so
`unit.TextUnitCompanion` and `text/unit.TextUnitCompanion` are the **same value**, not two
companions. The `next/` mirror (`compose/ui/next/unit`, `compose/ui/next/text/unit`) repeats
this: it is the *other* family's companion, independent of the non-`next` one.

---

```go
type TextStyle struct {
	color      Color
	fontSize   Dp
	fontWeight FontWeight
}

// 1. `TUnspecified` – sentinel is a singleton
var TextStyleUnspecified = &TextStyle{
	color:      ColorUnspecified,
	fontSize:   DpUnspecified,
	fontWeight: FontWeightUnspecified,
}

// 2. `IsSpecifiedT` Package-level helpers (NO methods on *TextStyle)
func IsSpecifiedTextStyle(style *TextStyle) bool {
	return style != nil && style != TextStyleUnspecified
}

// 3. `TakeOrElseT` Package-level helpers (NO methods on *TextStyle)
func TakeOrElseTextStyle(style, defaultStyle *TextStyle) *TextStyle {
	if style == nil || style == TextStyleUnspecified {
		return defaultStyle
	}
	return style
}

// 4. `MergeT` – composition merge, package-level function
func MergeTextStyle(a, b *TextStyle) *TextStyle {
	a = CoalesceTextStyle(a, TextStyleUnspecified)
	b = CoalesceTextStyle(b, TextStyleUnspecified)

	if a == TextStyleUnspecified { return b }
	if b == TextStyleUnspecified { return a }

	// Both are custom: allocate new merged style
	return &TextStyle{
		color:      b.color.TakeOrElse(a.color), // primitive struct
		fontSize:   b.fontSize.TakeOrElse(a.fontSize), // primitive struct
		fontWeight: b.fontWeight.TakeOrElse(a.fontWeight), // primitive struct
	}
}
// 5. `StringT` – stringification, package-level function  
func StringTextStyle(style *TextStyle) string {
	if !IsSpecifiedTextStyle(style) {
		return "TextStyle{Unspecified}"
	}

	return fmt.Sprintf(
		"TextStyle{color: %s, fontSize: %s, fontWeight: %s}",
		style.color,
		style.fontSize,
		style.fontWeight,
	)
}
// 6. `CoalesceT` – nil coalescing, package-level function
func CoalesceTextStyle(ptr, def *TextStyle) *TextStyle {
	if ptr == nil {
		return def
	}
	return ptr
}

// 7. `SameT` – identity, package-level function
func SameTextStyle(a, b *TextStyle) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return b == TextStyleUnspecified
	}
	if b == nil {
		return a == TextStyleUnspecified
	}
	return a == b
}

// 8. `SemanticEqualT` – semantic equality, package-level function
func SemanticEqualTextStyle(a, b *TextStyle) bool {

	a = CoalesceTextStyle(a, TextStyleUnspecified)
	b = CoalesceTextStyle(b, TextStyleUnspecified)

	return EqualColor(a.spanStyle, b.spanStyle) &&
		EqualParagraphStyle(a.paragraphStyle, b.paragraphStyle)
}

// 9. `EqualT` – semantic equality, package-level function
func EqualTextStyle(a, b *TextStyle) bool {
	if !SameTextStyle(a, b) {
		return SemanticEqualTextStyle(a, b)
	}
	return true
}

// 10. `CopyT` – copy, package-level function
func CopyTextStyle(ts *TextStyle, options ...TextStyleOption) *TextStyle {
	opt := *TextStyleUnspecified

	for _, option := range options {
		option(&opt)
	}
	return &TextStyle{
		color:      TakeOrElseColor(opt.color, ts.color),
		fontSize:   TakeOrElseDp(opt.fontSize, ts.fontSize),
		fontWeight: TakeOrElseFontWeight(opt.fontWeight, ts.fontWeight),
	}
}

```

**Memory layout**: pointer on stack (8 B) → singleton in `.data` → **zero heap bytes at call site**.

---

## 2. Semantic Levels of "Unspecified"

| Level | Value | Meaning | Typical Use |
|-------|-------|---------|-------------|
| **Object absent** | `nil` | "No style provided" | function parameter default , will be coalesced to unspecified |
| **Object present, all fields deferred** | `TextStyleUnspecified` | "Use theme for everything" | partial merge base |
| **Field absent** | `DpUnspecified` (inside struct) | "Use ambient for this field" | field-level override |

**Example**:
```go
Text("hi")                                    					// nil → use theme 100 %
Text("hi", WithTextStyle(TextStyleUnspecified))                 // singleton → use theme 100 %
Text("hi", WithTextStyle(&TextStyle{fontSize: 20}))             // partial → theme color + 20 sp
```

---

## 3. Public API for Complex Types (Nil-Safe, No Scattered Checks)

Expose **only** these four functions for complex types—**never** methods on `*T` (pointer receiver):

```go

// 2. `IsSpecifiedT`
func IsSpecifiedTextStyle(style *TextStyle) bool

// 3. `TakeOrElseT`
func TakeOrElseTextStyle(style, defaultStyle *TextStyle) *TextStyle

// 4. `MergeT`
func MergeTextStyle(a, b *TextStyle) *TextStyle

// 5. `StringT`
func StringTextStyle(style *TextStyle) string

// 6. `CoalesceT`
func CoalesceTextStyle(ptr, def *TextStyle) *TextStyle 

// 7. `SameT`
func SameTextStyle(a, b *TextStyle) bool 

// 8. `SemanticEqualT`
func SemanticEqualTextStyle(a, b *TextStyle) bool 

// 9. `EqualT`
func EqualTextStyle(a, b *TextStyle) bool

// Identity (2 ns)
func SameTextStyle(a, b *TextStyle) bool

// Semantic equality (field-by-field, 20 ns)
func EqualTextStyle(a, b *TextStyle) bool

// 10. `CopyT`
func CopyTextStyle(ts *TextStyle, options ...TextStyleOption) *TextStyle 
```

Business code **never writes `== nil`** outside these four helpers.

---


## 3.1 The Copy Pattern (Functional Options)

Use **functional options** backed by sentinel values to implement strictly typed partial updates (Copy). This solves the ambiguity between "set to zero" and "no change".

**Example**: `Shadow.Copy(WithBlurRadius(0))` vs `Shadow.Copy()`

```go
// 1. Options Struct (fields match the type, initialized to Unspecified)
type Shadow struct {
	Color      Color
	Offset     Offset
	BlurRadius float32
}
// 2. Functional Option
type ShadowOption func(*Shadow)

func WithBlurRadius(r float32) ShadowOption {
	return func(o *ShadowOptions) {
		o.BlurRadius = r
	}
}

// 3. Copy Method
func CopyShadow(s *Shadow, options ...ShadowOption) *Shadow {
	opt := *ShadowUnspecified // dereference unspecified

	for _, option := range options {
		option(&opt)
	}

	return &Shadow{
		Color:      opt.Color.TakeOrElse(s.Color),
		Offset:     opt.Offset.TakeOrElse(s.Offset),
		BlurRadius: floatutils.TakeOrElse(opt.BlurRadius, s.BlurRadius),
	}
}
```

## 4. Sentinel Patterns for `int`, `string`, `bool`

When the type cannot hold an actual `NaN`, reserve a **valid bit-pattern** that is **semantically impossible** in your domain.

### 5-A `int` (or `int32`, `int64`, `uint`, …)

```go
type Offset int32
const OffsetUnspecified Offset = math.MinInt32  // 0x80000000
func (o Offset) IsOffset() bool              { return o != OffsetUnspecified }

// Usage: 0 is valid → MinInt32 is sentinel
var off Offset = OffsetUnspecified  // register literal, 0 heap
```

### 5-B `string`

```go
type Profile string
const ProfileUnspecified Profile = ""  // zero value = sentinel (when empty is rare)
func (p Profile) IsSpecified() bool { return p != ProfileUnspecified }
```

If **empty string is meaningful**, reserve a **globally unique** magic value:

```go
const ProfileUnspecified Profile = "\x00unspecified"  // impossible in user data
```

Both are **zero-allocation**: the constant lives in `.rodata`, the field stores a **pointer to that data**.

### 5-C `bool`

```go
type Visible bool
const VisibleUnspecified Visible = false  // zero value = sentinel
func (v Visible) IsVisible() bool { return v != VisibleUnspecified }
```

If you **must** distinguish “explicitly false” from “unspecified”, promote to an enum:

```go
type Visible int8
const (
	VisibleUnspecified Visible = iota  // 0
	VisibleFalse
	VisibleTrue
)
```

---

## 5. Copy-Paste Templates

```go
// Primitive type
type MyUnit float32
const MyUnitUnspecified MyUnit = MyUnit(math.NaN())
func (u MyUnit) IsSpecified() bool { return !math.IsNaN(float64(u)) }

// Complex type
type MyStyle struct { field1 Color; field2 Dp }
var MyStyleUnspecified = &MyStyle{field1: ColorUnspecified, field2: DpUnspecified}

func IsSpecifiedMyStyle(s *MyStyle) bool {
	return s != nil && s != MyStyleUnspecified
}
func TakeOrElseMyStyle(s, def *MyStyle) *MyStyle {
	if s == nil || s == MyStyleUnspecified { return def }
	return s
}
```

---

## 6. Performance Contract

- **Primitive sentinel**: 0 heap bytes, 2 ns, inlined by compiler
- **Complex nil**: 0 heap bytes, 5 ns
- **Complex custom**: 24 B struct alloc, 20 ns (paid only when user creates new style)

Verify with:
```bash
go test -bench=. -gcflags="-m" 2>&1 | grep -E "does not escape|inlined"
```

---

## 7. Anti-Patterns (reject on sight)

```go
❌ type Color interface { IsSpecified() bool }        // interface forces alloc
❌ var UnspecifiedColor *Color = nil                  // nil interface → panic
❌ func (ts *TextStyle) IsSpecified() bool            // method on nil receiver
❌ func (ts *TextStyle) TakeOrElse(def *TextStyle) *TextStyle // REJECT: method on nil receiver
❌ func TakeOrElseColor(block func() Color) Color // lambda escapes in Go
❌ type T struct { isSpecified bool }                 // flag field = double mem
```

---

## 8. Function Signatures (Principle of Least Parameters)

### Public API (2-parameter fallback)
```go
// TakeOrElseTextStyle: public API - clear, no lambda, no 3rd param
func TakeOrElseTextStyle(style, defaultStyle *TextStyle) *TextStyle {
	if style == nil || style == TextStyleUnspecified {
		return defaultStyle
	}
	return style
}
```

### Private Helper (3-parameter generic)
```go
// takeOrElse: generic 3-param helper used INSIDE helpers to merge fields
func takeOrElse[T comparable](a, b, sentinel T) T {
	if a != sentinel {
		return a
	}
	return b
}
```

**Rule**: Public API always exposes **2-parameter** `TakeOrElseT`; the 3-param generic is an **implementation detail**.

---

## 9. Package-Level Contract (Machine-Readable)

```go
// UI_PACKAGE_CONTRACT
// For every exported type T in package, the following symbols MUST exist:
//   const/var TUnspecified  // Sentinel value or singleton pointer
//   func IsSpecifiedT(v *T) bool      // Package-level predicate (never method on *T)
//   func TakeOrElseT(a, b *T) *T // Package-level fallback (never method on *T)
//   func MergeT(a, b *T) *T      // Package-level composition (never method on *T)
//   func StringT(s *T) string      // Package-level predicate (never method on *T)
//   func CoalesceT(ptr, def *T) *T // Package-level fallback (never method on *T)
//   func SameT(a, b *T) bool      // Package-level predicate (never method on *T)
//   func SemanticEqualT(a, b *T) bool      // Package-level predicate (never method on *T)
//   func EqualT(a, b *T) bool      // Package-level predicate (never method on *T)
//   func CopyT(a, b *T) *T      // Package-level predicate (never method on *T)

// Additional symbols (abbreviations, helpers) are allowed but must be documented.
// All sentinel values must be compile-time constants or package-level variables.
// No function may accept func() parameters in hot paths (forces heap escape).
// No method may have a pointer receiver that checks for nil (anti-pattern).
// END_CONTRACT
```

---

## 10. Verification Commands

```bash
# Check for heap escapes (must show "does not escape")
go test -gcflags="-m" -run=ExampleUsage 2>&1 | grep -E "does not escape|escapes to heap"

# Check for inlining (must show "can inline" for hot paths)
go test -gcflags="-m" -run=ExampleUsage 2>&1 | grep "can inline"

# Benchmark allocation counts (should be 0 or 1 per operation)
go test -bench=. -benchmem

# Build with escape analysis for entire package
go build -gcflags="-m" ./... 2>&1 | grep "ui/"
```

---

**Document Version**: 2.2  
**Last Updated**: 2025-12-24