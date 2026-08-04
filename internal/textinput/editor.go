package textinput

import (
	"image/color"

	"gioui.org/io/input"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/widget"

	"github.com/zodimo/go-compose/internal/layoutnode"
	tw "github.com/zodimo/go-compose/internal/textwidget"
)

// Editor wraps gioui.org/widget.Editor to hide gioui types from the public API.
type Editor struct {
	E widget.Editor
}

// NewEditor creates a new wrapped editor with the given configuration.
func NewEditor(singleLine bool, submit bool, mask rune) *Editor {
	return &Editor{E: widget.Editor{SingleLine: singleLine, Submit: submit, Mask: mask}}
}

// --- Text content methods ---

// SetText replaces the entire content.
func (e *Editor) SetText(s string) { e.E.SetText(s) }

// Text returns the current content.
func (e *Editor) Text() string { return e.E.Text() }

// Len returns the number of characters.
func (e *Editor) Len() int { return e.E.Len() }

// Insert inserts text at the caret position.
func (e *Editor) Insert(s string) int { return e.E.Insert(s) }

// Delete deletes grapheme clusters from the caret position.
func (e *Editor) Delete(n int) int { return e.E.Delete(n) }

// SelectedText returns the currently selected text.
func (e *Editor) SelectedText() string { return e.E.SelectedText() }

// --- Selection methods ---

// Selection returns the cursor/selection range.
func (e *Editor) Selection() (start, end int) { return e.E.Selection() }

// SetCaret positions the cursor.
func (e *Editor) SetCaret(start, end int) { e.E.SetCaret(start, end) }

// ClearSelection clears the selection.
func (e *Editor) ClearSelection() { e.E.ClearSelection() }

// --- Configuration methods ---

// SetMaxLen sets the maximum content length.
func (e *Editor) SetMaxLen(n int) { e.E.MaxLen = n }

// SetSingleLine sets single-line mode.
func (e *Editor) SetSingleLine(b bool) { e.E.SingleLine = b }

// SetReadOnly sets read-only mode.
func (e *Editor) SetReadOnly(b bool) { e.E.ReadOnly = b }

// SetLineHeightScale sets the line height scale factor.
func (e *Editor) SetLineHeightScale(s float32) { e.E.LineHeightScale = s }

// SetSubmit enables submit mode.
func (e *Editor) SetSubmit(b bool) { e.E.Submit = b }

// SetMask sets the mask character.
func (e *Editor) SetMask(r rune) { e.E.Mask = r }

// SetFilter sets the allowed character filter.
func (e *Editor) SetFilter(f string) { e.E.Filter = f }

// --- Event processing ---

// Update processes all pending editor events. For each ChangeEvent, it calls
// the onChange callback. This preserves the original behavior of processing
// every event in a loop.
func (e *Editor) Update(gtx layoutnode.LayoutContext, onChange func()) {
	for {
		event, ok := e.E.Update(*gtx.ToGio())
		if !ok {
			break
		}
		if _, isChange := event.(widget.ChangeEvent); isChange {
			if onChange != nil {
				onChange()
			}
		}
	}
}

// --- State sync ---

// UpdateState syncs the editor from external state. If the editor text differs
// from stateText, it restores the editor content and clamps the caret.
func (e *Editor) UpdateState(stateText string) {
	if e.E.Text() != stateText {
		caretStart, caretEnd := e.E.Selection()
		e.E.SetText(stateText)
		textLen := e.E.Len()
		if caretStart > textLen {
			caretStart = textLen
		}
		if caretEnd > textLen {
			caretEnd = textLen
		}
		e.E.SetCaret(caretStart, caretEnd)
	}
}

// --- Color draw op helper ---

// NewColorDrawOp records a color draw operation and returns an opaque DrawOp.
func NewColorDrawOp(gtx layoutnode.LayoutContext, col color.NRGBA) tw.DrawOp {
	macro := op.Record(gtx.ToGio().Ops)
	paint.ColorOp{Color: col}.Add(gtx.ToGio().Ops)
	return tw.NewDrawOp(macro.Stop())
}

// FocusTag returns the pointer used as the focus tag for key events.
func (ed *Editor) FocusTag() interface{} { return &ed.E }

// SourceFocused reports whether the given source is focused on this editor.
func (ed *Editor) SourceFocused(src input.Source) bool {
	return src.Focused(&ed.E)
}

// ProcessEvents drains all pending editor events from the frame.
// Returns true if a submit event was received.
func (ed *Editor) ProcessEvents(gtx layoutnode.LayoutContext) bool {
	hasSubmit := false
	for {
		ev, ok := ed.E.Update(*gtx.ToGio())
		if !ok {
			break
		}
		if _, ok := ev.(widget.SubmitEvent); ok {
			hasSubmit = true
		}
	}
	return hasSubmit
}

// Disabled reports whether the editor is disabled (no source).
func Disabled(gtx layoutnode.LayoutContext) bool {
	return gtx.ToGio().Source == (input.Source{})
}
