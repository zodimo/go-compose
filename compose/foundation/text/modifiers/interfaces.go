// Package modifiers provides text modifier implementations for foundation text components.
//
// This package contains the selection controller and related types for handling
// text selection in compose text components.
package modifiers

import (
	"github.com/zodimo/go-compose/compose/ui/text"
)

// LayoutCoordinates provides access to layout coordinates for a composable.
// This is a subset of the full LayoutCoordinates interface, containing only
// the methods needed for selection handling.
//
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui/src/commonMain/kotlin/androidx/compose/ui/layout/LayoutCoordinates.kt
type LayoutCoordinates interface {
	// IsAttached returns true if this layout is currently attached.
	IsAttached() bool
}

// Selectable represents a selectable region of text.
//
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/foundation/foundation/src/commonMain/kotlin/androidx/compose/foundation/text/selection/Selectable.kt
type Selectable interface {
	// GetLastVisibleOffset returns the last visible character offset.
	GetLastVisibleOffset() int
}

// SelectionAnchorInfo contains information about a selection anchor point.
//
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/foundation/foundation/src/commonMain/kotlin/androidx/compose/foundation/text/selection/Selection.kt
type SelectionAnchorInfo struct {
	// Offset is the character offset of this anchor in the text.
	Offset int
	// Direction is the text direction at this anchor point.
	Direction int
	// SelectableId is the id of the selectable this anchor belongs to.
	SelectableId int64
}

// Selection represents a text selection with start and end anchors.
//
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/foundation/foundation/src/commonMain/kotlin/androidx/compose/foundation/text/selection/Selection.kt
type Selection struct {
	// Start is the start anchor of the selection.
	Start SelectionAnchorInfo
	// End is the end anchor of the selection.
	End SelectionAnchorInfo
	// HandlesCrossed indicates if the selection handles have crossed.
	// When true, the start handle is visually after the end handle.
	HandlesCrossed bool
}

// MultiWidgetSelectionDelegate is a selection delegate that coordinates selection
// across multiple text widgets.
//
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/foundation/foundation/src/commonMain/kotlin/androidx/compose/foundation/text/selection/MultiWidgetSelectionDelegate.kt
type MultiWidgetSelectionDelegate struct {
	// SelectableId is the unique identifier for this selectable.
	SelectableId int64
	// CoordinatesCallback returns the current layout coordinates.
	CoordinatesCallback func() LayoutCoordinates
	// LayoutResultCallback returns the current text layout result.
	LayoutResultCallback func() *text.TextLayoutResult
}

// SelectionRegistrar manages the registration and coordination of selectable text regions.
//
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/foundation/foundation/src/commonMain/kotlin/androidx/compose/foundation/text/selection/SelectionRegistrar.kt
type SelectionRegistrar interface {
	// Subscribe registers a selectable and returns the Selectable handle.
	Subscribe(delegate MultiWidgetSelectionDelegate) Selectable

	// Unsubscribe removes a selectable from the registrar.
	Unsubscribe(selectable Selectable)

	// NotifySelectableChange notifies the registrar that a selectable's content has changed.
	NotifySelectableChange(selectableId int64)

	// NotifyPositionChange notifies the registrar that a selectable's position has changed.
	NotifyPositionChange(selectableId int64)

	// Subselections returns the map of selection ID to Selection for active subselections.
	Subselections() map[int64]*Selection
}

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
