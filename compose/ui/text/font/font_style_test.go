package font

import "testing"

func TestFontStyle_Constants(t *testing.T) {
	if FontStyleNormal.Value() != 0 {
		t.Errorf("FontStyleNormal should be 0, got %d", FontStyleNormal.Value())
	}
	if FontStyleItalic.Value() != 1 {
		t.Errorf("FontStyleItalic should be 1, got %d", FontStyleItalic.Value())
	}
}

func TestFontStyle_String(t *testing.T) {
	if FontStyleNormal.String() != "Normal" {
		t.Errorf("Expected 'Normal', got %s", FontStyleNormal.String())
	}
	if FontStyleItalic.String() != "Italic" {
		t.Errorf("Expected 'Italic', got %s", FontStyleItalic.String())
	}
	invalid := FontStyle(99)
	if invalid.String() != "Invalid" {
		t.Errorf("Expected 'Invalid', got %s", invalid.String())
	}
}

func TestFontStyleValues(t *testing.T) {
	values := FontStyleValues()
	if len(values) != 2 {
		t.Errorf("Expected 2 values, got %d", len(values))
	}
	if values[0] != FontStyleNormal {
		t.Error("First value should be Normal")
	}
	if values[1] != FontStyleItalic {
		t.Error("Second value should be Italic")
	}
}
