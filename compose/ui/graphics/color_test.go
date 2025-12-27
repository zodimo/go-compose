package graphics_test

import (
	"testing"

	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/graphics/colorspace"
)

func TestNewColorSrgb(t *testing.T) {
	// Test basic sRGB color creation
	c := graphics.NewColorSrgb(255, 128, 64, 255)

	if c.Alpha() < 0.99 || c.Alpha() > 1.01 {
		t.Errorf("Alpha = %f, want 1.0", c.Alpha())
	}
	if c.Red() < 0.99 || c.Red() > 1.01 {
		t.Errorf("Red = %f, want 1.0", c.Red())
	}
	if abs(c.Green()-0.502) > 0.01 {
		t.Errorf("Green = %f, want ~0.502", c.Green())
	}
	if abs(c.Blue()-0.251) > 0.01 {
		t.Errorf("Blue = %f, want ~0.251", c.Blue())
	}
}

func TestColorConstants(t *testing.T) {
	tests := []struct {
		name       string
		color      graphics.Color
		r, g, b, a float32
	}{
		{"Black", graphics.ColorBlack, 0, 0, 0, 1},
		{"White", graphics.ColorWhite, 1, 1, 1, 1},
		{"Red", graphics.ColorRed, 1, 0, 0, 1},
		{"Green", graphics.ColorGreen, 0, 1, 0, 1},
		{"Blue", graphics.ColorBlue, 0, 0, 1, 1},
		{"Transparent", graphics.ColorTransparent, 0, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if abs(tt.color.Red()-tt.r) > 0.01 {
				t.Errorf("Red = %f, want %f", tt.color.Red(), tt.r)
			}
			if abs(tt.color.Green()-tt.g) > 0.01 {
				t.Errorf("Green = %f, want %f", tt.color.Green(), tt.g)
			}
			if abs(tt.color.Blue()-tt.b) > 0.01 {
				t.Errorf("Blue = %f, want %f", tt.color.Blue(), tt.b)
			}
			if abs(tt.color.Alpha()-tt.a) > 0.01 {
				t.Errorf("Alpha = %f, want %f", tt.color.Alpha(), tt.a)
			}
		})
	}
}

func TestColorSpace(t *testing.T) {
	// Test sRGB color space
	c := graphics.ColorRed
	cs := c.ColorSpace()
	if cs.Id() != colorspace.ColorSpaceSrgb {
		t.Errorf("ColorSpace Id = %d, want %d", cs.Id(), colorspace.ColorSpaceSrgb)
	}
	if !cs.IsSrgb() {
		t.Error("Expected IsSrgb() to be true for sRGB color")
	}
}

func TestNewColorNonSrgb(t *testing.T) {
	// Create a DisplayP3 color
	c := graphics.NewColor(0.5, 0.5, 0.5, 1.0, colorspace.DisplayP3)

	// Verify components
	if abs(c.Red()-0.5) > 0.01 {
		t.Errorf("Red = %f, want 0.5", c.Red())
	}
	if abs(c.Green()-0.5) > 0.01 {
		t.Errorf("Green = %f, want 0.5", c.Green())
	}
	if abs(c.Blue()-0.5) > 0.01 {
		t.Errorf("Blue = %f, want 0.5", c.Blue())
	}
	if abs(c.Alpha()-1.0) > 0.01 {
		t.Errorf("Alpha = %f, want 1.0", c.Alpha())
	}

	// Verify color space
	if c.ColorSpaceId() != colorspace.ColorSpaceDisplayP3 {
		t.Errorf("ColorSpaceId = %d, want %d", c.ColorSpaceId(), colorspace.ColorSpaceDisplayP3)
	}
}

func TestIsSpecified(t *testing.T) {
	if !graphics.ColorRed.IsSpecified() {
		t.Error("ColorRed should be specified")
	}
	if graphics.ColorUnspecified.IsSpecified() {
		t.Error("ColorUnspecified should not be specified")
	}
	if !graphics.ColorUnspecified.IsUnspecified() {
		t.Error("ColorUnspecified should be unspecified")
	}
}

func TestTakeOrElse(t *testing.T) {
	red := graphics.ColorRed
	blue := graphics.ColorBlue

	// Specified color should return itself
	result := red.TakeOrElse(blue)
	if result != red {
		t.Error("TakeOrElse should return self for specified color")
	}

	// Unspecified color should return the alternative
	result = graphics.ColorUnspecified.TakeOrElse(blue)
	if result != blue {
		t.Error("TakeOrElse should return alternative for unspecified color")
	}
}

func TestGetComponents(t *testing.T) {
	c := graphics.NewColorSrgb(255, 128, 64, 200)
	components := c.GetComponents()

	if abs(components[0]-c.Red()) > 0.001 {
		t.Errorf("Component[0] = %f, want %f", components[0], c.Red())
	}
	if abs(components[1]-c.Green()) > 0.001 {
		t.Errorf("Component[1] = %f, want %f", components[1], c.Green())
	}
	if abs(components[2]-c.Blue()) > 0.001 {
		t.Errorf("Component[2] = %f, want %f", components[2], c.Blue())
	}
	if abs(components[3]-c.Alpha()) > 0.001 {
		t.Errorf("Component[3] = %f, want %f", components[3], c.Alpha())
	}
}

func TestComponentAccessors(t *testing.T) {
	c := graphics.ColorRed

	if c.Component1() != c.Red() {
		t.Error("Component1 should equal Red")
	}
	if c.Component2() != c.Green() {
		t.Error("Component2 should equal Green")
	}
	if c.Component3() != c.Blue() {
		t.Error("Component3 should equal Blue")
	}
	if c.Component4() != c.Alpha() {
		t.Error("Component4 should equal Alpha")
	}
	if c.Component5().Id() != c.ColorSpaceId() {
		t.Error("Component5 should equal ColorSpace")
	}
}

func TestToArgb(t *testing.T) {
	// Pure red should be 0xFFFF0000
	argb := graphics.ColorRed.ToArgb()
	if argb != 0xFFFF0000 {
		t.Errorf("ColorRed.ToArgb() = 0x%08X, want 0xFFFF0000", argb)
	}

	// Pure blue should be 0xFF0000FF
	argb = graphics.ColorBlue.ToArgb()
	if argb != 0xFF0000FF {
		t.Errorf("ColorBlue.ToArgb() = 0x%08X, want 0xFF0000FF", argb)
	}

	// Transparent should be 0x00000000
	argb = graphics.ColorTransparent.ToArgb()
	if argb != 0x00000000 {
		t.Errorf("ColorTransparent.ToArgb() = 0x%08X, want 0x00000000", argb)
	}
}

func TestCopy(t *testing.T) {
	original := graphics.ColorRed
	copied := original.Copy(0.5, 0.1, 0.2, 0.3)

	if abs(copied.Alpha()-0.5) > 0.01 {
		t.Errorf("Copied Alpha = %f, want 0.5", copied.Alpha())
	}
	if abs(copied.Red()-0.1) > 0.01 {
		t.Errorf("Copied Red = %f, want 0.1", copied.Red())
	}
}

func TestLuminance(t *testing.T) {
	// White should have luminance close to 1.0
	lum := graphics.ColorWhite.Luminance()
	if abs(lum-1.0) > 0.01 {
		t.Errorf("White Luminance = %f, want ~1.0", lum)
	}

	// Black should have luminance close to 0.0
	lum = graphics.ColorBlack.Luminance()
	if abs(lum-0.0) > 0.01 {
		t.Errorf("Black Luminance = %f, want ~0.0", lum)
	}
}

func TestCompositeOver(t *testing.T) {
	// Semi-transparent red over white
	semiRed := graphics.NewColorSrgb(255, 0, 0, 128)
	result := semiRed.CompositeOver(graphics.ColorWhite)

	// Should be pinkish (mix of red and white)
	if result.Red() < 0.5 {
		t.Errorf("Composite Red = %f, expected > 0.5", result.Red())
	}
	if result.Alpha() < 0.99 {
		t.Errorf("Composite Alpha = %f, want ~1.0", result.Alpha())
	}
}

func TestLerp(t *testing.T) {
	// Lerp from black to white at 0.5 should give some gray
	result := graphics.Lerp(graphics.ColorBlack, graphics.ColorWhite, 0.5)

	// The result should have equal RGB components (gray)
	if abs(result.Red()-result.Green()) > 0.05 {
		t.Errorf("Lerp result should be gray, R=%f G=%f", result.Red(), result.Green())
	}
	if abs(result.Green()-result.Blue()) > 0.05 {
		t.Errorf("Lerp result should be gray, G=%f B=%f", result.Green(), result.Blue())
	}

	// Lerp at 0.0 should give approximately the start color
	result = graphics.Lerp(graphics.ColorRed, graphics.ColorBlue, 0.0)
	if abs(result.Red()-1.0) > 0.05 {
		t.Errorf("Lerp at 0.0: Red = %f, want ~1.0", result.Red())
	}

	// Lerp at 1.0 should give approximately the stop color
	result = graphics.Lerp(graphics.ColorRed, graphics.ColorBlue, 1.0)
	if abs(result.Blue()-1.0) > 0.05 {
		t.Errorf("Lerp at 1.0: Blue = %f, want ~1.0", result.Blue())
	}
}

func TestHsv(t *testing.T) {
	// Pure red: H=0, S=1, V=1
	red := graphics.Hsv(0, 1, 1, 1, nil)
	if abs(red.Red()-1.0) > 0.01 {
		t.Errorf("Hsv(0,1,1) Red = %f, want 1.0", red.Red())
	}
	if abs(red.Green()) > 0.01 {
		t.Errorf("Hsv(0,1,1) Green = %f, want 0.0", red.Green())
	}
	if abs(red.Blue()) > 0.01 {
		t.Errorf("Hsv(0,1,1) Blue = %f, want 0.0", red.Blue())
	}

	// Pure green: H=120, S=1, V=1
	green := graphics.Hsv(120, 1, 1, 1, nil)
	if abs(green.Green()-1.0) > 0.01 {
		t.Errorf("Hsv(120,1,1) Green = %f, want 1.0", green.Green())
	}
}

func TestHsl(t *testing.T) {
	// Pure red: H=0, S=1, L=0.5
	red := graphics.Hsl(0, 1, 0.5, 1, nil)
	if abs(red.Red()-1.0) > 0.01 {
		t.Errorf("Hsl(0,1,0.5) Red = %f, want 1.0", red.Red())
	}

	// White: H=0, S=0, L=1
	white := graphics.Hsl(0, 0, 1, 1, nil)
	if abs(white.Red()-1.0) > 0.01 || abs(white.Green()-1.0) > 0.01 || abs(white.Blue()-1.0) > 0.01 {
		t.Errorf("Hsl(0,0,1) = (%f,%f,%f), want (1,1,1)", white.Red(), white.Green(), white.Blue())
	}
}

func TestNewColorLong(t *testing.T) {
	// Test with a color that would have sign issues as int
	c := graphics.NewColorLong(0xFF000080) // Blue with full alpha

	if abs(c.Alpha()-1.0) > 0.01 {
		t.Errorf("Alpha = %f, want 1.0", c.Alpha())
	}
	if abs(c.Blue()-0.502) > 0.01 {
		t.Errorf("Blue = %f, want ~0.502", c.Blue())
	}
}

func TestConvert(t *testing.T) {
	// Convert sRGB red to DisplayP3 and back
	srgbRed := graphics.ColorRed
	p3Red := srgbRed.Convert(colorspace.DisplayP3)
	backToSrgb := p3Red.Convert(colorspace.Srgb)

	// Should be approximately the same after round trip
	if abs(srgbRed.Red()-backToSrgb.Red()) > 0.05 {
		t.Errorf("Round-trip Red = %f, want %f", backToSrgb.Red(), srgbRed.Red())
	}

	// Converting to same space should return identical color
	samespace := srgbRed.Convert(colorspace.Srgb)
	if samespace != srgbRed {
		t.Error("Converting to same color space should return identical color")
	}
}

func TestString(t *testing.T) {
	s := graphics.ColorRed.String()
	if len(s) == 0 {
		t.Error("String() should not be empty")
	}
	// Should contain "sRGB" in the output
	if !contains(s, "sRGB") {
		t.Errorf("String() = %s, should contain color space name", s)
	}
}

func TestColorToNRGBA(t *testing.T) {
	nrgba := graphics.ColorToNRGBA(graphics.ColorRed)
	if nrgba.R != 255 {
		t.Errorf("NRGBA.R = %d, want 255", nrgba.R)
	}
	if nrgba.G != 0 {
		t.Errorf("NRGBA.G = %d, want 0", nrgba.G)
	}
	if nrgba.B != 0 {
		t.Errorf("NRGBA.B = %d, want 0", nrgba.B)
	}
	if nrgba.A != 255 {
		t.Errorf("NRGBA.A = %d, want 255", nrgba.A)
	}
}

// Helper functions
func abs(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
