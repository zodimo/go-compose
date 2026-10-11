package style

import (
	"testing"
)

func TestTextAlignConstants(t *testing.T) {
	tests := []struct {
		name     string
		align    TextAlign
		expected int
	}{
		{"Unspecified", TextAlignUnspecified, 0},
		{"Left", TextAlignLeft, 1},
		{"Right", TextAlignRight, 2},
		{"Center", TextAlignCenter, 3},
		{"Justify", TextAlignJustify, 4},
		{"Start", TextAlignStart, 5},
		{"End", TextAlignEnd, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.align) != tt.expected {
				t.Errorf("TextAlign%s = %d, want %d", tt.name, int(tt.align), tt.expected)
			}
		})
	}
}

func TestTextAlign_String(t *testing.T) {
	tests := []struct {
		align    TextAlign
		expected string
	}{
		{TextAlignUnspecified, "Unspecified"},
		{TextAlignLeft, "Left"},
		{TextAlignRight, "Right"},
		{TextAlignCenter, "Center"},
		{TextAlignJustify, "Justify"},
		{TextAlignStart, "Start"},
		{TextAlignEnd, "End"},
		{TextAlign(99), "Invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.align.String(); got != tt.expected {
				t.Errorf("TextAlign.String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestTextAlign_IsSpecified(t *testing.T) {
	tests := []struct {
		name     string
		align    TextAlign
		expected bool
	}{
		{"Unspecified", TextAlignUnspecified, false},
		{"Left", TextAlignLeft, true},
		{"Right", TextAlignRight, true},
		{"Center", TextAlignCenter, true},
		{"Justify", TextAlignJustify, true},
		{"Start", TextAlignStart, true},
		{"End", TextAlignEnd, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.align.IsSpecified(); got != tt.expected {
				t.Errorf("TextAlign%s.IsSpecified() = %v, want %v", tt.name, got, tt.expected)
			}
		})
	}
}

func TestTextAlignValues(t *testing.T) {
	values := TextAlignValues()

	expected := []TextAlign{
		TextAlignLeft,
		TextAlignRight,
		TextAlignCenter,
		TextAlignJustify,
		TextAlignStart,
		TextAlignEnd,
	}

	if len(values) != len(expected) {
		t.Fatalf("TextAlignValues() returned %d values, want %d", len(values), len(expected))
	}

	for i, v := range values {
		if v != expected[i] {
			t.Errorf("TextAlignValues()[%d] = %v, want %v", i, v, expected[i])
		}
	}

	// Verify Unspecified is not in the list
	for _, v := range values {
		if v == TextAlignUnspecified {
			t.Error("TextAlignValues() should not contain Unspecified")
		}
	}
}

func TestTextAlignValueOf(t *testing.T) {
	tests := []struct {
		value       int
		expected    TextAlign
		expectError bool
	}{
		{0, TextAlignUnspecified, false},
		{1, TextAlignLeft, false},
		{2, TextAlignRight, false},
		{3, TextAlignCenter, false},
		{4, TextAlignJustify, false},
		{5, TextAlignStart, false},
		{6, TextAlignEnd, false},
		{-1, TextAlignUnspecified, true},
		{7, TextAlignUnspecified, true},
		{100, TextAlignUnspecified, true},
	}

	for _, tt := range tests {
		t.Run(tt.expected.String(), func(t *testing.T) {
			got, err := TextAlignValueOf(tt.value)

			if tt.expectError {
				if err == nil {
					t.Errorf("TextAlignValueOf(%d) expected error, got nil", tt.value)
				}
			} else {
				if err != nil {
					t.Errorf("TextAlignValueOf(%d) unexpected error: %v", tt.value, err)
				}
				if got != tt.expected {
					t.Errorf("TextAlignValueOf(%d) = %v, want %v", tt.value, got, tt.expected)
				}
			}
		})
	}
}

func TestTextAlign_TakeOrElse(t *testing.T) {
	defaultAlign := TextAlignLeft

	t.Run("specified value returns itself", func(t *testing.T) {
		align := TextAlignCenter
		result := align.TakeOrElse(defaultAlign)
		if result != TextAlignCenter {
			t.Errorf("TakeOrElse() = %v, want %v", result, TextAlignCenter)
		}
	})

	t.Run("unspecified value returns block result", func(t *testing.T) {
		align := TextAlignUnspecified
		result := align.TakeOrElse(defaultAlign)
		if result != defaultAlign {
			t.Errorf("TakeOrElse() = %v, want %v", result, defaultAlign)
		}
	})

	t.Run("block is not called when specified", func(t *testing.T) {
		align := TextAlignRight
		takeAlign := align.TakeOrElse(defaultAlign)
		if takeAlign != TextAlignRight {
			t.Errorf("TakeOrElse() = %v, want %v", takeAlign, TextAlignRight)
		}
	})

	t.Run("block is called when unspecified", func(t *testing.T) {
		align := TextAlignUnspecified
		takeAlign := align.TakeOrElse(defaultAlign)
		if takeAlign != defaultAlign {
			t.Errorf("TakeOrElse() = %v, want %v", takeAlign, defaultAlign)
		}
	})
}

func TestTextAlign_Equality(t *testing.T) {
	t.Run("same values are equal", func(t *testing.T) {
		a := TextAlignCenter
		b := TextAlignCenter
		if a != b {
			t.Errorf("TextAlignCenter != TextAlignCenter")
		}
	})

	t.Run("different values are not equal", func(t *testing.T) {
		a := TextAlignLeft
		b := TextAlignRight
		if a == b {
			t.Errorf("TextAlignLeft == TextAlignRight")
		}
	})
}
