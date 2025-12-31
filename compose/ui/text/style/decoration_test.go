package style

import (
	"testing"
)

func TestTextDecoration_Combine(t *testing.T) {
	combined := Combine([]*TextDecoration{TextDecorationUnderline, TextDecorationLineThrough})
	if !combined.Contains(TextDecorationUnderline) {
		t.Errorf("Expected combined decoration to contain Underline")
	}
	if !combined.Contains(TextDecorationLineThrough) {
		t.Errorf("Expected combined decoration to contain LineThrough")
	}
	if combined.mask != (TextDecorationMaskUnderline | TextDecorationMaskLineThrough) {
		t.Errorf("Expected mask to be %d, got %d", (TextDecorationMaskUnderline | TextDecorationMaskLineThrough), combined.mask)
	}
}

func TestTextDecoration_Plus(t *testing.T) {
	plus := TextDecorationUnderline.Plus(TextDecorationLineThrough)
	if !plus.Contains(TextDecorationUnderline) {
		t.Errorf("Expected plus decoration to contain Underline")
	}
	if !plus.Contains(TextDecorationLineThrough) {
		t.Errorf("Expected plus decoration to contain LineThrough")
	}
}

func TestTextDecoration_String(t *testing.T) {
	if StringTextDecoration(TextDecorationNone) != "TextDecoration.None" {
		t.Errorf("Expected TextDecoration.None, got %s", StringTextDecoration(TextDecorationNone))
	}
	if StringTextDecoration(TextDecorationUnderline) != "TextDecoration.Underline" {
		t.Errorf("Expected TextDecoration.Underline, got %s", StringTextDecoration(TextDecorationUnderline))
	}
	combined := TextDecorationUnderline.Plus(TextDecorationLineThrough)
	expected := "TextDecoration[Underline, LineThrough]"
	// Order depends on implementation, but here we append in specific order
	if StringTextDecoration(combined) != expected {
		t.Errorf("Expected %s, got %s", expected, StringTextDecoration(combined))
	}
}

func TestNewTextDecoration(t *testing.T) {
	// Valid masks
	NewTextDecoration(TextDecorationMaskNone)
	NewTextDecoration(TextDecorationMaskUnderline)
	// Invalid mask
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	NewTextDecoration(0x04) // Should panic 0b100
}
