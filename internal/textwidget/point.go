package textwidget

import (
	"gioui.org/f32"
)

// Point is a 2D float32 point without gioui dependency.
type Point struct {
	X, Y float32
}

// FromGioPoint converts a gioui f32.Point to the framework Point.
func FromGioPoint(p f32.Point) Point {
	return Point{X: p.X, Y: p.Y}
}

// ToGio converts this Point to a gioui f32.Point.
func (p Point) ToGio() f32.Point {
	return f32.Point{X: p.X, Y: p.Y}
}
