package graphics

import (
	"image"

	"github.com/zodimo/go-compose/compose/ui/geometry"
)

// PathFillType determines how the interior of a path is calculated.
//
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-graphics/src/commonMain/kotlin/androidx/compose/ui/graphics/PathFillType.kt
type PathFillType int

const (
	// PathFillTypeNonZero specifies that the path uses the non-zero winding rule.
	PathFillTypeNonZero PathFillType = iota

	// PathFillTypeEvenOdd specifies that the path uses the even-odd winding rule.
	PathFillTypeEvenOdd
)

// String returns the string representation of the PathFillType.
func (f PathFillType) String() string {
	switch f {
	case PathFillTypeNonZero:
		return "NonZero"
	case PathFillTypeEvenOdd:
		return "EvenOdd"
	default:
		return "Unknown"
	}
}

// PathDirection specifies how closed shapes are wound when added to a path.
type PathDirection int

const (
	// PathDirectionCounterClockwise means the shape is wound in counter-clockwise order.
	PathDirectionCounterClockwise PathDirection = iota

	// PathDirectionClockwise means the shape is wound in clockwise order.
	PathDirectionClockwise
)

// String returns the string representation of the PathDirection.
func (d PathDirection) String() string {
	switch d {
	case PathDirectionCounterClockwise:
		return "CounterClockwise"
	case PathDirectionClockwise:
		return "Clockwise"
	default:
		return "Unknown"
	}
}

// Path defines an interface for 2D geometric paths.
//
// A path contains zero or more subpaths, each of which contains zero or more
// straight line segments and/or bezier curve segments.
//
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-graphics/src/commonMain/kotlin/androidx/compose/ui/graphics/Path.kt
type Path interface {
	// FillType returns the path fill type.
	FillType() PathFillType

	// SetFillType sets the path fill type.
	SetFillType(fillType PathFillType)

	// IsConvex returns true if the path is convex.
	IsConvex() bool

	// IsEmpty returns true if the path is empty (contains no lines or curves).
	IsEmpty() bool

	// MoveTo starts a new subpath at the given coordinate.
	MoveTo(x, y float32)

	// RelativeMoveTo starts a new subpath at the given offset from the current point.
	RelativeMoveTo(dx, dy float32)

	// LineTo adds a straight line segment from the current point to the given point.
	LineTo(x, y float32)

	// RelativeLineTo adds a straight line segment from the current point to the point
	// at the given offset from the current point.
	RelativeLineTo(dx, dy float32)

	// QuadraticTo adds a quadratic bezier segment that curves from the current point
	// to (x2, y2), using (x1, y1) as the control point.
	QuadraticTo(x1, y1, x2, y2 float32)

	// RelativeQuadraticTo adds a quadratic bezier segment using relative coordinates.
	RelativeQuadraticTo(dx1, dy1, dx2, dy2 float32)

	// CubicTo adds a cubic bezier segment that curves from the current point to (x3, y3),
	// using (x1, y1) and (x2, y2) as control points.
	CubicTo(x1, y1, x2, y2, x3, y3 float32)

	// RelativeCubicTo adds a cubic bezier segment using relative coordinates.
	RelativeCubicTo(dx1, dy1, dx2, dy2, dx3, dy3 float32)

	// ArcTo adds an arc segment from the current point.
	ArcTo(rect geometry.Rect, startAngleDegrees, sweepAngleDegrees float32, forceMoveTo bool)

	// AddRect adds a rectangle as a new subpath.
	AddRect(rect geometry.Rect, direction PathDirection)

	// AddOval adds an oval (ellipse) as a new subpath.
	AddOval(oval geometry.Rect, direction PathDirection)

	// AddArc adds an arc segment as a new subpath.
	AddArc(oval geometry.Rect, startAngleDegrees, sweepAngleDegrees float32)

	// AddPath adds another path to this path with an optional offset.
	AddPath(path Path, offset geometry.Offset)

	// Close closes the current subpath.
	Close()

	// Reset clears all subpaths from the path.
	Reset()

	// Rewind clears lines and curves but keeps internal data structure for faster reuse.
	Rewind()

	// Translate translates all segments by the given offset.
	Translate(offset geometry.Offset)

	// GetBounds computes the bounds of the control points of the path.
	GetBounds() geometry.Rect

	// Op performs a boolean operation on two paths.
	Op(path1, path2 Path, operation PathOperation) bool

	// Contains returns true if the point is inside the path according to the fill type.
	// Uses winding or even-odd fill rule based on path's fill type.
	Contains(point image.Point) bool
}

// PathOperation specifies the boolean operation to perform on two paths.
type PathOperation int

const (
	// PathOperationDifference subtracts path2 from path1.
	PathOperationDifference PathOperation = iota

	// PathOperationIntersect returns the intersection of path1 and path2.
	PathOperationIntersect

	// PathOperationUnion returns the union of path1 and path2.
	PathOperationUnion

	// PathOperationXor returns the exclusive-or of path1 and path2.
	PathOperationXor

	// PathOperationReverseDifference subtracts path1 from path2.
	PathOperationReverseDifference
)

// String returns the string representation of the PathOperation.
func (o PathOperation) String() string {
	switch o {
	case PathOperationDifference:
		return "Difference"
	case PathOperationIntersect:
		return "Intersect"
	case PathOperationUnion:
		return "Union"
	case PathOperationXor:
		return "Xor"
	case PathOperationReverseDifference:
		return "ReverseDifference"
	default:
		return "Unknown"
	}
}
