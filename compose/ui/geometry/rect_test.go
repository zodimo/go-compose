package geometry

import (
	"testing"
)

func TestRectConstruction(t *testing.T) {
	r := NewRect(10, 20, 30, 40)
	if r.Left != 10 || r.Top != 20 || r.Right != 30 || r.Bottom != 40 {
		t.Errorf("NewRect(10, 20, 30, 40) = %v, want Rect(10, 20, 30, 40)", r)
	}

	r2 := RectFromOffsetSize(NewOffset(10, 20), NewSize(20, 20))
	if r2.Left != 10 || r2.Top != 20 || r2.Right != 30 || r2.Bottom != 40 {
		t.Errorf("RectFromOffsetSize(...) = %v, want Rect(10, 20, 30, 40)", r2)
	}

	r3 := RectFromTwoOffsets(NewOffset(10, 20), NewOffset(30, 40))
	if !r3.Equal(r) {
		t.Errorf("RectFromTwoOffsets(...) = %v, want %v", r3, r)
	}

	r4 := RectFromCircle(NewOffset(20, 30), 10)
	expectedCircleRect := NewRect(10, 20, 30, 40)
	if !r4.Equal(expectedCircleRect) {
		t.Errorf("RectFromCircle(...) = %v, want %v", r4, expectedCircleRect)
	}
}

func TestRectProperties(t *testing.T) {
	r := NewRect(10, 20, 30, 40)
	if r.Width() != 20 {
		t.Errorf("Width = %f, want 20", r.Width())
	}
	if r.Height() != 20 {
		t.Errorf("Height = %f, want 20", r.Height())
	}

	s := r.Size()
	if s.Width() != 20 || s.Height() != 20 {
		t.Errorf("Size = %v, want Size(20, 20)", s)
	}

	if r.IsEmpty() {
		t.Errorf("IsEmpty should be false")
	}

	emptyRect := NewRect(10, 20, 10, 20)
	if !emptyRect.IsEmpty() {
		t.Errorf("emptyRect.IsEmpty should be true")
	}
}

func TestRectTransformations(t *testing.T) {
	r := NewRect(10, 20, 30, 40)

	translated := r.Translate(NewOffset(10, 10))
	expectedTranslated := NewRect(20, 30, 40, 50)
	if !translated.Equal(expectedTranslated) {
		t.Errorf("Translate = %v, want %v", translated, expectedTranslated)
	}

	translatedXY := r.TranslateXY(10, 10)
	if !translatedXY.Equal(expectedTranslated) {
		t.Errorf("TranslateXY = %v, want %v", translatedXY, expectedTranslated)
	}

	inflated := r.Inflate(5)
	expectedInflated := NewRect(5, 15, 35, 45)
	if !inflated.Equal(expectedInflated) {
		t.Errorf("Inflate = %v, want %v", inflated, expectedInflated)
	}

	deflated := r.Deflate(5)
	expectedDeflated := NewRect(15, 25, 25, 35)
	if !deflated.Equal(expectedDeflated) {
		t.Errorf("Deflate = %v, want %v", deflated, expectedDeflated)
	}
}

func TestRectIntersection(t *testing.T) {
	r1 := NewRect(0, 0, 100, 100)
	r2 := NewRect(50, 50, 150, 150)

	intersect := r1.Intersect(r2)
	expected := NewRect(50, 50, 100, 100)
	if !intersect.Equal(expected) {
		t.Errorf("Intersect = %v, want %v", intersect, expected)
	}

	if !r1.Overlaps(r2) {
		t.Errorf("Overlaps should be true")
	}

	r3 := NewRect(200, 200, 300, 300)
	if r1.Overlaps(r3) {
		t.Errorf("Overlaps should be false")
	}
}

func TestRectContains(t *testing.T) {
	r := NewRect(0, 0, 100, 100)

	if !r.Contains(NewOffset(50, 50)) {
		t.Errorf("Contains(50, 50) should be true")
	}

	if r.Contains(NewOffset(150, 50)) {
		t.Errorf("Contains(150, 50) should be false")
	}

	// Boundary checks
	if !r.Contains(NewOffset(0, 0)) {
		t.Errorf("Contains(0, 0) should be true (inclusive top-left)")
	}
	if r.Contains(NewOffset(100, 100)) {
		t.Errorf("Contains(100, 100) should be false (exclusive bottom-right)")
	}
}

func TestRectOffsets(t *testing.T) {
	r := NewRect(10, 20, 110, 120) // 100x100

	checkOffset := func(name string, got, want Offset) {
		if !got.Equal(want) {
			t.Errorf("%s = %v, want %v", name, got, want)
		}
	}

	checkOffset("TopLeft", r.TopLeft(), NewOffset(10, 20))
	checkOffset("TopRight", r.TopRight(), NewOffset(110, 20))
	checkOffset("BottomLeft", r.BottomLeft(), NewOffset(10, 120))
	checkOffset("BottomRight", r.BottomRight(), NewOffset(110, 120))
	checkOffset("Center", r.Center(), NewOffset(60, 70))
	checkOffset("TopCenter", r.TopCenter(), NewOffset(60, 20))
	checkOffset("BottomCenter", r.BottomCenter(), NewOffset(60, 120))
	checkOffset("CenterLeft", r.CenterLeft(), NewOffset(10, 70))
	checkOffset("CenterRight", r.CenterRight(), NewOffset(110, 70))
}

func TestRectLerp(t *testing.T) {
	r1 := NewRect(0, 0, 10, 10)
	r2 := NewRect(100, 100, 110, 110)

	l := LerpRect(r1, r2, 0.5)
	expected := NewRect(50, 50, 60, 60)
	if !l.Equal(expected) {
		t.Errorf("LerpRect(...) = %v, want %v", l, expected)
	}
}
