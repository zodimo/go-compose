package icon

import (
	"bytes"
	"image"
	"image/color"
	"testing"

	"gioui.org/gpu/headless"
	"gioui.org/layout"
	"gioui.org/op"
)

// rasterise runs ops through a headless GPU and returns the resulting pixels.
//
// This is the only way to assert that an icon is actually PAINTED. Every other
// test in this package checks geometry or dimensions, both of which stay correct
// even when nothing is drawn — which is exactly how an empty-op cache replay
// shipped unnoticed.
func rasterise(t *testing.T, size image.Point, record func(*op.Ops)) *image.RGBA {
	t.Helper()

	window, err := headless.NewWindow(size.X, size.Y)
	if err != nil {
		t.Skipf("headless GPU unavailable: %v", err)
	}
	t.Cleanup(window.Release)

	var ops op.Ops
	record(&ops)
	window.Frame(&ops)

	img := image.NewRGBA(image.Rectangle{Max: size})
	if err := window.Screenshot(img); err != nil {
		t.Skipf("headless screenshot unavailable: %v", err)
	}
	return img
}

// nonBackground counts pixels that differ from the top-left pixel, which is
// outside the icon's ink and therefore shows the background.
func nonBackgroundPixels(img *image.RGBA) int {
	bg := img.RGBAAt(0, 0)
	count := 0
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			if img.RGBAAt(x, y) != bg {
				count++
			}
		}
	}
	return count
}

// TestSymbol_PaintsPixelsAcrossFrames is the regression test for the demo bug
// where every frame after the first drew nothing: the paint cache recorded into
// a private op.Ops but the widget drew into gtx.Ops, so the recorded macro was
// empty and every cache hit replayed emptiness.
//
// It renders the same icon configuration twice through the widget path and
// requires ink on BOTH frames, since only the first frame took the recording
// path.
func TestSymbol_PaintsPixelsAcrossFrames(t *testing.T) {
	freshOpsCache(t, 8)

	render, err := buildSymbolRender(SymbolOutlined, SymbolHome, 32, resolveAxes(DefaultIconOptions(), 32))
	if err != nil {
		t.Fatalf("buildSymbolRender: %v", err)
	}

	f := fill{color: color.NRGBA{A: 255}}
	size := image.Pt(48, 48)

	paintFrame := func(frame int) {
		t.Helper()

		widget := symbolWidgetConstructor(
			geometryKey{symbol: SymbolHome, style: SymbolOutlined, axes: resolveAxes(DefaultIconOptions(), 32), sizePx: 32},
			render, f, 0,
		).Make(nil)

		img := rasterise(t, size, func(ops *op.Ops) {
			gtx := layout.Context{Ops: ops, Constraints: layout.Exact(render.size)}
			widget(gtx)
		})

		if got := nonBackgroundPixels(img); got == 0 {
			t.Fatalf("frame %d painted no pixels; the icon is invisible", frame)
		}
	}

	// Frame 0 records into the cache; frame 1 replays the recorded ops.
	paintFrame(0)
	paintFrame(1)
}

// TestSymbol_PaintDiffersBetweenColors proves the recorded ops are actually
// replayed per paint key rather than being empty.
func TestSymbol_PaintDiffersBetweenColors(t *testing.T) {
	freshOpsCache(t, 8)

	render, err := buildSymbolRender(SymbolOutlined, SymbolHome, 32, resolveAxes(DefaultIconOptions(), 32))
	if err != nil {
		t.Fatalf("buildSymbolRender: %v", err)
	}

	geomKey := geometryKey{
		symbol: SymbolHome, style: SymbolOutlined,
		axes: resolveAxes(DefaultIconOptions(), 32), sizePx: 32,
	}
	size := image.Pt(48, 48)

	draw := func(c color.NRGBA) *image.RGBA {
		widget := symbolWidgetConstructor(geomKey, render, fill{color: c}, 0).Make(nil)
		return rasterise(t, size, func(ops *op.Ops) {
			gtx := layout.Context{Ops: ops, Constraints: layout.Exact(render.size)}
			widget(gtx)
			widget(gtx) // second call hits the ops cache
		})
	}

	red := draw(color.NRGBA{R: 255, A: 255})
	blue := draw(color.NRGBA{B: 255, A: 255})

	if nonBackgroundPixels(red) == 0 {
		t.Fatal("red icon painted nothing")
	}
	if nonBackgroundPixels(blue) == 0 {
		t.Fatal("blue icon painted nothing")
	}
	if bytes.Equal(red.Pix, blue.Pix) {
		t.Fatal("red and blue renders are byte-identical; paint is not being recorded")
	}
}
