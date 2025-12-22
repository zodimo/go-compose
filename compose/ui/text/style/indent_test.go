package style

import (
	"testing"

	"github.com/zodimo/go-compose/compose/ui/unit"
)

func TestTextIndentNone(t *testing.T) {
	none := TextIndentNone
	if none.FirstLine.Value() != 0 || !none.FirstLine.IsSp() {
		t.Errorf("TextIndentNone.FirstLine should be 0.sp, got %s", none.FirstLine)
	}
	if none.RestLine.Value() != 0 || !none.RestLine.IsSp() {
		t.Errorf("TextIndentNone.RestLine should be 0.sp, got %s", none.RestLine)
	}
}

func TestNewTextIndent(t *testing.T) {
	ti := NewTextIndent(unit.Sp(10), unit.Sp(5))
	if ti.FirstLine.Value() != 10 {
		t.Errorf("FirstLine should be 10, got %v", ti.FirstLine.Value())
	}
	if ti.RestLine.Value() != 5 {
		t.Errorf("RestLine should be 5, got %v", ti.RestLine.Value())
	}
}

func TestTextIndent_Copy(t *testing.T) {
	original := NewTextIndent(unit.Sp(10), unit.Sp(5))

	// Copy with no changes
	copy1 := original.Copy(nil, nil)
	if !copy1.Equals(original) {
		t.Errorf("Copy with nil args should equal original")
	}

	// Copy with FirstLine changed
	newFirstLine := unit.Sp(20)
	copy2 := original.Copy(&newFirstLine, nil)
	if copy2.FirstLine.Value() != 20 {
		t.Errorf("Copy FirstLine should be 20, got %v", copy2.FirstLine.Value())
	}
	if copy2.RestLine.Value() != 5 {
		t.Errorf("Copy RestLine should remain 5, got %v", copy2.RestLine.Value())
	}

	// Copy with RestLine changed
	newRestLine := unit.Sp(15)
	copy3 := original.Copy(nil, &newRestLine)
	if copy3.FirstLine.Value() != 10 {
		t.Errorf("Copy FirstLine should remain 10, got %v", copy3.FirstLine.Value())
	}
	if copy3.RestLine.Value() != 15 {
		t.Errorf("Copy RestLine should be 15, got %v", copy3.RestLine.Value())
	}
}

func TestTextIndent_Equals(t *testing.T) {
	ti1 := NewTextIndent(unit.Sp(10), unit.Sp(5))
	ti2 := NewTextIndent(unit.Sp(10), unit.Sp(5))
	ti3 := NewTextIndent(unit.Sp(10), unit.Sp(6))
	ti4 := NewTextIndent(unit.Sp(11), unit.Sp(5))

	if !ti1.Equals(ti2) {
		t.Error("Equal TextIndents should be equal")
	}
	if ti1.Equals(ti3) {
		t.Error("TextIndents with different RestLine should not be equal")
	}
	if ti1.Equals(ti4) {
		t.Error("TextIndents with different FirstLine should not be equal")
	}
}

func TestTextIndent_String(t *testing.T) {
	ti := NewTextIndent(unit.Sp(10), unit.Sp(5))
	str := ti.String()
	expected := "TextIndent(firstLine=10.sp, restLine=5.sp)"
	if str != expected {
		t.Errorf("String() = %q, want %q", str, expected)
	}
}

func TestTextIndent_HashCode(t *testing.T) {
	ti1 := NewTextIndent(unit.Sp(10), unit.Sp(5))
	ti2 := NewTextIndent(unit.Sp(10), unit.Sp(5))
	ti3 := NewTextIndent(unit.Sp(10), unit.Sp(6))

	if ti1.HashCode() != ti2.HashCode() {
		t.Error("Equal TextIndents should have equal hash codes")
	}
	if ti1.HashCode() == ti3.HashCode() {
		t.Error("Different TextIndents should likely have different hash codes")
	}
}

// Tests for lerpDiscrete
func TestLerpDiscrete(t *testing.T) {
	a := "start"
	b := "stop"

	// Fraction < 0.5 should return a
	if result := lerpDiscrete(a, b, 0.0); result != a {
		t.Errorf("lerpDiscrete at 0.0 should return a, got %v", result)
	}
	if result := lerpDiscrete(a, b, 0.25); result != a {
		t.Errorf("lerpDiscrete at 0.25 should return a, got %v", result)
	}
	if result := lerpDiscrete(a, b, 0.49); result != a {
		t.Errorf("lerpDiscrete at 0.49 should return a, got %v", result)
	}

	// Fraction >= 0.5 should return b
	if result := lerpDiscrete(a, b, 0.5); result != b {
		t.Errorf("lerpDiscrete at 0.5 should return b, got %v", result)
	}
	if result := lerpDiscrete(a, b, 0.75); result != b {
		t.Errorf("lerpDiscrete at 0.75 should return b, got %v", result)
	}
	if result := lerpDiscrete(a, b, 1.0); result != b {
		t.Errorf("lerpDiscrete at 1.0 should return b, got %v", result)
	}
}

// Tests for LerpTextUnitInheritable
func TestLerpTextUnitInheritable_BothSpecified(t *testing.T) {
	a := unit.Sp(10)
	b := unit.Sp(20)

	// Standard linear interpolation when both specified
	result := LerpTextUnitInheritable(a, b, 0.5)
	if result.Value() != 15 {
		t.Errorf("LerpTextUnitInheritable(10sp, 20sp, 0.5) = %v, want 15", result.Value())
	}

	result0 := LerpTextUnitInheritable(a, b, 0.0)
	if result0.Value() != 10 {
		t.Errorf("LerpTextUnitInheritable(10sp, 20sp, 0.0) = %v, want 10", result0.Value())
	}

	result1 := LerpTextUnitInheritable(a, b, 1.0)
	if result1.Value() != 20 {
		t.Errorf("LerpTextUnitInheritable(10sp, 20sp, 1.0) = %v, want 20", result1.Value())
	}
}

func TestLerpTextUnitInheritable_StartUnspecified(t *testing.T) {
	a := unit.TextUnitUnspecified
	b := unit.Sp(16)

	// Discrete interpolation: before 0.5 returns Unspecified
	resultBefore := LerpTextUnitInheritable(a, b, 0.25)
	if !resultBefore.IsUnspecified() {
		t.Errorf("LerpTextUnitInheritable(Unspecified, 16sp, 0.25) should be Unspecified, got %s", resultBefore)
	}

	// Discrete interpolation: at/after 0.5 returns the specified value
	resultAfter := LerpTextUnitInheritable(a, b, 0.5)
	if resultAfter.IsUnspecified() || resultAfter.Value() != 16 {
		t.Errorf("LerpTextUnitInheritable(Unspecified, 16sp, 0.5) should be 16sp, got %s", resultAfter)
	}
}

func TestLerpTextUnitInheritable_StopUnspecified(t *testing.T) {
	a := unit.Sp(16)
	b := unit.TextUnitUnspecified

	// Discrete interpolation: before 0.5 returns the specified value
	resultBefore := LerpTextUnitInheritable(a, b, 0.25)
	if resultBefore.IsUnspecified() || resultBefore.Value() != 16 {
		t.Errorf("LerpTextUnitInheritable(16sp, Unspecified, 0.25) should be 16sp, got %s", resultBefore)
	}

	// Discrete interpolation: at/after 0.5 returns Unspecified
	resultAfter := LerpTextUnitInheritable(a, b, 0.5)
	if !resultAfter.IsUnspecified() {
		t.Errorf("LerpTextUnitInheritable(16sp, Unspecified, 0.5) should be Unspecified, got %s", resultAfter)
	}
}

func TestLerpTextUnitInheritable_BothUnspecified(t *testing.T) {
	a := unit.TextUnitUnspecified
	b := unit.TextUnitUnspecified

	// Both unspecified: always returns Unspecified
	result := LerpTextUnitInheritable(a, b, 0.5)
	if !result.IsUnspecified() {
		t.Errorf("LerpTextUnitInheritable(Unspecified, Unspecified, 0.5) should be Unspecified, got %s", result)
	}
}

// Tests for LerpTextIndent
func TestLerpTextIndent_BothSpecified(t *testing.T) {
	start := NewTextIndent(unit.Sp(0), unit.Sp(0))
	stop := NewTextIndent(unit.Sp(20), unit.Sp(10))

	result := LerpTextIndent(start, stop, 0.5)
	if result.FirstLine.Value() != 10 {
		t.Errorf("LerpTextIndent FirstLine at 0.5 = %v, want 10", result.FirstLine.Value())
	}
	if result.RestLine.Value() != 5 {
		t.Errorf("LerpTextIndent RestLine at 0.5 = %v, want 5", result.RestLine.Value())
	}
}

func TestLerpTextIndent_WithUnspecified(t *testing.T) {
	start := NewTextIndent(unit.TextUnitUnspecified, unit.Sp(0))
	stop := NewTextIndent(unit.Sp(20), unit.Sp(10))

	// Before 0.5: FirstLine is Unspecified (from start)
	resultBefore := LerpTextIndent(start, stop, 0.25)
	if !resultBefore.FirstLine.IsUnspecified() {
		t.Errorf("LerpTextIndent FirstLine at 0.25 should be Unspecified, got %s", resultBefore.FirstLine)
	}
	// RestLine interpolates normally since both are specified
	if resultBefore.RestLine.Value() != 2.5 {
		t.Errorf("LerpTextIndent RestLine at 0.25 = %v, want 2.5", resultBefore.RestLine.Value())
	}

	// After 0.5: FirstLine snaps to specified value
	resultAfter := LerpTextIndent(start, stop, 0.75)
	if resultAfter.FirstLine.IsUnspecified() || resultAfter.FirstLine.Value() != 20 {
		t.Errorf("LerpTextIndent FirstLine at 0.75 should be 20sp, got %s", resultAfter.FirstLine)
	}
}

func TestLerpTextIndent_StartAndStop(t *testing.T) {
	start := NewTextIndent(unit.Sp(10), unit.Sp(5))
	stop := NewTextIndent(unit.Sp(30), unit.Sp(15))

	// At fraction 0, should equal start
	result0 := LerpTextIndent(start, stop, 0.0)
	if !result0.Equals(start) {
		t.Errorf("LerpTextIndent at 0.0 should equal start")
	}

	// At fraction 1, should equal stop
	result1 := LerpTextIndent(start, stop, 1.0)
	if !result1.Equals(stop) {
		t.Errorf("LerpTextIndent at 1.0 should equal stop")
	}
}
