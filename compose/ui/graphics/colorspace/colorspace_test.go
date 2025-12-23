package colorspace

import (
	"math"
	"testing"
)

func TestSrgb(t *testing.T) {
	srgb := Srgb
	if !srgb.IsSrgb() {
		t.Errorf("Expected IsSrgb to be true for Srgb colorspace")
	}
	if srgb.IsWideGamut() {
		t.Errorf("Expected IsWideGamut to be false for Srgb colorspace")
	}
	if srgb.Name() != "sRGB IEC61966-2.1" {
		t.Errorf("Unexpected name: %s", srgb.Name())
	}

	// Test ToXyz/FromXyz roundtrip (approximate)
	// Gray 0.5
	r, g, b := float32(0.5), float32(0.5), float32(0.5)
	xyz := srgb.ToXyz([]float32{r, g, b})
	rgb := srgb.FromXyz(xyz)
	r2, g2, b2 := rgb[0], rgb[1], rgb[2]

	if math.Abs(float64(r-r2)) > 1e-4 || math.Abs(float64(g-g2)) > 1e-4 || math.Abs(float64(b-b2)) > 1e-4 {
		t.Errorf("Roundtrip failed: %v,%v,%v -> %v,%v,%v -> %v,%v,%v", r, g, b, xyz[0], xyz[1], xyz[2], r2, g2, b2)
	}
}

func TestOklab(t *testing.T) {
	// Test Oklab
	oklab := OklabInstance
	if !oklab.IsWideGamut() {
		t.Errorf("Expected Oklab to be wide gamut")
	}
	if oklab.Name() != "Oklab" {
		t.Errorf("Expected Oklab name, got %s", oklab.Name())
	}

	// Test conversions if implemented
	// Since XyzaToColor returns placeholder, we might test ToXyz/FromXyz
	// Oklab ToXyz/FromXyz IS implemented in oklab.go

	// Test vector
	l, a, bArg := float32(0.6), float32(0.1), float32(-0.1)
	xyz := oklab.ToXyz([]float32{l, a, bArg})
	lab := oklab.FromXyz(xyz)
	l2, a2, b2 := lab[0], lab[1], lab[2]

	if math.Abs(float64(l-l2)) > 1e-4 || math.Abs(float64(a-a2)) > 1e-4 || math.Abs(float64(bArg-b2)) > 1e-4 {
		t.Errorf("Oklab Roundtrip failed: %v,%v,%v -> %v,%v,%v -> %v,%v,%v", l, a, bArg, xyz[0], xyz[1], xyz[2], l2, a2, b2)
	}
}

func TestGetColorSpace(t *testing.T) {
	srgb := Get(ColorSpaceSrgb)
	if srgb != Srgb {
		t.Errorf("Get(ColorSpaceSrgb) returned incorrect instance")
	}

	oklab := Get(ColorSpaceOklab)
	if oklab != OklabInstance {
		t.Errorf("Get(ColorSpaceOklab) returned incorrect instance")
	}
}
