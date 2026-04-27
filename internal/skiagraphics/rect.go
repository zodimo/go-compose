package skiagraphics

import (
	"github.com/zodimo/gio-skia/skia"
	"github.com/zodimo/go-compose/compose/ui/geometry"
)

func GeometryRectToSkiaRect(r geometry.Rect) skia.Rect {
	return skia.Rect{
		Left:   r.Left,
		Top:    r.Top,
		Right:  r.Right,
		Bottom: r.Bottom,
	}
}

func SkiaRectToGeometryRect(r skia.Rect) geometry.Rect {
	return geometry.Rect{
		Left:   r.Left,
		Top:    r.Top,
		Right:  r.Right,
		Bottom: r.Bottom,
	}
}
