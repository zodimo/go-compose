package undo

import (
	"strings"
	"time"

	"github.com/zodimo/go-compose/compose/ui/next/text"
)

// MergeTimeWindowMs is the maximum time between edits that allows merging.
const MergeTimeWindowMs = 2000

// TextUndoOperation represents a single atomic text change for undo/redo.
//
// It stores the text before and after a change, along with selection states,
// which allows the change to be applied or reversed.
//
// This is a port of androidx.compose.foundation.text.input.internal.undo.TextUndoOperation.
type TextUndoOperation struct {
	// Index is the start position of the change.
	Index int

	// PreText is the text that was replaced (deleted text).
	PreText string

	// PostText is the text that replaced it (inserted text).
	PostText string

	// PreSelection is the selection before the change.
	PreSelection text.TextRange

	// PostSelection is the selection after the change.
	PostSelection text.TextRange

	// TimeInMillis is when the change was committed.
	TimeInMillis int64

	// CanMerge indicates whether this operation can be merged with adjacent ones.
	CanMerge bool
}

// NewTextUndoOperation creates a new TextUndoOperation with the current timestamp.
func NewTextUndoOperation(
	index int,
	preText string,
	postText string,
	preSelection text.TextRange,
	postSelection text.TextRange,
) *TextUndoOperation {
	return &TextUndoOperation{
		Index:         index,
		PreText:       preText,
		PostText:      postText,
		PreSelection:  preSelection,
		PostSelection: postSelection,
		TimeInMillis:  time.Now().UnixMilli(),
		CanMerge:      true,
	}
}

// NewTextUndoOperationWithOptions creates a TextUndoOperation with explicit options.
func NewTextUndoOperationWithOptions(
	index int,
	preText string,
	postText string,
	preSelection text.TextRange,
	postSelection text.TextRange,
	timeInMillis int64,
	canMerge bool,
) *TextUndoOperation {
	return &TextUndoOperation{
		Index:         index,
		PreText:       preText,
		PostText:      postText,
		PreSelection:  preSelection,
		PostSelection: postSelection,
		TimeInMillis:  timeInMillis,
		CanMerge:      canMerge,
	}
}

// TextEditType returns the type of edit (Insert, Delete, or Replace).
func (op *TextUndoOperation) TextEditType() TextEditType {
	hasPreText := len(op.PreText) > 0
	hasPostText := len(op.PostText) > 0

	if !hasPreText && !hasPostText {
		// This shouldn't happen - one of them should have content
		return TextEditTypeReplace
	}
	if !hasPreText {
		return TextEditTypeInsert
	}
	if !hasPostText {
		return TextEditTypeDelete
	}
	return TextEditTypeReplace
}

// DeleteType returns the deletion direction (only valid for Delete type).
//
// This is used to determine if consecutive deletions can be merged.
func (op *TextUndoOperation) DeleteType() TextDeleteType {
	if op.TextEditType() != TextEditTypeDelete {
		return TextDeleteTypeNotByUser
	}

	// Post selection must be collapsed for user-initiated delete
	if !op.PostSelection.Collapsed() {
		return TextDeleteTypeNotByUser
	}

	if op.PreSelection.Collapsed() {
		// Single character delete
		if op.PreSelection.Start > op.PostSelection.Start {
			return TextDeleteTypeStart // Backspace
		}
		return TextDeleteTypeEnd // Delete key
	}

	// Selection was deleted
	if op.PreSelection.Start == op.PostSelection.Start && op.PreSelection.Start == op.Index {
		return TextDeleteTypeInner
	}

	return TextDeleteTypeNotByUser
}

// Merge attempts to merge this operation with the next operation.
//
// Returns the merged operation, or nil if they cannot be merged.
//
// Merge rules:
//  1. Both must have CanMerge=true
//  2. Time difference must be < 2 seconds
//  3. Newline insertions never merge
//  4. Only same TextEditType can merge
//  5. For insertions: next must extend this (cursor moving forward)
//  6. For deletions: must have same direction and share boundary
func (op *TextUndoOperation) Merge(next *TextUndoOperation) *TextUndoOperation {
	// Rule 1: Both must allow merging
	if !op.CanMerge || !next.CanMerge {
		return nil
	}

	// Rule 2: Time window check
	if next.TimeInMillis-op.TimeInMillis > MergeTimeWindowMs {
		return nil
	}

	// Rule 3: Newline insertions don't merge
	if strings.Contains(next.PostText, "\n") {
		return nil
	}

	// Rule 4: Same edit type required
	thisType := op.TextEditType()
	nextType := next.TextEditType()
	if thisType != nextType {
		return nil
	}

	switch thisType {
	case TextEditTypeInsert:
		return op.mergeInsert(next)
	case TextEditTypeDelete:
		return op.mergeDelete(next)
	default:
		// Replace operations don't merge
		return nil
	}
}

// mergeInsert merges two insert operations.
func (op *TextUndoOperation) mergeInsert(next *TextUndoOperation) *TextUndoOperation {
	// Next insert must start where this one ends
	thisEnd := op.Index + len(op.PostText)
	if next.Index != thisEnd {
		return nil
	}

	return NewTextUndoOperationWithOptions(
		op.Index,
		op.PreText,
		op.PostText+next.PostText,
		op.PreSelection,
		next.PostSelection,
		next.TimeInMillis,
		true,
	)
}

// mergeDelete merges two delete operations.
func (op *TextUndoOperation) mergeDelete(next *TextUndoOperation) *TextUndoOperation {
	// Must have same delete direction
	thisDeleteType := op.DeleteType()
	nextDeleteType := next.DeleteType()
	if thisDeleteType != nextDeleteType {
		return nil
	}

	switch thisDeleteType {
	case TextDeleteTypeStart:
		// Backspace: next deletes what's before current position
		if next.Index+len(next.PreText) != op.Index {
			return nil
		}
		return NewTextUndoOperationWithOptions(
			next.Index,
			next.PreText+op.PreText,
			"",
			op.PreSelection,
			next.PostSelection,
			next.TimeInMillis,
			true,
		)

	case TextDeleteTypeEnd:
		// Delete key: next deletes from same position
		if next.Index != op.Index {
			return nil
		}
		return NewTextUndoOperationWithOptions(
			op.Index,
			op.PreText+next.PreText,
			"",
			op.PreSelection,
			next.PostSelection,
			next.TimeInMillis,
			true,
		)

	case TextDeleteTypeInner:
		// Selection delete at same position
		if next.Index != op.Index {
			return nil
		}
		return NewTextUndoOperationWithOptions(
			op.Index,
			op.PreText+next.PreText,
			"",
			op.PreSelection,
			next.PostSelection,
			next.TimeInMillis,
			true,
		)

	default:
		return nil
	}
}
