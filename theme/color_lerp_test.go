package theme_test

import (
	"testing"

	"gioui.org/layout"
	"gioui.org/op"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/theme"
)

func TestColorLerp(t *testing.T) {
	tm := theme.GetThemeManager()
	gtx := layout.Context{
		Ops: new(op.Ops),
	}
	gtx = tm.Material3ThemeInit(gtx).(layout.Context)

	t.Run("Static Color Lerp", func(t *testing.T) {
		start := theme.ColorLerp(
			theme.ColorHelper.SpecificColor(graphics.ColorRed),
			theme.ColorHelper.SpecificColor(graphics.ColorBlue),
			0.5,
		)

		resolved := tm.ResolveColorDescriptor(start).AsNRGBA()

		// Expected: Purple (R: 127/128, G: 0, B: 127/128, A: 255)
		// lerp.Float32 logic: (1-0.5)*255 + 0.5*0 = 127.5 -> 127 // or 128 depending on rounding/casting
		// The current impl casts back to uint8, effectively truncating.

		expectedR := uint8(127) // or 128
		expectedB := uint8(127) // or 128

		if resolved.R < expectedR-1 || resolved.R > expectedR+2 {
			t.Errorf("Red channel mismatch: got %d, want approx %d", resolved.R, expectedR)
		}
		if resolved.B < expectedB-1 || resolved.B > expectedB+2 {
			t.Errorf("Blue channel mismatch: got %d, want approx %d", resolved.B, expectedB)
		}
		if resolved.G != 0 {
			t.Errorf("Green channel mismatch: got %d, want 0", resolved.G)
		}
		if resolved.A != 255 {
			t.Errorf("Alpha channel mismatch: got %d, want 255", resolved.A)
		}
	})

	t.Run("Role Lerp", func(t *testing.T) {
		// Primary (Default M3 Purple) to Secondary
		// We just verify that it produces something different from start and stop

		primaryDesc := theme.ColorHelper.ColorSelector().PrimaryRoles.Primary
		secondaryDesc := theme.ColorHelper.ColorSelector().SecondaryRoles.Secondary

		lerpedDesc := theme.ColorLerp(primaryDesc, secondaryDesc, 0.5)

		primary := tm.ResolveColorDescriptor(primaryDesc).AsNRGBA()
		secondary := tm.ResolveColorDescriptor(secondaryDesc).AsNRGBA()
		lerped := tm.ResolveColorDescriptor(lerpedDesc).AsNRGBA()

		if lerped == primary {
			t.Error("Lerped color should not be equal to Primary")
		}
		if lerped == secondary {
			t.Error("Lerped color should not be equal to Secondary")
		}

		// Verify interpolation logic for Red channel roughly
		expectedR := uint8((float32(primary.R) + float32(secondary.R)) / 2)
		if diff := int(lerped.R) - int(expectedR); diff < -2 || diff > 2 {
			t.Errorf("Lerped R channel mismatch: got %d, expected approx %d", lerped.R, expectedR)
		}
	})

	t.Run("Chained Lerp", func(t *testing.T) {
		start := theme.ColorHelper.SpecificColor(graphics.ColorRed)
		mid := theme.ColorHelper.SpecificColor(graphics.ColorGreen)
		end := theme.ColorHelper.SpecificColor(graphics.ColorBlue)

		// Start -> Mid (50%) = (127, 127, 0)
		// Then Lerp to End (50%) = (63, 63, 127)

		lerped := theme.ColorLerp(start, mid, 0.5)
		doubleLerped := theme.ColorLerp(lerped, end, 0.5)

		resolved := tm.ResolveColorDescriptor(doubleLerped).AsNRGBA()

		if resolved.R < 60 || resolved.R > 70 {
			t.Errorf("R mismatch: got %d", resolved.R)
		}
		if resolved.G < 60 || resolved.G > 70 {
			t.Errorf("G mismatch: got %d", resolved.G)
		}
		if resolved.B < 120 || resolved.B > 135 {
			t.Errorf("B mismatch: got %d", resolved.B)
		}
	})
}
