package icon

import (
	"image"
	"image/color"
	"testing"

	"gioui.org/op"
	"gioui.org/op/clip"

	"github.com/zodimo/go-compose/compose/ui/geometry"
	"github.com/zodimo/go-compose/compose/ui/graphics"
)

// gradientFixture paints a linear gradient across a solid rectangle and reports
// the colour at three points along its axis. Rasterising is the only way to prove
// the bands are composited in order; op counts say nothing about which band wins.
func gradientFixture(t *testing.T, brush graphics.LinearGradient) (left, mid, right color.RGBA) {
	t.Helper()

	size := image.Pt(120, 40)
	f := fill{brush: graphics.Brush(brush), color: color.NRGBA{A: 255}}
	img := rasterise(t, size, func(ops *op.Ops) {
		f.paint(ops, image.Rectangle{Max: size}, clip.Rect{Max: size}.Op())
	})

	// Sample exactly at the gradient's own stops: its axis runs 0..100 here, so
	// x=0, 50 and 100 sit on stop 0, 0.5 and 1.
	return img.RGBAAt(0, size.Y/2), img.RGBAAt(50, size.Y/2), img.RGBAAt(100, size.Y/2)
}

// TestFill_MultiStopGradientKeepsFirstStopAtTheStart is the regression test for
// bands painted unclipped: every band covered the whole shape, so the last band
// won and the start of the gradient showed the second colour instead of the
// first.
func TestFill_MultiStopGradientKeepsFirstStopAtTheStart(t *testing.T) {
	brush := graphics.LinearGradient{
		Colors: []graphics.Color{
			graphics.NewColorSrgb(255, 0, 0, 255),
			graphics.NewColorSrgb(0, 255, 0, 255),
			graphics.NewColorSrgb(0, 0, 255, 255),
		},
		Stops: []float32{0, 0.5, 1},
		Start: geometry.NewOffset(0, 0),
		End:   geometry.NewOffset(100, 0),
	}

	left, mid, right := gradientFixture(t, brush)
	t.Logf("left=%v mid=%v right=%v", left, mid, right)

	if left.R < 200 || left.G > 60 || left.B > 60 {
		t.Errorf("left edge = %v, want the first stop colour (red)", left)
	}
	if right.B < 200 || right.R > 60 || right.G > 60 {
		t.Errorf("right edge = %v, want the last stop colour (blue)", right)
	}
	if mid.R > 60 || mid.B > 60 || mid.G < 150 {
		t.Errorf("midpoint = %v, want the middle stop colour (green)", mid)
	}
}

// TestFill_MultiStopGradientCoversWholeShape guards the other half of the
// clipping change: a band clip that is too narrow would leave unpainted gaps
// where nothing is drawn at all.
func TestFill_MultiStopGradientCoversWholeShape(t *testing.T) {
	brush := graphics.LinearGradient{
		Colors: []graphics.Color{
			graphics.NewColorSrgb(255, 0, 0, 255),
			graphics.NewColorSrgb(0, 255, 0, 255),
			graphics.NewColorSrgb(0, 0, 255, 255),
			graphics.NewColorSrgb(255, 255, 255, 255),
		},
		Stops: []float32{0, 0.33, 0.66, 1},
		Start: geometry.NewOffset(0, 0),
		End:   geometry.NewOffset(100, 0),
	}

	size := image.Pt(120, 40)
	f := fill{brush: graphics.Brush(brush), color: color.NRGBA{A: 255}}
	img := rasterise(t, size, func(ops *op.Ops) {
		f.paint(ops, image.Rectangle{Max: size}, clip.Rect{Max: size}.Op())
	})

	// Every column inside the rectangle must be painted; a zero alpha column is
	// a gap between bands.
	for x := 0; x < size.X; x++ {
		if got := img.RGBAAt(x, size.Y/2); got.A == 0 {
			t.Fatalf("column %d was left unpainted (%v)", x, got)
		}
	}
}

// TestFill_HardStopPaintsBothSides covers duplicate stops, which express a hard
// colour transition. The zero-width band must fill its slice rather than
// abandoning the remaining bands.
func TestFill_HardStopPaintsBothSides(t *testing.T) {
	brush := graphics.LinearGradient{
		Colors: []graphics.Color{
			graphics.NewColorSrgb(255, 0, 0, 255),
			graphics.NewColorSrgb(255, 0, 0, 255),
			graphics.NewColorSrgb(0, 0, 255, 255),
		},
		Stops: []float32{0, 0.5, 1},
		Start: geometry.NewOffset(0, 0),
		End:   geometry.NewOffset(100, 0),
	}

	size := image.Pt(120, 40)
	f := fill{brush: graphics.Brush(brush), color: color.NRGBA{A: 255}}
	img := rasterise(t, size, func(ops *op.Ops) {
		f.paint(ops, image.Rectangle{Max: size}, clip.Rect{Max: size}.Op())
	})

	if got := img.RGBAAt(10, size.Y/2); got.R < 200 || got.B > 60 {
		t.Errorf("before the hard stop = %v, want red", got)
	}
	if got := img.RGBAAt(110, size.Y/2); got.B < 200 || got.R > 60 {
		t.Errorf("after the hard stop = %v, want blue", got)
	}
	for x := 0; x < size.X; x++ {
		if got := img.RGBAAt(x, size.Y/2); got.A == 0 {
			t.Fatalf("column %d was left unpainted (%v)", x, got)
		}
	}
}
