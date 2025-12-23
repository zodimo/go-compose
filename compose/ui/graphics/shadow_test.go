package graphics

import (
	"testing"

	"github.com/zodimo/go-compose/compose/ui/geometry"
)

func TestShadow_SentinelPattern(t *testing.T) {
	// Test EmptyShadow singleton
	if EmptyShadow == nil {
		t.Error("EmptyShadow should not be nil")
	}
	if EmptyShadow.IsSpecified() {
		t.Error("EmptyShadow should not be specified")
	}

	// Test TakeOrElse
	t.Run("TakeOrElse", func(t *testing.T) {
		s1 := NewShadow(ColorBlack, geometry.NewOffset(1, 1), 10)
		s1Ptr := &s1
		defaultShadow := NewShadow(ColorRed, geometry.NewOffset(2, 2), 5)
		defaultPtr := &defaultShadow

		// Case 1: s is specified
		res := TakeOrElse(s1Ptr, defaultPtr)
		if res != s1Ptr {
			t.Error("TakeOrElse should return s when s is specified")
		}

		// Case 2: s is nil
		res = TakeOrElse(nil, defaultPtr)
		if res != defaultPtr {
			t.Error("TakeOrElse should return default when s is nil")
		}

		// Case 3: s is EmptyShadow
		res = TakeOrElse(EmptyShadow, defaultPtr)
		if res != defaultPtr {
			t.Error("TakeOrElse should return default when s is EmptyShadow")
		}
	})

	// Test IsSpecified
	t.Run("IsSpecified", func(t *testing.T) {
		s1 := NewShadow(ColorBlack, geometry.NewOffset(1, 1), 10)
		if !s1.IsSpecified() {
			t.Error("Regular shadow should be specified")
		}

		// Manually construct unspecified shadow (should match EmptyShadow content)
		// We can't use ShadowUnspecified directly anymore as it was removed,
		// but we can copy EmptyShadow which serves the same purpose of holding unspecified values.
		sUnspec := *EmptyShadow
		if sUnspec.IsSpecified() {
			t.Error("Shadow with unspecified fields should not be specified")
		}
	})

	// Test String()
	t.Run("String", func(t *testing.T) {
		if EmptyShadow.String() != "EmptyShadow" {
			t.Errorf("Expected 'EmptyShadow', got '%s'", EmptyShadow.String())
		}
	})
}
