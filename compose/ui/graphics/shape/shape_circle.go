package shape

import (
	"image"

	"gioui.org/op"
	"gioui.org/op/clip"
	gioUnit "gioui.org/unit"
)

// Deprecated: Use CircleShape instead
var ShapeCircle Shape = circleShape{}

// CircleShape is a shape describing a circle.
var CircleShape Shape = circleShape{}

// CircleShape
type circleShape struct{}

func (c circleShape) CreateOutline(size image.Point, metric gioUnit.Metric) Outline {
	return ellipseOutline{clip.Ellipse{Max: size}}
}

type ellipseOutline struct {
	clip.Ellipse
}

func (e ellipseOutline) Push(ops *op.Ops) clip.Stack {
	return e.Ellipse.Push(ops)
}

func (e ellipseOutline) Op(ops *op.Ops) clip.Op {
	return e.Ellipse.Op(ops)
}

// Ellipse.Path takes ops argument? No, Ellipse.Path(ops) returns PathSpec.
// checking docs/source... Ellipse.Path(ops) -> PathSpec.
func (e ellipseOutline) Path(ops *op.Ops) clip.PathSpec {
	return e.Ellipse.Path(ops)
}
