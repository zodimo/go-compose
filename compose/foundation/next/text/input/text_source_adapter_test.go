package input

import (
	"testing"
)

func TestTextSourceAdapter_ReadOnly(t *testing.T) {
	adapter := NewTextSourceAdapterFromString("hello world")

	if adapter.IsEditable() {
		t.Error("should not be editable")
	}
	if adapter.Text() != "hello world" {
		t.Errorf("expected 'hello world', got %q", adapter.Text())
	}
	if adapter.Size() != 11 {
		t.Errorf("expected size 11, got %d", adapter.Size())
	}
}

func TestTextSourceAdapter_ReadAt(t *testing.T) {
	adapter := NewTextSourceAdapterFromString("hello world")

	buf := make([]byte, 5)
	n, err := adapter.ReadAt(buf, 0)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if n != 5 {
		t.Errorf("expected 5, got %d", n)
	}
	if string(buf) != "hello" {
		t.Errorf("expected 'hello', got %q", string(buf))
	}

	n, err = adapter.ReadAt(buf, 6)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(buf) != "world" {
		t.Errorf("expected 'world', got %q", string(buf))
	}
}

func TestTextSourceAdapter_Changed(t *testing.T) {
	adapter := NewTextSourceAdapterFromString("hello")

	// First call should return false (no change yet)
	if adapter.Changed() {
		t.Error("should not be changed initially after construction")
	}

	// After SetText, should be changed
	adapter.SetText("world")
	if !adapter.Changed() {
		t.Error("should be changed after SetText")
	}

	// After reading changed, should not be changed
	if adapter.Changed() {
		t.Error("should not be changed after reading Changed()")
	}
}

func TestTextSourceAdapter_ReplaceRunes_ReadOnly(t *testing.T) {
	adapter := NewTextSourceAdapterFromString("hello")

	// ReplaceRunes should be a no-op for read-only adapters
	adapter.ReplaceRunes(0, 1, "X")

	if adapter.Text() != "hello" {
		t.Errorf("text should be unchanged for read-only adapter, got %q", adapter.Text())
	}
}

func TestTextSourceAdapter_FromState(t *testing.T) {
	state := NewTextFieldState("hello world")
	adapter := NewTextSourceAdapterFromState(state)

	if !adapter.IsEditable() {
		t.Error("should be editable")
	}
	if adapter.State() != state {
		t.Error("State() should return the original state")
	}
	if adapter.Text() != "hello world" {
		t.Errorf("expected 'hello world', got %q", adapter.Text())
	}
}

func TestTextSourceAdapter_ReplaceRunes_Editable(t *testing.T) {
	state := NewTextFieldState("hello")
	adapter := NewTextSourceAdapterFromState(state)

	// Replace "h" with "X" (first byte/rune)
	adapter.ReplaceRunes(0, 1, "X")

	if adapter.Text() != "Xello" {
		t.Errorf("expected 'Xello', got %q", adapter.Text())
	}

	if !adapter.Changed() {
		t.Error("should be changed after ReplaceRunes")
	}
}

func TestTextSourceAdapter_SetText(t *testing.T) {
	adapter := NewTextSourceAdapterFromString("hello")

	adapter.SetText("world")

	if adapter.Text() != "world" {
		t.Errorf("expected 'world', got %q", adapter.Text())
	}
	if adapter.Size() != 5 {
		t.Errorf("expected size 5, got %d", adapter.Size())
	}
}
