package unit

import (
	"math"
	"testing"
)

func TestTextUnitPacking(t *testing.T) {
	t.Run("Sp packing", func(t *testing.T) {
		sp := Sp(12.0)
		if !sp.IsSp() {
			t.Errorf("Expected IsSp() to be true")
		}
		if sp.IsEm() {
			t.Errorf("Expected IsEm() to be false")
		}
		if sp.IsUnspecified() {
			t.Errorf("Expected IsUnspecified() to be false")
		}
		if sp.Type() != TextUnitTypeSp {
			t.Errorf("Expected Type() to be TextUnitTypeSp, got %v", sp.Type())
		}
		if sp.Value() != 12.0 {
			t.Errorf("Expected Value() to be 12.0, got %v", sp.Value())
		}
	})

	t.Run("Em packing", func(t *testing.T) {
		em := Em(1.5)
		if !em.IsEm() {
			t.Errorf("Expected IsEm() to be true")
		}
		if em.IsSp() {
			t.Errorf("Expected IsSp() to be false")
		}
		if em.IsUnspecified() {
			t.Errorf("Expected IsUnspecified() to be false")
		}
		if em.Type() != TextUnitTypeEm {
			t.Errorf("Expected Type() to be TextUnitTypeEm, got %v", em.Type())
		}
		if em.Value() != 1.5 {
			t.Errorf("Expected Value() to be 1.5, got %v", em.Value())
		}
	})

	t.Run("Unspecified packing", func(t *testing.T) {
		u := TextUnitUnspecified
		if !u.IsUnspecified() {
			t.Errorf("Expected IsUnspecified() to be true")
		}
		if u.IsSpecified() {
			t.Errorf("Expected IsSpecified() to be false")
		}
		if u.Type() != TextUnitTypeUnspecified {
			t.Errorf("Expected Type() to be TextUnitTypeUnspecified, got %v", u.Type())
		}
		// Value of Unspecified is NaN, so direct comparison fails.
		if !math.IsNaN(float64(u.Value())) {
			t.Errorf("Expected Value() to be NaN, got %v", u.Value())
		}
	})
}

func TestTextUnitArithmetic(t *testing.T) {
	sp10 := Sp(10)
	sp20 := Sp(20)

	t.Run("UnaryMinus", func(t *testing.T) {
		neg := sp10.UnaryMinus()
		if neg.Value() != -10 {
			t.Errorf("Expected -10, got %v", neg.Value())
		}
		if neg.Type() != TextUnitTypeSp {
			t.Errorf("Expected Type Sp, got %v", neg.Type())
		}
	})

	t.Run("Times", func(t *testing.T) {
		doubled := sp10.Times(2)
		if doubled.Value() != 20 {
			t.Errorf("Expected 20, got %v", doubled.Value())
		}
	})

	t.Run("Div", func(t *testing.T) {
		halved := sp20.Div(2)
		if halved.Value() != 10 {
			t.Errorf("Expected 10, got %v", halved.Value())
		}
	})

	t.Run("Panic on Unspecified", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		TextUnitUnspecified.Times(2)
	})
}

func TestTextUnitComparison(t *testing.T) {
	sp10 := Sp(10)
	sp20 := Sp(20)
	sp10Clone := Sp(10)
	em10 := Em(10)

	if sp10.Compare(sp20) >= 0 {
		t.Errorf("Expected 10.sp < 20.sp")
	}
	if sp20.Compare(sp10) <= 0 {
		t.Errorf("Expected 20.sp > 10.sp")
	}
	if sp10.Compare(sp10Clone) != 0 {
		t.Errorf("Expected 10.sp == 10.sp")
	}

	// Test Equals
	if !sp10.Equals(sp10Clone) {
		t.Errorf("Expected Equals to be true")
	}
	if sp10.Equals(em10) {
		t.Errorf("Expected different types to be unequal")
	}
	if sp10.Equals(sp20) {
		t.Errorf("Expected different values to be unequal")
	}
}

func TestLerpTextUnit(t *testing.T) {
	start := Sp(10)
	stop := Sp(20)

	mid := LerpTextUnit(start, stop, 0.5)
	if mid.Value() != 15 {
		t.Errorf("Expected 15, got %v", mid.Value())
	}
	if mid.Type() != TextUnitTypeSp {
		t.Errorf("Expected type Sp, got %v", mid.Type())
	}

	t.Run("Lerp panic on mismatch", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		LerpTextUnit(Sp(10), Em(20), 0.5)
	})
}

func TestTextUnitString(t *testing.T) {
	if Sp(12.5).String() != "12.5.sp" {
		t.Errorf("Expected 12.5.sp, got %v", Sp(12.5).String())
	}
	// Note: Float string formatting might vary slightly, but 12.5 is exact.
	if Em(2).String() != "2.em" {
		t.Errorf("Expected 2.em, got %v", Em(2).String())
	}
	if TextUnitUnspecified.String() != "Unspecified" {
		t.Errorf("Expected Unspecified, got %v", TextUnitUnspecified.String())
	}
}
