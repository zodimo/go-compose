// Package input provides text field state management and editing infrastructure.
// This file provides the TextLayoutController that manages text layout state
// for bridging between Compose input APIs and Gio widget rendering.

package input

import (
	"gioui.org/font"
	"gioui.org/layout"
	gioOp "gioui.org/op"
	gioText "gioui.org/text"
	"gioui.org/unit"

	"github.com/zodimo/go-compose/compose/foundation/next/text/widget"
	"github.com/zodimo/go-compose/compose/ui/next/text"
	uiFont "github.com/zodimo/go-compose/compose/ui/next/text/font"
	"github.com/zodimo/go-compose/compose/ui/next/text/style"
)

// TextLayoutController manages text layout state and bridges between
// Compose input APIs and Gio widget rendering.
//
// It wraps a widget.TextView internally and provides a higher-level API
// for BasicText and BasicTextField composables.
type TextLayoutController struct {
	// Internal text view for rendering
	view widget.TextView

	// Source adapter that provides text data
	source *TextSourceAdapter

	// Current text style
	textStyle *text.TextStyle

	// Configuration
	maxLines   int
	minLines   int
	softWrap   bool
	singleLine bool
	truncator  string
	alignment  gioText.Alignment
	wrapPolicy gioText.WrapPolicy
}

// NewTextLayoutController creates a new TextLayoutController.
func NewTextLayoutController(source *TextSourceAdapter) *TextLayoutController {
	c := &TextLayoutController{
		source:     source,
		maxLines:   1,
		minLines:   1,
		softWrap:   true,
		singleLine: false,
		truncator:  "…",
	}
	c.view.SetSource(source)
	return c
}

// SetTextStyle configures the text style.
func (c *TextLayoutController) SetTextStyle(textStyle *text.TextStyle) {
	c.textStyle = textStyle
}

// SetMaxLines sets the maximum number of lines.
func (c *TextLayoutController) SetMaxLines(maxLines int) {
	c.maxLines = maxLines
	c.view.MaxLines = maxLines
}

// SetMinLines sets the minimum number of lines.
func (c *TextLayoutController) SetMinLines(minLines int) {
	c.minLines = minLines
}

// SetSoftWrap enables or disables soft wrapping.
func (c *TextLayoutController) SetSoftWrap(softWrap bool) {
	c.softWrap = softWrap
}

// SetSingleLine enables or disables single line mode.
func (c *TextLayoutController) SetSingleLine(singleLine bool) {
	c.singleLine = singleLine
	c.view.SingleLine = singleLine
}

// SetTruncator sets the truncation string.
func (c *TextLayoutController) SetTruncator(truncator string) {
	c.truncator = truncator
	c.view.Truncator = truncator
}

// SetAlignment sets text alignment.
func (c *TextLayoutController) SetAlignment(alignment style.TextAlign) {
	c.alignment = textAlignToGioAlign(alignment)
	c.view.Alignment = c.alignment
}

// SetWrapPolicy sets the line wrap policy.
func (c *TextLayoutController) SetWrapPolicy(wrapPolicy gioText.WrapPolicy) {
	c.wrapPolicy = wrapPolicy
	c.view.WrapPolicy = wrapPolicy
}

// SetLineHeight sets the line height.
func (c *TextLayoutController) SetLineHeight(lineHeight unit.Sp) {
	c.view.LineHeight = lineHeight
}

// SetLineHeightScale sets the line height scale.
func (c *TextLayoutController) SetLineHeightScale(scale float32) {
	c.view.LineHeightScale = scale
}

// Layout performs text layout and returns dimensions.
func (c *TextLayoutController) Layout(gtx layout.Context, shaper *gioText.Shaper, gioFont font.Font, size unit.Sp) layout.Dimensions {
	c.view.Layout(gtx, shaper, gioFont, size)
	return c.view.Dimensions()
}

// PaintText clips and paints the text glyphs using the provided material.
func (c *TextLayoutController) PaintText(gtx layout.Context, textMaterial gioOp.CallOp) {
	c.view.PaintText(gtx, textMaterial)
}

// LayoutAndPaint performs layout and paints the text in one call.
// This is the main entry point for rendering text.
func (c *TextLayoutController) LayoutAndPaint(gtx layout.Context, shaper *gioText.Shaper, textMaterial gioOp.CallOp) layout.Dimensions {
	gioFont := c.GetFont()
	size := c.GetFontSize()
	c.view.Layout(gtx, shaper, gioFont, size)
	c.PaintText(gtx, textMaterial)
	return c.view.Dimensions()
}

// Len returns the length of the text in runes.
func (c *TextLayoutController) Len() int {
	return c.view.Len()
}

// Selection returns the current selection range in runes.
func (c *TextLayoutController) Selection() (start, end int) {
	return c.view.Selection()
}

// SetCaret sets the caret position and selection.
func (c *TextLayoutController) SetCaret(start, end int) {
	c.view.SetCaret(start, end)
}

// Truncated returns whether the text is truncated.
func (c *TextLayoutController) Truncated() bool {
	return c.view.Truncated()
}

// TextView returns the underlying widget.TextView for advanced use cases.
func (c *TextLayoutController) TextView() *widget.TextView {
	return &c.view
}

// Source returns the text source adapter.
func (c *TextLayoutController) Source() *TextSourceAdapter {
	return c.source
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
		// Gio doesn't support justify, fall back to Start
		return gioText.Start
	default:
		return gioText.Start
	}
}

// ConfigureFromTextStyle applies settings from a TextStyle.
func (c *TextLayoutController) ConfigureFromTextStyle(ts *text.TextStyle) {
	if ts == nil {
		return
	}
	c.textStyle = ts
	c.SetAlignment(ts.TextAlign())
	c.SetLineHeight(unit.Sp(ts.LineHeight().Value()))
	c.SetWrapPolicy(lineBreakToGioWrapPolicy(ts.LineBreak()))
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

// GetFont returns a Gio font from the current text style.
func (c *TextLayoutController) GetFont() font.Font {
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
func (c *TextLayoutController) GetFontSize() unit.Sp {
	if c.textStyle == nil {
		return 14 // Default font size
	}
	return unit.Sp(c.textStyle.FontSize().Value())
}
