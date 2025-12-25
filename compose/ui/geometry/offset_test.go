package geometry

import (
	"math"
	"testing"
)

func TestOffsetZero(t *testing.T) {
	if OffsetZero.X() != 0 || OffsetZero.Y() != 0 {
		t.Errorf("OffsetZero = %v, want (0, 0)", OffsetZero)
	}
}

func TestOffsetInfinite(t *testing.T) {
	if !math.IsInf(float64(OffsetInfinite.X()), 1) || !math.IsInf(float64(OffsetInfinite.Y()), 1) {
		t.Errorf("OffsetInfinite = %v, want (+Inf, +Inf)", OffsetInfinite)
	}
}

func TestOffsetEqualInfinite(t *testing.T) {
	inf1 := OffsetInfinite
	inf2 := OffsetInfinite

	if !inf1.Equal(inf2) {
		t.Errorf("OffsetInfinite should equal OffsetInfinite")
	}
}

func TestOffsetUnspecified(t *testing.T) {
	if !math.IsNaN(float64(OffsetUnspecified.X())) || !math.IsNaN(float64(OffsetUnspecified.Y())) {
		t.Errorf("OffsetUnspecified = %v, want (NaN, NaN)", OffsetUnspecified)
	}
}

func TestOffsetIsValid(t *testing.T) {
	tests := []struct {
		offset Offset
		valid  bool
	}{
		{NewOffset(0, 0), true},
		{NewOffset(10, 20), true},
		{NewOffset(float32(math.Inf(1)), 0), true}, // Infinite is valid in Kotlin
		{NewOffset(0, float32(math.Inf(1))), true},
		{NewOffset(float32(math.NaN()), 0), false},
		{NewOffset(0, float32(math.NaN())), false},
		{OffsetUnspecified, false},
		{OffsetInfinite, true},
	}

	for _, tt := range tests {
		if got := tt.offset.IsValid(); got != tt.valid {
			t.Errorf("Offset{%v}.IsValid() = %v, want %v", tt.offset, got, tt.valid)
		}
	}
}

func TestOffsetIsSpecified(t *testing.T) {
	if OffsetUnspecified.IsOffset() {
		t.Error("OffsetUnspecified should not be specified")
	}
	if !OffsetZero.IsOffset() {
		t.Error("OffsetZero should be specified")
	}
	// Test partial NaN
	partial := NewOffset(float32(math.NaN()), 1.0)
	if !partial.IsOffset() {
		t.Error("Offset{NaN, 1.0} should be specified")
	}
}

func TestOffsetIsUnspecified(t *testing.T) {
	if !OffsetUnspecified.IsUnspecified() {
		t.Error("OffsetUnspecified should be unspecified")
	}
	if OffsetZero.IsUnspecified() {
		t.Error("OffsetZero should not be unspecified")
	}
	// Test partial NaN
	partial := NewOffset(float32(math.NaN()), 1.0)
	if partial.IsUnspecified() {
		t.Error("Offset{NaN, 1.0} should NOT be unspecified")
	}
}

func TestOffsetIsFinite(t *testing.T) {
	tests := []struct {
		offset Offset
		finite bool
	}{
		{NewOffset(0, 0), true},
		{NewOffset(10, 20), true},
		{NewOffset(float32(math.Inf(1)), 0), false},
		{NewOffset(0, float32(math.Inf(1))), false},
		{NewOffset(float32(math.NaN()), 0), false},
		{OffsetInfinite, false},
		{OffsetUnspecified, false},
	}

	for _, tt := range tests {
		if got := tt.offset.IsFinite(); got != tt.finite {
			t.Errorf("Offset{%v}.IsFinite() = %v, want %v", tt.offset, got, tt.finite)
		}
	}
}

func TestOffsetOps(t *testing.T) {
	a := NewOffset(10, 20)
	b := NewOffset(30, 40)

	// Plus
	sum := a.Plus(b)
	if sum.X() != 40 || sum.Y() != 60 {
		t.Errorf("Plus mismatched: got %v, want (40, 60)", sum)
	}

	// Minus
	diff := a.Minus(b)
	if diff.X() != -20 || diff.Y() != -20 {
		t.Errorf("Minus mismatched: got %v, want (-20, -20)", diff)
	}

	// Times
	scaled := a.Times(2)
	if scaled.X() != 20 || scaled.Y() != 40 {
		t.Errorf("Times mismatched: got %v, want (20, 40)", scaled)
	}

	// Div
	div := a.Div(2)
	if div.X() != 5 || div.Y() != 10 {
		t.Errorf("Div mismatched: got %v, want (5, 10)", div)
	}

	// Rem
	rem := NewOffset(11, 21).Rem(10)
	if rem.X() != 1 || rem.Y() != 1 {
		t.Errorf("Rem mismatched: got %v, want (1, 1)", rem)
	}

	// UnaryMinus
	neg := a.UnaryMinus()
	if neg.X() != -10 || neg.Y() != -20 {
		t.Errorf("UnaryMinus mismatched: got %v, want (-10, -20)", neg)
	}
}

func TestOffsetDistance(t *testing.T) {
	o := NewOffset(3, 4)
	if d := o.GetDistance(); d != 5 {
		t.Errorf("GetDistance() = %v, want 5", d)
	}
	if d2 := o.GetDistanceSquared(); d2 != 25 {
		t.Errorf("GetDistanceSquared() = %v, want 25", d2)
	}
}

func TestLerp(t *testing.T) {
	start := NewOffset(0, 0)
	end := NewOffset(10, 20)

	mid := LerpOffset(start, end, 0.5)
	if mid.X() != 5 || mid.Y() != 10 {
		t.Errorf("Lerp(0.5) = %v, want (5, 10)", mid)
	}

	// Extrapolation
	extra := LerpOffset(start, end, 2.0)
	if extra.X() != 20 || extra.Y() != 40 {
		t.Errorf("Lerp(2.0) = %v, want (20, 40)", extra)
	}
}

func TestOffsetEqual(t *testing.T) {
	o1 := NewOffset(1, 2)
	o2 := NewOffset(1, 2)
	o3 := NewOffset(1, 3)

	if !o1.Equal(o2) {
		t.Error("Equal() should be true for identical offsets")
	}
	if o1.Equal(o3) {
		t.Error("Equal() should be false for different offsets")
	}

}
