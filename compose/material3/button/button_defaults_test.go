package button_test

import (
	"testing"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/material3/button"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

// Mock Composer to satisfy interface
type mockComposer struct {
	compose.Composer
}

// Basic checks for ButtonDefaults values
func TestButtonDefaults_ContentPadding(t *testing.T) {
	padding := button.ButtonDefaults.ContentPadding()
	expectedHorizontal := unit.Dp(24)
	expectedVertical := unit.Dp(8)

	if padding.Start != expectedHorizontal {
		t.Errorf("Expected Start padding %v, got %v", expectedHorizontal, padding.Start)
	}
	if padding.End != expectedHorizontal {
		t.Errorf("Expected End padding %v, got %v", expectedHorizontal, padding.End)
	}
	if padding.Top != expectedVertical {
		t.Errorf("Expected Top padding %v, got %v", expectedVertical, padding.Top)
	}
	if padding.Bottom != expectedVertical {
		t.Errorf("Expected Bottom padding %v, got %v", expectedVertical, padding.Bottom)
	}
}

// Note: Testing methods requiring Composer (OutlinedButtonColors, etc.) is complex without a full Composition setup.
// For unit testing defaults, we primarily verify static values or simple logic where possible.
// Since OutlinedButtonColors relies on material3.Theme(c).ColorScheme(), integration tests or a mock Composer with CompositionLocals are needed.
// For now, we skip those in this basic unit test and rely on the fact that we're delegating to the Theme system.
