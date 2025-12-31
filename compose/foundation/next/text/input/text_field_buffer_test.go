package input

import (
	"testing"

	"github.com/zodimo/go-compose/compose/ui/next/text"
)

func TestTextFieldBuffer_New(t *testing.T) {
	tfcs := NewTextFieldCharSequence("Hello", text.NewTextRange(5, 5))
	buffer := NewTextFieldBuffer(tfcs)

	if buffer.String() != "Hello" {
		t.Errorf("expected 'Hello', got '%s'", buffer.String())
	}
	if buffer.Length() != 5 {
		t.Errorf("expected length 5, got %d", buffer.Length())
	}
}

func TestTextFieldBuffer_Replace(t *testing.T) {
	tfcs := NewTextFieldCharSequence("Hello", text.NewTextRange(5, 5))
	buffer := NewTextFieldBuffer(tfcs)

	buffer.Replace(5, 5, " World")
	if buffer.String() != "Hello World" {
		t.Errorf("expected 'Hello World', got '%s'", buffer.String())
	}
}

func TestTextFieldBuffer_Append(t *testing.T) {
	tfcs := NewTextFieldCharSequence("Hello", text.NewTextRange(5, 5))
	buffer := NewTextFieldBuffer(tfcs)

	buffer.Append(" World")
	if buffer.String() != "Hello World" {
		t.Errorf("expected 'Hello World', got '%s'", buffer.String())
	}
}

func TestTextFieldBuffer_Insert(t *testing.T) {
	tfcs := NewTextFieldCharSequence("Hllo", text.NewTextRange(1, 1))
	buffer := NewTextFieldBuffer(tfcs)

	buffer.Insert(1, "e")
	if buffer.String() != "Hello" {
		t.Errorf("expected 'Hello', got '%s'", buffer.String())
	}
}

func TestTextFieldBuffer_Delete(t *testing.T) {
	tfcs := NewTextFieldCharSequence("Helllo", text.NewTextRange(4, 4))
	buffer := NewTextFieldBuffer(tfcs)

	buffer.Delete(3, 4)
	if buffer.String() != "Hello" {
		t.Errorf("expected 'Hello', got '%s'", buffer.String())
	}
}

func TestTextFieldBuffer_SelectionAdjustment(t *testing.T) {
	// Selection after insertion point should shift
	tfcs := NewTextFieldCharSequence("Hello", text.NewTextRange(5, 5))
	buffer := NewTextFieldBuffer(tfcs)

	buffer.Insert(0, "Say ")
	sel := buffer.Selection()
	if sel.Start != 9 || sel.End != 9 {
		t.Errorf("expected selection (9,9), got (%d,%d)", sel.Start, sel.End)
	}
}

func TestTextFieldBuffer_PlaceCursorAtEnd(t *testing.T) {
	tfcs := NewTextFieldCharSequence("Hello", text.NewTextRange(0, 0))
	buffer := NewTextFieldBuffer(tfcs)

	buffer.PlaceCursorAtEnd()
	sel := buffer.Selection()
	if sel.Start != 5 || sel.End != 5 {
		t.Errorf("expected selection (5,5), got (%d,%d)", sel.Start, sel.End)
	}
}

func TestTextFieldBuffer_SelectAll(t *testing.T) {
	tfcs := NewTextFieldCharSequence("Hello", text.NewTextRange(0, 0))
	buffer := NewTextFieldBuffer(tfcs)

	buffer.SelectAll()
	sel := buffer.Selection()
	if sel.Start != 0 || sel.End != 5 {
		t.Errorf("expected selection (0,5), got (%d,%d)", sel.Start, sel.End)
	}
}

func TestTextFieldBuffer_DeleteSelectedText(t *testing.T) {
	tfcs := NewTextFieldCharSequence("Hello World", text.NewTextRange(5, 11))
	buffer := NewTextFieldBuffer(tfcs)

	deleted := buffer.DeleteSelectedText()
	if !deleted {
		t.Error("expected DeleteSelectedText to return true")
	}
	if buffer.String() != "Hello" {
		t.Errorf("expected 'Hello', got '%s'", buffer.String())
	}
}

func TestTextFieldBuffer_DeleteSelectedText_NoSelection(t *testing.T) {
	tfcs := NewTextFieldCharSequence("Hello", text.NewTextRange(2, 2))
	buffer := NewTextFieldBuffer(tfcs)

	deleted := buffer.DeleteSelectedText()
	if deleted {
		t.Error("expected DeleteSelectedText to return false for collapsed selection")
	}
	if buffer.String() != "Hello" {
		t.Errorf("expected unchanged 'Hello', got '%s'", buffer.String())
	}
}

func TestTextFieldBuffer_RevertAllChanges(t *testing.T) {
	tfcs := NewTextFieldCharSequence("Hello", text.NewTextRange(5, 5))
	buffer := NewTextFieldBuffer(tfcs)

	buffer.Append(" World")
	buffer.Replace(0, 5, "Goodbye")
	buffer.RevertAllChanges()

	if buffer.String() != "Hello" {
		t.Errorf("expected 'Hello', got '%s'", buffer.String())
	}
	sel := buffer.Selection()
	if sel.Start != 5 || sel.End != 5 {
		t.Errorf("expected selection (5,5), got (%d,%d)", sel.Start, sel.End)
	}
}

func TestTextFieldBuffer_ChangesTracked(t *testing.T) {
	tfcs := NewTextFieldCharSequence("Hello", text.NewTextRange(5, 5))
	buffer := NewTextFieldBuffer(tfcs)

	buffer.Append(" World")

	changes := buffer.Changes()
	if changes.ChangeCount() == 0 {
		t.Error("expected changes to be tracked")
	}
}

func TestTextFieldBuffer_HasSelection(t *testing.T) {
	tests := []struct {
		selection text.TextRange
		want      bool
	}{
		{text.NewTextRange(0, 0), false},
		{text.NewTextRange(5, 5), false},
		{text.NewTextRange(0, 5), true},
		{text.NewTextRange(2, 4), true},
	}

	for _, tt := range tests {
		tfcs := NewTextFieldCharSequence("Hello", tt.selection)
		buffer := NewTextFieldBuffer(tfcs)

		got := buffer.HasSelection()
		if got != tt.want {
			t.Errorf("HasSelection() for selection %v = %v, want %v", tt.selection, got, tt.want)
		}
	}
}

func TestTextFieldBuffer_Clear(t *testing.T) {
	tfcs := NewTextFieldCharSequence("Hello World", text.NewTextRange(0, 0))
	buffer := NewTextFieldBuffer(tfcs)

	buffer.Clear()
	if buffer.String() != "" {
		t.Errorf("expected empty string, got '%s'", buffer.String())
	}
	if buffer.Length() != 0 {
		t.Errorf("expected length 0, got %d", buffer.Length())
	}
}

func TestTextFieldBuffer_ReplaceAll(t *testing.T) {
	tfcs := NewTextFieldCharSequence("Hello Hello Hello", text.NewTextRange(0, 0))
	buffer := NewTextFieldBuffer(tfcs)

	count := buffer.ReplaceAll("Hello", "Hi")
	if count != 3 {
		t.Errorf("expected 3 replacements, got %d", count)
	}
	if buffer.String() != "Hi Hi Hi" {
		t.Errorf("expected 'Hi Hi Hi', got '%s'", buffer.String())
	}
}
