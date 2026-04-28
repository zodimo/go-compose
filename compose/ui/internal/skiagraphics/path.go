package skiagraphics

import (
	"image"

	"github.com/zodimo/gio-skia/skia"
	"github.com/zodimo/go-compose/compose/ui/geometry"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-skia-support/skia/enums"
	"github.com/zodimo/go-skia-support/skia/models"
)

func NewPath() graphics.Path {
	return &pathWrapper{path: skia.NewPath()}
}

var _ graphics.Path = (*pathWrapper)(nil)

type pathWrapper struct {
	path skia.SkPath
}

// FillType returns the path fill type.
func (p *pathWrapper) FillType() graphics.PathFillType {
	return SkiaFillPathTypeToGraphicsFillType(p.path.FillType())
}

// SetFillType sets the path fill type.
func (p *pathWrapper) SetFillType(fillType graphics.PathFillType) {
	p.path.SetFillType(GraphicsFillPathTypeToSkiaFillPathType(fillType))
}

// IsConvex returns true if the path is convex.
func (p *pathWrapper) IsConvex() bool {
	return p.path.IsConvex()
}

// IsEmpty returns true if the path is empty (contains no lines or curves).
func (p *pathWrapper) IsEmpty() bool {
	return p.path.IsEmpty()
}

// MoveTo starts a new subpath at the given coordinate.
func (p *pathWrapper) MoveTo(x, y float32) {
	p.path.MoveTo(x, y)
}

// RelativeMoveTo starts a new subpath at the given offset from the current point.
func (p *pathWrapper) RelativeMoveTo(dx, dy float32) {
	p.path.RMoveTo(dx, dy)
}

// LineTo adds a straight line segment from the current point to the given point.
func (p *pathWrapper) LineTo(x, y float32) {
	p.path.LineTo(x, y)
}

// RelativeLineTo adds a straight line segment from the current point to the point
// at the given offset from the current point.
func (p *pathWrapper) RelativeLineTo(dx, dy float32) {
	p.path.RLineTo(dx, dy)
}

// QuadraticTo adds a quadratic bezier segment that curves from the current point
// to (x2, y2), using (x1, y1) as the control point.
func (p *pathWrapper) QuadraticTo(x1, y1, x2, y2 float32) {
	p.path.QuadTo(x1, y1, x2, y2)
}

// RelativeQuadraticTo adds a quadratic bezier segment using relative coordinates.
func (p *pathWrapper) RelativeQuadraticTo(dx1, dy1, dx2, dy2 float32) {
	p.path.RQuadTo(dx1, dy1, dx2, dy2)
}

// CubicTo adds a cubic bezier segment that curves from the current point to (x3, y3),
// using (x1, y1) and (x2, y2) as control points.
func (p *pathWrapper) CubicTo(x1, y1, x2, y2, x3, y3 float32) {
	p.path.CubicTo(x1, y1, x2, y2, x3, y3)
}

// RelativeCubicTo adds a cubic bezier segment using relative coordinates.
func (p *pathWrapper) RelativeCubicTo(dx1, dy1, dx2, dy2, dx3, dy3 float32) {
	p.path.RCubicTo(dx1, dy1, dx2, dy2, dx3, dy3)
}

// ArcTo adds an arc segment from the current point.
func (p *pathWrapper) ArcTo(rect geometry.Rect, startAngleDegrees, sweepAngleDegrees float32, forceMoveTo bool) {
	p.path.ArcTo(geometryRectToSkiaRect(rect), startAngleDegrees, sweepAngleDegrees, forceMoveTo)
}

// AddRect adds a rectangle as a new subpath.
func (p *pathWrapper) AddRect(rect geometry.Rect, direction graphics.PathDirection) {
	p.path.AddRect(geometryRectToSkiaRect(rect), graphicsPathDirectionToSkiaPathDirection(direction), 0)
}

// AddOval adds an oval (ellipse) as a new subpath.
func (p *pathWrapper) AddOval(oval geometry.Rect, direction graphics.PathDirection) {
	p.path.AddOval(geometryRectToSkiaRect(oval), graphicsPathDirectionToSkiaPathDirection(direction))
}

// AddArc adds an arc segment as a new subpath.
func (p *pathWrapper) AddArc(oval geometry.Rect, startAngleDegrees, sweepAngleDegrees float32) {
	p.path.AddArc(geometryRectToSkiaRect(oval), startAngleDegrees, sweepAngleDegrees)
}

// AddPath adds another path to this path with an optional offset.
func (p *pathWrapper) AddPath(path graphics.Path, offset geometry.Offset) {
	addMode := enums.AddPathModeAppend
	p.path.AddPath(path.(*pathWrapper).path, offset.X(), offset.Y(), addMode)
}

// Close closes the current subpath.
func (p *pathWrapper) Close() {
	p.path.Close()
}

// Reset clears all subpaths from the path.
func (p *pathWrapper) Reset() {
	p.path.Reset()
}

// Rewind clears lines and curves but keeps internal data structure for faster reuse.
func (p *pathWrapper) Rewind() {
	p.path.Reset()
}

// Translate translates all segments by the given offset.
func (p *pathWrapper) Translate(offset geometry.Offset) {
	p.path.Offset(offset.X(), offset.Y())
}

// GetBounds computes the bounds of the control points of the path.
func (p *pathWrapper) GetBounds() geometry.Rect {
	return skiaRectToGeometryRect(p.path.Bounds())
}

// Op performs a boolean operation on two paths.
func (p *pathWrapper) Op(path1, path2 graphics.Path, operation graphics.PathOperation) bool {
	p1 := path1.(*pathWrapper)
	p2 := path2.(*pathWrapper)
	result := p1.path.Op(p2.path, graphicsPathOperationToSkiaPathOp(operation))
	p1.path = result
	return true
}

// Contains returns true if the point is inside the path according to the fill type.
// Uses winding or even-odd fill rule based on path's fill type.
func (p *pathWrapper) Contains(point image.Point) bool {
	return p.path.Contains(models.Point{X: float32(point.X), Y: float32(point.Y)})
}
