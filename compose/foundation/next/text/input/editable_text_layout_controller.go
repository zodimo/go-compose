// Package input provides text field state management and editing infrastructure.
// This file provides the EditableTextLayoutController that manages editable text layout state
// for bridging between Compose input APIs and Gio widget.Editor rendering.

package input

import (
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	gioOp "gioui.org/op"
	"gioui.org/op/paint"
	gioText "gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/zodimo/go-compose/compose/ui/text"
	uiFont "github.com/zodimo/go-compose/compose/ui/text/font"
	"github.com/zodimo/go-compose/compose/ui/text/style"
)

// EditableTextLayoutController manages editable text layout state and bridges between
// Compose input APIs and Gio widget.Editor rendering.
//
// Unlike TextLayoutController which uses widget.TextView (read-only), this controller
// uses widget.Editor which provides full input handling including:
// - Keyboard input
// - Cursor positioning and blinking
// - Text selection
// - Undo/redo
type EditableTextLayoutController struct {
	// Internal editor for input handling and rendering
	editor widget.Editor

	// TextFieldState for Compose-style state management
	state *TextFieldState

	// Current text style
	textStyle *text.TextStyle

	// Track last known text to detect external state changes
	lastStateText string

	// Configuration
	maxLines   int
	minLines   int
	softWrap   bool
	singleLine bool
	readOnly   bool
	alignment  gioText.Alignment
	wrapPolicy gioText.WrapPolicy

	// Input transformation to apply on changes
	inputTransformation InputTransformation

	// OnValueChange callback when text changes
	onValueChange func(string)

	// Colors
	selectionColor color.NRGBA
}

// NewEditableTextLayoutController creates a new EditableTextLayoutController.
func NewEditableTextLayoutController(state *TextFieldState) *EditableTextLayoutController {
	c := &EditableTextLayoutController{
		state:          state,
		maxLines:       0, // 0 means unlimited
		minLines:       1,
		softWrap:       true,
		singleLine:     false,
		readOnly:       false,
		lastStateText:  state.Text(),
		selectionColor: color.NRGBA{R: 100, G: 149, B: 237, A: 128}, // Cornflower blue
	}
	// Initialize editor with current state text
	c.editor.SetText(state.Text())
	return c
}

// SetTextStyle configures the text style.
func (c *EditableTextLayoutController) SetTextStyle(textStyle *text.TextStyle) {
	c.textStyle = textStyle
}

// SetMaxLines sets the maximum number of lines.
func (c *EditableTextLayoutController) SetMaxLines(maxLines int) {
	c.maxLines = maxLines
	c.editor.MaxLen = 0 // MaxLen is in runes, not lines - we don't limit this way
}

// SetMinLines sets the minimum number of lines.
func (c *EditableTextLayoutController) SetMinLines(minLines int) {
	c.minLines = minLines
}

// SetSoftWrap enables or disables soft wrapping.
func (c *EditableTextLayoutController) SetSoftWrap(softWrap bool) {
	c.softWrap = softWrap
}

// SetSingleLine enables or disables single line mode.
func (c *EditableTextLayoutController) SetSingleLine(singleLine bool) {
	c.singleLine = singleLine
	c.editor.SingleLine = singleLine
}

// SetReadOnly enables or disables read-only mode.
func (c *EditableTextLayoutController) SetReadOnly(readOnly bool) {
	c.readOnly = readOnly
	c.editor.ReadOnly = readOnly
}

// SetAlignment sets text alignment.
func (c *EditableTextLayoutController) SetAlignment(alignment style.TextAlign) {
	c.alignment = textAlignToGioAlign(alignment)
	c.editor.Alignment = c.alignment
}

// SetWrapPolicy sets the line wrap policy.
func (c *EditableTextLayoutController) SetWrapPolicy(wrapPolicy gioText.WrapPolicy) {
	c.wrapPolicy = wrapPolicy
	c.editor.WrapPolicy = wrapPolicy
}

// SetLineHeight sets the line height.
func (c *EditableTextLayoutController) SetLineHeight(lineHeight unit.Sp) {
	c.editor.LineHeight = lineHeight
}

// SetLineHeightScale sets the line height scale.
func (c *EditableTextLayoutController) SetLineHeightScale(scale float32) {
	c.editor.LineHeightScale = scale
}

// SetInputTransformation sets the input transformation.
func (c *EditableTextLayoutController) SetInputTransformation(t InputTransformation) {
	c.inputTransformation = t
}

// SetSelectionColor sets the selection highlight color.
func (c *EditableTextLayoutController) SetSelectionColor(color color.NRGBA) {
	c.selectionColor = color
}

// SetOnValueChange sets the callback for text changes.
func (c *EditableTextLayoutController) SetOnValueChange(callback func(string)) {
	c.onValueChange = callback
}

// ConfigureFromTextStyle applies settings from a TextStyle.
func (c *EditableTextLayoutController) ConfigureFromTextStyle(ts *text.TextStyle) {
	if ts == nil {
		return
	}
	c.textStyle = ts
	c.SetAlignment(ts.TextAlign())
	c.SetLineHeight(unit.Sp(ts.LineHeight().Value()))
	c.SetWrapPolicy(lineBreakToGioWrapPolicy(ts.LineBreak()))
}

// Update processes input events and syncs state.
// Should be called before Layout.
//
// This implements the controlled component pattern from Jetpack Compose:
// 1. Process editor events and call onValueChange with new text
// 2. Sync editor FROM state - if state wasn't updated by callback, editor reverts
func (c *EditableTextLayoutController) Update(gtx layout.Context) {
	// Process editor events and notify via callback
	for {
		event, ok := c.editor.Update(gtx)
		if !ok {
			break
		}
		if _, isChange := event.(widget.ChangeEvent); isChange {
			newText := c.editor.Text()
			// Apply input transformation if present
			if c.inputTransformation != nil {
				buffer := NewTextFieldBuffer(NewTextFieldCharSequence(newText, c.state.Selection()))
				c.inputTransformation.TransformInput(buffer)
				newText = buffer.String()
			}
			// Call the change callback - this is where the user updates state
			if c.onValueChange != nil {
				c.onValueChange(newText)
			}
		}
	}

	// Controlled component pattern: Always sync Editor FROM TextFieldState
	// This ensures if the callback didn't update state, the editor reverts
	stateText := c.state.Text()
	if c.editor.Text() != stateText {
		// Save cursor position before SetText (which resets it)
		caretStart, caretEnd := c.editor.Selection()

		c.editor.SetText(stateText)

		// Restore cursor position, clamped to valid range
		textLen := c.editor.Len()
		if caretStart > textLen {
			caretStart = textLen
		}
		if caretEnd > textLen {
			caretEnd = textLen
		}
		c.editor.SetCaret(caretStart, caretEnd)
	}
}

// Layout performs text layout and returns dimensions.
func (c *EditableTextLayoutController) Layout(gtx layout.Context, shaper *gioText.Shaper, textMaterial, selectMaterial gioOp.CallOp) layout.Dimensions {
	gioFont := c.GetFont()
	size := c.GetFontSize()
	return c.editor.Layout(gtx, shaper, gioFont, size, textMaterial, selectMaterial)
}

// LayoutAndPaint performs update, layout and paints the text in one call.
// This is the main entry point for rendering editable text.
func (c *EditableTextLayoutController) LayoutAndPaint(gtx layout.Context, shaper *gioText.Shaper, textMaterial gioOp.CallOp) layout.Dimensions {
	// Update state first
	c.Update(gtx)

	// Create selection material
	selectionColorMacro := gioOp.Record(gtx.Ops)
	paint.ColorOp{Color: c.selectionColor}.Add(gtx.Ops)
	selectMaterial := selectionColorMacro.Stop()

	return c.Layout(gtx, shaper, textMaterial, selectMaterial)
}

// Len returns the length of the text in runes.
func (c *EditableTextLayoutController) Len() int {
	return c.editor.Len()
}

// Selection returns the current selection range in runes.
func (c *EditableTextLayoutController) Selection() (start, end int) {
	return c.editor.Selection()
}

// SetCaret sets the caret position and selection.
func (c *EditableTextLayoutController) SetCaret(start, end int) {
	c.editor.SetCaret(start, end)
}

// Text returns the current text content.
func (c *EditableTextLayoutController) Text() string {
	return c.editor.Text()
}

// SetText sets the text content.
func (c *EditableTextLayoutController) SetText(s string) {
	c.editor.SetText(s)
}

// SelectedText returns the currently selected text.
func (c *EditableTextLayoutController) SelectedText() string {
	return c.editor.SelectedText()
}

// ClearSelection clears the selection.
func (c *EditableTextLayoutController) ClearSelection() {
	c.editor.ClearSelection()
}

// Insert inserts text at the current caret position.
func (c *EditableTextLayoutController) Insert(s string) int {
	return c.editor.Insert(s)
}

// Delete deletes runes from the caret position.
func (c *EditableTextLayoutController) Delete(graphemeClusters int) int {
	return c.editor.Delete(graphemeClusters)
}

// Editor returns the underlying widget.Editor for advanced use cases.
func (c *EditableTextLayoutController) Editor() *widget.Editor {
	return &c.editor
}

// State returns the underlying TextFieldState.
func (c *EditableTextLayoutController) State() *TextFieldState {
	return c.state
}

// GetFont returns a Gio font from the current text style.
func (c *EditableTextLayoutController) GetFont() font.Font {
	if c.textStyle == nil {
		return font.Font{}
	}
	return uiFont.ToGioFont(
		c.textStyle.FontFamily(),
		c.textStyle.FontWeight(),
		c.textStyle.FontStyle(),
	)
}

// GetFontSize returns the font size in Sp.
func (c *EditableTextLayoutController) GetFontSize() unit.Sp {
	if c.textStyle == nil {
		return 14 // Default font size
	}
	return unit.Sp(c.textStyle.FontSize().Value())
}
