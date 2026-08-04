package background

import (
	"image"

	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	"github.com/zodimo/go-compose/internal/render"
)

// toRenderColor converts a framework Color to a render.Color (linear sRGB [0,1]).
func toRenderColor(c graphics.Color) render.Color {
	nrgba := graphics.ColorToNRGBA(c)
	return render.Color{
		R: float32(nrgba.R) / 255,
		G: float32(nrgba.G) / 255,
		B: float32(nrgba.B) / 255,
		A: float32(nrgba.A) / 255,
	}
}

// toRenderShape converts a framework Shape + size to a render.Shape for the
// software backend. It handles Rectangle, RoundedRectangle, and Circle shapes.
func toRenderShape(s shape.Shape, size image.Point, metric shape.Metric) render.Shape {
	if s == nil {
		return render.Shape{
			Kind:   render.ShapeRectangle,
			Bounds: render.Rect{Min: render.Point{}, Max: pointToRenderPoint(size)},
		}
	}
	switch v := s.(type) {
	case *shape.RoundedCornerShape:
		return roundedCornerToRenderShape(v, size, metric)
	default:
		// Rectangle, Circle, CutCorner → treat as rectangle for FillShape.
		// Circle and CutCorner would need their own ShapeKinds for exact
		// rasterization, but for golden-test DrawCall assertions the bounds
		// are sufficient.
		return render.Shape{
			Kind:   render.ShapeRectangle,
			Bounds: render.Rect{Min: render.Point{}, Max: pointToRenderPoint(size)},
		}
	}
}

func roundedCornerToRenderShape(r *shape.RoundedCornerShape, size image.Point, metric shape.Metric) render.Shape {
	nw, ne, se, sw := roundedCornerRadii(r, metric)
	if nw == 0 && ne == 0 && se == 0 && sw == 0 {
		return render.Shape{
			Kind:   render.ShapeRectangle,
			Bounds: render.Rect{Min: render.Point{}, Max: pointToRenderPoint(size)},
		}
	}
	return render.Shape{
		Kind:   render.ShapeRoundedRectangle,
		Bounds: render.Rect{Min: render.Point{}, Max: pointToRenderPoint(size)},
		Radius: render.CornerRadius{
			TL: float32(nw),
			TR: float32(ne),
			BL: float32(sw),
			BR: float32(se),
		},
	}
}

func roundedCornerRadii(r *shape.RoundedCornerShape, metric shape.Metric) (nw, ne, se, sw int) {
	if r.Radius.IsSpecified() {
		radius := int(float32(r.Radius) * metric.PxPerDp)
		return radius, radius, radius, radius
	}
	topStart := r.TopStart.TakeOrElse(0)
	topEnd := r.TopEnd.TakeOrElse(0)
	bottomEnd := r.BottomEnd.TakeOrElse(0)
	bottomStart := r.BottomStart.TakeOrElse(0)
	nw = int(float32(topStart) * metric.PxPerDp)
	ne = int(float32(topEnd) * metric.PxPerDp)
	se = int(float32(bottomEnd) * metric.PxPerDp)
	sw = int(float32(bottomStart) * metric.PxPerDp)
	return
}

func pointToRenderPoint(p image.Point) render.Point {
	return render.Point{X: float32(p.X), Y: float32(p.Y)}
}
