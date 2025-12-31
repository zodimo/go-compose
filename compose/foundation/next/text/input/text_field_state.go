package input

import (
	"github.com/zodimo/go-compose/compose/ui/next/text"
)

// TextFieldState is the editable text state for a text field.
//
// This is the main state holder that manages:
//   - Current text content
//   - Cursor/selection position
//   - Composition (IME input in progress)
//   - Undo/redo history
//
// Use Edit() to make atomic changes to the text. The edit block receives
// a TextFieldBuffer that can be modified. After the block returns, changes
// are committed and undo history is updated.
//
// This is a port of androidx.compose.foundation.text.input.TextFieldState.
type TextFieldState struct {
	mainBuffer      *TextFieldBuffer
	textUndoManager *TextUndoManager
	value           *TextFieldCharSequence
	isEditing       bool

	// UndoState provides public access to undo/redo functionality.
	UndoState *UndoState
}

// NewTextFieldState creates a new TextFieldState with the given initial text.
// The cursor is placed at the end of the text.
func NewTextFieldState(initialText string) *TextFieldState {
	return NewTextFieldStateWithSelection(
		initialText,
		text.NewTextRangeCollapsed(len(initialText)),
	)
}

// NewTextFieldStateWithSelection creates a TextFieldState with specified text and selection.
func NewTextFieldStateWithSelection(initialText string, initialSelection text.TextRange) *TextFieldState {
	value := NewTextFieldCharSequence(initialText, initialSelection)
	state := &TextFieldState{
		textUndoManager: NewTextUndoManager(),
		value:           value,
	}
	state.mainBuffer = NewTextFieldBuffer(value)
	state.UndoState = newUndoState(state)
	return state
}

// Text returns the current text content.
func (s *TextFieldState) Text() string {
	return s.value.Text()
}

// Selection returns the current selection range.
// If collapsed, this is the cursor position.
func (s *TextFieldState) Selection() text.TextRange {
	return s.value.Selection()
}

// Composition returns the current IME composition range.
// Returns nil when there is no active composition.
func (s *TextFieldState) Composition() *text.TextRange {
	return s.value.Composition()
}

// Value returns the current immutable snapshot of the text field contents.
func (s *TextFieldState) Value() *TextFieldCharSequence {
	return s.value
}

// Edit performs an atomic edit operation.
//
// The provided block receives a TextFieldBuffer that reflects the current
// state. Any changes made to the buffer are committed when the block returns.
// Changes are recorded in the undo history.
//
// Example:
//
//	state.Edit(func(buffer *TextFieldBuffer) {
//	    buffer.Append("Hello")
//	    buffer.PlaceCursorAtEnd()
//	})
func (s *TextFieldState) Edit(block func(*TextFieldBuffer)) {
	s.editInternal(block, true, UndoBehaviorMergeIfPossible)
}

// EditWithBehavior performs an edit with specific undo behavior.
func (s *TextFieldState) EditWithBehavior(block func(*TextFieldBuffer), behavior TextFieldEditUndoBehavior) {
	s.editInternal(block, true, behavior)
}

// editInternal is the core edit implementation.
func (s *TextFieldState) editInternal(
	block func(*TextFieldBuffer),
	recordUndo bool,
	undoBehavior TextFieldEditUndoBehavior,
) {
	if s.isEditing {
		panic("cannot call edit() from within another edit() block")
	}

	s.isEditing = true
	defer func() { s.isEditing = false }()

	// Create buffer from current state
	buffer := NewTextFieldBuffer(s.value)

	// Save pre-edit state for undo
	preValue := s.value

	// Execute the edit block
	block(buffer)

	// Create new value from buffer
	newValue := buffer.ToTextFieldCharSequence()

	// Only update if something changed
	if !s.value.Equals(newValue) {
		// Record undo if requested and content changed
		if recordUndo && !preValue.ContentEquals(newValue.Text()) {
			s.textUndoManager.RecordChanges(preValue, newValue, buffer.Changes(), undoBehavior == UndoBehaviorMergeIfPossible)
		}

		s.value = newValue
		s.mainBuffer = NewTextFieldBuffer(newValue)
	}
}

// editWithNoSideEffects performs an edit without recording undo.
// This is used internally for undo/redo operations.
func (s *TextFieldState) editWithNoSideEffects(block func(*TextFieldBuffer)) {
	s.editInternal(block, false, UndoBehaviorMergeIfPossible)
}

// SetTextAndPlaceCursorAtEnd replaces all text and places the cursor at the end.
//
// This is a convenience method equivalent to:
//
//	state.Edit(func(buffer *TextFieldBuffer) {
//	    buffer.Replace(0, buffer.Length(), newText)
//	    buffer.PlaceCursorAtEnd()
//	})
func (s *TextFieldState) SetTextAndPlaceCursorAtEnd(newText string) {
	s.Edit(func(buffer *TextFieldBuffer) {
		buffer.Replace(0, buffer.Length(), newText)
		buffer.PlaceCursorAtEnd()
	})
}

// SetTextAndSelectAll replaces all text and selects it.
func (s *TextFieldState) SetTextAndSelectAll(newText string) {
	s.Edit(func(buffer *TextFieldBuffer) {
		buffer.Replace(0, buffer.Length(), newText)
		buffer.SelectAll()
	})
}

// ClearText removes all text from the field.
func (s *TextFieldState) ClearText() {
	s.Edit(func(buffer *TextFieldBuffer) {
		buffer.Delete(0, buffer.Length())
	})
}

// SelectAll selects all text in the field.
func (s *TextFieldState) SelectAll() {
	s.Edit(func(buffer *TextFieldBuffer) {
		buffer.SelectAll()
	})
}

// PlaceCursorAtEnd moves the cursor to the end of the text.
func (s *TextFieldState) PlaceCursorAtEnd() {
	s.Edit(func(buffer *TextFieldBuffer) {
		buffer.PlaceCursorAtEnd()
	})
}

// Undo reverts the last edit operation.
func (s *TextFieldState) Undo() {
	s.UndoState.Undo()
}

// Redo re-applies the last undone operation.
func (s *TextFieldState) Redo() {
	s.UndoState.Redo()
}

// CanUndo returns true if undo is possible.
func (s *TextFieldState) CanUndo() bool {
	return s.UndoState.CanUndo()
}

// CanRedo returns true if redo is possible.
func (s *TextFieldState) CanRedo() bool {
	return s.UndoState.CanRedo()
}

// ClearHistory removes all undo/redo history.
func (s *TextFieldState) ClearHistory() {
	s.UndoState.ClearHistory()
}

// Length returns the length of the current text.
func (s *TextFieldState) Length() int {
	return s.value.Len()
}

// IsEmpty returns true if the text field is empty.
func (s *TextFieldState) IsEmpty() bool {
	return s.value.Len() == 0
}
