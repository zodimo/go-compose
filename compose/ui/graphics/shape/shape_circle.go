package shape

import (
	"image"

	"gioui.org/op/clip"
	"github.com/zodimo/go-compose/internal/clipconvert"
)

// Deprecated: Use CircleShape instead
var ShapeCircle Shape = &circleShape{}

// CircleShape is a shape describing a circle.
var CircleShape Shape = &circleShape{}

// CircleShape
type circleShape struct{}

func (c *circleShape) CreateOutline(size image.Point, metric Metric) Outline {
	radius := min(size.X, size.Y) / 2
	return rrectOutline{clip.RRect{
		Rect: image.Rectangle{Max: size},
		SE:   radius,
		SW:   radius,
		NW:   radius,
		NE:   radius,
	}}
}

// Sentinel Private methods
func (c *circleShape) mergeShape(other Shape) Shape {
	if otherCircle, ok := other.(*circleShape); ok {
		return otherCircle
	}
	return c
}
func (c *circleShape) sameShape(other Shape) bool {
	if _, ok := other.(*circleShape); ok {
		return true
	}
	return false
}
func (c *circleShape) semanticEqualShape(other Shape) bool {
	if _, ok := other.(*circleShape); ok {
		return true
	}
	return false
}
func (c *circleShape) copyShape(options ...ShapeOption) Shape {
	copy := *c
	if len(options) > 0 {
		// should we panic here ?
		return &copy
	}
	return &copy
}
func (c *circleShape) stringShape() string {
	return "CircleShape"
}

type ellipseOutline struct {
	clip.Ellipse
}

func (e ellipseOutline) Push(ops *clipconvert.Ops) clipconvert.Stack {
	return clipconvert.NewStack(e.Ellipse.Push(ops.ToGio()))
}

func (e ellipseOutline) ClipOp(ops *clipconvert.Ops) clipconvert.ClipOp {
	return clipconvert.NewClipOp(e.Ellipse.Op(ops.ToGio()))
}

func (e ellipseOutline) Path(ops *clipconvert.Ops) clipconvert.PathSpec {
	return clipconvert.NewPathSpec(e.Ellipse.Path(ops.ToGio()))
}
