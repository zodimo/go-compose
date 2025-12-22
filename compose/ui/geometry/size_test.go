package geometry

import (
	"testing"
)

func TestSizeConstruction(t *testing.T) {
	s := NewSize(10, 20)
	if s.Width != 10 || s.Height != 20 {
		t.Errorf("NewSize(10, 20) = %v, want Size(10, 20)", s)
	}

	if SizeZero.Width != 0 || SizeZero.Height != 0 {
		t.Errorf("SizeZero = %v, want Size(0, 0)", SizeZero)
	}
}

func TestSizeUnspecified(t *testing.T) {
	s := SizeUnspecified
	if !s.IsUnspecified() {
		t.Errorf("SizeUnspecified.IsUnspecified() = false, want true")
	}
	if s.IsSpecified() {
		t.Errorf("SizeUnspecified.IsSpecified() = true, want false")
	}

	s2 := NewSize(10, 20)
	if s2.IsUnspecified() {
		t.Errorf("Size(10, 20).IsUnspecified() = true, want false")
	}
	if !s2.IsSpecified() {
		t.Errorf("Size(10, 20).IsSpecified() = false, want true")
	}
}

func TestSizeIsEmpty(t *testing.T) {
	tests := []struct {
		s    Size
		want bool
	}{
		{NewSize(10, 20), false},
		{NewSize(0, 20), true},
		{NewSize(10, 0), true},
		{NewSize(-10, 20), true},
		{NewSize(10, -20), true},
		{SizeZero, true},
		{SizeUnspecified, true},
	}

	for _, tt := range tests {
		if got := tt.s.IsEmpty(); got != tt.want {
			t.Errorf("%v.IsEmpty() = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestSizeArithmetic(t *testing.T) {
	s := NewSize(10, 20)

	// Times
	scaled := s.Times(2)
	expectedTimes := NewSize(20, 40)
	if !scaled.Equal(expectedTimes) {
		t.Errorf("%v.Times(2) = %v, want %v", s, scaled, expectedTimes)
	}

	// Div
	divided := s.Div(2)
	expectedDiv := NewSize(5, 10)
	if !divided.Equal(expectedDiv) {
		t.Errorf("%v.Div(2) = %v, want %v", s, divided, expectedDiv)
	}
}

func TestSizeDimensions(t *testing.T) {
	s := NewSize(10, 20)
	if s.MinDimension() != 10 {
		t.Errorf("%v.MinDimension() = %f, want 10", s, s.MinDimension())
	}
	if s.MaxDimension() != 20 {
		t.Errorf("%v.MaxDimension() = %f, want 20", s, s.MaxDimension())
	}

	sNeg := NewSize(-10, -30)
	if sNeg.MinDimension() != 10 {
		t.Errorf("%v.MinDimension() = %f, want 10", sNeg, sNeg.MinDimension())
	}
	if sNeg.MaxDimension() != 30 {
		t.Errorf("%v.MaxDimension() = %f, want 30", sNeg, sNeg.MaxDimension())
	}
}

func TestSizeCenter(t *testing.T) {
	s := NewSize(20, 40)
	center := s.Center()
	expected := NewOffset(10, 20)
	if !center.Equal(expected) {
		t.Errorf("%v.Center() = %v, want %v", s, center, expected)
	}
}

func TestSizeLerp(t *testing.T) {
	s1 := NewSize(0, 0)
	s2 := NewSize(100, 200)

	l := LerpSize(s1, s2, 0.5)
	expected := NewSize(50, 100)
	if !l.Equal(expected) {
		t.Errorf("LerpSize(%v, %v, 0.5) = %v, want %v", s1, s2, l, expected)
	}
}

func TestSizeString(t *testing.T) {
	s := NewSize(10.5, 20.0)
	if s.String() == "Size.Unspecified" {
		t.Errorf("NewSize(10.5, 20.0).String() should not be Size.Unspecified")
	}

	sUnspec := SizeUnspecified
	if sUnspec.String() != "Size.Unspecified" {
		t.Errorf("SizeUnspecified.String() = %v, want Size.Unspecified", sUnspec.String())
	}
}

func TestSizeTakeOrElse(t *testing.T) {
	s := NewSize(10, 20)
	res := s.TakeOrElse(func() Size { return NewSize(30, 40) })
	if !res.Equal(s) {
		t.Errorf("TakeOrElse on specified size should return itself")
	}

	sUnspec := SizeUnspecified
	fallback := NewSize(30, 40)
	res2 := sUnspec.TakeOrElse(func() Size { return fallback })
	if !res2.Equal(fallback) {
		t.Errorf("TakeOrElse on unspecified size should return fallback")
	}
}

func TestSizeEquality(t *testing.T) {
	s1 := NewSize(10, 20)
	s2 := NewSize(10, 20)
	s3 := NewSize(10, 20.0001) // Small diff

	if !s1.Equal(s2) {
		t.Errorf("%v should equal %v", s1, s2)
	}

	if s1.Equal(s3) {
		t.Errorf("%v should not equal %v (diff)", s1, s3)
	}

	// Assuming floatutils.Float32EqualityThreshold is small enough
	// If it's very strict, this might fail, but let's test exact inequality first
	s4 := NewSize(10, 21)
	if s1.Equal(s4) {
		t.Errorf("%v should not equal %v", s1, s4)
	}

	u1 := SizeUnspecified
	u2 := SizeUnspecified
	if !u1.Equal(u2) {
		t.Errorf("SizeUnspecified should equal SizeUnspecified")
	}
}
