# The Unspecified Pattern: Complete Go Reference  
*Declarative UI Composition with Zero-Allocation Semantics*

---

> **AGENT CONTRACT**: This is a single-file, copy-paste-ready reference. Every exported type `T` must provide exactly **4 symbols**:  
> 1. `TUnspecified` – sentinel (constant for values, singleton for structs)  
> 2. `IsT` – predicate, package-level function  
> 3. `TakeOrElseT` – 2-param fallback, package-level function  
> 4. `MergeT` – composition merge, package-level function  
> No other public symbols are required. All verification commands are in Section 8.

---

## 0. Quick Glance Cheat-Sheet

| Type Family | Go Declaration | Sentinel | Is-Specified Check | Alloc/Cost | Use-Case |
|-------------|----------------|----------|--------------------|------------|----------|
| **Primitive** | `type Color uint64` | `const ColorUnspecified = 0` | `c != ColorUnspecified` | **0 B / 2 ns** | color, dp, unit |
| **Complex** | `type TextStyle struct{ ... }` | `var TextStyleUnspecified = &TextStyle{ ... }` | `IsTextStyle(style)` | **0 B if nil / 24 B if custom** | style, modifier |

---

## 1. The Two Mutually Exclusive Patterns

Choose **once per type**—never mix.

### 1-A Primitive / Value Types (zero-allocation sentinel)

```go
// 1. Declare the value type
type Color uint64

// 2. Pick a sentinel bit-pattern that is never a valid explicit value
const ColorUnspecified Color = 0          // transparent black
const ColorRed       Color = 0xFFFF0000   // normal colours

// 3. Zero-defensive method (value receiver → never nil)
func (c Color) IsColor() bool {
	return c != ColorUnspecified
}

// 4. Zero-allocation fallback
func (c Color) TakeOrElseColor(defaultColor Color) Color {
	if c.IsColor() {
		return c
	}
	return defaultColor
}
```

**Sentinel choice guide**  
- `0` → OK when 0 is rare (Color)  
- `NaN` → use when 0 is common (Dp, TextUnit)  
- `MinInt32` → for packed structs (IntOffset)  
- `iota 0` → for enums (Alignment, FontWeight)

**Guarantees**: lives in `.rodata`, copied into register/stack → **zero heap bytes**.

---

### 1-B Complex Object Types (nullable pointer + singleton)

```go
// 1. Multi-field struct
type TextStyle struct {
	color      Color
	fontSize   Dp
	fontWeight FontWeight
}

// 2. ONE global singleton (static data segment)
var TextStyleUnspecified = &TextStyle{
	color:      ColorUnspecified,
	fontSize:   DpUnspecified,
	fontWeight: FontWeightUnspecified,
}

// 3. Package-level helpers (NO methods on *TextStyle)
func IsTextStyle(style *TextStyle) bool {
	return style != nil && style != TextStyleUnspecified
}

// 4. Package-level helpers (NO methods on *TextStyle)
func TakeOrElseTextStyle(style, defaultStyle *TextStyle) *TextStyle {
	if style == nil || style == TextStyleUnspecified {
		return defaultStyle
	}
	return style
}

// 5. Generic 3-param helper (used inside helpers)
func takeOrElse[T comparable](a, b, sentinel T) T {
	if a != sentinel {
		return a
	}
	return b
}

// 6. Merge / equality helpers
func MergeTextStyle(a, b *TextStyle) *TextStyle {
	a = CoalesceTextStyle(a, TextStyleUnspecified)
	b = CoalesceTextStyle(b, TextStyleUnspecified)

	if a == TextStyleUnspecified { return b }
	if b == TextStyleUnspecified { return a }

	// Both are custom: allocate new merged style
	return &TextStyle{
		color:      takeOrElse(a.color, b.color, ColorUnspecified),
		fontSize:   takeOrElse(a.fontSize, b.fontSize, DpUnspecified),
		fontWeight: takeOrElse(a.fontWeight, b.fontWeight, FontWeightUnspecified),
	}
}

func CoalesceTextStyle(ptr, def *TextStyle) *TextStyle {
	if ptr == nil {
		return def
	}
	return ptr
}
```

**Memory layout**: pointer on stack (8 B) → singleton in `.data` → **zero heap bytes at call site**.

---

## 2. Semantic Levels of "Unspecified"

| Level | Value | Meaning | Typical Use |
|-------|-------|---------|-------------|
| **Object absent** | `nil` | "No style provided" | function parameter default |
| **Object present, all fields deferred** | `TextStyleUnspecified` | "Use theme for everything" | partial merge base |
| **Field absent** | `ColorUnspecified` (inside struct) | "Use ambient for this field" | field-level override |

**Example**:
```go
Text("hi", nil)                                    // nil → use theme 100 %
Text("hi", TextStyleUnspecified)                   // singleton → use theme 100 %
Text("hi", &TextStyle{fontSize: 20})               // partial → theme color + 20 sp
```

---

## 3. Public API for Complex Types (Nil-Safe, No Scattered Checks)

Expose **only** these four functions for complex types—**never** methods on `*T`:

```go
// Identity (2 ns)
func SameStyle(a, b *TextStyle) bool

// Semantic equality (field-by-field, 20 ns)
func EqualStyle(a, b *TextStyle) bool

// Merge (0 or 1 alloc)
func MergeTextStyle(a, b *TextStyle) *TextStyle

// Coalesce nil → singleton
func CoalesceTextStyle(ptr, def *TextStyle) *TextStyle
```

Business code **never writes `== nil`** outside these four helpers.

---

## 4. Sentinel Choice Reference

| Domain | Sentinel | Reason |
|--------|----------|--------|
| `Color` | `0` | transparent black is rare |
| `Dp` | `math.NaN()` | 0 dp is common |
| `TextUnit` | `math.NaN()` | 0 sp is valid |
| `IntOffset` | `(math.MinInt32, math.MinInt32)` | extremes unused |
| `FontWeight` | `0` (iota) | zero = unspecified |
| `Alignment` | `0` (iota) | zero = unspecified |

---

## 5. Sentinel Patterns for `int`, `string`, `bool`

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
func (p Profile) IsProfile() bool { return p != ProfileUnspecified }
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

## 6. Copy-Paste Templates

```go
// Primitive type
type MyUnit float32
const MyUnitUnspecified MyUnit = MyUnit(math.NaN())
func (u MyUnit) IsMyUnit() bool { return !math.IsNaN(float64(u)) }

// Complex type
type MyStyle struct { field1 Color; field2 Dp }
var MyStyleUnspecified = &MyStyle{field1: ColorUnspecified, field2: DpUnspecified}

func IsMyStyle(s *MyStyle) bool {
	return s != nil && s != MyStyleUnspecified
}
func TakeOrElseMyStyle(s, def *MyStyle) *MyStyle {
	if s == nil || s == MyStyleUnspecified { return def }
	return s
}
```

---

## 7. Performance Contract

- **Primitive sentinel**: 0 heap bytes, 2 ns, inlined by compiler
- **Complex nil**: 0 heap bytes, 5 ns
- **Complex custom**: 24 B struct alloc, 20 ns (paid only when user creates new style)

Verify with:
```bash
go test -bench=. -gcflags="-m" 2>&1 | grep -E "does not escape|inlined"
```

---

## 8. Anti-Patterns (reject on sight)

```go
❌ type Color interface { IsSpecified() bool }        // interface forces alloc
❌ var UnspecifiedColor *Color = nil                  // nil interface → panic
❌ func (ts *TextStyle) IsTextStyle() bool            // method on nil receiver
❌ func (ts *TextStyle) TakeOrElseTextStyle(block *TextStyle) *TextStyle // REJECT: method on nil receiver
❌ func TakeOrElseColor(block func() Color) Color // lambda escapes in Go
❌ type T struct { isSpecified bool }                 // flag field = double mem
```

---

## 9. Function Signatures (Principle of Least Parameters)

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

## 10. Package-Level Contract (Machine-Readable)

```go
// UI_PACKAGE_CONTRACT
// For every exported type T in package ui, the following symbols MUST exist:
//   const/var TUnspecified  // Sentinel value or singleton pointer
//   func IsT(v *T) bool      // Package-level predicate (never method on *T)
//   func TakeOrElseT(a, b *T) *T // Package-level fallback (never method on *T)
//   func MergeT(a, b *T) *T      // Package-level composition (never method on *T)
// Additional symbols (abbreviations, helpers) are allowed but must be documented.
// All sentinel values must be compile-time constants or package-level variables.
// No function may accept func() parameters in hot paths (forces heap escape).
// No method may have a pointer receiver that checks for nil (anti-pattern).
// END_CONTRACT
```

---

## 11. Verification Commands

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