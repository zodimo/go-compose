package theme_test

import (
	"testing"

	"gioui.org/layout"
	"gioui.org/op"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/theme"
)

func TestSurfaceColorResolution(t *testing.T) {
	// Initialize ThemeManager and Theme
	tm := theme.GetThemeManager()
	gtx := layout.Context{
		Ops: new(op.Ops),
	}
	// Force initialization (might need to reset execution state if singleton persists across tests, but simple init usually works)
	// Material3ThemeInit checks if theme is nil.
	gtx = tm.Material3ThemeInit(gtx).(layout.Context)

	// Get Surface Descriptor
	surfaceDesc := theme.ColorHelper.ColorSelector().SurfaceRoles.Surface

	// Resolve it
	resolvedColor := tm.ResolveColorDescriptor(surfaceDesc)
	rgba := resolvedColor.AsNRGBA()

	// Expected M3 Baseline Light Surface Color: #FFFBFE (R:255, G:251, B:254) ?
	// Or based on my logs: #FDF7FF (R:253, G:247, B:255) ?
	// Let's print it first or check against what we saw.
	// If the test fails, we know what the canonical color is.

	t.Logf("Resolved Surface Color: %+v", rgba)

	// Verify it is NOT transparent
	if rgba.A == 0 {
		t.Error("Surface color resolved to Transparent")
	}

	// Verify it is NOT Primary (Purple 40: #6750A4 roughly)
	// Primary: R:103 G:80 B:164
	if rgba.R == 103 && rgba.G == 80 && rgba.B == 164 {
		t.Error("Surface color resolved to Primary!")
	}
}

func TestSpecificColorPreservation(t *testing.T) {
	tm := theme.GetThemeManager()
	targetColor := graphics.ColorRed

	desc := theme.ColorHelper.SpecificColor(targetColor)
	resolved := tm.ResolveColorDescriptor(desc).AsNRGBA()

	if resolved != graphics.ColorToNRGBA(targetColor) {
		t.Errorf("SpecificColor failed: got %+v, want %+v", resolved, targetColor)
	}
}

func TestTransparencyResolution(t *testing.T) {
	tm := theme.GetThemeManager()
	transparent := graphics.ColorTransparent

	desc := theme.ColorHelper.SpecificColor(transparent)
	resolved := tm.ResolveColorDescriptor(desc).AsNRGBA()

	if resolved != graphics.ColorToNRGBA(transparent) {
		t.Errorf("Transparent color failed: got %+v, want %+v", resolved, transparent)
	}
}
