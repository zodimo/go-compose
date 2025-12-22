package unit

import (
	"math"
	"testing"
)

func TestTextUnit_Creation(t *testing.T) {
	sp := Sp(10)
	if !sp.IsSp() {
		t.Errorf("Expected Sp unit, got %v", sp.Type())
	}
	if sp.Value() != 10 {
		t.Errorf("Expected 10, got %v", sp.Value())
	}

	em := Em(2.5)
	if !em.IsEm() {
		t.Errorf("Expected Em unit, got %v", em.Type())
	}
	if em.Value() != 2.5 {
		t.Errorf("Expected 2.5, got %v", em.Value())
	}

	unspecified := TextUnitUnspecified
	if !unspecified.IsUnspecified() {
		t.Errorf("Expected Unspecified unit, got %v", unspecified.Type())
	}
}

func TestTextUnit_String(t *testing.T) {
	tests := []struct {
		unit     TextUnit
		expected string
	}{
		{Sp(12), "12.sp"},
		{Em(1.5), "1.5.em"},
		{TextUnitUnspecified, "Unspecified"},
	}

	for _, tt := range tests {
		if got := tt.unit.String(); got != tt.expected {
			t.Errorf("TextUnit.String() = %v, want %v", got, tt.expected)
		}
	}
}

func TestTextUnit_Helpers(t *testing.T) {
	sp := Sp(1)
	if sp.IsUnspecified() {
		t.Error("Sp should not be Unspecified")
	}
	if !sp.IsSpecified() {
		t.Error("Sp should be Specified")
	}

	unspec := TextUnitUnspecified
	if !unspec.IsUnspecified() {
		t.Error("Unspecified should be Unspecified")
	}
	if unspec.IsSpecified() {
		t.Error("Unspecified should not be Specified")
	}
}

func TestTextUnit_TakeOrElse(t *testing.T) {
	sp := Sp(10)
	fallback := Em(20)

	result := sp.TakeOrElse(func() TextUnit { return fallback })
	if result != sp {
		t.Errorf("Expected original %v, got %v", sp, result)
	}

	result = TextUnitUnspecified.TakeOrElse(func() TextUnit { return fallback })
	if result != fallback {
		t.Errorf("Expected fallback %v, got %v", fallback, result)
	}
}

func TestTextUnit_Arithmetic_Valid(t *testing.T) {
	val := Sp(10)

	// UnaryMinus
	minus := val.UnaryMinus()
	if minus.Value() != -10 || !minus.IsSp() {
		t.Errorf("UnaryMinus failed: %v", minus)
	}

	// Times
	times := val.Times(2)
	if times.Value() != 20 || !times.IsSp() {
		t.Errorf("Times failed: %v", times)
	}

	// Div
	div := val.Div(2)
	if div.Value() != 5 || !div.IsSp() {
		t.Errorf("Div failed: %v", div)
	}
}

func TestTextUnit_Compare(t *testing.T) {
	v1 := Sp(10)
	v2 := Sp(20)
	v3 := Sp(10)

	if v1.Compare(v2) >= 0 {
		t.Error("Compare 10.sp vs 20.sp should be negative")
	}
	if v2.Compare(v1) <= 0 {
		t.Error("Compare 20.sp vs 10.sp should be positive")
	}
	if v1.Compare(v3) != 0 {
		t.Error("Compare 10.sp vs 10.sp should be 0")
	}
}

func TestTextUnit_Lerp(t *testing.T) {
	v1 := Sp(10)
	v2 := Sp(20)

	mid := LerpTextUnit(v1, v2, 0.5)
	if mid.Value() != 15 || !mid.IsSp() {
		t.Errorf("Lerp 0.5 failed: %v", mid)
	}

	start := LerpTextUnit(v1, v2, 0)
	if start.Value() != 10 {
		t.Errorf("Lerp 0 failed: %v", start)
	}

	end := LerpTextUnit(v1, v2, 1)
	if end.Value() != 20 {
		t.Errorf("Lerp 1 failed: %v", end)
	}
}

func TestTextUnit_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("The code did not panic")
		}
	}()
	TextUnitUnspecified.Times(2)
}

func TestTextUnit_Compare_Panic_Unspecified(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("The code did not panic")
		}
	}()
	Sp(10).Compare(TextUnitUnspecified)
}

func TestTextUnit_Compare_Panic_TypeMismatch(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("The code did not panic")
		}
	}()
	Sp(10).Compare(Em(10))
}

func TestTextUnit_Lerp_Panic_MixedTypes(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("The code did not panic")
		}
	}()
	LerpTextUnit(Sp(10), Em(20), 0.5)
}

// Ensure NaN handling in Unspecified
func TestTextUnit_Unspecified_IsNaN(t *testing.T) {
	if !math.IsNaN(float64(TextUnitUnspecified.Value())) {
		t.Error("TextUnitUnspecified value should be NaN")
	}
}
