package unit

import (
	"fmt"
	"math"

	"github.com/zodimo/go-compose/pkg/floatutils"
)

// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-unit/src/commonMain/kotlin/androidx/compose/ui/unit/TextUnit.kt;drc=4970f6e96cdb06089723da0ab8ec93ae3f067c7a

type TextUnitType int

const (
	TextUnitTypeUnspecified TextUnitType = iota
	TextUnitTypeSp
	TextUnitTypeEm
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

type TextUnit struct {
	value    float32
	unitType TextUnitType
}

var TextUnitUnspecified = TextUnit{value: float32(math.NaN()), unitType: TextUnitTypeUnspecified}

func NewTextUnit(value float32, unitType TextUnitType) TextUnit {
	return TextUnit{
		value:    value,
		unitType: unitType,
	}
}

func Sp(value float32) TextUnit {
	return NewTextUnit(value, TextUnitTypeSp)
}

func Em(value float32) TextUnit {
	return NewTextUnit(value, TextUnitTypeEm)
}

func (tu TextUnit) Type() TextUnitType {
	return tu.unitType
}

func (tu TextUnit) Value() float32 {
	return tu.value
}

func (tu TextUnit) IsSp() bool {
	return tu.unitType == TextUnitTypeSp
}

func (tu TextUnit) IsEm() bool {
	return tu.unitType == TextUnitTypeEm
}

func (tu TextUnit) IsUnspecified() bool {
	return tu.unitType == TextUnitTypeUnspecified
}

func (tu TextUnit) IsSpecified() bool {
	return !tu.IsUnspecified()
}

func (tu TextUnit) TakeOrElse(block func() TextUnit) TextUnit {
	if tu.IsSpecified() {
		return tu
	}
	return block()
}

func (tu TextUnit) UnaryMinus() TextUnit {
	if tu.IsUnspecified() {
		panic("Cannot perform operation for Unspecified type.")
	}
	return TextUnit{
		value:    -tu.value,
		unitType: tu.unitType,
	}
}

func (tu TextUnit) Div(other float32) TextUnit {
	if tu.IsUnspecified() {
		panic("Cannot perform operation for Unspecified type.")
	}
	return TextUnit{
		value:    tu.value / other,
		unitType: tu.unitType,
	}
}

func (tu TextUnit) Times(other float32) TextUnit {
	if tu.IsUnspecified() {
		panic("Cannot perform operation for Unspecified type.")
	}
	return TextUnit{
		value:    tu.value * other,
		unitType: tu.unitType,
	}
}

func (tu TextUnit) Compare(other TextUnit) int {
	if tu.IsUnspecified() || other.IsUnspecified() {
		panic("Cannot perform operation for Unspecified type.")
	}
	if tu.unitType != other.unitType {
		panic(fmt.Sprintf("Cannot perform operation for %s and %s", tu.unitType, other.unitType))
	}
	if floatutils.Float32Equals(tu.value, other.value, floatutils.Float32EqualityThreshold) {
		return 0
	}
	if tu.value < other.value {
		return -1
	} else {
		return 1
	}
}

func (tu TextUnit) String() string {
	switch tu.unitType {
	case TextUnitTypeUnspecified:
		return "Unspecified"
	case TextUnitTypeSp:
		return fmt.Sprintf("%v.sp", tu.value)
	case TextUnitTypeEm:
		return fmt.Sprintf("%v.em", tu.value)
	default:
		return "Invalid"
	}
}

func (tu TextUnit) Equals(other TextUnit) bool {
	return float32Equals(tu.value, other.value, float32EqualityThreshold) && tu.unitType == other.unitType
}

// LerpTextUnit linearly interpolates between two TextUnits.
func LerpTextUnit(start, stop TextUnit, fraction float32) TextUnit {
	if start.IsUnspecified() || stop.IsUnspecified() {
		panic("Cannot perform operation for Unspecified type.")
	}
	if start.unitType != stop.unitType {
		panic(fmt.Sprintf("Cannot perform operation for %s and %s", start.unitType, stop.unitType))
	}

	val := lerpBetween(start.value, stop.value, float64(fraction))
	return TextUnit{
		value:    val,
		unitType: start.unitType,
	}
}
