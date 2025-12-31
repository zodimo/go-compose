package font

import "testing"

func TestFontLoadingStrategy_Constants(t *testing.T) {
	if FontLoadingStrategyBlocking.Value() != 0 {
		t.Errorf("Blocking should be 0, got %d", FontLoadingStrategyBlocking.Value())
	}
	if FontLoadingStrategyOptionalLocal.Value() != 1 {
		t.Errorf("OptionalLocal should be 1, got %d", FontLoadingStrategyOptionalLocal.Value())
	}
	if FontLoadingStrategyAsync.Value() != 2 {
		t.Errorf("Async should be 2, got %d", FontLoadingStrategyAsync.Value())
	}
}

func TestFontLoadingStrategy_String(t *testing.T) {
	tests := []struct {
		strategy FontLoadingStrategy
		expected string
	}{
		{FontLoadingStrategyBlocking, "Blocking"},
		{FontLoadingStrategyOptionalLocal, "OptionalLocal"},
		{FontLoadingStrategyAsync, "Async"},
	}
	for _, tt := range tests {
		if tt.strategy.String() != tt.expected {
			t.Errorf("Expected %s, got %s", tt.expected, tt.strategy.String())
		}
	}
}

func TestFontLoadingStrategy_String_Invalid(t *testing.T) {
	invalid := FontLoadingStrategy(99)
	expected := "Invalid(value=99)"
	if invalid.String() != expected {
		t.Errorf("Expected %s, got %s", expected, invalid.String())
	}
}
