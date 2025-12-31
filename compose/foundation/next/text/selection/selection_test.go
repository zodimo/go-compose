package selection

import (
	"testing"

	"github.com/zodimo/go-compose/compose/ui/next/text/style"
)

func TestAnchorInfo_Creation(t *testing.T) {
	anchor := AnchorInfo{
		Direction:    style.ResolvedTextDirectionLtr,
		Offset:       10,
		SelectableId: 42,
	}

	if anchor.Direction != style.ResolvedTextDirectionLtr {
		t.Errorf("Expected Direction to be TextDirectionLtr")
	}
	if anchor.Offset != 10 {
		t.Errorf("Expected Offset to be 10, got %d", anchor.Offset)
	}
	if anchor.SelectableId != 42 {
		t.Errorf("Expected SelectableId to be 42, got %d", anchor.SelectableId)
	}
}

func TestSelection_NewSelection(t *testing.T) {
	start := AnchorInfo{Offset: 5, SelectableId: 1}
	end := AnchorInfo{Offset: 10, SelectableId: 1}

	sel := NewSelection(start, end, false)

	if sel.Start != start {
		t.Error("Expected Start to match")
	}
	if sel.End != end {
		t.Error("Expected End to match")
	}
	if sel.HandlesCrossed {
		t.Error("Expected HandlesCrossed to be false")
	}
}

func TestSelection_IsCollapsed(t *testing.T) {
	// Collapsed selection - same offset and same selectable
	start := AnchorInfo{Offset: 5, SelectableId: 1}
	end := AnchorInfo{Offset: 5, SelectableId: 1}
	sel := NewSelection(start, end, false)

	if !sel.IsCollapsed() {
		t.Error("Expected selection to be collapsed")
	}

	// Non-collapsed - different offsets
	end2 := AnchorInfo{Offset: 10, SelectableId: 1}
	sel2 := NewSelection(start, end2, false)

	if sel2.IsCollapsed() {
		t.Error("Expected selection to not be collapsed")
	}

	// Non-collapsed - different selectables
	end3 := AnchorInfo{Offset: 5, SelectableId: 2}
	sel3 := NewSelection(start, end3, false)

	if sel3.IsCollapsed() {
		t.Error("Expected selection with different selectables to not be collapsed")
	}
}

func TestSelection_Copy(t *testing.T) {
	start := AnchorInfo{Offset: 5, SelectableId: 1}
	end := AnchorInfo{Offset: 10, SelectableId: 2}
	sel := NewSelection(start, end, true)

	copied := sel.Copy()

	if copied.Start != sel.Start {
		t.Error("Expected copied Start to match")
	}
	if copied.End != sel.End {
		t.Error("Expected copied End to match")
	}
	if copied.HandlesCrossed != sel.HandlesCrossed {
		t.Error("Expected copied HandlesCrossed to match")
	}
}

func TestSelection_ToTextRange(t *testing.T) {
	start := AnchorInfo{Offset: 5}
	end := AnchorInfo{Offset: 10}
	sel := NewSelection(start, end, false)

	textRange := sel.ToTextRange()

	if textRange.Start != 5 {
		t.Errorf("Expected Start to be 5, got %d", textRange.Start)
	}
	if textRange.End != 10 {
		t.Errorf("Expected End to be 10, got %d", textRange.End)
	}
}

func TestSelection_Merge_NilOther(t *testing.T) {
	start := AnchorInfo{Offset: 5, SelectableId: 1}
	end := AnchorInfo{Offset: 10, SelectableId: 1}
	sel := NewSelection(start, end, false)

	result := sel.Merge(nil)

	if result != sel {
		t.Error("Expected Merge with nil to return original selection")
	}
}

func TestSelection_Merge_NoHandlesCrossed(t *testing.T) {
	// First selection: 5-10
	start1 := AnchorInfo{Offset: 5, SelectableId: 1}
	end1 := AnchorInfo{Offset: 10, SelectableId: 1}
	sel1 := NewSelection(start1, end1, false)

	// Second selection: 15-20
	start2 := AnchorInfo{Offset: 15, SelectableId: 2}
	end2 := AnchorInfo{Offset: 20, SelectableId: 2}
	sel2 := NewSelection(start2, end2, false)

	result := sel1.Merge(&sel2)

	// When neither is crossed: start from sel1, end from sel2
	if result.Start != start1 {
		t.Error("Expected merged Start to be from first selection")
	}
	if result.End != end2 {
		t.Error("Expected merged End to be from second selection")
	}
	if result.HandlesCrossed {
		t.Error("Expected merged HandlesCrossed to be false")
	}
}

func TestSelection_Merge_WithHandlesCrossed(t *testing.T) {
	// First selection with handles crossed
	start1 := AnchorInfo{Offset: 10, SelectableId: 1}
	end1 := AnchorInfo{Offset: 5, SelectableId: 1}
	sel1 := NewSelection(start1, end1, true)

	// Second selection not crossed
	start2 := AnchorInfo{Offset: 15, SelectableId: 2}
	end2 := AnchorInfo{Offset: 20, SelectableId: 2}
	sel2 := NewSelection(start2, end2, false)

	result := sel1.Merge(&sel2)

	// When sel1.HandlesCrossed: newEnd = sel1.End
	// When !other.HandlesCrossed: newStart = other.End
	if result.Start != end2 {
		t.Error("Expected merged Start to be other.End when other is not crossed")
	}
	if result.End != end1 {
		t.Error("Expected merged End to be sel1.End when sel1 is crossed")
	}
	if !result.HandlesCrossed {
		t.Error("Expected merged HandlesCrossed to be true")
	}
}
