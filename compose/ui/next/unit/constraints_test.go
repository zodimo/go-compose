package unit

import (
	"testing"
)

func TestConstraints_CreationAndGetters(t *testing.T) {
	tests := []struct {
		name      string
		minW      int
		maxW      int
		minH      int
		maxH      int
		expectErr bool
	}{
		{"Simple", 10, 20, 30, 40, false},
		{"Fixed", 50, 50, 60, 60, false},
		{"Unbounded", 0, Infinity, 0, Infinity, false},
		{"Mixed Unbounded", 10, Infinity, 20, 50, false},
		{"Invalid Min > Max", 20, 10, 30, 40, true},
		{"Invalid Negative", -1, 10, 0, 10, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tt.expectErr {
					if r == nil {
						t.Errorf("Expected panic but did not panic")
					}
				} else {
					if r != nil {
						t.Errorf("Unexpected panic: %v", r)
					}
				}
			}()

			c := NewConstraints(tt.minW, tt.maxW, tt.minH, tt.maxH)

			if !tt.expectErr {
				if c.MinWidth() != tt.minW {
					t.Errorf("Expected MinWidth %d, got %d", tt.minW, c.MinWidth())
				}
				if c.MaxWidth() != tt.maxW {
					t.Errorf("Expected MaxWidth %d, got %d", tt.maxW, c.MaxWidth())
				}
				if c.MinHeight() != tt.minH {
					t.Errorf("Expected MinHeight %d, got %d", tt.minH, c.MinHeight())
				}
				if c.MaxHeight() != tt.maxH {
					t.Errorf("Expected MaxHeight %d, got %d", tt.maxH, c.MaxHeight())
				}
			}
		})
	}
}

func TestConstraints_BitPacking(t *testing.T) {
	// Test boundary values for bit packing
	// Note: These values depend on the internal constants in constraints.go

	// Case 1: Small dimensions (should fit in any focus)
	c1 := NewConstraints(100, 200, 100, 200)
	if c1.MinWidth() != 100 || c1.MaxWidth() != 200 {
		t.Errorf("Small dims failed")
	}

	// Case 2: Large Width (requires MaxFocusWidth or similar)
	largeW := 100000 // > 65535, fits in 18 bits (262143)
	c2 := NewConstraints(largeW, largeW, 10, 100)
	if c2.MinWidth() != largeW {
		t.Errorf("Large width min failed: got %d, want %d", c2.MinWidth(), largeW)
	}
	if c2.MaxWidth() != largeW {
		t.Errorf("Large width max failed: got %d, want %d", c2.MaxWidth(), largeW)
	}

	// Case 3: Large Height
	largeH := 100000
	c3 := NewConstraints(10, 100, largeH, largeH)
	if c3.MinHeight() != largeH {
		t.Errorf("Large height min failed: got %d, want %d", c3.MinHeight(), largeH)
	}

	// Case 4: Max allowed values
	// 18 bits max is 262143.
	// 13 bits max is 8191.

	// Try a width that needs 18 bits and a height that fits in 13 bits
	w18 := 200000
	h13 := 5000
	c4 := NewConstraints(w18, w18, h13, h13)
	if c4.MinWidth() != w18 || c4.MinHeight() != h13 {
		t.Errorf("18-bit width failed")
	}

	// Try a height that needs 18 bits and a width that fits in 13 bits
	c5 := NewConstraints(h13, h13, w18, w18)
	if c5.MinWidth() != h13 || c5.MinHeight() != w18 {
		t.Errorf("18-bit height failed")
	}
}

func TestConstraints_Operations(t *testing.T) {
	c := NewConstraints(10, 100, 20, 200)

	// HasBounded
	if !c.HasBoundedWidth() || !c.HasBoundedHeight() {
		t.Errorf("Expected bounded")
	}

	inf := NewConstraints(0, Infinity, 0, Infinity)
	if inf.HasBoundedWidth() || inf.HasBoundedHeight() {
		t.Errorf("Expected unbounded")
	}

	// Fixed
	fixed := Fixed(50, 60)
	if !fixed.HasFixedWidth() || !fixed.HasFixedHeight() {
		t.Errorf("Expected fixed")
	}
	if fixed.MinWidth() != 50 || fixed.MaxWidth() != 50 {
		t.Errorf("Fixed width mismatch")
	}

	// IsSatisfiedBy
	size := IntSize{Width: 50, Height: 100}
	if !c.IsSatisfiedBy(size) {
		t.Errorf("Expected satisfied")
	}

	badSize := IntSize{Width: 5, Height: 100}
	if c.IsSatisfiedBy(badSize) {
		t.Errorf("Expected not satisfied (width too small)")
	}

	// Constrain
	constrainedSize := c.ConstrainSize(IntSize{Width: 500, Height: 500})
	if constrainedSize.Width != 100 || constrainedSize.Height != 200 {
		t.Errorf("ConstrainSize failed: got %v", constrainedSize)
	}

	constrainedSizeLow := c.ConstrainSize(IntSize{Width: 1, Height: 1})
	if constrainedSizeLow.Width != 10 || constrainedSizeLow.Height != 20 {
		t.Errorf("ConstrainSize low failed: got %v", constrainedSizeLow)
	}
}

func TestConstraints_Offset(t *testing.T) {
	c := NewConstraints(10, 100, 20, 200)

	// Positive offset
	res := c.Offset(5, 10)
	// minW: 10+5=15, maxW: 100+5=105
	// minH: 20+10=30, maxH: 200+10=210
	if res.MinWidth() != 15 || res.MaxWidth() != 105 {
		t.Errorf("Offset width failed")
	}
	if res.MinHeight() != 30 || res.MaxHeight() != 210 {
		t.Errorf("Offset height failed")
	}

	// Negative offset
	res2 := c.Offset(-5, -5)
	// minW: 10-5=5, maxW: 100-5=95
	if res2.MinWidth() != 5 || res2.MaxWidth() != 95 {
		t.Errorf("Negative offset failed")
	}

	// Offset leading to negative (should clamp to 0)
	res3 := c.Offset(-20, -30)
	if res3.MinWidth() != 0 || res3.MinHeight() != 0 {
		t.Errorf("Offset clamping failed")
	}
}
