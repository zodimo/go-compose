package input

import (
	"testing"

	"github.com/zodimo/go-compose/compose/ui/next/text"
)

func TestTextFieldCharSequence_New(t *testing.T) {
	tfcs := NewTextFieldCharSequence("Hello", text.NewTextRange(2, 4))

	if tfcs.Text() != "Hello" {
		t.Errorf("expected 'Hello', got '%s'", tfcs.Text())
	}
	if tfcs.Len() != 5 {
		t.Errorf("expected length 5, got %d", tfcs.Len())
	}
	if tfcs.Selection().Start != 2 || tfcs.Selection().End != 4 {
		t.Errorf("expected selection (2,4), got (%d,%d)", tfcs.Selection().Start, tfcs.Selection().End)
	}
	if tfcs.Composition() != nil {
		t.Error("expected nil composition")
	}
}

func TestTextFieldCharSequence_CoerceSelection(t *testing.T) {
	// Selection should be coerced to valid range
	tfcs := NewTextFieldCharSequence("Hello", text.NewTextRange(0, 100))

	sel := tfcs.Selection()
	if sel.End != 5 {
		t.Errorf("expected selection end 5, got %d", sel.End)
	}
}

func TestTextFieldCharSequence_WithComposition(t *testing.T) {
	comp := text.NewTextRange(2, 4)
	tfcs := NewTextFieldCharSequenceWithComposition("Hello", text.NewTextRange(0, 0), &comp)

	if tfcs.Composition() == nil {
		t.Fatal("expected non-nil composition")
	}
	if tfcs.Composition().Start != 2 || tfcs.Composition().End != 4 {
		t.Errorf("expected composition (2,4), got (%d,%d)", tfcs.Composition().Start, tfcs.Composition().End)
	}
}

func TestTextFieldCharSequence_GetSelectedText(t *testing.T) {
	tfcs := NewTextFieldCharSequence("Hello World", text.NewTextRange(6, 11))
	selected := tfcs.GetSelectedText()

	if selected != "World" {
		t.Errorf("expected 'World', got '%s'", selected)
	}
}

func TestTextFieldCharSequence_GetTextBeforeSelection(t *testing.T) {
	tfcs := NewTextFieldCharSequence("Hello World", text.NewTextRange(6, 11))
	before := tfcs.GetTextBeforeSelection(3)

	if before != "lo " {
		t.Errorf("expected 'lo ', got '%s'", before)
	}
}

func TestTextFieldCharSequence_GetTextAfterSelection(t *testing.T) {
	tfcs := NewTextFieldCharSequence("Hello World", text.NewTextRange(0, 5))
	after := tfcs.GetTextAfterSelection(3)

	if after != " Wo" {
		t.Errorf("expected ' Wo', got '%s'", after)
	}
}

func TestTextFieldCharSequence_ContentEquals(t *testing.T) {
	tfcs := NewTextFieldCharSequence("Hello", text.NewTextRange(0, 0))

	if !tfcs.ContentEquals("Hello") {
		t.Error("expected ContentEquals true for same text")
	}
	if tfcs.ContentEquals("World") {
		t.Error("expected ContentEquals false for different text")
	}
}

func TestTextFieldCharSequence_Equals(t *testing.T) {
	tfcs1 := NewTextFieldCharSequence("Hello", text.NewTextRange(0, 0))
	tfcs2 := NewTextFieldCharSequence("Hello", text.NewTextRange(0, 0))
	tfcs3 := NewTextFieldCharSequence("Hello", text.NewTextRange(1, 1))
	tfcs4 := NewTextFieldCharSequence("World", text.NewTextRange(0, 0))

	if !tfcs1.Equals(tfcs2) {
		t.Error("expected equal for same text and selection")
	}
	if tfcs1.Equals(tfcs3) {
		t.Error("expected not equal for different selection")
	}
	if tfcs1.Equals(tfcs4) {
		t.Error("expected not equal for different text")
	}
}

func TestTextFieldCharSequence_ShouldShowSelection(t *testing.T) {
	tfcs := NewTextFieldCharSequence("Hello", text.NewTextRange(0, 0))
	if !tfcs.ShouldShowSelection() {
		t.Error("expected ShouldShowSelection true when no highlight")
	}

	highlight := &TextHighlight{
		Type:  TextHighlightTypeHandwritingSelectPreview,
		Range: text.NewTextRange(0, 5),
	}
	tfcsWithHighlight := NewTextFieldCharSequenceFull("Hello", text.NewTextRange(0, 0), nil, highlight, nil, nil)
	if tfcsWithHighlight.ShouldShowSelection() {
		t.Error("expected ShouldShowSelection false when highlight exists")
	}
}
