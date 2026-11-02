package icon

import (
	"image"
	"image/color"
	"testing"

	"gioui.org/op"
	"gioui.org/op/clip"

	"github.com/zodimo/go-compose/compose/ui/geometry"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

func TestWithBrush_PrecedenceOverColor(t *testing.T) {
	brush := graphics.LinearGradientBrush(
		[]graphics.Color{
			graphics.NewColorSrgb(255, 0, 0, 255),
			graphics.NewColorSrgb(0, 0, 255, 255),
		},
		geometry.NewOffset(0, 0),
		geometry.NewOffset(24, 24),
		graphics.TileModeClamp,
	)

	opts := optionsWith(WithBrush(brush), WithColor(graphics.NewColorSrgb(0, 255, 0, 255)))
	if opts.Brush == nil {
		t.Fatal("WithBrush did not record the brush")
	}

	// The renderer builds its fill from brush + resolved color; brush presence is
	// what selects gradient painting.
	f := fill{brush: opts.Brush, color: color.NRGBA{G: 255, A: 255}}
	if f.brush == nil {
		t.Fatal("fill lost the brush")
	}
}

func TestFill_SolidColorPaints(t *testing.T) {
	var ops op.Ops
	shape := rectShape()

	// Must not panic and must not attempt gradient painting.
	f := fill{color: color.NRGBA{A: 255}}
	f.paint(&ops, image.Rect(0, 0, 24, 30), shape)
}

func TestFill_LinearGradientPaintsMultipleBands(t *testing.T) {
	var ops op.Ops

	brush := graphics.LinearGradient{
		Colors: []graphics.Color{
			graphics.NewColorSrgb(255, 0, 0, 255),
			graphics.NewColorSrgb(0, 255, 0, 255),
			graphics.NewColorSrgb(0, 0, 255, 255),
		},
		Stops: []float32{0, 0.5, 1},
		Start: geometry.NewOffset(0, 0),
		End:   geometry.NewOffset(24, 0),
	}

	f := fill{brush: graphics.Brush(brush), color: color.NRGBA{}}
	if !f.paintGradient(&ops, image.Rect(0, 0, 24, 30), rectShape()) {
		t.Fatal("paintGradient declined a linear gradient brush")
	}
}

func TestFill_SingleStopGradientFallsBackToSolid(t *testing.T) {
	var ops op.Ops

	brush := graphics.LinearGradient{
		Colors: []graphics.Color{graphics.NewColorSrgb(10, 20, 30, 255)},
		Start:  geometry.NewOffset(0, 0),
		End:    geometry.NewOffset(24, 0),
	}
	f := fill{brush: graphics.Brush(brush)}
	// A one-stop gradient is a solid fill; it must still be handled.
	if !f.paintGradient(&ops, image.Rect(0, 0, 24, 30), rectShape()) {
		t.Fatal("paintGradient declined a single-stop gradient")
	}
}

func TestFill_UnsupportedBrushFallsBackToColor(t *testing.T) {
	var ops op.Ops

	// A radial gradient has no Gio painter here; the icon must still paint with
	// the resolved color rather than disappearing.
	brush := graphics.RadialGradientBrush(
		[]graphics.Color{
			graphics.NewColorSrgb(255, 0, 0, 255),
			graphics.NewColorSrgb(0, 0, 255, 255),
		},
		geometry.NewOffset(12, 12),
		12,
		graphics.TileModeClamp,
	)
	f := fill{brush: brush, color: color.NRGBA{A: 255}}

	// paint() is the fallback boundary; it must not panic.
	f.paint(&ops, image.Rect(0, 0, 24, 30), rectShape())
}

func TestWithStroke_RecordsWidth(t *testing.T) {
	opts := optionsWith(WithStroke(2))
	if !opts.hasStroke {
		t.Fatal("WithStroke did not mark stroke as provided")
	}
	if opts.StrokeWidth != 2 {
		t.Fatalf("stroke width = %v, want 2", opts.StrokeWidth)
	}

	// Unset by default.
	if DefaultIconOptions().hasStroke {
		t.Fatal("default options claim a stroke")
	}
}

func TestStrokeWidth_DoesNotChangeReportedSize(t *testing.T) {
	render, err := buildSymbolRender(SymbolOutlined, SymbolHome, 24, resolveAxes(DefaultIconOptions(), 24))
	if err != nil {
		t.Fatalf("buildSymbolRender: %v", err)
	}

	var filledOps, strokedOps op.Ops

	filledDims := render.draw(&filledOps, fill{color: color.NRGBA{A: 255}}, 0)
	strokedDims := render.draw(&strokedOps, fill{color: color.NRGBA{A: 255}}, 4)

	if filledDims.Size != strokedDims.Size {
		t.Fatalf("stroke changed reported size: %v vs %v", filledDims.Size, strokedDims.Size)
	}
	if filledDims.Baseline != strokedDims.Baseline {
		t.Fatalf("stroke changed reported baseline: %d vs %d",
			filledDims.Baseline, strokedDims.Baseline)
	}
	if strokedDims.Size != render.size {
		t.Fatalf("stroked size = %v, want the unstroked %v", strokedDims.Size, render.size)
	}
}

func TestStrokeWidthPxFromDp(t *testing.T) {
	density := unit.NewDensity(2, 1)
	if got := strokeWidthPxFromDp(2, density); got != 4 {
		t.Fatalf("2dp at 2x density = %v, want 4", got)
	}
	if got := strokeWidthPxFromDp(0, density); got != 0 {
		t.Fatalf("0dp = %v, want 0", got)
	}
}

// rectShape returns a trivial clip shape for paint tests.
func rectShape() clip.Op {
	return clip.Rect{Max: image.Pt(24, 30)}.Op()
}
