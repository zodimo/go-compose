// Package input provides text field state management and editing infrastructure.
// This file provides the TextLayoutController that manages text layout state
// for bridging between Compose input APIs and Gio widget rendering.

package input

import (
	"gioui.org/font"

	"github.com/zodimo/go-compose/compose/foundation/next/text/widget"
	"github.com/zodimo/go-compose/compose/ui/next/text"
	"github.com/zodimo/go-compose/compose/ui/next/text/style"
	wstyle "github.com/zodimo/go-compose/compose/ui/text/style"
	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/internal/textconvert"
	"github.com/zodimo/go-compose/internal/textinput"
	tw "github.com/zodimo/go-compose/internal/textwidget"
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
	alignment  style.TextAlign
	wrapPolicy style.LineBreak
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
	c.alignment = alignment
	// Convert next/text/style.TextAlign to widget's style.TextAlign
	c.view.Alignment = wstyle.TextAlign(alignment)
}

// SetWrapPolicy sets the line wrap policy.
func (c *TextLayoutController) SetWrapPolicy(wrapPolicy style.LineBreak) {
	c.wrapPolicy = wrapPolicy
	// Convert next/text/style.LineBreak to widget's style.LineBreak
	c.view.WrapPolicy = wstyle.LineBreak(wrapPolicy)
}

// SetLineHeight sets the line height.
func (c *TextLayoutController) SetLineHeight(lineHeight unit.TextUnit) {
	c.view.LineHeight = lineHeight
}

// SetLineHeightScale sets the line height scale.
func (c *TextLayoutController) SetLineHeightScale(scale float32) {
	c.view.LineHeightScale = scale
}

// Layout performs text layout and returns dimensions.
func (c *TextLayoutController) Layout(gtx layoutnode.LayoutContext, shaper *tw.Shaper) layoutnode.LayoutDimensions {
	fontSpec := c.getFontSpec()
	size := c.GetFontSize()
	c.view.Layout(gtx, shaper, fontSpec, size)
	return c.view.Dimensions()
}

// PaintText clips and paints the text glyphs using the provided material.
func (c *TextLayoutController) PaintText(gtx layoutnode.LayoutContext, textMaterial tw.DrawOp) {
	c.view.PaintText(gtx, textMaterial)
}

// LayoutAndPaint performs layout and paints the text in one call.
// This is the main entry point for rendering text.
func (c *TextLayoutController) LayoutAndPaint(gtx layoutnode.LayoutContext, shaper *tw.Shaper, textMaterial tw.DrawOp) layoutnode.LayoutDimensions {
	fontSpec := c.getFontSpec()
	size := c.GetFontSize()
	c.view.Layout(gtx, shaper, fontSpec, size)
	c.view.PaintText(gtx, textMaterial)
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

// ConfigureFromTextStyle applies settings from a TextStyle.
func (c *TextLayoutController) ConfigureFromTextStyle(ts *text.TextStyle) {
	if ts == nil {
		return
	}
	c.textStyle = ts
	c.SetAlignment(ts.TextAlign())
	c.SetLineHeight(ts.LineHeight())
	c.SetWrapPolicy(ts.LineBreak())
}

// getFontSpec returns a FontSpec from the current text style.
func (c *TextLayoutController) getFontSpec() tw.FontSpec {
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
func (c *TextLayoutController) getFont() font.Font {
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
func (c *TextLayoutController) GetFontSize() unit.TextUnit {
	if c.textStyle == nil {
		return unit.Sp(14) // Default font size
	}
	return c.textStyle.FontSize()
}
