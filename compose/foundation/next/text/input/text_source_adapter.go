// Package input provides text field state management and editing infrastructure.
// This file provides the TextSourceAdapter that bridges input package types
// to the widget.TextSource interface for rendering.

package input

import (
	"io"
	"unicode/utf8"
)

// TextSourceAdapter adapts TextFieldCharSequence to implement the TextSource interface
// from the widget package for text rendering. This is the bridge between Compose-style
// text state management and Gio-based text rendering.
//
// TextSourceAdapter can be used in two modes:
//   - Read-only mode (from AnnotatedString): for BasicText
//   - Editable mode (from TextFieldState): for BasicTextField
type TextSourceAdapter struct {
	// For read-only text (BasicText)
	readOnlyText string

	// For editable text (BasicTextField)
	state *TextFieldState

	// Tracks whether contents have changed since last Changed() call
	changed      bool
	lastSnapshot string
}

// NewTextSourceAdapterFromString creates a read-only TextSourceAdapter from a string.
// This is suitable for use with BasicText.
func NewTextSourceAdapterFromString(text string) *TextSourceAdapter {
	return &TextSourceAdapter{
		readOnlyText: text,
		lastSnapshot: text,
	}
}

// NewTextSourceAdapterFromState creates an editable TextSourceAdapter from a TextFieldState.
// This is suitable for use with BasicTextField.
func NewTextSourceAdapterFromState(state *TextFieldState) *TextSourceAdapter {
	return &TextSourceAdapter{
		state:        state,
		lastSnapshot: state.Text(),
	}
}

// Size returns the total length of the text data in bytes.
func (a *TextSourceAdapter) Size() int64 {
	return int64(len(a.currentText()))
}

// Changed returns whether the contents have changed since the last call to Changed.
func (a *TextSourceAdapter) Changed() bool {
	current := a.currentText()
	changed := current != a.lastSnapshot || a.changed
	a.lastSnapshot = current
	a.changed = false
	return changed
}

// ReadAt reads len(p) bytes into p starting at byte offset off.
// It implements io.ReaderAt.
func (a *TextSourceAdapter) ReadAt(p []byte, off int64) (int, error) {
	text := a.currentText()
	if off >= int64(len(text)) {
		return 0, io.EOF
	}
	n := copy(p, text[off:])
	if n < len(p) {
		return n, io.EOF
	}
	return n, nil
}

// ReplaceRunes replaces runeCount runes starting at byteOffset within the
// data with the provided string.
// For read-only adapters, this is a no-op.
// For editable adapters, this forwards to the TextFieldState.
func (a *TextSourceAdapter) ReplaceRunes(byteOffset int64, runeCount int64, replacement string) {
	if a.state == nil {
		// Read-only mode - no-op
		return
	}

	// Convert byte offset to rune offset
	text := a.state.Text()
	runeOffset := a.byteToRuneOffset(text, int(byteOffset))

	a.state.Edit(func(buffer *TextFieldBuffer) {
		endRune := runeOffset + int(runeCount)
		if endRune > buffer.Length() {
			endRune = buffer.Length()
		}
		buffer.Replace(runeOffset, endRune, replacement)
	})

	a.changed = true
}

// currentText returns the current text content.
func (a *TextSourceAdapter) currentText() string {
	if a.state != nil {
		return a.state.Text()
	}
	return a.readOnlyText
}

// byteToRuneOffset converts a byte offset to a rune offset.
func (a *TextSourceAdapter) byteToRuneOffset(text string, byteOffset int) int {
	if byteOffset <= 0 {
		return 0
	}
	if byteOffset >= len(text) {
		return utf8.RuneCountInString(text)
	}
	return utf8.RuneCountInString(text[:byteOffset])
}

// SetText updates the text content for read-only adapters.
// For editable adapters, use the TextFieldState directly.
func (a *TextSourceAdapter) SetText(text string) {
	if a.state != nil {
		a.state.SetTextAndPlaceCursorAtEnd(text)
	} else {
		a.readOnlyText = text
	}
	a.changed = true
}

// Text returns the current text content.
func (a *TextSourceAdapter) Text() string {
	return a.currentText()
}

// IsEditable returns true if this adapter is backed by a TextFieldState.
func (a *TextSourceAdapter) IsEditable() bool {
	return a.state != nil
}

// State returns the underlying TextFieldState, or nil for read-only adapters.
func (a *TextSourceAdapter) State() *TextFieldState {
	return a.state
}
