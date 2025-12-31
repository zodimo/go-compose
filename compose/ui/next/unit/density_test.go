package unit

import (
	"testing"

	"github.com/zodimo/go-compose/compose/ui/geometry"
)

func TestDensityValues(t *testing.T) {
	density := NewDensity(2.0, 3.0)
	if density.Density() != 2.0 {
		t.Errorf("Expected density 2.0, got %f", density.Density())
	}
	if density.FontScale() != 3.0 {
		t.Errorf("Expected fontScale 3.0, got %f", density.FontScale())
	}
}

func TestDpToPx(t *testing.T) {
	density := NewDensity(2.0, 3.0)
	px := density.DpToPx(NewDp(1.0))
	if px != 2.0 {
		t.Errorf("Expected 2.0, got %f", px)
	}
}

func TestDpRoundToPx(t *testing.T) {
	density := NewDensity(2.0, 3.0)

	val1 := density.DpRoundToPx(NewDp(4.95))
	if val1 != 10 {
		t.Errorf("Expected 10, got %d", val1)
	}

	val2 := density.DpRoundToPx(NewDp(4.75))
	if val2 != 10 {
		t.Errorf("Expected 10, got %d", val2)
	}

	val3 := density.DpRoundToPx(NewDp(4.74))
	if val3 != 9 {
		t.Errorf("Expected 9, got %d", val3)
	}
}

func TestIntToDp(t *testing.T) {
	density := NewDensity(2.0, 3.0)
	dp := density.IntToDp(2)
	expected := NewDp(1.0)
	if dp.CompareTo(expected) != 0 {
		t.Errorf("Expected %v, got %v", expected, dp)
	}
}

func TestFloatToDp(t *testing.T) {
	density := NewDensity(2.0, 3.0)
	dp := density.FloatToDp(2.0)
	expected := NewDp(1.0)
	if dp.CompareTo(expected) != 0 {
		t.Errorf("Expected %v, got %v", expected, dp)
	}
}

func TestDpRectToRect(t *testing.T) {
	density := NewDensity(2.0, 3.0)
	dpRect := NewDpRect(NewDp(1.0), NewDp(2.0), NewDp(3.0), NewDp(4.0))
	rect := density.DpRectToRect(dpRect)

	if rect.Left != 2.0 {
		t.Errorf("Expected left 2.0, got %f", rect.Left)
	}
	if rect.Top != 4.0 {
		t.Errorf("Expected top 4.0, got %f", rect.Top)
	}
	if rect.Right != 6.0 {
		t.Errorf("Expected right 6.0, got %f", rect.Right)
	}
	if rect.Bottom != 8.0 {
		t.Errorf("Expected bottom 8.0, got %f", rect.Bottom)
	}
}

func TestDpSizeToSize(t *testing.T) {
	density := NewDensity(2.0, 3.0)
	dpSize := NewDpSize(NewDp(1.0), NewDp(3.0))
	size := density.DpSizeToSize(dpSize)

	expected := geometry.NewSize(2.0, 6.0)
	if !size.Equal(expected) {
		t.Errorf("Expected %v, got %v", expected, size)
	}
}

func TestSizeToDpSize(t *testing.T) {
	density := NewDensity(2.0, 3.0)
	size := geometry.NewSize(2.0, 6.0)
	dpSize := density.SizeToDpSize(size)

	expected := NewDpSize(NewDp(1.0), NewDp(3.0))
	// Need DpSize equality check or manual
	if dpSize.Width.CompareTo(expected.Width) != 0 || dpSize.Height.CompareTo(expected.Height) != 0 {
		t.Errorf("Expected %v, got %v", expected, dpSize)
	}
}

func TestDpSizeUnspecifiedToSize(t *testing.T) {
	density := NewDensity(2.0, 3.0)
	size := density.DpSizeToSize(DpSizeUnspecified)
	if !size.IsUnspecified() {
		t.Errorf("Expected Unspecified, got %v", size)
	}
}

func TestSizeUnspecifiedToDpSize(t *testing.T) {
	density := NewDensity(2.0, 3.0)
	dpSize := density.SizeToDpSize(geometry.SizeUnspecified)
	if !dpSize.Width.IsUnspecified() || !dpSize.Height.IsUnspecified() {
		t.Errorf("Expected valid Unspecified parts, got %v", dpSize)
	}
	// Strict equality might be tricky with NaN, but IsUnspecified checks are sufficient.
}

func TestTextUnitToPx(t *testing.T) {
	density := NewDensity(2.0, 3.0)
	// Sp -> Dp -> Px
	// Sp(10) -> Dp(10 * 3) = Dp(30) -> Px(30 * 2) = 60
	sp := Sp(10)
	px := density.TextUnitToPx(sp)
	if px != 60.0 {
		t.Errorf("Expected 60.0, got %f", px)
	}
}

func TestTextUnitRoundToPx(t *testing.T) {
	density := NewDensity(2.0, 3.0)
	// Sp(10.1) -> Dp(30.3) -> Px(60.6) -> Round(61)
	sp := Sp(10.1)
	px := density.TextUnitRoundToPx(sp)
	if px != 61 {
		t.Errorf("Expected 61, got %d", px)
	}
}
