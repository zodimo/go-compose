// Package input provides text field state management and editing infrastructure.
// This file provides the EditableTextLayoutController that manages editable text layout state
// for bridging between Compose input APIs and Gio widget.Editor rendering.

package input
// Engine-bound (Door 2): this component's public API is gioui-free; its
// implementation uses the engine directly behind the seam (internal/render).


import (
	"image/color"

	"gioui.org/font"
	gioText "gioui.org/text"

	"github.com/zodimo/go-compose/compose/ui/unit"

	"github.com/zodimo/go-compose/compose/ui/next/text"
	"github.com/zodimo/go-compose/compose/ui/next/text/style"
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/internal/textconvert"
	"github.com/zodimo/go-compose/internal/textinput"
	tw "github.com/zodimo/go-compose/internal/textwidget"
	"github.com/zodimo/go-compose/internal/unitconvert"
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
	editor *textinput.Editor

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
	alignment  style.TextAlign
	wrapPolicy style.LineBreak

	// Input transformation to apply on changes
	inputTransformation InputTransformation

	// OnValueChange callback when text changes
	onValueChange func(string)

	// Colors
	selectionColor color.NRGBA
}

// NewEditableTextLayoutController creates a new EditableTextLayoutController.
func NewEditableTextLayoutController(state *TextFieldState) *EditableTextLayoutController {
	editor := textinput.NewEditor(false, false, 0)
	editor.SetText(state.Text())
	c := &EditableTextLayoutController{
		editor:          editor,
		state:           state,
		maxLines:        0, // 0 means unlimited
		minLines:        1,
		softWrap:        true,
		singleLine:      false,
		readOnly:        false,
		lastStateText:   state.Text(),
		selectionColor:  color.NRGBA{R: 100, G: 149, B: 237, A: 128}, // Cornflower blue
	}
	return c
}

// SetTextStyle configures the text style.
func (c *EditableTextLayoutController) SetTextStyle(textStyle *text.TextStyle) {
	c.textStyle = textStyle
}

// SetMaxLines sets the maximum number of lines.
func (c *EditableTextLayoutController) SetMaxLines(maxLines int) {
	c.maxLines = maxLines
	c.editor.SetMaxLen(0) // MaxLen is in runes, not lines - we don't limit this way
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
	c.editor.SetSingleLine(singleLine)
}

// SetReadOnly enables or disables read-only mode.
func (c *EditableTextLayoutController) SetReadOnly(readOnly bool) {
	c.readOnly = readOnly
	c.editor.SetReadOnly(readOnly)
}

// SetAlignment sets text alignment.
func (c *EditableTextLayoutController) SetAlignment(alignment style.TextAlign) {
	c.alignment = alignment
	c.editor.E.Alignment = textAlignToGioAlign(alignment)
}

// SetWrapPolicy sets the line wrap policy.
func (c *EditableTextLayoutController) SetWrapPolicy(wrapPolicy style.LineBreak) {
	c.wrapPolicy = wrapPolicy
	c.editor.E.WrapPolicy = lineBreakToGioWrapPolicy(wrapPolicy)
}

// SetLineHeight sets the line height.
func (c *EditableTextLayoutController) SetLineHeight(lineHeight unit.TextUnit) {
	c.editor.E.LineHeight = unitconvert.TextUnitToGioSpUnsafe(lineHeight)
}

// SetLineHeightScale sets the line height scale.
func (c *EditableTextLayoutController) SetLineHeightScale(scale float32) {
	c.editor.SetLineHeightScale(scale)
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
	c.SetLineHeight(ts.LineHeight())
	c.SetWrapPolicy(ts.LineBreak())
}

// Update processes input events and syncs state.
// Should be called before Layout.
//
// This implements the controlled component pattern from Jetpack Compose:
// 1. Process editor events and call onValueChange with new text
// 2. Sync editor FROM state - if state wasn't updated by callback, editor reverts
func (c *EditableTextLayoutController) Update(gtx layoutnode.LayoutContext) {
	// Process editor events and notify via callback
	c.editor.Update(gtx, func() {
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
	})

	// Controlled component pattern: Always sync Editor FROM TextFieldState
	// This ensures if the callback didn't update state, the editor reverts
	c.editor.UpdateState(c.state.Text())
}

// Layout performs text layout and returns dimensions.
func (c *EditableTextLayoutController) Layout(gtx layoutnode.LayoutContext, shaper *tw.Shaper, textMaterial, selectMaterial tw.DrawOp) layoutnode.LayoutDimensions {
	fontSpec := c.getFontSpec()
	size := c.GetFontSize()
	dims := c.editor.E.Layout(*gtx.ToGio(), shaper.S, tw.ToGioFont(fontSpec), unitconvert.TextUnitToGioSpUnsafe(size), textMaterial.O, selectMaterial.O)
	return layoutnode.FromGioDimensions(dims)
}

// LayoutAndPaint performs update, layout and paints the text in one call.
// This is the main entry point for rendering editable text.
func (c *EditableTextLayoutController) LayoutAndPaint(gtx layoutnode.LayoutContext, shaper *tw.Shaper, textMaterial tw.DrawOp) layoutnode.LayoutDimensions {
	// Update state first
	c.Update(gtx)

	// Create selection material
	selectMaterial := textinput.NewColorDrawOp(gtx, c.selectionColor)

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

// Editor returns the internal editor wrapper for advanced use cases.
func (c *EditableTextLayoutController) Editor() *textinput.Editor {
	return c.editor
}

// State returns the underlying TextFieldState.
func (c *EditableTextLayoutController) State() *TextFieldState {
	return c.state
}

// getFontSpec returns a FontSpec from the current text style.
func (c *EditableTextLayoutController) getFontSpec() tw.FontSpec {
	if c.textStyle == nil {
		return tw.DefaultFontSpec()
	}
	return textinput.ToFontSpecNext(
		c.textStyle.FontFamily(),
		c.textStyle.FontWeight(),
		c.textStyle.FontStyle(),
	)
}

// getFont returns a Gio font from the current text style (for internal use).
func (c *EditableTextLayoutController) getFont() font.Font {
	if c.textStyle == nil {
		return font.Font{}
	}
	return textconvert.ToGioFontNext(
		c.textStyle.FontFamily(),
		c.textStyle.FontWeight(),
		c.textStyle.FontStyle(),
	)
}

// GetFontSize returns the font size.
func (c *EditableTextLayoutController) GetFontSize() unit.TextUnit {
	if c.textStyle == nil {
		return unit.Sp(14) // Default font size
	}
	return c.textStyle.FontSize()
}

// textAlignToGioAlign converts compose TextAlign to gio text.Alignment.
func textAlignToGioAlign(textAlign style.TextAlign) gioText.Alignment {
	switch textAlign {
	case style.TextAlignLeft:
		return gioText.Start
	case style.TextAlignCenter:
		return gioText.Middle
	case style.TextAlignRight:
		return gioText.End
	case style.TextAlignJustify:
		return gioText.Start
	default:
		return gioText.Start
	}
}

// lineBreakToGioWrapPolicy converts compose LineBreak to gio WrapPolicy.
func lineBreakToGioWrapPolicy(lineBreak style.LineBreak) gioText.WrapPolicy {
	switch lineBreak {
	case style.LineBreakSimple:
		return gioText.WrapGraphemes
	case style.LineBreakHeading:
		return gioText.WrapWords
	case style.LineBreakParagraph:
		return gioText.WrapHeuristically
	default:
		return gioText.WrapWords
	}
}
