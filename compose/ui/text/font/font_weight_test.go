package font

import (
	"testing"
)

func TestNewFontWeight_Valid(t *testing.T) {
	tests := []int{1, 100, 400, 700, 1000}
	for _, weight := range tests {
		fw := NewFontWeight(weight)
		if fw.Weight() != weight {
			t.Errorf("Expected weight %d, got %d", weight, fw.Weight())
		}
	}
}

func TestNewFontWeight_InvalidLow(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for weight 0")
		}
	}()
	NewFontWeight(0)
}

func TestNewFontWeight_InvalidHigh(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for weight 1001")
		}
	}()
	NewFontWeight(1001)
}

func TestFontWeight_Compare(t *testing.T) {
	if FontWeightNormal.Compare(FontWeightBold) != -1 {
		t.Error("Normal should be less than Bold")
	}
	if FontWeightBold.Compare(FontWeightNormal) != 1 {
		t.Error("Bold should be greater than Normal")
	}
	if FontWeightNormal.Compare(FontWeightW400) != 0 {
		t.Error("Normal should equal W400")
	}
}

func TestFontWeight_Equals(t *testing.T) {
	if !FontWeightNormal.Equals(FontWeightW400) {
		t.Error("Normal should equal W400")
	}
	if FontWeightNormal.Equals(FontWeightBold) {
		t.Error("Normal should not equal Bold")
	}
}

func TestFontWeight_String(t *testing.T) {
	expected := "FontWeight(weight=400)"
	if FontWeightNormal.String() != expected {
		t.Errorf("Expected %s, got %s", expected, FontWeightNormal.String())
	}
}

func TestFontWeight_Constants(t *testing.T) {
	tests := []struct {
		fw     FontWeight
		weight int
	}{
		{FontWeightThin, 100},
		{FontWeightExtraLight, 200},
		{FontWeightLight, 300},
		{FontWeightNormal, 400},
		{FontWeightMedium, 500},
		{FontWeightSemiBold, 600},
		{FontWeightBold, 700},
		{FontWeightExtraBold, 800},
		{FontWeightBlack, 900},
	}
	for _, tt := range tests {
		if tt.fw.Weight() != tt.weight {
			t.Errorf("Expected weight %d, got %d", tt.weight, tt.fw.Weight())
		}
	}
}

func TestFontWeightValues(t *testing.T) {
	values := FontWeightValues()
	if len(values) != 9 {
		t.Errorf("Expected 9 values, got %d", len(values))
	}
	expectedWeights := []int{100, 200, 300, 400, 500, 600, 700, 800, 900}
	for i, fw := range values {
		if fw.Weight() != expectedWeights[i] {
			t.Errorf("Expected weight %d at index %d, got %d", expectedWeights[i], i, fw.Weight())
		}
	}
}

func TestLerpFontWeight(t *testing.T) {
	// Test at 0 returns start
	result := LerpFontWeight(FontWeightNormal, FontWeightBold, 0)
	if result.Weight() != 400 {
		t.Errorf("Lerp at 0 should return start, got %d", result.Weight())
	}

	// Test at 1 returns stop
	result = LerpFontWeight(FontWeightNormal, FontWeightBold, 1)
	if result.Weight() != 700 {
		t.Errorf("Lerp at 1 should return stop, got %d", result.Weight())
	}

	// Test at 0.5 returns middle
	result = LerpFontWeight(FontWeightNormal, FontWeightBold, 0.5)
	expected := 550 // (400 + 700) / 2
	if result.Weight() != expected {
		t.Errorf("Lerp at 0.5 should return %d, got %d", expected, result.Weight())
	}

	// Test clamping to max
	result = LerpFontWeight(FontWeightW900, FontWeightW900, 2.0)
	if result.Weight() > 1000 {
		t.Errorf("Lerp should clamp to 1000, got %d", result.Weight())
	}
}
