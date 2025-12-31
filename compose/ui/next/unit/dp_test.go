package unit

import (
	"testing"
)

func TestDpOperations(t *testing.T) {
	d1 := NewDp(10)
	d2 := NewDp(20)

	if d1.Add(d2) != NewDp(30) {
		t.Errorf("Add failed")
	}
	if d2.Subtract(d1) != NewDp(10) {
		t.Errorf("Subtract failed")
	}
	if d1.Times(2) != NewDp(20) {
		t.Errorf("Times failed")
	}
	if d2.Div(2) != NewDp(10) {
		t.Errorf("Div failed")
	}
}

func TestDpMinMax(t *testing.T) {
	d1 := NewDp(10)
	d2 := NewDp(20)

	if MinDp(d1, d2) != d1 {
		t.Errorf("Min failed")
	}
	if MaxDp(d1, d2) != d2 {
		t.Errorf("Max failed")
	}
}

func TestDpCoerce(t *testing.T) {
	d := NewDp(10)
	min := NewDp(5)
	max := NewDp(15)

	if d.CoerceIn(min, max) != d {
		t.Errorf("CoerceIn failed in range")
	}
	if NewDp(0).CoerceIn(min, max) != min {
		t.Errorf("CoerceIn failed below min")
	}
	if NewDp(20).CoerceIn(min, max) != max {
		t.Errorf("CoerceIn failed above max")
	}
}

func TestDpOffset(t *testing.T) {
	o1 := NewDpOffset(NewDp(10), NewDp(20))
	o2 := NewDpOffset(NewDp(5), NewDp(5))

	sum := o1.Add(o2)
	if sum.X != NewDp(15) || sum.Y != NewDp(25) {
		t.Errorf("DpOffset Add failed")
	}

	diff := o1.Subtract(o2)
	if diff.X != NewDp(5) || diff.Y != NewDp(15) {
		t.Errorf("DpOffset Subtract failed")
	}
}

func TestDpSize(t *testing.T) {
	s1 := NewDpSize(NewDp(100), NewDp(200))
	center := s1.Center()

	if center.X != NewDp(50) || center.Y != NewDp(100) {
		t.Errorf("DpSize Center failed")
	}
}

func TestDpRectFromOriginSize(t *testing.T) {
	origin := NewDpOffset(NewDp(10), NewDp(10))
	size := NewDpSize(NewDp(20), NewDp(30))
	rect := NewDpRectFromOriginSize(origin, size)

	if rect.Left != NewDp(10) || rect.Top != NewDp(10) || rect.Right != NewDp(30) || rect.Bottom != NewDp(40) {
		t.Errorf("NewDpRectFromOriginSize failed: %v", rect)
	}

	if rect.Width() != size.Width {
		t.Errorf("Rect Width failed")
	}
	if rect.Height() != size.Height {
		t.Errorf("Rect Height failed")
	}
}
