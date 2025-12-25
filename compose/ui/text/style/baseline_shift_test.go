package style

import (
	"math"
	"testing"
)

func TestBaselineShift_Constants(t *testing.T) {
	// Verify constant values match Kotlin's definitions
	if BaselineShiftSuperscript != 0.5 {
		t.Errorf("Expected Superscript to be 0.5, got %v", BaselineShiftSuperscript)
	}
	if BaselineShiftSubscript != -0.5 {
		t.Errorf("Expected Subscript to be -0.5, got %v", BaselineShiftSubscript)
	}
	if BaselineShiftNone != 0.0 {
		t.Errorf("Expected None to be 0.0, got %v", BaselineShiftNone)
	}
	if !math.IsNaN(float64(BaselineShiftUnspecified)) {
		t.Errorf("Expected Unspecified to be NaN, got %v", BaselineShiftUnspecified)
	}
}

func TestBaselineShift_IsBaselineShift(t *testing.T) {
	tests := []struct {
		name     string
		bs       BaselineShift
		expected bool
	}{
		{"Superscript is specified", BaselineShiftSuperscript, true},
		{"Subscript is specified", BaselineShiftSubscript, true},
		{"None is specified", BaselineShiftNone, true},
		{"Custom value is specified", NewBaselineShift(0.25), true},
		{"Unspecified is not specified", BaselineShiftUnspecified, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsBaselineShift(tt.bs); got != tt.expected {
				t.Errorf("IsBaselineShift() = %v, want %v", got, tt.expected)
			}
		})
	}

}

func TestTakeOrElseBaselineShift(t *testing.T) {
	fallback := NewBaselineShift(0.75)

	t.Run("returns self when specified", func(t *testing.T) {
		result := TakeOrElseBaselineShift(BaselineShiftSuperscript, fallback)
		if result != BaselineShiftSuperscript {
			t.Errorf("Expected Superscript, got %v", result)
		}
	})

	t.Run("returns block result when unspecified", func(t *testing.T) {
		result := TakeOrElseBaselineShift(BaselineShiftUnspecified, fallback)
		if result != fallback {
			t.Errorf("Expected fallback %v, got %v", fallback, result)
		}
	})

	t.Run("None returns self (not fallback)", func(t *testing.T) {
		result := TakeOrElseBaselineShift(BaselineShiftNone, fallback)
		if result != BaselineShiftNone {
			t.Errorf("Expected None, got %v", result)
		}
	})
}

func TestMergeBaselineShift(t *testing.T) {
	t.Run("returns first when specified", func(t *testing.T) {
		result := MergeBaselineShift(BaselineShiftSuperscript, BaselineShiftSubscript)
		if result != BaselineShiftSuperscript {
			t.Errorf("Expected Superscript, got %v", result)
		}
	})

	t.Run("returns second when first is unspecified", func(t *testing.T) {
		result := MergeBaselineShift(BaselineShiftUnspecified, BaselineShiftSubscript)
		if result != BaselineShiftSubscript {
			t.Errorf("Expected Subscript, got %v", result)
		}
	})

	t.Run("returns unspecified when both are unspecified", func(t *testing.T) {
		result := MergeBaselineShift(BaselineShiftUnspecified, BaselineShiftUnspecified)
		if IsBaselineShift(result) {
			t.Errorf("Expected Unspecified, got %v", result)
		}
	})
}

func TestBaselineShift_Multiplier(t *testing.T) {
	tests := []struct {
		name     string
		bs       BaselineShift
		expected float32
	}{
		{"Superscript", BaselineShiftSuperscript, 0.5},
		{"Subscript", BaselineShiftSubscript, -0.5},
		{"None", BaselineShiftNone, 0.0},
		{"Custom", NewBaselineShift(1.5), 1.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bs.Multiplier(); got != tt.expected {
				t.Errorf("Multiplier() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNewBaselineShift(t *testing.T) {
	bs := NewBaselineShift(0.33)
	if bs.Multiplier() != 0.33 {
		t.Errorf("Expected multiplier 0.33, got %v", bs.Multiplier())
	}
	if !IsBaselineShift(bs) {
		t.Errorf("Expected new BaselineShift to be specified")
	}
}

func TestBaselineShift_String(t *testing.T) {
	tests := []struct {
		name     string
		bs       BaselineShift
		expected string
	}{
		{"Superscript", BaselineShiftSuperscript, "BaselineShift.Superscript"},
		{"Subscript", BaselineShiftSubscript, "BaselineShift.Subscript"},
		{"None", BaselineShiftNone, "BaselineShift.None"},
		{"Unspecified", BaselineShiftUnspecified, "BaselineShift.Unspecified"},
		{"Custom", NewBaselineShift(0.25), "BaselineShift(multiplier=0.25)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringBaselineShift(tt.bs); got != tt.expected {
				t.Errorf("StringBaselineShift() = %v, want %v", got, tt.expected)
			}
			// Also verify the String() method still works
			if got := tt.bs.String(); got != tt.expected {
				t.Errorf("String() = %v, want %v", got, tt.expected)
			}
		})
	}

}

func TestLerpBaselineShift(t *testing.T) {
	t.Run("lerp between None and Superscript", func(t *testing.T) {
		result := LerpBaselineShift(BaselineShiftNone, BaselineShiftSuperscript, 0.5)
		expected := BaselineShift(0.25)
		if result != expected {
			t.Errorf("LerpBaselineShift(None, Superscript, 0.5) = %v, want %v", result, expected)
		}
	})

	t.Run("lerp at fraction 0 returns start", func(t *testing.T) {
		result := LerpBaselineShift(BaselineShiftSubscript, BaselineShiftSuperscript, 0.0)
		if result != BaselineShiftSubscript {
			t.Errorf("LerpBaselineShift at fraction 0 = %v, want %v", result, BaselineShiftSubscript)
		}
	})

	t.Run("lerp at fraction 1 returns stop", func(t *testing.T) {
		result := LerpBaselineShift(BaselineShiftSubscript, BaselineShiftSuperscript, 1.0)
		if result != BaselineShiftSuperscript {
			t.Errorf("LerpBaselineShift at fraction 1 = %v, want %v", result, BaselineShiftSuperscript)
		}
	})

	t.Run("lerp between Subscript and Superscript", func(t *testing.T) {
		// -0.5 to 0.5, at 0.5 fraction should be 0
		result := LerpBaselineShift(BaselineShiftSubscript, BaselineShiftSuperscript, 0.5)
		if result != BaselineShiftNone {
			t.Errorf("LerpBaselineShift(Subscript, Superscript, 0.5) = %v, want %v", result, BaselineShiftNone)
		}
	})

	t.Run("lerp with Unspecified start produces Unspecified", func(t *testing.T) {
		result := LerpBaselineShift(BaselineShiftUnspecified, BaselineShiftSuperscript, 0.5)
		if IsBaselineShift(result) {
			t.Errorf("Expected Unspecified result when start is Unspecified, got %v", result)
		}
	})

	t.Run("lerp with Unspecified stop produces Unspecified", func(t *testing.T) {
		result := LerpBaselineShift(BaselineShiftSubscript, BaselineShiftUnspecified, 0.5)
		if IsBaselineShift(result) {
			t.Errorf("Expected Unspecified result when stop is Unspecified, got %v", result)
		}
	})

	t.Run("lerp with custom values", func(t *testing.T) {
		start := NewBaselineShift(0.0)
		stop := NewBaselineShift(1.0)
		result := LerpBaselineShift(start, stop, 0.75)
		expected := BaselineShift(0.75)
		if result != expected {
			t.Errorf("LerpBaselineShift(0.0, 1.0, 0.75) = %v, want %v", result, expected)
		}
	})
}
