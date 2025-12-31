package selection

import (
	"github.com/zodimo/go-compose/compose/ui/next/text"
)

// TextFieldSelectionState manages selection gestures and handles for text fields.
//
// This is a stub interface that documents the operational semantics.
// A full implementation requires integration with gesture detection and
// the compositor for rendering selection handles.
//
// Operational Semantics:
//
// Cursor Handle State:
//   - The cursor handle is shown when:
//   - The text field is focused
//   - There is no selection (cursor is collapsed)
//   - The cursor is within the visible viewport
//   - The handle provides visual feedback and drag-to-move functionality
//
// Selection Handles:
//   - Start and end handles appear when text is selected
//   - Dragging handles adjusts the selection range
//   - The start handle marks the anchor, end handle marks the focus
//
// Gesture Recognition:
//   - Single tap: Place cursor at tap position
//   - Double tap: Select word at position
//   - Triple tap: Select paragraph at position
//   - Long press: Begin selection, show magnifier
//   - Drag from handle: Expand/contract selection
//
// Text Toolbar:
//   - Shows cut/copy/paste/select-all when selection is active
//   - Position is calculated based on selection bounds
//   - Dismissed when selection changes or focus is lost
//
// This is a port of androidx.compose.foundation.text.input.internal.selection.TextFieldSelectionState.
type TextFieldSelectionState interface {
	// GetCursorHandleState returns the cursor handle visibility and position.
	//
	// Parameters:
	//   - includePosition: If false, omit position from result (for derived state optimization)
	//
	// Returns cursor handle state including:
	//   - Whether the handle is visible
	//   - Position in decoration box coordinates (if includePosition is true)
	GetCursorHandleState(includePosition bool) CursorHandleState

	// IsCursorVisible returns whether the cursor should be drawn.
	// The cursor blinks, so this value changes over time when focused.
	IsCursorVisible() bool

	// IsCursorHandleInVisibleBounds returns whether the cursor handle is
	// within visible bounds. Used to determine if the handle should be drawn.
	IsCursorHandleInVisibleBounds() bool

	// GetCursorRect returns the rectangle area the cursor occupies.
	// Returns zero rect if there is no cursor (i.e., selection is not collapsed).
	GetCursorRect() Rect

	// GetFocusRect returns where the focus should be for panning/scrolling.
	// Returns the cursor or selection rect in decorator coordinates.
	GetFocusRect() Rect

	// Cut cuts the selected text to clipboard.
	// No-op if nothing is selected.
	Cut()

	// Copy copies the selected text to clipboard.
	// No-op if nothing is selected.
	Copy()

	// Paste pastes text from clipboard at cursor/selection.
	// Replaces selection if present, otherwise inserts at cursor.
	Paste()

	// SelectAll selects all text in the field.
	SelectAll()
}

// CursorHandleState represents the state of the cursor handle.
type CursorHandleState struct {
	// Visible indicates whether the cursor handle should be shown.
	Visible bool

	// Position is the cursor handle position in decoration box coordinates.
	// This is where the visual handle should be rendered.
	Position Offset
}

// SelectionHandleState represents the state of a selection handle.
type SelectionHandleState struct {
	// Visible indicates whether this handle should be shown.
	Visible bool

	// Position is the handle position in decoration box coordinates.
	Position Offset

	// HandleDirection indicates whether this is the start or end handle.
	HandleDirection HandleDirection
}

// HandleDirection indicates whether a handle is at the start or end of selection.
type HandleDirection int

const (
	HandleDirectionStart HandleDirection = iota
	HandleDirectionEnd
)

// TextToolbarState represents the state of the floating text toolbar.
type TextToolbarState int

const (
	// TextToolbarStateNone: No toolbar shown
	TextToolbarStateNone TextToolbarState = iota
	// TextToolbarStateCursor: Show toolbar for cursor (paste only)
	TextToolbarStateCursor
	// TextToolbarStateSelection: Show toolbar for selection (cut/copy/paste)
	TextToolbarStateSelection
)

// Offset represents a 2D position.
type Offset struct {
	X float64
	Y float64
}

// Rect represents a rectangle.
type Rect struct {
	Left   float64
	Top    float64
	Right  float64
	Bottom float64
}

// Zero returns a zero rect.
var ZeroRect = Rect{}

// IsEmpty returns true if the rect has zero area.
func (r Rect) IsEmpty() bool {
	return r.Right <= r.Left || r.Bottom <= r.Top
}

// Width returns the width of the rect.
func (r Rect) Width() float64 {
	return r.Right - r.Left
}

// Height returns the height of the rect.
func (r Rect) Height() float64 {
	return r.Bottom - r.Top
}

// Clipboard stub interface for clipboard operations.
type Clipboard interface {
	// GetText returns text from the clipboard, or nil if unavailable.
	GetText() *string

	// SetText sets text to the clipboard.
	SetText(text string)

	// HasText returns true if the clipboard contains text.
	HasText() bool
}

// HapticFeedback stub interface for haptic feedback.
type HapticFeedback interface {
	// PerformHapticFeedback performs haptic feedback for the given type.
	PerformHapticFeedback(feedbackType HapticFeedbackType)
}

// HapticFeedbackType represents types of haptic feedback.
type HapticFeedbackType int

const (
	HapticFeedbackTypeLongPress HapticFeedbackType = iota
	HapticFeedbackTypeTextHandleMove
)

// Density represents display density.
type Density struct {
	Density   float64
	FontScale float64
}

// TextRange alias for convenience
type TextRange = text.TextRange
