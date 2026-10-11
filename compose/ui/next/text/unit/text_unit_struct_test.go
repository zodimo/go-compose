package unit

import (
	"math"
	"testing"
)

// TextUnitStruct is a struct-based alternative to TextUnit (which is int64)
// This test file proves that both approaches have identical allocation behavior.
type TextUnitStruct struct {
	packed int64
}

var TextUnitStructUnspecified = TextUnitStruct{packed: int64(packStruct(int64(TextUnitTypeUnspecified), float32(math.NaN())))}

func packStruct(unitType int64, v float32) int64 {
	valBits := int64(math.Float32bits(v)) & 0xFFFFFFFF
	return unitType | valBits
}

func SpStruct(value float32) TextUnitStruct {
	return TextUnitStruct{packed: packStruct(int64(TextUnitTypeSp), value)}
}

func (tu TextUnitStruct) rawType() int64 {
	return tu.packed & unitMask
}

func (tu TextUnitStruct) IsUnspecified() bool {
	return tu.rawType() == int64(TextUnitTypeUnspecified)
}

func (tu TextUnitStruct) IsSpecified() bool {
	return !tu.IsUnspecified()
}

func (tu TextUnitStruct) Value() float32 {
	return math.Float32frombits(uint32(tu.packed & 0xFFFFFFFF))
}

func (tu TextUnitStruct) TakeOrElse(block TextUnitStruct) TextUnitStruct {
	if tu.IsSpecified() {
		return tu
	}
	return block
}

// === Benchmarks to prove zero allocation ===

// BenchmarkTextUnitInt64_Create tests the current int64-based TextUnit
func BenchmarkTextUnitInt64_Create(b *testing.B) {
	var result TextUnit
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		result = Sp(float32(i))
	}
	_ = result
}

// BenchmarkTextUnitStruct_Create tests the struct-based alternative
func BenchmarkTextUnitStruct_Create(b *testing.B) {
	var result TextUnitStruct
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		result = SpStruct(float32(i))
	}
	_ = result
}

// BenchmarkTextUnitInt64_TakeOrElse tests TakeOrElse with int64-based TextUnit
func BenchmarkTextUnitInt64_TakeOrElse(b *testing.B) {
	specified := Sp(24)
	unspecified := TextUnitUnspecified
	defaultVal := Sp(14)
	var result TextUnit
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			result = specified.TakeOrElse(defaultVal)
		} else {
			result = unspecified.TakeOrElse(defaultVal)
		}
	}
	_ = result
}

// BenchmarkTextUnitStruct_TakeOrElse tests TakeOrElse with struct-based TextUnit
func BenchmarkTextUnitStruct_TakeOrElse(b *testing.B) {
	specified := SpStruct(24)
	unspecified := TextUnitStructUnspecified
	defaultVal := SpStruct(14)
	var result TextUnitStruct
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			result = specified.TakeOrElse(defaultVal)
		} else {
			result = unspecified.TakeOrElse(defaultVal)
		}
	}
	_ = result
}

// BenchmarkTextUnitInt64_Value tests Value() with int64-based TextUnit
func BenchmarkTextUnitInt64_Value(b *testing.B) {
	tu := Sp(24)
	var result float32
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		result = tu.Value()
	}
	_ = result
}

// BenchmarkTextUnitStruct_Value tests Value() with struct-based TextUnit
func BenchmarkTextUnitStruct_Value(b *testing.B) {
	tu := SpStruct(24)
	var result float32
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		result = tu.Value()
	}
	_ = result
}

// TestTextUnitStruct_Equivalence verifies struct behaves identically to int64
func TestTextUnitStruct_Equivalence(t *testing.T) {
	// Test Sp creation
	int64Unit := Sp(24)
	structUnit := SpStruct(24)

	if int64Unit.Value() != structUnit.Value() {
		t.Errorf("Value mismatch: int64=%v, struct=%v", int64Unit.Value(), structUnit.Value())
	}

	if int64Unit.IsSpecified() != structUnit.IsSpecified() {
		t.Errorf("IsSpecified mismatch: int64=%v, struct=%v", int64Unit.IsSpecified(), structUnit.IsSpecified())
	}

	// Test unspecified
	if TextUnitUnspecified.IsUnspecified() != TextUnitStructUnspecified.IsUnspecified() {
		t.Errorf("IsUnspecified mismatch")
	}

	// Test TakeOrElse
	defaultInt64 := Sp(14)
	defaultStruct := SpStruct(14)

	resultInt64 := TextUnitUnspecified.TakeOrElse(defaultInt64)
	resultStruct := TextUnitStructUnspecified.TakeOrElse(defaultStruct)

	if resultInt64.Value() != resultStruct.Value() {
		t.Errorf("TakeOrElse mismatch: int64=%v, struct=%v", resultInt64.Value(), resultStruct.Value())
	}
}
