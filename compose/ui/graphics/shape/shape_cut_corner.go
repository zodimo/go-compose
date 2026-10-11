package shape

import (
	"image"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	gioUnit "gioui.org/unit"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

var cutCornerShapeUnspecified = &CutCornerShape{
	Radius: unit.DpCompanion.Unspecified(),
}
var _ Shape = (*CutCornerShape)(nil)

// CutCornerShape
type CutCornerShape struct {
	Radius unit.Dp
}

func (c *CutCornerShape) CreateOutline(size image.Point, metric gioUnit.Metric) Outline {
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
			Radius: unit.DpCompanion.TakeOrElse(otherCutCorner.Radius, c.Radius),
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

func (c *cutCornerOutline) generatePath(ops *op.Ops) clip.PathSpec {
	w, h := float32(c.size.X), float32(c.size.Y)
	r := c.radius

	var p clip.Path
	p.Begin(ops)
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

func (c *cutCornerOutline) Push(ops *op.Ops) clip.Stack {
	return clip.Outline{Path: c.generatePath(ops)}.Op().Push(ops)
}

func (c *cutCornerOutline) Op(ops *op.Ops) clip.Op {
	return clip.Outline{Path: c.generatePath(ops)}.Op()
}

func (c *cutCornerOutline) Path(ops *op.Ops) clip.PathSpec {
	return c.generatePath(ops)
}
