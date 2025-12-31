package internal

import (
	"github.com/zodimo/go-compose/compose/ui/next/text"
)

// ChangeList represents a list of non-overlapping, ordered text changes.
//
// This interface is implemented by ChangeTracker and used by InputTransformation
// to inspect what changes were made during an edit session.
//
// Changes are ordered by their position in the current text. Each change
// has both a "current" range (position in the modified text) and an
// "original" range (position in the text before any changes).
type ChangeList interface {
	// ChangeCount returns the number of discrete changes.
	ChangeCount() int

	// GetRange returns the range in the current (modified) text for the change at index.
	GetRange(changeIndex int) text.TextRange

	// GetOriginalRange returns the range in the original (unmodified) text that was replaced.
	GetOriginalRange(changeIndex int) text.TextRange
}

// Change represents a single tracked text change.
type Change struct {
	// OriginalStart is the start position in the original text.
	OriginalStart int
	// OriginalEnd is the end position in the original text.
	OriginalEnd int
	// CurrentStart is the start position in the current (modified) text.
	CurrentStart int
	// CurrentEnd is the end position in the current (modified) text.
	CurrentEnd int
}

// ChangeTracker implements ChangeList and records changes as they occur.
//
// It maintains a list of non-overlapping changes ordered by position.
// Adjacent or overlapping changes are merged automatically.
//
// The tracker handles the complexity of maintaining consistent coordinate
// mappings between the original text and the current text as edits accumulate.
//
// This is a port of androidx.compose.foundation.text.input.internal.ChangeTracker.
type ChangeTracker struct {
	changes []Change
}

// NewChangeTracker creates a new ChangeTracker.
// If initialChanges is provided, it copies the changes from that tracker.
func NewChangeTracker(initialChanges *ChangeTracker) *ChangeTracker {
	ct := &ChangeTracker{}
	if initialChanges != nil {
		ct.changes = make([]Change, len(initialChanges.changes))
		copy(ct.changes, initialChanges.changes)
	}
	return ct
}

// ChangeCount returns the number of tracked changes.
func (c *ChangeTracker) ChangeCount() int {
	return len(c.changes)
}

// GetRange returns the range in the current text for change at index.
func (c *ChangeTracker) GetRange(changeIndex int) text.TextRange {
	change := c.changes[changeIndex]
	return text.NewTextRange(change.CurrentStart, change.CurrentEnd)
}

// GetOriginalRange returns the range in the original text for change at index.
func (c *ChangeTracker) GetOriginalRange(changeIndex int) text.TextRange {
	change := c.changes[changeIndex]
	return text.NewTextRange(change.OriginalStart, change.OriginalEnd)
}

// TrackChange records a change at the given position.
//
// Parameters:
//   - preStart: Start position in pre-change text (before this change)
//   - preEnd: End position in pre-change text (before this change)
//   - postLength: Length of the replacement text
//
// The change is merged with existing tracking data if it overlaps or is adjacent.
func (c *ChangeTracker) TrackChange(preStart, preEnd, postLength int) {
	if preStart == preEnd && postLength == 0 {
		return // No-op change
	}

	postEnd := preStart + postLength

	if len(c.changes) == 0 {
		c.changes = append(c.changes, Change{
			OriginalStart: preStart,
			OriginalEnd:   preEnd,
			CurrentStart:  preStart,
			CurrentEnd:    postEnd,
		})
		return
	}

	// Find affected changes (those that overlap or are adjacent to the edit)
	// and merge them into a single change.
	firstAffected := -1
	lastAffected := -1

	for i, ch := range c.changes {
		// Check if this change is affected by the edit
		// A change is affected if the edit range overlaps or is adjacent to its current range
		if preEnd >= ch.CurrentStart && preStart <= ch.CurrentEnd {
			if firstAffected == -1 {
				firstAffected = i
			}
			lastAffected = i
		}
	}

	if firstAffected == -1 {
		// No overlapping changes - insert a new one at the right position
		c.insertChange(preStart, preEnd, postEnd)
		return
	}

	// Merge with affected changes
	first := c.changes[firstAffected]
	last := c.changes[lastAffected]

	// Calculate merged original range
	mergedOrigStart := first.OriginalStart
	if preStart < first.CurrentStart {
		// Edit extends before first affected change
		mergedOrigStart = preStart
	}

	mergedOrigEnd := last.OriginalEnd
	if preEnd > last.CurrentEnd {
		// Edit extends after last affected change
		expandAmount := preEnd - last.CurrentEnd
		mergedOrigEnd += expandAmount
	}

	// Calculate merged current range
	mergedCurStart := first.CurrentStart
	if preStart < first.CurrentStart {
		mergedCurStart = preStart
	}

	mergedCurEnd := postEnd
	if last.CurrentEnd > preEnd {
		// Some of the last change extends past the edit
		excess := last.CurrentEnd - preEnd
		mergedCurEnd += excess
	}

	// Calculate the delta (how much the text length changed)
	delta := postLength - (preEnd - preStart)

	// Update changes after the affected region
	for i := lastAffected + 1; i < len(c.changes); i++ {
		c.changes[i].CurrentStart += delta
		c.changes[i].CurrentEnd += delta
	}

	// Replace affected changes with merged change
	mergedChange := Change{
		OriginalStart: mergedOrigStart,
		OriginalEnd:   mergedOrigEnd,
		CurrentStart:  mergedCurStart,
		CurrentEnd:    mergedCurEnd,
	}

	// Remove the range [firstAffected+1, lastAffected] and update firstAffected
	if firstAffected == lastAffected {
		c.changes[firstAffected] = mergedChange
	} else {
		c.changes[firstAffected] = mergedChange
		c.changes = append(c.changes[:firstAffected+1], c.changes[lastAffected+1:]...)
	}
}

// insertChange adds a new change at the correct sorted position.
func (c *ChangeTracker) insertChange(preStart, preEnd, postEnd int) {
	newChange := Change{
		OriginalStart: preStart,
		OriginalEnd:   preEnd,
		CurrentStart:  preStart,
		CurrentEnd:    postEnd,
	}

	delta := postEnd - preEnd

	// Find insertion point
	insertAt := len(c.changes)
	for i, ch := range c.changes {
		if preStart < ch.CurrentStart {
			insertAt = i
			break
		}
	}

	// Update all changes after insertion point
	for i := insertAt; i < len(c.changes); i++ {
		c.changes[i].CurrentStart += delta
		c.changes[i].CurrentEnd += delta
	}

	// Insert the new change
	c.changes = append(c.changes[:insertAt], append([]Change{newChange}, c.changes[insertAt:]...)...)
}

// ClearChanges removes all tracked changes.
func (c *ChangeTracker) ClearChanges() {
	c.changes = c.changes[:0]
}

// ForEachChange calls the given function for each tracked change.
func (c *ChangeTracker) ForEachChange(fn func(originalRange, currentRange text.TextRange)) {
	for _, ch := range c.changes {
		fn(
			text.NewTextRange(ch.OriginalStart, ch.OriginalEnd),
			text.NewTextRange(ch.CurrentStart, ch.CurrentEnd),
		)
	}
}
