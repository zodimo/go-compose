package shape

import (
	"image"

	"gioui.org/op/clip"
	"github.com/zodimo/go-compose/internal/clipconvert"
)

var ShapeRectangle Shape = &rectangleShape{}

// RectangleShape
type rectangleShape struct{}

func (r *rectangleShape) CreateOutline(size image.Point, metric Metric) Outline {
	return rectOutline{clip.Rect{Max: size}}
}

func (r *rectangleShape) mergeShape(other Shape) Shape {
	return other
}
func (r *rectangleShape) sameShape(other Shape) bool {
	if _, ok := other.(*rectangleShape); ok {
		return true
	}
	return false
}
func (r *rectangleShape) semanticEqualShape(other Shape) bool {
	if _, ok := other.(*rectangleShape); ok {
		return true
	}
	return false
}
func (r *rectangleShape) copyShape(options ...ShapeOption) Shape {
	copy := *r
	if len(options) > 0 {
		panic("copyShape: options not supported")
	}
	return &copy
}

func (r *rectangleShape) stringShape() string {
	return "RectangleShape"
}

type rectOutline struct {
	clip.Rect
}

func (r rectOutline) Push(ops *clipconvert.Ops) clipconvert.Stack {
	return clipconvert.NewStack(r.Rect.Push(ops.ToGio()))
}

func (r rectOutline) ClipOp(ops *clipconvert.Ops) clipconvert.ClipOp {
	return clipconvert.NewClipOp(r.Rect.Op())
}

func (r rectOutline) Path(ops *clipconvert.Ops) clipconvert.PathSpec {
	return clipconvert.NewPathSpec(r.Rect.Path())
}

type rrectOutline struct {
	clip.RRect
}

func (r rrectOutline) Push(ops *clipconvert.Ops) clipconvert.Stack {
	return clipconvert.NewStack(r.RRect.Push(ops.ToGio()))
}

func (r rrectOutline) ClipOp(ops *clipconvert.Ops) clipconvert.ClipOp {
	return clipconvert.NewClipOp(r.RRect.Op(ops.ToGio()))
}

func (r rrectOutline) Path(ops *clipconvert.Ops) clipconvert.PathSpec {
	return clipconvert.NewPathSpec(r.RRect.Path(ops.ToGio()))
}
