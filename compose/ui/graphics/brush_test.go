package graphics_test

import (
	"testing"

	"github.com/zodimo/go-compose/compose/ui/geometry"
	"github.com/zodimo/go-compose/compose/ui/graphics"
)

func TestSolidColor_Equal(t *testing.T) {
	c1 := graphics.NewSolidColor(graphics.ColorBlack)
	c2 := graphics.NewSolidColor(graphics.ColorBlack)
	c3 := graphics.NewSolidColor(graphics.ColorTransparent)

	if !graphics.EqualBrush(c1, c2) {
		t.Errorf("SolidColor Black should equal Black")
	}
	if graphics.EqualBrush(c1, c3) {
		t.Errorf("SolidColor Black should not equal Transparent")
	}
}

func TestLinearGradient_Equal(t *testing.T) {
	colors := []graphics.Color{graphics.ColorBlack, graphics.ColorTransparent}
	g1 := graphics.LinearGradientBrush(colors, geometry.OffsetZero, geometry.OffsetInfinite, graphics.TileModeClamp)
	g2 := graphics.LinearGradientBrush(colors, geometry.OffsetZero, geometry.OffsetInfinite, graphics.TileModeClamp)
	g3 := graphics.LinearGradientBrush(colors, geometry.OffsetZero, geometry.NewOffset(10, 10), graphics.TileModeClamp)

	if !graphics.EqualBrush(g1, g2) {
		t.Errorf("LinearGradient should equal identical instance")
	}
	if graphics.EqualBrush(g1, g3) {
		t.Errorf("LinearGradient should not equal different end offset")
	}
}

func TestRadialGradient_Equal(t *testing.T) {
	colors := []graphics.Color{graphics.ColorBlack, graphics.ColorTransparent}
	g1 := graphics.RadialGradientBrush(colors, geometry.OffsetZero, 100, graphics.TileModeClamp)
	g2 := graphics.RadialGradientBrush(colors, geometry.OffsetZero, 100, graphics.TileModeClamp)
	g3 := graphics.RadialGradientBrush(colors, geometry.OffsetZero, 200, graphics.TileModeClamp)

	if !graphics.EqualBrush(g1, g2) {
		t.Errorf("RadialGradient should equal identical instance")
	}
	if graphics.EqualBrush(g1, g3) {
		t.Errorf("RadialGradient should not equal different radius")
	}
}

func TestHorizontalGradient(t *testing.T) {
	colors := []graphics.Color{graphics.ColorBlack, graphics.ColorTransparent}
	g := graphics.HorizontalGradient(colors, 0, 100, graphics.TileModeMirror)

	if g.Start.X() != 0 || g.Start.Y() != 0 {
		t.Errorf("HorizontalGradient start should be (0,0)")
	}
	if g.End.X() != 100 || g.End.Y() != 0 {
		t.Errorf("HorizontalGradient end should be (100,0)")
	}
	if g.TileMode != graphics.TileModeMirror {
		t.Errorf("HorizontalGradient tile mode mismatch")
	}
}

func TestLerpBrush(t *testing.T) {
	// Case 1: SolidColor <-> SolidColor
	c1 := graphics.NewSolidColor(graphics.ColorBlack)
	c2 := graphics.NewSolidColor(graphics.ColorTransparent)
	// Lerp at 0.0 -> Black
	l1 := graphics.LerpBrush(c1, c2, 0.0)
	if !graphics.EqualBrush(l1, c1) {
		t.Errorf("LerpBrush(0.0) should be equal to start")
	}
	// Lerp at 1.0 -> Transparent
	l2 := graphics.LerpBrush(c1, c2, 1.0)
	if !graphics.EqualBrush(l2, c2) {
		t.Errorf("LerpBrush(1.0) should be equal to stop")
	}

	// Case 2: LinearGradient <-> LinearGradient
	g1 := graphics.HorizontalGradient([]graphics.Color{graphics.ColorBlack, graphics.ColorTransparent}, 0, 100, graphics.TileModeClamp)
	g2 := graphics.HorizontalGradient([]graphics.Color{graphics.ColorBlack, graphics.ColorTransparent}, 0, 200, graphics.TileModeClamp)

	// Lerp at 0.5 -> EndX should be 150
	l3 := graphics.LerpBrush(g1, g2, 0.5)
	lg3, ok := l3.(*graphics.LinearGradient)
	if !ok {
		t.Fatalf("Lerp of LinearGradient should be LinearGradient")
	}
	if lg3.End.X() != 150 {
		t.Errorf("LerpBrush(0.5) LinearGradient End.X = %v, want 150", lg3.End.X())
	}

	// Case 3: Mixed (Solid <-> Gradient)
	// < 0.5 -> start
	l4 := graphics.LerpBrush(c1, g1, 0.4)
	if !graphics.EqualBrush(l4, c1) {
		t.Errorf("LerpBrush(0.4) mixed should be start")
	}
	// >= 0.5 -> stop
	l5 := graphics.LerpBrush(c1, g1, 0.6)
	if !graphics.EqualBrush(l5, g1) {
		t.Errorf("LerpBrush(0.6) mixed should be stop")
	}
}
