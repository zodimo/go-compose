package shape

import (
	"image"

	"gioui.org/op/clip"
	gioUnit "gioui.org/unit"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

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

func (r RoundedCornerShape) CreateOutline(size image.Point, metric gioUnit.Metric) Outline {
	// Determine if per-corner radius is being used
	hasPerCorner := r.TopStart > 0 || r.TopEnd > 0 || r.BottomEnd > 0 || r.BottomStart > 0

	var nw, ne, se, sw int
	if hasPerCorner {
		nw = metric.Dp(unit.DpToGioUnit(r.TopStart))
		ne = metric.Dp(unit.DpToGioUnit(r.TopEnd))
		se = metric.Dp(unit.DpToGioUnit(r.BottomEnd))
		sw = metric.Dp(unit.DpToGioUnit(r.BottomStart))
	} else {
		radius := metric.Dp(unit.DpToGioUnit(r.Radius))
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
