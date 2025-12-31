package undo

import (
	"testing"
	"time"

	"github.com/zodimo/go-compose/compose/ui/next/text"
)

func TestTextUndoOperation_TextEditType(t *testing.T) {
	tests := []struct {
		name     string
		preText  string
		postText string
		want     TextEditType
	}{
		{"insert", "", "abc", TextEditTypeInsert},
		{"delete", "abc", "", TextEditTypeDelete},
		{"replace", "abc", "xyz", TextEditTypeReplace},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := NewTextUndoOperation(0, tt.preText, tt.postText, text.NewTextRange(0, 0), text.NewTextRange(0, 0))
			got := op.TextEditType()
			if got != tt.want {
				t.Errorf("TextEditType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTextUndoOperation_DeleteType(t *testing.T) {
	tests := []struct {
		name          string
		index         int
		preText       string
		preSelection  text.TextRange
		postSelection text.TextRange
		want          TextDeleteType
	}{
		{
			name:          "backspace",
			index:         4,
			preText:       "x",
			preSelection:  text.NewTextRange(5, 5),
			postSelection: text.NewTextRange(4, 4),
			want:          TextDeleteTypeStart,
		},
		{
			name:          "delete key",
			index:         5,
			preText:       "x",
			preSelection:  text.NewTextRange(5, 5),
			postSelection: text.NewTextRange(5, 5),
			want:          TextDeleteTypeEnd,
		},
		{
			name:          "selection delete",
			index:         2,
			preText:       "abc",
			preSelection:  text.NewTextRange(2, 5),
			postSelection: text.NewTextRange(2, 2),
			want:          TextDeleteTypeInner,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := &TextUndoOperation{
				Index:         tt.index,
				PreText:       tt.preText,
				PostText:      "",
				PreSelection:  tt.preSelection,
				PostSelection: tt.postSelection,
			}
			got := op.DeleteType()
			if got != tt.want {
				t.Errorf("DeleteType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTextUndoOperation_Merge_Insert(t *testing.T) {
	now := time.Now().UnixMilli()

	op1 := NewTextUndoOperationWithOptions(0, "", "a", text.NewTextRange(0, 0), text.NewTextRange(1, 1), now, true)
	op2 := NewTextUndoOperationWithOptions(1, "", "b", text.NewTextRange(1, 1), text.NewTextRange(2, 2), now+100, true)

	merged := op1.Merge(op2)
	if merged == nil {
		t.Fatal("expected merge to succeed")
	}

	if merged.PostText != "ab" {
		t.Errorf("expected merged PostText 'ab', got '%s'", merged.PostText)
	}
	if merged.Index != 0 {
		t.Errorf("expected merged Index 0, got %d", merged.Index)
	}
}

func TestTextUndoOperation_Merge_TimeWindow(t *testing.T) {
	now := time.Now().UnixMilli()

	op1 := NewTextUndoOperationWithOptions(0, "", "a", text.NewTextRange(0, 0), text.NewTextRange(1, 1), now, true)
	op2 := NewTextUndoOperationWithOptions(1, "", "b", text.NewTextRange(1, 1), text.NewTextRange(2, 2), now+3000, true) // 3 seconds later

	merged := op1.Merge(op2)
	if merged != nil {
		t.Error("expected merge to fail due to time window")
	}
}

func TestTextUndoOperation_Merge_Newline(t *testing.T) {
	now := time.Now().UnixMilli()

	op1 := NewTextUndoOperationWithOptions(0, "", "a", text.NewTextRange(0, 0), text.NewTextRange(1, 1), now, true)
	op2 := NewTextUndoOperationWithOptions(1, "", "\n", text.NewTextRange(1, 1), text.NewTextRange(2, 2), now+100, true)

	merged := op1.Merge(op2)
	if merged != nil {
		t.Error("expected merge to fail for newline")
	}
}

func TestTextUndoOperation_Merge_CanMergeFalse(t *testing.T) {
	now := time.Now().UnixMilli()

	op1 := NewTextUndoOperationWithOptions(0, "", "a", text.NewTextRange(0, 0), text.NewTextRange(1, 1), now, false)
	op2 := NewTextUndoOperationWithOptions(1, "", "b", text.NewTextRange(1, 1), text.NewTextRange(2, 2), now+100, true)

	merged := op1.Merge(op2)
	if merged != nil {
		t.Error("expected merge to fail when first op has CanMerge=false")
	}
}

func TestTextUndoOperation_Merge_DifferentTypes(t *testing.T) {
	now := time.Now().UnixMilli()

	op1 := NewTextUndoOperationWithOptions(0, "", "a", text.NewTextRange(0, 0), text.NewTextRange(1, 1), now, true)     // Insert
	op2 := NewTextUndoOperationWithOptions(0, "a", "", text.NewTextRange(1, 1), text.NewTextRange(0, 0), now+100, true) // Delete

	merged := op1.Merge(op2)
	if merged != nil {
		t.Error("expected merge to fail for different edit types")
	}
}
