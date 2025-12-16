package shape

import (
	"image"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/unit"
)

// https://developer.android.com/reference/kotlin/androidx/compose/ui/graphics/Shape

type Outline interface {
	Push(ops *op.Ops) clip.Stack
	Op(ops *op.Ops) clip.Op
	Path(ops *op.Ops) clip.PathSpec
}

type Shape interface {
	CreateOutline(size image.Point, metric unit.Metric) Outline
}

// RectangleShape
type rectangleShape struct{}

func (r rectangleShape) CreateOutline(size image.Point, metric unit.Metric) Outline {
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

var ShapeRectangle Shape = rectangleShape{}

// CircleShape
type circleShape struct{}

func (c circleShape) CreateOutline(size image.Point, metric unit.Metric) Outline {
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

var ShapeCircle Shape = circleShape{}

// RoundedCornerShape supports uniform radius via Radius, or per-corner radius
// via TopStart (NW), TopEnd (NE), BottomEnd (SE), BottomStart (SW).
// Per-corner fields take precedence when any are non-zero.
// This follows Jetpack Compose's RoundedCornerShape API.
type RoundedCornerShape struct {
	// Uniform radius applied to all corners (used when per-corner fields are all zero)
	Radius unit.Dp

	// Per-corner radius (following LTR layout direction):
	// TopStart = NW (top-left), TopEnd = NE (top-right)
	// BottomEnd = SE (bottom-right), BottomStart = SW (bottom-left)
	TopStart    unit.Dp
	TopEnd      unit.Dp
	BottomEnd   unit.Dp
	BottomStart unit.Dp
}

func (r RoundedCornerShape) CreateOutline(size image.Point, metric unit.Metric) Outline {
	// Determine if per-corner radius is being used
	hasPerCorner := r.TopStart > 0 || r.TopEnd > 0 || r.BottomEnd > 0 || r.BottomStart > 0

	var nw, ne, se, sw int
	if hasPerCorner {
		nw = metric.Dp(r.TopStart)
		ne = metric.Dp(r.TopEnd)
		se = metric.Dp(r.BottomEnd)
		sw = metric.Dp(r.BottomStart)
	} else {
		radius := metric.Dp(r.Radius)
		if radius == 0 {
			return rectOutline{clip.Rect{Max: size}}
		}
		nw, ne, se, sw = radius, radius, radius, radius
	}

	return rrectOutline{clip.RRect{
		Rect: image.Rectangle{Max: size},
		NW:   nw,
		NE:   ne,
		SE:   se,
		SW:   sw,
	}}
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

// CutCornerShape
type CutCornerShape struct {
	Radius unit.Dp
}

func (c CutCornerShape) CreateOutline(size image.Point, metric unit.Metric) Outline {
	radius := float32(metric.Dp(c.Radius))
	if radius <= 0 {
		return rectOutline{clip.Rect{Max: size}}
	}

	return &cutCornerOutline{
		size:   size,
		radius: radius,
	}
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
