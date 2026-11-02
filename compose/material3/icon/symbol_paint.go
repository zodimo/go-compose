package icon

import (
	"encoding/binary"
	"gioui.org/op"
	"hash/fnv"
	"image"
	"image/color"
	"math"
	"reflect"

	"gioui.org/f32"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"github.com/zodimo/go-compose/compose/ui/geometry"
	"github.com/zodimo/go-compose/compose/ui/graphics"
)

// fill describes how an icon's geometry is painted.
type fill struct {
	// brush, when non-nil, wins over color: a gradient is painted across the
	// icon box.
	brush graphics.Brush
	// color is used when brush is nil.
	color color.NRGBA
}

// paint fills shape inside gtx.
//
// Gio has no arbitrary-stop gradient op, only a two-stop LinearGradientOp, so a
// multi-stop brush is painted as one two-stop band per adjacent stop pair. Bands
// are clamped to the shape's bounds, so the result is an ordinary N-stop
// gradient rather than an approximation.
func (f fill) paint(ops *op.Ops, box image.Rectangle, shape clip.Op) {
	if f.brush == nil {
		paint.FillShape(ops, f.color, shape)
		return
	}

	if !f.paintGradient(ops, box, shape) {
		// Unsupported brush kind: fall back to the resolved color so the icon is
		// still visible rather than silently dropped.
		paint.FillShape(ops, f.color, shape)
	}
}

// paintGradient paints a linear gradient brush and reports whether it handled
// the brush. Rectangle-based brushes (solid color) are handled as a plain fill.
func (f fill) paintGradient(ops *op.Ops, box image.Rectangle, shape clip.Op) bool {
	defer paint.PushOpacity(ops, 1).Pop()

	switch brush := f.brush.(type) {
	case *graphics.SolidColor:
		paint.FillShape(ops, graphics.ColorToNRGBA(brush.Value), shape)
		return true
	case graphics.LinearGradient:
		f.paintLinearGradient(ops, box, shape, brush.Start, brush.End, brush.Colors, brush.Stops)
		return true
	case *graphics.LinearGradient:
		f.paintLinearGradient(ops, box, shape, brush.Start, brush.End, brush.Colors, brush.Stops)
		return true
	default:
		return false
	}
}

// paintLinearGradient paints a multi-stop linear gradient by drawing one band per
// adjacent stop pair, clipped to the shape and to that band's slice.
//
// Clipping to the band is what makes the result an N-stop gradient rather than
// "last band wins": Gio's LinearGradientOp extrapolates beyond its two stops
// with the endpoint colours, so an unclipped band paints its start colour over
// the whole shape and the final band erases every earlier one.
func (f fill) paintLinearGradient(
	ops *op.Ops,
	box image.Rectangle,
	shape clip.Op,
	start, end geometry.Offset,
	colors []graphics.Color,
	stops []float32,
) {
	n := len(colors)
	if n == 0 {
		return
	}

	// A single color, or no stops at all, is just a solid fill.
	if n == 1 || len(stops) < n {
		paint.FillShape(ops, graphics.ColorToNRGBA(colors[0]), shape)
		return
	}

	from := resolveGradientPoint(start, box)
	to := resolveGradientPoint(end, box)
	if from == to {
		paint.FillShape(ops, graphics.ColorToNRGBA(colors[0]), shape)
		return
	}

	// Clip every band to the icon shape first, then to its own slice.
	shapeStack := shape.Push(ops)
	defer shapeStack.Pop()

	// The first and last bands extend past the gradient's endpoints: with clamp
	// tiling, the endpoint colour covers everything beyond them, which is how a
	// gradient whose axis is narrower than the box still fills the box.
	axisLen := float32(math.Hypot(float64(to.X-from.X), float64(to.Y-from.Y)))
	overrun := (float32(box.Dx()+box.Dy())+1)/axisLen + 1

	for i := 0; i+1 < n; i++ {
		t0, t1 := stops[i], stops[i+1]
		if t1 < t0 {
			t0, t1 = t1, t0
		}

		// The clip covers the (possibly extended) band slice...
		clipFrom, clipTo := t0, t1
		if i == 0 {
			clipFrom -= overrun
		}
		if i+2 == n {
			clipTo += overrun
		}
		// ...while the gradient's own stops stay inside the axis, so each band
		// interpolates between its two colours rather than across the whole axis.
		stopFrom := lerpPoint(from, to, clampAxis(t0, 0, 1))
		stopTo := lerpPoint(from, to, clampAxis(t1, 0, 1))

		bandStack := gradientBandClip(ops, from, to, clipFrom, clipTo, box).Push(ops)

		if stopFrom == stopTo {
			// A zero-width band is a hard colour transition (duplicate stops):
			// fill its slice with the band's colour instead of interpolating.
			paint.FillShape(ops, graphics.ColorToNRGBA(colors[i]), shape)
		} else {
			paint.LinearGradientOp{
				Stop1:  stopFrom,
				Color1: graphics.ColorToNRGBA(colors[i]),
				Stop2:  stopTo,
				Color2: graphics.ColorToNRGBA(colors[i+1]),
			}.Add(ops)

			paint.PaintOp{}.Add(ops)
		}

		bandStack.Pop()
	}
}

// gradientBandClip builds a clip covering the slice of the gradient between the
// normalized positions t0 and t1, wide enough perpendicular to the gradient to
// cover box.
func gradientBandClip(ops *op.Ops, from, to f32.Point, t0, t1 float32, box image.Rectangle) clip.Op {
	d := to.Sub(from)
	// A degenerate gradient was rejected by the caller; guard anyway so the
	// normal below is never the zero vector.
	if d.X == 0 && d.Y == 0 {
		return clip.Rect{Max: box.Max}.Op()
	}

	// Perpendicular half-width: the box diagonal is enough to cover the box from
	// any gradient direction.
	halfWidth := float32(box.Dx()+box.Dy()) + 1
	len2 := d.X*d.X + d.Y*d.Y
	invLen := 1 / float32(math.Sqrt(float64(len2)))
	dir := f32.Pt(d.X*invLen, d.Y*invLen)
	normal := f32.Pt(-dir.Y, dir.X).Mul(halfWidth)

	at := func(t float32) f32.Point { return from.Add(d.Mul(t)) }

	var path clip.Path
	path.Begin(ops)
	path.MoveTo(at(t0).Sub(normal))
	path.LineTo(at(t0).Add(normal))
	path.LineTo(at(t1).Add(normal))
	path.LineTo(at(t1).Sub(normal))
	path.Close()

	return clip.Outline{Path: path.End()}.Op()
}

// resolveGradientPoint maps a brush endpoint onto absolute pixels within box.
// Gradient brushes express infinite endpoints as +Inf, meaning "the far edge".
func resolveGradientPoint(o geometry.Offset, box image.Rectangle) f32.Point {
	x, y := o.X(), o.Y()

	if math.IsInf(float64(x), 1) {
		x = float32(box.Max.X)
	}
	if math.IsInf(float64(y), 1) {
		y = float32(box.Max.Y)
	}
	if math.IsNaN(float64(x)) {
		x = 0
	}
	if math.IsNaN(float64(y)) {
		y = 0
	}

	// Brush coordinates are box-relative; the icon is drawn at the origin.
	return f32.Pt(x+float32(box.Min.X), y+float32(box.Min.Y))
}

// lerpPoint returns the point along from->to at t.
func lerpPoint(from, to f32.Point, t float32) f32.Point {
	return f32.Pt(
		from.X+(to.X-from.X)*t,
		from.Y+(to.Y-from.Y)*t,
	)
}

// strokeShape builds a stroke outline for geometry at the given pixel width.
func strokeShape(geometryPath clip.PathSpec, widthPx float32) clip.Op {
	return clip.Stroke{Path: geometryPath, Width: widthPx}.Op()
}

// pointerOf returns the identity of a pointer-like value for cache keying.
//
// The brush is an interface, so a type switch covers the pointer-backed kinds;
// value-backed brushes return a constant, which makes them cache as a single
// entry rather than one per call.
func pointerOf(b any) uintptr {
	switch v := b.(type) {
	case *graphics.SolidColor:
		return uintptr(reflect.ValueOf(v).Pointer())
	case *graphics.LinearGradient:
		return uintptr(reflect.ValueOf(v).Pointer())
	case graphics.LinearGradient:
		// A value brush has no stable identity; fold its shape into a hash so two
		// different gradients do not share an entry but an identical one does.
		return uintptr(hashBrush(v))
	default:
		return 0
	}
}

// hashBrush produces a stable identity for a value-typed gradient.
func hashBrush(g graphics.LinearGradient) uint64 {
	h := fnv.New64a()
	for _, c := range g.Colors {
		n := graphics.ColorToNRGBA(c)
		h.Write([]byte{byte(n.R), byte(n.G), byte(n.B), byte(n.A)})
	}
	for _, s := range g.Stops {
		var buf [4]byte
		binary.LittleEndian.PutUint32(buf[:], math.Float32bits(s))
		h.Write(buf[:])
	}
	for _, o := range []geometry.Offset{g.Start, g.End} {
		for _, f := range []float32{o.X(), o.Y()} {
			var buf [4]byte
			binary.LittleEndian.PutUint32(buf[:], math.Float32bits(f))
			h.Write(buf[:])
		}
	}
	return h.Sum64()
}
