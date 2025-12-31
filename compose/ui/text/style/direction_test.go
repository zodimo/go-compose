package style

import (
	"testing"
)

func TestTextDirectionConstants(t *testing.T) {
	tests := []struct {
		name      string
		direction TextDirection
		expected  int
	}{
		{"Unspecified", TextDirectionUnspecified, 0},
		{"Ltr", TextDirectionLtr, 1},
		{"Rtl", TextDirectionRtl, 2},
		{"Content", TextDirectionContent, 3},
		{"ContentOrLtr", TextDirectionContentOrLtr, 4},
		{"ContentOrRtl", TextDirectionContentOrRtl, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.direction) != tt.expected {
				t.Errorf("TextDirection%s = %d, want %d", tt.name, int(tt.direction), tt.expected)
			}
		})
	}
}

func TestTextDirection_String(t *testing.T) {
	tests := []struct {
		direction TextDirection
		expected  string
	}{
		{TextDirectionUnspecified, "Unspecified"},
		{TextDirectionLtr, "Ltr"},
		{TextDirectionRtl, "Rtl"},
		{TextDirectionContent, "Content"},
		{TextDirectionContentOrLtr, "ContentOrLtr"},
		{TextDirectionContentOrRtl, "ContentOrRtl"},
		{TextDirection(99), "Invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.direction.String(); got != tt.expected {
				t.Errorf("TextDirection.String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestTextDirection_IsSpecified(t *testing.T) {
	tests := []struct {
		name      string
		direction TextDirection
		expected  bool
	}{
		{"Unspecified", TextDirectionUnspecified, false},
		{"Ltr", TextDirectionLtr, true},
		{"Rtl", TextDirectionRtl, true},
		{"Content", TextDirectionContent, true},
		{"ContentOrLtr", TextDirectionContentOrLtr, true},
		{"ContentOrRtl", TextDirectionContentOrRtl, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.direction.IsSpecified(); got != tt.expected {
				t.Errorf("TextDirection%s.IsSpecified() = %v, want %v", tt.name, got, tt.expected)
			}
		})
	}
}

func TestTextDirectionValues(t *testing.T) {
	values := TextDirectionValues()

	expected := []TextDirection{
		TextDirectionLtr,
		TextDirectionRtl,
		TextDirectionContent,
		TextDirectionContentOrLtr,
		TextDirectionContentOrRtl,
	}

	if len(values) != len(expected) {
		t.Fatalf("TextDirectionValues() returned %d values, want %d", len(values), len(expected))
	}

	for i, v := range values {
		if v != expected[i] {
			t.Errorf("TextDirectionValues()[%d] = %v, want %v", i, v, expected[i])
		}
	}

	// Verify Unspecified is not in the list
	for _, v := range values {
		if v == TextDirectionUnspecified {
			t.Error("TextDirectionValues() should not contain Unspecified")
		}
	}
}

func TestTextDirectionValueOf(t *testing.T) {
	tests := []struct {
		value       int
		expected    TextDirection
		expectError bool
	}{
		{0, TextDirectionUnspecified, false},
		{1, TextDirectionLtr, false},
		{2, TextDirectionRtl, false},
		{3, TextDirectionContent, false},
		{4, TextDirectionContentOrLtr, false},
		{5, TextDirectionContentOrRtl, false},
		{-1, TextDirectionUnspecified, true},
		{6, TextDirectionUnspecified, true},
		{100, TextDirectionUnspecified, true},
	}

	for _, tt := range tests {
		t.Run(tt.expected.String(), func(t *testing.T) {
			got, err := TextDirectionValueOf(tt.value)

			if tt.expectError {
				if err == nil {
					t.Errorf("TextDirectionValueOf(%d) expected error, got nil", tt.value)
				}
			} else {
				if err != nil {
					t.Errorf("TextDirectionValueOf(%d) unexpected error: %v", tt.value, err)
				}
				if got != tt.expected {
					t.Errorf("TextDirectionValueOf(%d) = %v, want %v", tt.value, got, tt.expected)
				}
			}
		})
	}
}

func TestTextDirection_TakeOrElse(t *testing.T) {
	defaultDir := TextDirectionLtr

	t.Run("specified value returns itself", func(t *testing.T) {
		dir := TextDirectionRtl
		result := dir.TakeOrElse(defaultDir)
		if result != TextDirectionRtl {
			t.Errorf("TakeOrElse() = %v, want %v", result, TextDirectionRtl)
		}
	})

	t.Run("unspecified value returns block result", func(t *testing.T) {
		dir := TextDirectionUnspecified
		result := dir.TakeOrElse(defaultDir)
		if result != defaultDir {
			t.Errorf("TakeOrElse() = %v, want %v", result, defaultDir)
		}
	})

	t.Run("block is not called when specified", func(t *testing.T) {
		dir := TextDirectionContent
		takeDir := dir.TakeOrElse(defaultDir)
		if takeDir != TextDirectionContent {
			t.Errorf("TakeOrElse() = %v, want %v", takeDir, TextDirectionContent)
		}
	})

	t.Run("block is called when unspecified", func(t *testing.T) {
		dir := TextDirectionUnspecified
		takeDir := dir.TakeOrElse(defaultDir)
		if takeDir != defaultDir {
			t.Errorf("TakeOrElse() = %v, want %v", takeDir, defaultDir)
		}
	})
}

func TestTextDirection_Equality(t *testing.T) {
	t.Run("same values are equal", func(t *testing.T) {
		a := TextDirection(3)
		b := TextDirection(3)
		if a != b {
			t.Errorf("TextDirection(3) != TextDirection(3)")
		}
	})

	t.Run("different values are not equal", func(t *testing.T) {
		a := TextDirection(1)
		b := TextDirection(2)
		if a == b {
			t.Errorf("TextDirection(1) == TextDirection(2)")
		}
	})

	t.Run("constants are equal to equivalent values", func(t *testing.T) {
		if TextDirectionContent != TextDirection(3) {
			t.Errorf("TextDirectionContent != TextDirection(3)")
		}
	})
}
