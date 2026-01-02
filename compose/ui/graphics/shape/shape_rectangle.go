package shape

import (
	"image"

	"gioui.org/op"
	"gioui.org/op/clip"
	gioUnit "gioui.org/unit"
)

var ShapeRectangle Shape = rectangleShape{}

// RectangleShape
type rectangleShape struct{}

func (r rectangleShape) CreateOutline(size image.Point, metric gioUnit.Metric) Outline {
	return rectOutline{clip.Rect{Max: size}}
}

type rectOutline struct {
	clip.Rect
}

func (r rectOutline) Push(ops *op.Ops) clip.Stack {
	return r.Rect.Push(ops)
}

func (r rectOutline) Op(ops *op.Ops) clip.Op {
	return r.Rect.Op()
}

func (r rectOutline) Path(ops *op.Ops) clip.PathSpec {
	return r.Rect.Path()
}

type rrectOutline struct {
	clip.RRect
}

func (r rrectOutline) Push(ops *op.Ops) clip.Stack {
	return r.RRect.Push(ops)
}

func (r rrectOutline) Op(ops *op.Ops) clip.Op {
	return r.RRect.Op(ops)
}

func (r rrectOutline) Path(ops *op.Ops) clip.PathSpec {
	return r.RRect.Path(ops)
}
