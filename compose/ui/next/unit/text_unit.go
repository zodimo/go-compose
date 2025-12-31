package unit

import (
	"fmt"
	"math"

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

// TextUnit encodes the unit information and float value into a single 64-bit integer.
// The higher 32 bits represent the metadata (unit type), and the lower 32 bits represent
// the bit representation of the float value.
type TextUnit int64

// TextUnitUnspecified is the sentinel value for an unspecified TextUnit.
// Kotlin uses a packed value with NaN for Unspecified.
// val Unspecified = pack(UNIT_TYPE_UNSPECIFIED, Float.NaN)
var TextUnitUnspecified = pack(int64(TextUnitTypeUnspecified), float32(math.NaN()))

// Pack packs a unit type and a float value into a TextUnit.
func pack(unitType int64, v float32) TextUnit {
	// unitType is expected to be already shifted (e.g., TextUnitTypeSp)
	// v.toRawBits().toLong() and 0xFFFF_FFFFL
	valBits := int64(math.Float32bits(v)) & 0xFFFFFFFF
	return TextUnit(unitType | valBits)
}

// NewTextUnit creates a new TextUnit.
// Note: In Kotlin this is constructor `TextUnit(value: Float, type: TextUnitType)`.
func NewTextUnit(value float32, unitType TextUnitType) TextUnit {
	return pack(int64(unitType), value)
}

// Sp creates a SP unit TextUnit.
func Sp(value float32) TextUnit {
	return pack(int64(TextUnitTypeSp), value)
}

// Em creates an EM unit TextUnit.
func Em(value float32) TextUnit {
	return pack(int64(TextUnitTypeEm), value)
}

// RawType returns the raw type bits.
func (tu TextUnit) rawType() int64 {
	return int64(tu) & unitMask
}

// Type returns the TextUnitType of this TextUnit.
func (tu TextUnit) Type() TextUnitType {
	// Kotlin: TextUnitTypes[(rawType ushr 32).toInt()]
	// But we can just mask it since our TextUnitType constants are already shifted.
	// However, TextUnitType is defined as the shifted value in Kotlin?
	// Kotlin: value class TextUnitType(internal val type: Long)
	// Companion object vals are shifted: val Sp = TextUnitType(UNIT_TYPE_SP)
	// So returning the masked value directly as TextUnitType is correct if TextUnitType is defined as the shifted value.
	return TextUnitType(tu.rawType())
}

// Value returns the float value of this TextUnit.
func (tu TextUnit) Value() float32 {
	return math.Float32frombits(uint32(int64(tu) & 0xFFFFFFFF))
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

// IsSpecified returns true if this is a specified unit type.
func (tu TextUnit) IsSpecified() bool {
	return !tu.IsUnspecified()
}

// TakeOrElse returns this TextUnit if specified, otherwise executes the block.
func (tu TextUnit) TakeOrElse(block TextUnit) TextUnit {
	if tu.IsSpecified() {
		return tu
	}
	return block
}

// UnaryMinus returns the negation of this TextUnit.
func (tu TextUnit) UnaryMinus() TextUnit {
	checkArithmetic(tu)
	return pack(tu.rawType(), -tu.Value())
}

// Div divides a TextUnit by a scalar.
func (tu TextUnit) Div(other float32) TextUnit {
	checkArithmetic(tu)
	return pack(tu.rawType(), tu.Value()/other)
}

// Times multiplies a TextUnit by a scalar.
func (tu TextUnit) Times(other float32) TextUnit {
	checkArithmetic(tu)
	return pack(tu.rawType(), tu.Value()*other)
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

// String returns the string representation of the TextUnit.
func (tu TextUnit) String() string {
	switch tu.Type() {
	case TextUnitTypeUnspecified:
		return "Unspecified"
	case TextUnitTypeSp:
		return fmt.Sprintf("%v.sp", tu.Value())
	case TextUnitTypeEm:
		return fmt.Sprintf("%v.em", tu.Value())
	default:
		return "Invalid"
	}
}

// Equals checks if two TextUnits are equal.
func (tu TextUnit) Equals(other TextUnit) bool {
	// Kotlin relies on value class equality which checks the underlying Long.
	// However, floating point equality might be tricky if we just compare the bits directly.
	// But since we pack it deterministically, direct comparison might work for same values.
	// Kotlin's implementation doesn't override equals for TextUnit, so it uses default value class equality (comparing the underlying long).
	// But wait, floating point NaN != NaN.
	// If we want exact Kotlin behavior, we should look at how value classes compare.
	// For now, let's use our explicit comparison logic for safety.

	// A nicer implementation might be to just compare the raw bits if we trust the packing.
	// But float arithmetic results can vary slightly.
	// Let's stick to checking type and fuzzy float equality.
	if tu.Type() != other.Type() {
		return false
	}
	// If both are Unspecified, they are equal regardless of the float payload (which is NaN)
	if tu.IsUnspecified() {
		return other.IsUnspecified()
	}
	return floatutils.Float32Equals(tu.Value(), other.Value(), floatutils.Float32EqualityThreshold)
}

// LerpTextUnit linearly interpolates between two TextUnits.
func LerpTextUnit(start, stop TextUnit, fraction float32) TextUnit {
	checkArithmetic2(start, stop)
	val := lerp.Float32(start.Value(), stop.Value(), fraction)
	return pack(start.rawType(), val)
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
