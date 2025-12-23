package geometry

import (
	"fmt"
	"math"

	"github.com/zodimo/go-compose/pkg/floatutils"
	"github.com/zodimo/go-compose/pkg/floatutils/lerp"
)

// Rect is an immutable, 2D, axis-aligned, floating-point rectangle whose coordinates are relative to a given origin.
type Rect struct {
	Left   float32
	Top    float32
	Right  float32
	Bottom float32
}

// RectZero is a rectangle with left, top, right, and bottom edges all at zero.
var RectZero = Rect{Left: 0, Top: 0, Right: 0, Bottom: 0}

// NewRect constructs a Rect from the given left, top, right, and bottom values.
func NewRect(left, top, right, bottom float32) Rect {
	return Rect{Left: left, Top: top, Right: right, Bottom: bottom}
}

// RectFromOffsetSize constructs a rectangle from its left and top edges as well as its width and height.
func RectFromOffsetSize(offset Offset, size Size) Rect {
	return Rect{
		Left:   offset.X(),
		Top:    offset.Y(),
		Right:  offset.X() + size.Width(),
		Bottom: offset.Y() + size.Height(),
	}
}

// RectFromTwoOffsets constructs the smallest rectangle that encloses the given offsets, treating them as vectors from the origin.
func RectFromTwoOffsets(topLeft, bottomRight Offset) Rect {
	return Rect{
		Left:   topLeft.X(),
		Top:    topLeft.Y(),
		Right:  bottomRight.X(),
		Bottom: bottomRight.Y(),
	}
}

// RectFromCircle constructs a rectangle that bounds the given circle.
func RectFromCircle(center Offset, radius float32) Rect {
	return Rect{
		Left:   center.X() - radius,
		Top:    center.Y() - radius,
		Right:  center.X() + radius,
		Bottom: center.Y() + radius,
	}
}

// Width returns the distance between the left and right edges of this rectangle.
func (r Rect) Width() float32 {
	return r.Right - r.Left
}

// Height returns the distance between the top and bottom edges of this rectangle.
func (r Rect) Height() float32 {
	return r.Bottom - r.Top
}

// Size returns the distance between the upper-left corner and the lower-right corner of this rectangle.
func (r Rect) Size() Size {
	return NewSize(r.Width(), r.Height())
}

// IsInfinite returns true if any of the coordinates of this rectangle are equal to positive infinity.
func (r Rect) IsInfinite() bool {
	return math.IsInf(float64(r.Left), 1) ||
		math.IsInf(float64(r.Top), 1) ||
		math.IsInf(float64(r.Right), 1) ||
		math.IsInf(float64(r.Bottom), 1)
}

// IsFinite returns true if all coordinates of this rectangle are finite.
func (r Rect) IsFinite() bool {
	return !math.IsInf(float64(r.Left), 0) && !math.IsNaN(float64(r.Left)) &&
		!math.IsInf(float64(r.Top), 0) && !math.IsNaN(float64(r.Top)) &&
		!math.IsInf(float64(r.Right), 0) && !math.IsNaN(float64(r.Right)) &&
		!math.IsInf(float64(r.Bottom), 0) && !math.IsNaN(float64(r.Bottom))
}

// IsEmpty returns true if this rectangle encloses a non-zero area. Negative areas are considered empty.
func (r Rect) IsEmpty() bool {
	return r.Left >= r.Right || r.Top >= r.Bottom
}

// Translate returns a new rectangle translated by the given offset.
func (r Rect) Translate(offset Offset) Rect {
	return Rect{
		Left:   r.Left + offset.X(),
		Top:    r.Top + offset.Y(),
		Right:  r.Right + offset.X(),
		Bottom: r.Bottom + offset.Y(),
	}
}

// TranslateXY returns a new rectangle with translateX added to the x components and translateY added to the y components.
func (r Rect) TranslateXY(translateX, translateY float32) Rect {
	return Rect{
		Left:   r.Left + translateX,
		Top:    r.Top + translateY,
		Right:  r.Right + translateX,
		Bottom: r.Bottom + translateY,
	}
}

// Inflate returns a new rectangle with edges moved outwards by the given delta.
func (r Rect) Inflate(delta float32) Rect {
	return Rect{
		Left:   r.Left - delta,
		Top:    r.Top - delta,
		Right:  r.Right + delta,
		Bottom: r.Bottom + delta,
	}
}

// Deflate returns a new rectangle with edges moved inwards by the given delta.
func (r Rect) Deflate(delta float32) Rect {
	return r.Inflate(-delta)
}

// Intersect returns a new rectangle that is the intersection of the given rectangle and this rectangle.
func (r Rect) Intersect(other Rect) Rect {
	return Rect{
		Left:   float32(math.Max(float64(r.Left), float64(other.Left))),
		Top:    float32(math.Max(float64(r.Top), float64(other.Top))),
		Right:  float32(math.Min(float64(r.Right), float64(other.Right))),
		Bottom: float32(math.Min(float64(r.Bottom), float64(other.Bottom))),
	}
}

// IntersectCoords returns a new rectangle that is the intersection of the given rectangle coordinates and this rectangle.
func (r Rect) IntersectCoords(otherLeft, otherTop, otherRight, otherBottom float32) Rect {
	return Rect{
		Left:   float32(math.Max(float64(r.Left), float64(otherLeft))),
		Top:    float32(math.Max(float64(r.Top), float64(otherTop))),
		Right:  float32(math.Min(float64(r.Right), float64(otherRight))),
		Bottom: float32(math.Min(float64(r.Bottom), float64(otherBottom))),
	}
}

// Overlaps returns whether `other` has a nonzero area of overlap with this rectangle.
func (r Rect) Overlaps(other Rect) bool {
	return r.Left < other.Right && other.Left < r.Right &&
		r.Top < other.Bottom && other.Top < r.Bottom
}

// MinDimension returns the lesser of the magnitudes of the width and the height of this rectangle.
func (r Rect) MinDimension() float32 {
	return float32(math.Min(math.Abs(float64(r.Width())), math.Abs(float64(r.Height()))))
}

// MaxDimension returns the greater of the magnitudes of the width and the height of this rectangle.
func (r Rect) MaxDimension() float32 {
	return float32(math.Max(math.Abs(float64(r.Width())), math.Abs(float64(r.Height()))))
}

// TopLeft returns the offset to the intersection of the top and left edges of this rectangle.
func (r Rect) TopLeft() Offset {
	return NewOffset(r.Left, r.Top)
}

// TopCenter returns the offset to the center of the top edge of this rectangle.
func (r Rect) TopCenter() Offset {
	return NewOffset(r.Left+r.Width()/2, r.Top)
}

// TopRight returns the offset to the intersection of the top and right edges of this rectangle.
func (r Rect) TopRight() Offset {
	return NewOffset(r.Right, r.Top)
}

// CenterLeft returns the offset to the center of the left edge of this rectangle.
func (r Rect) CenterLeft() Offset {
	return NewOffset(r.Left, r.Top+r.Height()/2)
}

// Center returns the offset to the point halfway between the left and right and the top and bottom edges.
func (r Rect) Center() Offset {
	return NewOffset(r.Left+r.Width()/2, r.Top+r.Height()/2)
}

// CenterRight returns the offset to the center of the right edge of this rectangle.
func (r Rect) CenterRight() Offset {
	return NewOffset(r.Right, r.Top+r.Height()/2)
}

// BottomLeft returns the offset to the intersection of the bottom and left edges of this rectangle.
func (r Rect) BottomLeft() Offset {
	return NewOffset(r.Left, r.Bottom)
}

// BottomCenter returns the offset to the center of the bottom edge of this rectangle.
func (r Rect) BottomCenter() Offset {
	return NewOffset(r.Left+r.Width()/2, r.Bottom)
}

// BottomRight returns the offset to the intersection of the bottom and right edges of this rectangle.
func (r Rect) BottomRight() Offset {
	return NewOffset(r.Right, r.Bottom)
}

// Contains returns whether the point specified by the given offset lies between the left and right and the top and bottom edges.
// Rectangles include their top and left edges but exclude their bottom and right edges.
func (r Rect) Contains(offset Offset) bool {
	return offset.X() >= r.Left && offset.X() < r.Right && offset.Y() >= r.Top && offset.Y() < r.Bottom
}

// String returns a string representation of the object.
func (r Rect) String() string {
	return fmt.Sprintf("Rect.fromLTRB(%.1f, %.1f, %.1f, %.1f)", r.Left, r.Top, r.Right, r.Bottom)
}

// Equal checks equality with another Rect.
func (r Rect) Equal(other Rect) bool {
	return floatutils.Float32Equals(r.Left, other.Left, floatutils.Float32EqualityThreshold) &&
		floatutils.Float32Equals(r.Top, other.Top, floatutils.Float32EqualityThreshold) &&
		floatutils.Float32Equals(r.Right, other.Right, floatutils.Float32EqualityThreshold) &&
		floatutils.Float32Equals(r.Bottom, other.Bottom, floatutils.Float32EqualityThreshold)
}

// LerpRect linearly interpolates between two rectangles.
func LerpRect(start, stop Rect, fraction float32) Rect {
	return Rect{
		Left:   lerp.Between32(start.Left, stop.Left, fraction),
		Top:    lerp.Between32(start.Top, stop.Top, fraction),
		Right:  lerp.Between32(start.Right, stop.Right, fraction),
		Bottom: lerp.Between32(start.Bottom, stop.Bottom, fraction),
	}
}
