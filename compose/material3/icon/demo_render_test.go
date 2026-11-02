package icon_test

import (
	"image"
	"testing"

	"gioui.org/gpu/headless"
	"gioui.org/op"

	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	m3icon "github.com/zodimo/go-compose/compose/material3/icon"
	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/internal/screenshot"
	"github.com/zodimo/go-compose/pkg/api"
)

// demoLikeFixture mirrors the iconsymbols demo: several icons with different
// options stacked in a column.
func demoLikeFixture() api.Composable {
	return func(c api.Composer) api.Composer {
		return column.Column(
			c.Sequence(
				m3icon.Icon(m3icon.SymbolFavorite, m3icon.WithFill(0), m3icon.WithSize(unit.Dp(40))),
				m3icon.Icon(m3icon.SymbolFavorite, m3icon.WithFill(1), m3icon.WithSize(unit.Dp(40))),
				m3icon.Icon(m3icon.SymbolStar, m3icon.WithWeight(700), m3icon.WithSize(unit.Dp(40))),
				m3icon.Icon(m3icon.SymbolSearch, m3icon.WithSize(unit.Dp(32))),
				m3icon.Icon(m3icon.SymbolHome, m3icon.WithStroke(unit.Dp(2)), m3icon.WithSize(unit.Dp(40))),
			),
			column.WithAlignment(column.Start),
		)(c)
	}
}

func TestDemoLikeIconColumnPaintsAcrossFrames(t *testing.T) {
	w, h := 200, 320

	render := func(frame int) *image.RGBA {
		result := screenshot.Render(demoLikeFixture(), w, h)

		window, err := headless.NewWindow(w, h)
		if err != nil {
			t.Skipf("headless GPU unavailable: %v", err)
		}
		defer window.Release()

		var ops op.Ops
		result.DrawOp.Add(&ops)
		window.Frame(&ops)

		img := image.NewRGBA(image.Rect(0, 0, w, h))
		if err := window.Screenshot(img); err != nil {
			t.Skipf("headless screenshot unavailable: %v", err)
		}
		return img
	}

	for frame := 0; frame < 3; frame++ {
		img := render(frame)
		bg := img.RGBAAt(0, 0)
		ink := 0
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				if img.RGBAAt(x, y) != bg {
					ink++
				}
			}
		}
		if ink == 0 {
			t.Fatalf("frame %d: the icon column painted nothing", frame)
		}
		t.Logf("frame %d: %d non-background pixels", frame, ink)
	}
}
