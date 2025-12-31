// Package modifiers provides text modifier implementations for foundation text components.
//
// This package contains the selection controller and related types for handling
// text selection in compose text components.
package modifiers

import (
	"github.com/zodimo/go-compose/compose/foundation/next/text/selection"
	"github.com/zodimo/go-compose/compose/ui/geometry"
	"github.com/zodimo/go-compose/compose/ui/layout"
	"github.com/zodimo/go-compose/compose/ui/next/text"
)

// Re-export types from selection package for backward compatibility and convenience.
// The canonical definitions are in the selection package.

// LayoutCoordinates provides access to layout coordinates for a composable.
// Use layout.LayoutCoordinates for new code.
type LayoutCoordinates = layout.LayoutCoordinates

// SelectionRegistrar is an alias to the selection package's SelectionRegistrar.
// Use selection.SelectionRegistrar for new code.
type SelectionRegistrar = selection.SelectionRegistrar

// Selection is an alias to the selection package's Selection.
// Use selection.Selection for new code.
type Selection = selection.Selection

// Selectable is an alias to the selection package's Selectable.
// Use selection.Selectable for new code.
type Selectable = selection.Selectable

// MultiWidgetSelectionDelegate is a selection delegate that coordinates selection
// across multiple text widgets. It implements the Selectable interface.
//
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/foundation/foundation/src/commonMain/kotlin/androidx/compose/foundation/text/selection/MultiWidgetSelectionDelegate.kt
type MultiWidgetSelectionDelegate struct {
	// selectableId is the unique identifier for this selectable.
	selectableId int64
	// CoordinatesCallback returns the current layout coordinates.
	CoordinatesCallback func() layout.LayoutCoordinates
	// LayoutResultCallback returns the current text layout result.
	LayoutResultCallback func() *text.TextLayoutResult
}

// NewMultiWidgetSelectionDelegate creates a new MultiWidgetSelectionDelegate.
func NewMultiWidgetSelectionDelegate(
	selectableId int64,
	coordinatesCallback func() layout.LayoutCoordinates,
	layoutResultCallback func() *text.TextLayoutResult,
) *MultiWidgetSelectionDelegate {
	return &MultiWidgetSelectionDelegate{
		selectableId:         selectableId,
		CoordinatesCallback:  coordinatesCallback,
		LayoutResultCallback: layoutResultCallback,
	}
}

// SelectableId implements Selectable.SelectableId.
func (d *MultiWidgetSelectionDelegate) SelectableId() int64 {
	return d.selectableId
}

// AppendSelectableInfoToBuilder implements Selectable.AppendSelectableInfoToBuilder.
func (d *MultiWidgetSelectionDelegate) AppendSelectableInfoToBuilder(builder selection.SelectionLayoutBuilder) {
	// TODO: Implement when SelectionLayoutBuilder is fully defined
}

// GetSelectAllSelection implements Selectable.GetSelectAllSelection.
func (d *MultiWidgetSelectionDelegate) GetSelectAllSelection() *selection.Selection {
	result := d.LayoutResultCallback()
	if result == nil {
		return nil
	}
	textLen := len(result.LayoutInput().Text.Text())
	if textLen == 0 {
		return nil
	}
	return &selection.Selection{
		Start: selection.AnchorInfo{
			Offset:       0,
			SelectableId: d.selectableId,
		},
		End: selection.AnchorInfo{
			Offset:       textLen,
			SelectableId: d.selectableId,
		},
		HandlesCrossed: false,
	}
}

// GetHandlePosition implements Selectable.GetHandlePosition.
func (d *MultiWidgetSelectionDelegate) GetHandlePosition(sel selection.Selection, isStartHandle bool) geometry.Offset {
	// TODO: Implement proper handle position calculation using TextLayoutResult
	return geometry.OffsetZero
}

// GetLayoutCoordinates implements Selectable.GetLayoutCoordinates.
func (d *MultiWidgetSelectionDelegate) GetLayoutCoordinates() layout.LayoutCoordinates {
	return d.CoordinatesCallback()
}

// TextLayoutResult implements Selectable.TextLayoutResult.
func (d *MultiWidgetSelectionDelegate) TextLayoutResult() *text.TextLayoutResult {
	return d.LayoutResultCallback()
}

// GetText implements Selectable.GetText.
func (d *MultiWidgetSelectionDelegate) GetText() text.AnnotatedString {
	result := d.LayoutResultCallback()
	if result == nil {
		return text.NewAnnotatedString("", nil, nil)
	}
	return result.LayoutInput().Text
}

// GetBoundingBox implements Selectable.GetBoundingBox.
func (d *MultiWidgetSelectionDelegate) GetBoundingBox(offset int) geometry.Rect {
	// TODO: Implement using TextLayoutResult
	return geometry.RectZero
}

// GetLineLeft implements Selectable.GetLineLeft.
func (d *MultiWidgetSelectionDelegate) GetLineLeft(offset int) float32 {
	// TODO: Implement using TextLayoutResult
	return 0
}

// GetLineRight implements Selectable.GetLineRight.
func (d *MultiWidgetSelectionDelegate) GetLineRight(offset int) float32 {
	// TODO: Implement using TextLayoutResult
	return 0
}

// GetCenterYForOffset implements Selectable.GetCenterYForOffset.
func (d *MultiWidgetSelectionDelegate) GetCenterYForOffset(offset int) float32 {
	// TODO: Implement using TextLayoutResult
	return 0
}

// GetRangeOfLineContaining implements Selectable.GetRangeOfLineContaining.
func (d *MultiWidgetSelectionDelegate) GetRangeOfLineContaining(offset int) text.TextRange {
	// TODO: Implement using TextLayoutResult
	return text.TextRangeZero
}

// GetLastVisibleOffset implements Selectable.GetLastVisibleOffset.
func (d *MultiWidgetSelectionDelegate) GetLastVisibleOffset() int {
	result := d.LayoutResultCallback()
	if result == nil {
		return 0
	}
	return len(result.LayoutInput().Text.Text())
}

// GetLineHeight implements Selectable.GetLineHeight.
func (d *MultiWidgetSelectionDelegate) GetLineHeight(offset int) float32 {
	// TODO: Implement using TextLayoutResult
	return 0
}

// Verify MultiWidgetSelectionDelegate implements Selectable at compile time.
var _ selection.Selectable = (*MultiWidgetSelectionDelegate)(nil)

// RememberObserver is an interface for objects that need lifecycle callbacks
// when remembered in composition.
//
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/runtime/runtime/src/commonMain/kotlin/androidx/compose/runtime/RememberObserver.kt
type RememberObserver interface {
	// OnRemembered is called when the object is successfully stored by remember.
	OnRemembered()

	// OnForgotten is called when the object is no longer being remembered.
	OnForgotten()

	// OnAbandoned is called when the remember call was not committed to the composition.
	OnAbandoned()
}
