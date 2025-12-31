package modifiers

import (
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/next/text"
	"github.com/zodimo/go-compose/internal/modifier"
)

// SelectionController manages text selection for a text composable.
// It handles registration with the selection system, updates on text/position changes,
// and drawing of selection highlights.
//
// This is essentially a Modifier.Node moved into remember because we need pointerInput.
//
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/foundation/foundation/src/commonMain/kotlin/androidx/compose/foundation/text/modifiers/SelectionController.kt
type SelectionController struct {
	selectableId             int64
	selectionRegistrar       SelectionRegistrar
	backgroundSelectionColor graphics.Color
	params                   StaticTextSelectionParams
	selectable               Selectable
	modifier                 modifier.Modifier
}

// NewSelectionController creates a new SelectionController.
func NewSelectionController(
	selectableId int64,
	selectionRegistrar SelectionRegistrar,
	backgroundSelectionColor graphics.Color,
) *SelectionController {
	sc := &SelectionController{
		selectableId:             selectableId,
		selectionRegistrar:       selectionRegistrar,
		backgroundSelectionColor: backgroundSelectionColor,
		params:                   EmptyStaticTextSelectionParams(),
		selectable:               nil,
	}

	// Build the modifier chain
	// In Kotlin: selectionRegistrar.makeSelectionModifier(...).pointerHoverIcon(PointerIcon.Text)
	// For now, we create an empty modifier as the full implementation requires
	// platform-specific pointer input handling
	sc.modifier = modifier.EmptyModifier

	return sc
}

// Modifier returns the modifier for this selection controller.
// This modifier should be applied to the text composable to enable selection.
func (sc *SelectionController) Modifier() modifier.Modifier {
	return sc.modifier
}

// SelectableId returns the unique identifier for this selectable.
func (sc *SelectionController) SelectableId() int64 {
	return sc.selectableId
}

// BackgroundSelectionColor returns the background color used for selection highlights.
func (sc *SelectionController) BackgroundSelectionColor() graphics.Color {
	return sc.backgroundSelectionColor
}

// OnRemembered is called when the controller is remembered in composition.
// It subscribes to the selection registrar.
func (sc *SelectionController) OnRemembered() {
	delegate := NewMultiWidgetSelectionDelegate(
		sc.selectableId,
		func() LayoutCoordinates {
			return sc.params.LayoutCoordinatesValue()
		},
		func() *text.TextLayoutResult {
			return sc.params.TextLayoutResultValue()
		},
	)
	sc.selectable = sc.selectionRegistrar.Subscribe(delegate)
}

// OnForgotten is called when the controller is no longer remembered.
// It unsubscribes from the selection registrar.
func (sc *SelectionController) OnForgotten() {
	if sc.selectable != nil {
		sc.selectionRegistrar.Unsubscribe(sc.selectable)
		sc.selectable = nil
	}
}

// OnAbandoned is called when the remember was abandoned before being committed.
// It unsubscribes from the selection registrar.
func (sc *SelectionController) OnAbandoned() {
	if sc.selectable != nil {
		sc.selectionRegistrar.Unsubscribe(sc.selectable)
		sc.selectable = nil
	}
}

// UpdateTextLayout updates the text layout result.
// If the text content has changed, it notifies the selection registrar.
func (sc *SelectionController) UpdateTextLayout(textLayoutResult *text.TextLayoutResult) {
	prevTextLayoutResult := sc.params.TextLayoutResultValue()

	// Don't notify on nil. We don't want every new Text that enters composition to
	// notify a selectable change. It was already handled when it was created.
	if prevTextLayoutResult != nil && textLayoutResult != nil {
		prevInput := prevTextLayoutResult.LayoutInput()
		newInput := textLayoutResult.LayoutInput()
		if prevInput.Text.String() != newInput.Text.String() {
			// Text content changed, notify selection to update itself.
			sc.selectionRegistrar.NotifySelectableChange(sc.selectableId)
		}
	}

	sc.params = sc.params.CopyWithTextLayoutResult(textLayoutResult)
}

// UpdateGlobalPosition updates the layout coordinates and notifies the selection registrar.
func (sc *SelectionController) UpdateGlobalPosition(coordinates LayoutCoordinates) {
	sc.params = sc.params.CopyWithLayoutCoordinates(coordinates)
	sc.selectionRegistrar.NotifyPositionChange(sc.selectableId)
}

// Draw draws the selection highlight using the provided draw function.
// The drawPath function is called with the selection path and background color.
//
// This is a simplified version that takes a draw function instead of DrawScope,
// as DrawScope depends on platform-specific graphics systems.
func (sc *SelectionController) Draw(drawPath func(path graphics.Path, color graphics.Color, shouldClip bool)) {
	subselections := sc.selectionRegistrar.Subselections()
	selection, exists := subselections[sc.selectableId]
	if !exists || selection == nil {
		return
	}

	var start, end int
	if !selection.HandlesCrossed {
		start = selection.Start.Offset
		end = selection.End.Offset
	} else {
		start = selection.End.Offset
		end = selection.Start.Offset
	}

	if start == end {
		return
	}

	lastOffset := 0
	if sc.selectable != nil {
		lastOffset = sc.selectable.GetLastVisibleOffset()
	}

	clippedStart := min(start, lastOffset)
	clippedEnd := min(end, lastOffset)

	selectionPath := sc.params.GetPathForRange(clippedStart, clippedEnd)
	if selectionPath == nil {
		return
	}

	drawPath(selectionPath, sc.backgroundSelectionColor, sc.params.ShouldClip())
}

// min returns the smaller of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Verify SelectionController implements RememberObserver at compile time.
var _ RememberObserver = (*SelectionController)(nil)
