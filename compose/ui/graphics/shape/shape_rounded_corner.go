package shape

import (
	"fmt"
	"image"

	"gioui.org/op/clip"
	gioUnit "gioui.org/unit"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

var roundedCornerShapeUnspecified = &RoundedCornerShape{
	Radius: unit.DpCompanion.Unspecified(),

	TopStart:    unit.DpCompanion.Unspecified(),
	TopEnd:      unit.DpCompanion.Unspecified(),
	BottomEnd:   unit.DpCompanion.Unspecified(),
	BottomStart: unit.DpCompanion.Unspecified(),
}

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

func (r *RoundedCornerShape) CreateOutline(size image.Point, metric gioUnit.Metric) Outline {
	rValid := coalesceRoundedCornerShape(r, roundedCornerShapeUnspecified)

	var nw, ne, se, sw int

	// Check if uniform Radius is specified first (this is the most common case)
	// Note: Dp(0) is technically "specified" but means no rounding, which is correct.
	// The uniform Radius field takes precedence when set.
	if rValid.Radius.IsSpecified() {
		radius := metric.Dp(unit.DpToGioUnitUnsafe(rValid.Radius))
		nw, ne, se, sw = radius, radius, radius, radius
	} else {
		// Fall back to per-corner values if uniform Radius is not set or is 0
		// Per-corner values are only used if at least one is > 0
		topStart := unit.DpCompanion.TakeOrElse(rValid.TopStart, unit.Dp(0))
		topEnd := unit.DpCompanion.TakeOrElse(rValid.TopEnd, unit.Dp(0))
		bottomEnd := unit.DpCompanion.TakeOrElse(rValid.BottomEnd, unit.Dp(0))
		bottomStart := unit.DpCompanion.TakeOrElse(rValid.BottomStart, unit.Dp(0))

		if topStart > 0 || topEnd > 0 || bottomEnd > 0 || bottomStart > 0 {
			nw = metric.Dp(unit.DpToGioUnitUnsafe(topStart))
			ne = metric.Dp(unit.DpToGioUnitUnsafe(topEnd))
			se = metric.Dp(unit.DpToGioUnitUnsafe(bottomEnd))
			sw = metric.Dp(unit.DpToGioUnitUnsafe(bottomStart))
		}
	}

	// If all corners are 0, use a simple rectangle
	if nw == 0 && ne == 0 && se == 0 && sw == 0 {
		return rectOutline{clip.Rect{Max: size}}
	}

	return rrectOutline{clip.RRect{
		Rect: image.Rectangle{Max: size},
		NW:   nw,
		NE:   ne,
		SE:   se,
		SW:   sw,
	}}
}

func (r *RoundedCornerShape) mergeShape(other Shape) Shape {
	if otherRoundedCorner, ok := other.(*RoundedCornerShape); ok {
		return &RoundedCornerShape{
			Radius:      unit.DpCompanion.TakeOrElse(otherRoundedCorner.Radius, r.Radius),
			TopStart:    unit.DpCompanion.TakeOrElse(otherRoundedCorner.TopStart, r.TopStart),
			TopEnd:      unit.DpCompanion.TakeOrElse(otherRoundedCorner.TopEnd, r.TopEnd),
			BottomEnd:   unit.DpCompanion.TakeOrElse(otherRoundedCorner.BottomEnd, r.BottomEnd),
			BottomStart: unit.DpCompanion.TakeOrElse(otherRoundedCorner.BottomStart, r.BottomStart),
		}
	}
	panic(fmt.Sprintf("RoundedCornerShape.mergeShape: cannot merge with %s", other.stringShape()))
}
func (r *RoundedCornerShape) sameShape(other Shape) bool {
	if _, ok := other.(*RoundedCornerShape); ok {
		return true
	}
	return false
}
func (r *RoundedCornerShape) semanticEqualShape(other Shape) bool {
	if otherShape, ok := other.(*RoundedCornerShape); ok {
		return otherShape.Radius == r.Radius
	}
	return false
}
func (r *RoundedCornerShape) copyShape(options ...ShapeOption) Shape {

	rValid := coalesceRoundedCornerShape(r, roundedCornerShapeUnspecified)

	copy := *rValid
	for _, option := range options {
		option(&copy)
	}
	return &copy
}

func (r *RoundedCornerShape) stringShape() string {
	if !isSpecifiedRoundedCornerShape(r) {
		return "RoundedCornerShape{Unspecified}"
	}
	return fmt.Sprintf(
		"RoundedCornerShape{Radius: %s, TopStart: %s, TopEnd: %s, BottomEnd: %s, BottomStart: %s}",
		r.Radius.String(), r.TopStart.String(), r.TopEnd.String(), r.BottomEnd.String(), r.BottomStart.String(),
	)
}

func isSpecifiedRoundedCornerShape(r *RoundedCornerShape) bool {
	return r != nil && r != roundedCornerShapeUnspecified
}

func coalesceRoundedCornerShape(ptr, def *RoundedCornerShape) *RoundedCornerShape {
	if ptr == nil {
		return def
	}
	return ptr
}
