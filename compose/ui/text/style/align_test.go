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
		{"Unspecified", TextAlignUnspecified, 99},
		{"Start", TextAlignStart, 0},
		{"End", TextAlignEnd, 1},
		{"Middle", TextAlignMiddle, 2},
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
		{TextAlignStart, "Start"},
		{TextAlignEnd, "End"},
		{TextAlignMiddle, "Middle"},
		{TextAlign(100), "Invalid"},
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
		{"Start", TextAlignStart, true},
		{"End", TextAlignEnd, true},
		{"Middle", TextAlignMiddle, true},
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
		TextAlignStart,
		TextAlignEnd,
		TextAlignMiddle,
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
		{99, TextAlignUnspecified, false},
		{0, TextAlignStart, false},
		{1, TextAlignEnd, false},
		{2, TextAlignMiddle, false},
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
	defaultAlign := TextAlignStart

	t.Run("specified value returns itself", func(t *testing.T) {
		align := TextAlignMiddle
		result := align.TakeOrElse(defaultAlign)
		if result != TextAlignMiddle {
			t.Errorf("TakeOrElse() = %v, want %v", result, TextAlignMiddle)
		}
	})

	t.Run("unspecified value returns block result", func(t *testing.T) {
		align := TextAlignUnspecified
		result := align.TakeOrElse(defaultAlign)
		if result != defaultAlign {
			t.Errorf("TakeOrElse() = %v, want %v", result, defaultAlign)
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
		a := TextAlignMiddle
		b := TextAlignMiddle
		if a != b {
			t.Errorf("TextAlignMiddle != TextAlignMiddle")
		}
	})

	t.Run("different values are not equal", func(t *testing.T) {
		a := TextAlignStart
		b := TextAlignEnd
		if a == b {
			t.Errorf("TextAlignStart == TextAlignEnd")
		}
	})
}
