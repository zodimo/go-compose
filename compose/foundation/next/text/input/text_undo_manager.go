package input

import (
	"github.com/zodimo/go-compose/compose/foundation/next/text/input/internal"
	"github.com/zodimo/go-compose/compose/foundation/next/text/input/internal/undo"
	"github.com/zodimo/go-compose/compose/ui/next/text"
)

// TextUndoCapacity is the default undo capacity for text fields.
const TextUndoCapacity = 100

// TextUndoManager manages undo/redo for text editing.
//
// It wraps UndoManager with text-specific merge logic, maintaining a
// "staging" operation that hasn't been committed yet. New edits attempt
// to merge with staging before flushing it to the undo stack.
//
// This is a port of androidx.compose.foundation.text.input.TextUndoManager.
type TextUndoManager struct {
	stagingUndo *undo.TextUndoOperation
	undoManager *undo.UndoManager[*undo.TextUndoOperation]
}

// NewTextUndoManager creates a new TextUndoManager.
func NewTextUndoManager() *TextUndoManager {
	return &TextUndoManager{
		undoManager: undo.NewUndoManager[*undo.TextUndoOperation](TextUndoCapacity),
	}
}

// NewTextUndoManagerWithStaging creates a TextUndoManager with initial staging operation.
func NewTextUndoManagerWithStaging(stagingUndo *undo.TextUndoOperation) *TextUndoManager {
	return &TextUndoManager{
		stagingUndo: stagingUndo,
		undoManager: undo.NewUndoManager[*undo.TextUndoOperation](TextUndoCapacity),
	}
}

// CanUndo returns true if undo is possible.
func (m *TextUndoManager) CanUndo() bool {
	return m.stagingUndo != nil || m.undoManager.CanUndo()
}

// CanRedo returns true if redo is possible.
func (m *TextUndoManager) CanRedo() bool {
	return m.undoManager.CanRedo()
}

// Undo reverts the last edit and applies it to the state.
func (m *TextUndoManager) Undo(state *TextFieldState) {
	if !m.CanUndo() {
		return
	}

	var op *undo.TextUndoOperation

	// Flush staging first if exists
	if m.stagingUndo != nil {
		m.flush()
	}

	// Pop from undo stack
	op = m.undoManager.Undo()

	// Apply the undo operation
	state.editWithNoSideEffects(func(buffer *TextFieldBuffer) {
		buffer.Replace(op.Index, op.Index+len(op.PostText), op.PreText)
		buffer.SetSelection(op.PreSelection.Start, op.PreSelection.End)
	})
}

// Redo re-applies a previously undone edit.
func (m *TextUndoManager) Redo(state *TextFieldState) {
	if !m.CanRedo() {
		return
	}

	// Flush any staging first
	if m.stagingUndo != nil {
		m.flush()
	}

	op := m.undoManager.Redo()

	// Apply the redo operation
	state.editWithNoSideEffects(func(buffer *TextFieldBuffer) {
		buffer.Replace(op.Index, op.Index+len(op.PreText), op.PostText)
		buffer.SetSelection(op.PostSelection.Start, op.PostSelection.End)
	})
}

// Record adds a new undo operation.
//
// The operation is staged for potential merging with subsequent operations.
func (m *TextUndoManager) Record(op *undo.TextUndoOperation) {
	if m.stagingUndo == nil {
		m.stagingUndo = op
		return
	}

	// Attempt to merge with staging
	merged := m.stagingUndo.Merge(op)
	if merged != nil {
		m.stagingUndo = merged
	} else {
		// Flush staging and start new staging
		m.flush()
		m.stagingUndo = op
	}
}

// RecordWithBehavior adds an undo operation with specific behavior.
func (m *TextUndoManager) RecordWithBehavior(op *undo.TextUndoOperation, behavior undo.TextFieldEditUndoBehavior) {
	switch behavior {
	case undo.TextFieldEditUndoBehaviorClearHistory:
		m.ClearHistory()
		return

	case undo.TextFieldEditUndoBehaviorNeverMerge:
		m.flush()
		op.CanMerge = false
		m.undoManager.Record(op)
		return

	default: // MergeIfPossible
		m.Record(op)
	}
}

// RecordChanges converts changes to an undo operation and records it.
//
// Parameters:
//   - pre: The text state before the changes
//   - post: The text state after the changes
//   - changes: The list of changes that were made
//   - allowMerge: Whether this operation can be merged with the previous one
func (m *TextUndoManager) RecordChanges(
	pre, post *TextFieldCharSequence,
	changes internal.ChangeList,
	allowMerge bool,
) {
	if changes.ChangeCount() == 0 {
		return
	}

	// For now, we record the full text change as a single operation.
	// More sophisticated implementations could record individual changes.
	preText := pre.Text()
	postText := post.Text()

	// Find the changed region
	start, preEnd, postEnd := findChangeRegion(preText, postText)
	if start == preEnd && start == postEnd {
		return // No actual change
	}

	op := undo.NewTextUndoOperationWithOptions(
		start,
		preText[start:preEnd],
		postText[start:postEnd],
		pre.Selection(),
		post.Selection(),
		0, // Will be set by NewTextUndoOperation
		allowMerge,
	)

	if allowMerge {
		m.Record(op)
	} else {
		m.RecordWithBehavior(op, undo.TextFieldEditUndoBehaviorNeverMerge)
	}
}

// flush moves the staging operation to the undo stack.
func (m *TextUndoManager) flush() {
	if m.stagingUndo == nil {
		return
	}
	m.undoManager.Record(m.stagingUndo)
	m.stagingUndo = nil
}

// ClearHistory removes all undo and redo history.
func (m *TextUndoManager) ClearHistory() {
	m.stagingUndo = nil
	m.undoManager.ClearHistory()
}

// ForceFlush flushes any staging operation to the undo stack.
// This ensures the staging operation is recorded even if no new operation follows.
func (m *TextUndoManager) ForceFlush() {
	m.flush()
}

// findChangeRegion finds the smallest region that differs between two strings.
// Returns (start, endInA, endInB) where:
//   - start is the first differing position
//   - endInA is the end of the differing region in string A
//   - endInB is the end of the differing region in string B
func findChangeRegion(a, b string) (start, endA, endB int) {
	runesA := []rune(a)
	runesB := []rune(b)
	lenA := len(runesA)
	lenB := len(runesB)

	// Find common prefix
	start = 0
	minLen := lenA
	if lenB < minLen {
		minLen = lenB
	}
	for start < minLen && runesA[start] == runesB[start] {
		start++
	}

	// Find common suffix (from the end)
	endA = lenA
	endB = lenB
	for endA > start && endB > start && runesA[endA-1] == runesB[endB-1] {
		endA--
		endB--
	}

	return start, endA, endB
}

// UndoState provides public access to undo/redo functionality.
//
// This is the public API for controlling undo/redo from outside the
// text field internals.
//
// This is a port of androidx.compose.foundation.text.input.UndoState.
type UndoState struct {
	state *TextFieldState
}

// newUndoState creates a new UndoState for the given TextFieldState.
func newUndoState(state *TextFieldState) *UndoState {
	return &UndoState{state: state}
}

// CanUndo returns true if undo is possible.
func (u *UndoState) CanUndo() bool {
	return u.state.textUndoManager.CanUndo()
}

// CanRedo returns true if redo is possible.
func (u *UndoState) CanRedo() bool {
	return u.state.textUndoManager.CanRedo()
}

// Undo reverts the latest edit action or group of merged actions.
// Calling repeatedly continues undoing previous actions.
func (u *UndoState) Undo() {
	u.state.textUndoManager.Undo(u.state)
}

// Redo re-applies a change previously reverted via Undo.
func (u *UndoState) Redo() {
	u.state.textUndoManager.Redo(u.state)
}

// ClearHistory removes all undo and redo history.
func (u *UndoState) ClearHistory() {
	u.state.textUndoManager.ClearHistory()
}

// Helper: adjustTextRange for undo/redo operations (already defined in text_field_buffer.go)

// SetSelectionCoerced is an extension function for TextFieldBuffer to set selection
// with values coerced to valid range.
func SetSelectionCoerced(buffer *TextFieldBuffer, start, end int) {
	buffer.SetSelection(start, end)
}

// Helper type for undo behavior
type TextFieldEditUndoBehavior = undo.TextFieldEditUndoBehavior

// Re-export undo behavior constants for convenience
const (
	UndoBehaviorMergeIfPossible = undo.TextFieldEditUndoBehaviorMergeIfPossible
	UndoBehaviorClearHistory    = undo.TextFieldEditUndoBehaviorClearHistory
	UndoBehaviorNeverMerge      = undo.TextFieldEditUndoBehaviorNeverMerge
)

// TextRange helper that creates a TextRange
func textRange(start, end int) text.TextRange {
	return text.NewTextRange(start, end)
}
