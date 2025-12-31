package selection

import (
	"github.com/zodimo/go-compose/compose/ui/geometry"
	"github.com/zodimo/go-compose/compose/ui/layout"
	"github.com/zodimo/go-compose/compose/ui/next/text"
)

// Selectable provides Selection information for a composable to SelectionContainer.
// Composables that can be selected should subscribe to SelectionRegistrar using this interface.
//
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/foundation/foundation/src/commonMain/kotlin/androidx/compose/foundation/text/selection/Selectable.kt
type Selectable interface {
	// SelectableId returns an ID used by SelectionRegistrar to identify this Selectable.
	// This value should not be InvalidSelectableId.
	SelectableId() int64

	// AppendSelectableInfoToBuilder adds SelectableInfo representing this Selectable
	// to the SelectionLayoutBuilder.
	AppendSelectableInfoToBuilder(builder SelectionLayoutBuilder)

	// GetSelectAllSelection returns selectAll Selection information for this selectable.
	// Returns nil if no selection can be provided.
	GetSelectAllSelection() *Selection

	// GetHandlePosition returns the Offset of a SelectionHandle.
	// isStartHandle is true for the start handle, false for the end handle.
	GetHandlePosition(selection Selection, isStartHandle bool) geometry.Offset

	// GetLayoutCoordinates returns the LayoutCoordinates of the Selectable.
	// This could be nil if called before composing.
	GetLayoutCoordinates() layout.LayoutCoordinates

	// TextLayoutResult returns the TextLayoutResult of the selectable.
	// This could be nil if called before composing.
	TextLayoutResult() *text.TextLayoutResult

	// GetText returns the text content as AnnotatedString of the Selectable.
	GetText() text.AnnotatedString

	// GetBoundingBox returns the bounding box of the character for given offset.
	// Returns Rect.Zero if the selectable is empty.
	GetBoundingBox(offset int) geometry.Rect

	// GetLineLeft returns the left x coordinate of the line for the given offset.
	GetLineLeft(offset int) float32

	// GetLineRight returns the right x coordinate of the line for the given offset.
	GetLineRight(offset int) float32

	// GetCenterYForOffset returns the center y coordinate of the line
	// on which the specified text offset appears.
	GetCenterYForOffset(offset int) float32

	// GetRangeOfLineContaining returns the offsets of the start and end of the line
	// containing the given offset, or TextRange.Zero if the selectable is empty.
	GetRangeOfLineContaining(offset int) text.TextRange

	// GetLastVisibleOffset returns the last visible character's offset.
	// Some lines can be hidden due to maxLines or Constraints.maxHeight.
	GetLastVisibleOffset() int

	// GetLineHeight returns the text line height for the given offset.
	GetLineHeight(offset int) float32
}

// SelectionLayoutBuilder is used to build selection layout information.
// This is a placeholder interface that will be expanded as needed.
//
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/foundation/foundation/src/commonMain/kotlin/androidx/compose/foundation/text/selection/SelectionLayout.kt
type SelectionLayoutBuilder interface {
	// TODO: Add methods as needed for selection layout building
}
