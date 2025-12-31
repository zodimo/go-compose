package internal

import (
	"testing"

	"github.com/zodimo/go-compose/compose/ui/next/text"
)

func TestChangeTracker_Empty(t *testing.T) {
	ct := NewChangeTracker(nil)
	if ct.ChangeCount() != 0 {
		t.Errorf("expected 0 changes, got %d", ct.ChangeCount())
	}
}

func TestChangeTracker_SingleChange(t *testing.T) {
	ct := NewChangeTracker(nil)
	ct.TrackChange(0, 0, 5) // Insert 5 chars at position 0

	if ct.ChangeCount() != 1 {
		t.Errorf("expected 1 change, got %d", ct.ChangeCount())
	}

	origRange := ct.GetOriginalRange(0)
	if origRange.Start != 0 || origRange.End != 0 {
		t.Errorf("expected original range (0,0), got (%d,%d)", origRange.Start, origRange.End)
	}

	curRange := ct.GetRange(0)
	if curRange.Start != 0 || curRange.End != 5 {
		t.Errorf("expected current range (0,5), got (%d,%d)", curRange.Start, curRange.End)
	}
}

func TestChangeTracker_MultipleNonOverlapping(t *testing.T) {
	ct := NewChangeTracker(nil)
	ct.TrackChange(0, 0, 5)   // Insert "Hello" at start
	ct.TrackChange(10, 10, 5) // Insert " more" at position 10

	if ct.ChangeCount() != 2 {
		t.Errorf("expected 2 changes, got %d", ct.ChangeCount())
	}
}

func TestChangeTracker_AdjacentMerge(t *testing.T) {
	ct := NewChangeTracker(nil)
	ct.TrackChange(0, 0, 1) // Insert 'H'
	ct.TrackChange(1, 1, 1) // Insert 'e'
	ct.TrackChange(2, 2, 1) // Insert 'l'

	// These should be merged into a single change
	if ct.ChangeCount() != 1 {
		t.Errorf("expected 1 merged change, got %d", ct.ChangeCount())
	}

	curRange := ct.GetRange(0)
	if curRange.Start != 0 || curRange.End != 3 {
		t.Errorf("expected current range (0,3), got (%d,%d)", curRange.Start, curRange.End)
	}
}

func TestChangeTracker_ClearChanges(t *testing.T) {
	ct := NewChangeTracker(nil)
	ct.TrackChange(0, 0, 5)
	ct.ClearChanges()

	if ct.ChangeCount() != 0 {
		t.Errorf("expected 0 changes after clear, got %d", ct.ChangeCount())
	}
}

func TestChangeTracker_ForEachChange(t *testing.T) {
	ct := NewChangeTracker(nil)
	ct.TrackChange(0, 0, 5)
	ct.TrackChange(10, 15, 3)

	var ranges []text.TextRange
	ct.ForEachChange(func(originalRange, currentRange text.TextRange) {
		ranges = append(ranges, currentRange)
	})

	if len(ranges) != 2 {
		t.Errorf("expected 2 changes in ForEach, got %d", len(ranges))
	}
}

func TestChangeTracker_NoOpChange(t *testing.T) {
	ct := NewChangeTracker(nil)
	ct.TrackChange(5, 5, 0) // No-op: insert nothing

	if ct.ChangeCount() != 0 {
		t.Errorf("expected 0 changes for no-op, got %d", ct.ChangeCount())
	}
}
