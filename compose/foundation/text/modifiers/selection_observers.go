package modifiers

import "github.com/zodimo/go-compose/compose/ui/geometry"

// SelectionAdjustment determines how selection should be adjusted during gestures.
//
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/foundation/foundation/src/commonMain/kotlin/androidx/compose/foundation/text/selection/SelectionAdjustment.kt
type SelectionAdjustment int

const (
	// SelectionAdjustmentNone means no adjustment to selection.
	SelectionAdjustmentNone SelectionAdjustment = iota

	// SelectionAdjustmentCharacter adjusts selection to character boundaries.
	SelectionAdjustmentCharacter

	// SelectionAdjustmentWord adjusts selection to word boundaries.
	SelectionAdjustmentWord

	// SelectionAdjustmentParagraph adjusts selection to paragraph boundaries.
	SelectionAdjustmentParagraph
)

// String returns the string representation of SelectionAdjustment.
func (s SelectionAdjustment) String() string {
	switch s {
	case SelectionAdjustmentNone:
		return "None"
	case SelectionAdjustmentCharacter:
		return "Character"
	case SelectionAdjustmentWord:
		return "Word"
	case SelectionAdjustmentParagraph:
		return "Paragraph"
	default:
		return "Unknown"
	}
}

// TextDragObserver observes text drag gestures for selection.
// This is used for long-press-drag selection on touch devices.
//
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/foundation/foundation/src/commonMain/kotlin/androidx/compose/foundation/text/TextDragObserver.kt
type TextDragObserver interface {
	// OnDown is called when the pointer first touches down.
	// Not supported for long-press-drag selection.
	OnDown(point geometry.Offset)

	// OnUp is called when the pointer is released.
	OnUp()

	// OnStart is called when a drag gesture starts (after long-press is detected).
	// The selectionAdjustment determines how the selection should be adjusted.
	OnStart(startPoint geometry.Offset, selectionAdjustment SelectionAdjustment)

	// OnDrag is called during a drag gesture with the delta from the previous position.
	OnDrag(delta geometry.Offset)

	// OnStop is called when the drag gesture ends normally.
	OnStop()

	// OnCancel is called when the drag gesture is cancelled.
	OnCancel()
}

// MouseSelectionObserver observes mouse selection gestures.
// This is used for click-and-drag selection with a mouse.
//
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/foundation/foundation/src/commonMain/kotlin/androidx/compose/foundation/text/selection/MouseSelectionObserver.kt
type MouseSelectionObserver interface {
	// OnExtend is called when the selection should be extended to the given position.
	// Returns true if the extension was successful.
	OnExtend(downPosition geometry.Offset) bool

	// OnExtendDrag is called during a drag to extend the selection.
	// Returns true if the drag should continue.
	OnExtendDrag(dragPosition geometry.Offset) bool

	// OnStart is called when a mouse selection gesture starts.
	// The adjustment determines how the selection should be adjusted (e.g., word on double-click).
	// The clickCount indicates single, double, or triple click.
	// Returns true if the selection was started successfully.
	OnStart(downPosition geometry.Offset, adjustment SelectionAdjustment, clickCount int) bool

	// OnDrag is called during a mouse drag gesture.
	// The adjustment determines how the selection should be adjusted.
	// Returns true if the drag should continue.
	OnDrag(dragPosition geometry.Offset, adjustment SelectionAdjustment) bool

	// OnDragDone is called when the mouse drag gesture ends.
	OnDragDone()
}

// SelectionRegistrarExtended extends SelectionRegistrar with methods needed for
// the default selection modifier implementation.
//
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/foundation/foundation/src/commonMain/kotlin/androidx/compose/foundation/text/selection/SelectionRegistrar.kt
type SelectionRegistrarExtended interface {
	SelectionRegistrar

	// HasSelection returns true if there is an active selection for the given selectableId.
	HasSelection(selectableId int64) bool

	// NotifySelectionUpdateStart notifies the registrar that a selection gesture has started.
	NotifySelectionUpdateStart(
		layoutCoordinates LayoutCoordinates,
		startPosition geometry.Offset,
		adjustment SelectionAdjustment,
		isInTouchMode bool,
	)

	// NotifySelectionUpdate notifies the registrar that the selection has been updated.
	// Returns true if the update was consumed.
	NotifySelectionUpdate(
		layoutCoordinates LayoutCoordinates,
		previousPosition geometry.Offset,
		newPosition geometry.Offset,
		isStartHandle bool,
		adjustment SelectionAdjustment,
		isInTouchMode bool,
	) bool

	// NotifySelectionUpdateEnd notifies the registrar that a selection gesture has ended.
	NotifySelectionUpdateEnd()

	// MakeSelectionModifier creates a modifier for handling selection gestures.
	// This is platform-specific and may return different implementations on different platforms.
	MakeSelectionModifier(
		selectableId int64,
		layoutCoordinates func() LayoutCoordinates,
	) interface{} // Returns Modifier, but typed as interface{} to avoid circular dependencies
}
