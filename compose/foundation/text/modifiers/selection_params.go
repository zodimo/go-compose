package modifiers

import (
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/text"
	"github.com/zodimo/go-compose/compose/ui/text/style"
)

// StaticTextSelectionParams holds the parameters needed for static text selection.
// This is used to track layout coordinates and text layout results for selection handling.
//
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/foundation/foundation/src/commonMain/kotlin/androidx/compose/foundation/text/modifiers/SelectionController.kt
type StaticTextSelectionParams struct {
	layoutCoordinates LayoutCoordinates
	textLayoutResult  *text.TextLayoutResult
}

// EmptyStaticTextSelectionParams returns an empty StaticTextSelectionParams instance.
func EmptyStaticTextSelectionParams() StaticTextSelectionParams {
	return StaticTextSelectionParams{
		layoutCoordinates: nil,
		textLayoutResult:  nil,
	}
}

// NewStaticTextSelectionParams creates a new StaticTextSelectionParams.
func NewStaticTextSelectionParams(
	layoutCoordinates LayoutCoordinates,
	textLayoutResult *text.TextLayoutResult,
) StaticTextSelectionParams {
	return StaticTextSelectionParams{
		layoutCoordinates: layoutCoordinates,
		textLayoutResult:  textLayoutResult,
	}
}

// LayoutCoordinatesValue returns the layout coordinates.
func (p StaticTextSelectionParams) LayoutCoordinatesValue() LayoutCoordinates {
	return p.layoutCoordinates
}

// TextLayoutResultValue returns the text layout result.
func (p StaticTextSelectionParams) TextLayoutResultValue() *text.TextLayoutResult {
	return p.textLayoutResult
}

// GetPathForRange returns the path for the given text range, or nil if not available.
func (p StaticTextSelectionParams) GetPathForRange(start, end int) graphics.Path {
	if p.textLayoutResult == nil {
		return nil
	}
	return p.textLayoutResult.GetPathForRange(start, end)
}

// ShouldClip returns true if the selection should be clipped to the layout bounds.
// This is true when the text has visual overflow and the overflow mode is not Visible.
func (p StaticTextSelectionParams) ShouldClip() bool {
	if p.textLayoutResult == nil {
		return false
	}
	layoutInput := p.textLayoutResult.LayoutInput()
	return layoutInput.Overflow != style.OverFlowVisible && p.textLayoutResult.HasVisualOverflow()
}

// Copy creates a copy of StaticTextSelectionParams with optional overrides.
func (p StaticTextSelectionParams) Copy(
	layoutCoordinates LayoutCoordinates,
	textLayoutResult *text.TextLayoutResult,
) StaticTextSelectionParams {
	newLayoutCoordinates := p.layoutCoordinates
	if layoutCoordinates != nil {
		newLayoutCoordinates = layoutCoordinates
	}

	newTextLayoutResult := p.textLayoutResult
	if textLayoutResult != nil {
		newTextLayoutResult = textLayoutResult
	}

	return StaticTextSelectionParams{
		layoutCoordinates: newLayoutCoordinates,
		textLayoutResult:  newTextLayoutResult,
	}
}

// CopyWithLayoutCoordinates creates a copy with updated layout coordinates.
func (p StaticTextSelectionParams) CopyWithLayoutCoordinates(layoutCoordinates LayoutCoordinates) StaticTextSelectionParams {
	return StaticTextSelectionParams{
		layoutCoordinates: layoutCoordinates,
		textLayoutResult:  p.textLayoutResult,
	}
}

// CopyWithTextLayoutResult creates a copy with updated text layout result.
func (p StaticTextSelectionParams) CopyWithTextLayoutResult(textLayoutResult *text.TextLayoutResult) StaticTextSelectionParams {
	return StaticTextSelectionParams{
		layoutCoordinates: p.layoutCoordinates,
		textLayoutResult:  textLayoutResult,
	}
}
