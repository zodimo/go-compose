package graphics

import (
	"testing"

	"github.com/zodimo/go-compose/compose/ui/geometry"
)

func TestShadow_SentinelPattern(t *testing.T) {
	// Test EmptyShadow singleton
	if ShadowUnspecified == nil {
		t.Error("EmptyShadow should not be nil")
	}
	if IsSpecifiedShadow(ShadowUnspecified) {
		t.Error("EmptyShadow should not be specified")
	}

	// Test TakeOrElse
	t.Run("TakeOrElse", func(t *testing.T) {
		s1 := NewShadow(ColorBlack, geometry.NewOffset(1, 1), 10)
		defaultShadow := NewShadow(ColorRed, geometry.NewOffset(2, 2), 5)

		// Case 1: s is specified
		res := TakeOrElseShadow(s1, defaultShadow)
		if res != s1 {
			t.Error("TakeOrElse should return s when s is specified")
		}

		// Case 2: s is nil
		res = TakeOrElseShadow(nil, defaultShadow)
		if res != defaultShadow {
			t.Error("TakeOrElse should return default when s is nil")
		}

		// Case 3: s is EmptyShadow
		res = TakeOrElseShadow(ShadowUnspecified, defaultShadow)
		if res != defaultShadow {
			t.Error("TakeOrElse should return default when s is EmptyShadow")
		}
	})

	// Test String()
	t.Run("String", func(t *testing.T) {
		if StringShadow(ShadowUnspecified) != "EmptyShadow" {
			t.Errorf("Expected 'EmptyShadow', got '%s'", StringShadow(ShadowUnspecified))
		}
	})
}

func TestShadow_Copy(t *testing.T) {
	initialColor := ColorRed
	initialOffset := geometry.NewOffset(10, 20)
	initialBlur := float32(5.0)

	original := NewShadow(initialColor, initialOffset, initialBlur)

	t.Run("Identity Copy", func(t *testing.T) {
		copy := original.Copy()
		if !EqualShadow(copy, original) {
			t.Errorf("Identity copy failed. Expected %v, got %v", original, copy)
		}
	})

	t.Run("Copy WithColor", func(t *testing.T) {
		newColor := ColorBlue
		copy := original.Copy(WithColor(newColor))

		expected := NewShadow(newColor, initialOffset, initialBlur)
		if !EqualShadow(copy, expected) {
			t.Errorf("Copy with Color failed. Expected %v, got %v", expected, copy)
		}
		// Verify original unchanged
		if !EqualShadow(original, NewShadow(initialColor, initialOffset, initialBlur)) {
			t.Error("Original shadow was modified")
		}
	})

	t.Run("Copy WithOffset", func(t *testing.T) {
		newOffset := geometry.NewOffset(30, 40)
		copy := original.Copy(WithOffset(newOffset))

		expected := NewShadow(initialColor, newOffset, initialBlur)
		if !EqualShadow(copy, expected) {
			t.Errorf("Copy with Offset failed. Expected %v, got %v", expected, copy)
		}
	})

	t.Run("Copy WithBlurRadius", func(t *testing.T) {
		newBlur := float32(8.0)
		copy := original.Copy(WithBlurRadius(newBlur))

		expected := NewShadow(initialColor, initialOffset, newBlur)
		if !EqualShadow(copy, expected) {
			t.Errorf("Copy with BlurRadius failed. Expected %v, got %v", expected, copy)
		}
	})

	t.Run("Copy All", func(t *testing.T) {
		newColor := ColorGreen
		newOffset := geometry.NewOffset(50, 60)
		newBlur := float32(12.0)

		copy := original.Copy(
			WithColor(newColor),
			WithOffset(newOffset),
			WithBlurRadius(newBlur),
		)

		expected := NewShadow(newColor, newOffset, newBlur)
		if !EqualShadow(copy, expected) {
			t.Errorf("Copy all failed. Expected %v, got %v", expected, copy)
		}
	})
}
