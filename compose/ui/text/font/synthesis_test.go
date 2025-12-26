package font

import "testing"

func TestFontSynthesis_Constants(t *testing.T) {
	if FontSynthesisNone.Value() != 0 {
		t.Errorf("None should be 0, got %d", FontSynthesisNone.Value())
	}
	if FontSynthesisWeight.Value() != 1 {
		t.Errorf("Weight should be 1, got %d", FontSynthesisWeight.Value())
	}
	if FontSynthesisStyle.Value() != 2 {
		t.Errorf("Style should be 2, got %d", FontSynthesisStyle.Value())
	}
	if FontSynthesisAll.Value() != 0xffff {
		t.Errorf("All should be 0xffff, got %d", FontSynthesisAll.Value())
	}
}

func TestFontSynthesis_IsWeightOn(t *testing.T) {
	if FontSynthesisNone.IsWeightOn() {
		t.Error("None should not have weight on")
	}
	if !FontSynthesisWeight.IsWeightOn() {
		t.Error("Weight should have weight on")
	}
	if FontSynthesisStyle.IsWeightOn() {
		t.Error("Style should not have weight on")
	}
	if !FontSynthesisAll.IsWeightOn() {
		t.Error("All should have weight on")
	}
}

func TestFontSynthesis_IsStyleOn(t *testing.T) {
	if FontSynthesisNone.IsStyleOn() {
		t.Error("None should not have style on")
	}
	if FontSynthesisWeight.IsStyleOn() {
		t.Error("Weight should not have style on")
	}
	if !FontSynthesisStyle.IsStyleOn() {
		t.Error("Style should have style on")
	}
	if !FontSynthesisAll.IsStyleOn() {
		t.Error("All should have style on")
	}
}

func TestFontSynthesis_String(t *testing.T) {
	tests := []struct {
		fs       *FontSynthesis
		expected string
	}{
		{FontSynthesisNone, "None"},
		{FontSynthesisWeight, "Weight"},
		{FontSynthesisStyle, "Style"},
		{FontSynthesisAll, "All"},
	}
	for _, tt := range tests {
		if StringFontSynthesis(tt.fs) != tt.expected {
			t.Errorf("Expected %s, got %s", tt.expected, StringFontSynthesis(tt.fs))
		}
	}
}

func TestFontSynthesisValueOf_Valid(t *testing.T) {
	tests := []struct {
		value    int
		expected *FontSynthesis
	}{
		{0, FontSynthesisNone},
		{1, FontSynthesisWeight},
		{2, FontSynthesisStyle},
		{0xffff, FontSynthesisAll},
	}
	for _, tt := range tests {
		result, err := FontSynthesisValueOf(tt.value)
		if err != nil {
			t.Errorf("Unexpected error for value %d: %v", tt.value, err)
		}
		if !EqualFontSynthesis(result, tt.expected) {
			t.Errorf("Expected %v for value %d, got %v", tt.expected, tt.value, result)
		}
	}
}

func TestFontSynthesisValueOf_Invalid(t *testing.T) {
	invalidValues := []int{3, 4, 5, 100, -1}
	for _, value := range invalidValues {
		_, err := FontSynthesisValueOf(value)
		if err == nil {
			t.Errorf("Expected error for value %d", value)
		}
	}
}

func TestFontSynthesis_Equals(t *testing.T) {
	if !EqualFontSynthesis(FontSynthesisNone, FontSynthesisNone) {
		t.Error("None should equal None")
	}
	if EqualFontSynthesis(FontSynthesisNone, FontSynthesisWeight) {
		t.Error("None should not equal Weight")
	}
}
