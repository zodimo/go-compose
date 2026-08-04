package shape

import (
	"image"

	"gioui.org/f32"
	"gioui.org/op/clip"
	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/internal/clipconvert"
)

var CutCornerShapeUnspecified = &CutCornerShape{
	Radius: unit.DpUnspecified,
}
var _ Shape = (*CutCornerShape)(nil)

// CutCornerShape
type CutCornerShape struct {
	Radius unit.Dp
}

func (c *CutCornerShape) CreateOutline(size image.Point, metric Metric) Outline {
	radius := float32(c.Radius) * metric.PxPerDp
	if radius <= 0 {
		return rectOutline{clip.Rect{Max: size}}
	}

	return &cutCornerOutline{
		size:   size,
		radius: radius,
	}
}

// Sentinel private functions
func (c *CutCornerShape) mergeShape(other Shape) Shape {
	if otherCutCorner, ok := other.(*CutCornerShape); ok {
		return &CutCornerShape{
			Radius: otherCutCorner.Radius.TakeOrElse(c.Radius),
		}

	}

	// should we panic here ?
	return &CutCornerShape{
		Radius: c.Radius,
	}
}
func (c *CutCornerShape) sameShape(other Shape) bool {
	if _, ok := other.(*CutCornerShape); ok {
		return true
	}
	return false
}
func (c *CutCornerShape) semanticEqualShape(other Shape) bool {
	if otherShape, ok := other.(*CutCornerShape); ok {
		return otherShape.Radius == c.Radius
	}
	return false
}
func (c *CutCornerShape) copyShape(options ...ShapeOption) Shape {
	copy := *c
	for _, option := range options {
		option(&copy)
	}
	return &copy
}
func (c *CutCornerShape) stringShape() string {
	return "CutCornerShape{Radius: " + c.Radius.String() + "}"
}

type cutCornerOutline struct {
	size   image.Point
	radius float32
}

func (c *cutCornerOutline) generatePath(ops *clipconvert.Ops) clip.PathSpec {
	gtx := ops.ToGio()
	w, h := float32(c.size.X), float32(c.size.Y)
	r := c.radius

	var p clip.Path
	p.Begin(gtx)
	p.MoveTo(f32.Pt(r, 0))
	p.LineTo(f32.Pt(w-r, 0))
	p.LineTo(f32.Pt(w, r))
	p.LineTo(f32.Pt(w, h-r))
	p.LineTo(f32.Pt(w-r, h))
	p.LineTo(f32.Pt(r, h))
	p.LineTo(f32.Pt(0, h-r))
	p.LineTo(f32.Pt(0, r))
	p.Close()
	return p.End()
}

func (c *cutCornerOutline) Push(ops *clipconvert.Ops) clipconvert.Stack {
	return clipconvert.NewStack(clip.Outline{Path: c.generatePath(ops)}.Op().Push(ops.ToGio()))
}

func (c *cutCornerOutline) ClipOp(ops *clipconvert.Ops) clipconvert.ClipOp {
	return clipconvert.NewClipOp(clip.Outline{Path: c.generatePath(ops)}.Op())
}

func (c *cutCornerOutline) Path(ops *clipconvert.Ops) clipconvert.PathSpec {
	return clipconvert.NewPathSpec(c.generatePath(ops))
}
