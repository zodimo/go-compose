package unit

import (
	"fmt"
	"math"

	gioUnit "gioui.org/unit"
	"github.com/zodimo/go-compose/pkg/floatutils"
	"github.com/zodimo/go-compose/pkg/floatutils/lerp"
)

// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-unit/src/commonMain/kotlin/androidx/compose/ui/unit/TextUnit.kt;drc=4970f6e96cdb06089723da0ab8ec93ae3f067c7a

type TextUnitType int64

const (
	TextUnitTypeUnspecified TextUnitType = 0x00 << 32
	TextUnitTypeSp          TextUnitType = 0x01 << 32
	TextUnitTypeEm          TextUnitType = 0x02 << 32
)

// Internal constants
const (
	unitMask = 0xFF << 32
)

func (t TextUnitType) String() string {
	switch t {
	case TextUnitTypeUnspecified:
		return "Unspecified"
	case TextUnitTypeSp:
		return "Sp"
	case TextUnitTypeEm:
		return "Em"
	default:
		return "Invalid"
	}
}

// TextUnit packs unit type (Sp/Em) and float value into 64 bits.
// Struct wrapper prevents `TextUnit(24)` - must use `Sp(24)` or `Em(1.5)`.
type TextUnit struct {
	packed int64
}

// 1. `TUnspecified` – sentinel (value, not pointer)
// TextUnitUnspecified is the sentinel value for an unspecified TextUnit.
// Kotlin uses a packed value with NaN for Unspecified.
// val Unspecified = pack(UNIT_TYPE_UNSPECIFIED, Float.NaN)
var TextUnitUnspecified = TextUnit{packed: packRaw(int64(TextUnitTypeUnspecified), float32(math.NaN()))}

// Internal packing function
func packRaw(unitType int64, v float32) int64 {
	valBits := int64(math.Float32bits(v)) & 0xFFFFFFFF
	return unitType | valBits
}

// NewTextUnit creates a new TextUnit.
// Note: In Kotlin this is constructor `TextUnit(value: Float, type: TextUnitType)`.
func NewTextUnit(value float32, unitType TextUnitType) TextUnit {
	return TextUnit{packed: packRaw(int64(unitType), value)}
}

// Sp creates a SP unit TextUnit.
func Sp(value float32) TextUnit {
	return TextUnit{packed: packRaw(int64(TextUnitTypeSp), value)}
}

// Em creates an EM unit TextUnit.
func Em(value float32) TextUnit {
	return TextUnit{packed: packRaw(int64(TextUnitTypeEm), value)}
}

// rawType returns the raw type bits.
func (tu TextUnit) rawType() int64 {
	return tu.packed & unitMask
}

// Type returns the TextUnitType of this TextUnit.
func (tu TextUnit) Type() TextUnitType {
	return TextUnitType(tu.rawType())
}

// Value returns the float value of this TextUnit.
func (tu TextUnit) Value() float32 {
	return math.Float32frombits(uint32(tu.packed & 0xFFFFFFFF))
}

// IsSp returns true if this is a SP unit type.
func (tu TextUnit) IsSp() bool {
	return tu.rawType() == int64(TextUnitTypeSp)
}

// IsEm returns true if this is an EM unit type.
func (tu TextUnit) IsEm() bool {
	return tu.rawType() == int64(TextUnitTypeEm)
}

// IsUnspecified returns true if this is an unspecified unit type.
func (tu TextUnit) IsUnspecified() bool {
	return tu.rawType() == int64(TextUnitTypeUnspecified)
}

// 2. `IsSpecified` – predicate (method on value receiver)
// IsSpecified returns true if this is a specified unit type.
func (tu TextUnit) IsSpecified() bool {
	return tu.rawType() != int64(TextUnitTypeUnspecified)
}

// 3. `TakeOrElse` – 2-param fallback (method on value receiver)
// TakeOrElse returns this TextUnit if specified, otherwise returns the default.
func (tu TextUnit) TakeOrElse(def TextUnit) TextUnit {
	if tu.IsSpecified() {
		return tu
	}
	return def
}

// 4. `Merge` – whole-value replacement for atomic packed types (method on value receiver)
// Merge returns other if specified, otherwise returns tu.
func (tu TextUnit) Merge(other TextUnit) TextUnit {
	if other.IsSpecified() {
		return other
	}
	return tu
}

// 5. `String` – stringification (method on value receiver)
// String returns the string representation of the TextUnit.
func (tu TextUnit) String() string {
	if !tu.IsSpecified() {
		return "TextUnit{Unspecified}"
	}
	return fmt.Sprintf("TextUnit{%v %s}", tu.Value(), tu.Type())
}

// 6. `Coalesce` – N/A for value types (no nil possible)

// 7-9. Equality – Equals method
// Equals checks if two TextUnits are equal.
func (tu TextUnit) Equals(other TextUnit) bool {
	if tu.Type() != other.Type() {
		return false
	}
	// If both are Unspecified, they are equal regardless of the float payload (which is NaN)
	if tu.IsUnspecified() {
		return other.IsUnspecified()
	}
	return floatutils.Float32Equals(tu.Value(), other.Value(), floatutils.Float32EqualityThreshold)
}

// 10. `Copy` – identity for immutable value types (just use assignment)
// For TextUnit, copy is a no-op as it's a value type

// UnaryMinus returns the negation of this TextUnit.
func (tu TextUnit) UnaryMinus() TextUnit {
	checkArithmetic(tu)
	return TextUnit{packed: packRaw(tu.rawType(), -tu.Value())}
}

// Div divides a TextUnit by a scalar.
func (tu TextUnit) Div(other float32) TextUnit {
	checkArithmetic(tu)
	return TextUnit{packed: packRaw(tu.rawType(), tu.Value()/other)}
}

// Times multiplies a TextUnit by a scalar.
func (tu TextUnit) Times(other float32) TextUnit {
	checkArithmetic(tu)
	return TextUnit{packed: packRaw(tu.rawType(), tu.Value()*other)}
}

// Compare compares this TextUnit with another.
func (tu TextUnit) Compare(other TextUnit) int {
	checkArithmetic2(tu, other)
	diff := tu.Value() - other.Value()
	if floatutils.Float32Equals(tu.Value(), other.Value(), floatutils.Float32EqualityThreshold) {
		return 0
	}
	if diff < 0 {
		return -1
	}
	return 1
}

func (tu TextUnit) AsGioSp() gioUnit.Sp {
	if !tu.IsSpecified() {
		return gioUnit.Sp(0)
	}
	if tu.IsEm() {
		panic("TextUnit is an EM unit, cannot convert to Sp")
	}
	return gioUnit.Sp(tu.Value())
}

// LerpTextUnit linearly interpolates between two TextUnits.
func LerpTextUnit(start, stop TextUnit, fraction float32) TextUnit {
	checkArithmetic2(start, stop)
	val := lerp.Float32(start.Value(), stop.Value(), fraction)
	return TextUnit{packed: packRaw(start.rawType(), val)}
}

func LerpTextUnitInheritable(a, b TextUnit, t float32) TextUnit {
	if a.IsUnspecified() || b.IsUnspecified() {
		return lerp.LerpDiscrete(a, b, t)
	}
	return LerpTextUnit(a, b, t)
}

func checkArithmetic(a TextUnit) {
	requirePrecondition(!a.IsUnspecified(), "Cannot perform operation for Unspecified type.")
}

func checkArithmetic2(a, b TextUnit) {

	requirePrecondition(!a.IsUnspecified() && !b.IsUnspecified(), "Cannot perform operation for Unspecified type.")
	requirePrecondition(a.Type() == b.Type(), fmt.Sprintf("Cannot perform operation for %s and %s", a.Type(), b.Type()))
}

func requirePrecondition(cond bool, message string) {
	if !cond {
		panic(message)
	}
}
