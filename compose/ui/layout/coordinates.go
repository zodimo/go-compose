package layout

import (
	"github.com/zodimo/go-compose/compose/ui/geometry"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

// LayoutCoordinates provides access to the position and size of a layout.
//
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui/src/commonMain/kotlin/androidx/compose/ui/layout/LayoutCoordinates.kt
type LayoutCoordinates interface {
	// IsAttached returns true if the layout is currently attached to the composition.
	IsAttached() bool

	// Size returns the size of this layout in pixels.
	Size() unit.IntSize

	// PositionInRoot returns the position of this layout relative to the root.
	PositionInRoot() geometry.Offset

	// PositionInWindow returns the position of this layout in window coordinates.
	PositionInWindow() geometry.Offset

	// LocalPositionOf converts a position from another LayoutCoordinates to this one.
	LocalPositionOf(sourceCoordinates LayoutCoordinates, relativeToSource geometry.Offset) geometry.Offset

	// VisibleBounds returns the visible bounds of this layout in local coordinates.
	VisibleBounds() geometry.Rect

	// BoundsInWindow returns the bounds of this layout in window coordinates.
	BoundsInWindow() geometry.Rect
}
