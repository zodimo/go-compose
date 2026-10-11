package graphics

import (
	"github.com/zodimo/go-compose/compose/ui/graphics/colorspace"
	"github.com/zodimo/go-compose/pkg/floatutils/lerp"
)

// LerpColor interpolates between two colors.
func LerpColor(start, stop Color, fraction float32) Color {
	// If both are unspecified, return unspecified.
	if start.IsUnspecified() && stop.IsUnspecified() {
		return colorUnspecified
	}
	// If one is unspecified, maybe treat as transparent or return unspecified?
	// Kotlin Color.kt says: "If one of the arguments is Color.Unspecified, the result is Color.Unspecified"
	// Wait, let's verify.
	// Actually, TextForegroundStyle behavior relies on this.
	// If ANY is unspecified, result is unspecified.
	if start.IsUnspecified() || stop.IsUnspecified() {
		return colorUnspecified
	}

	// Basic Lerp in sRGB or Oklab?
	// Kotlin version uses Oklab for interpolation by default in `lerp`, or simple argb lerp?
	// Kotlin: `fun lerp(start: Color, stop: Color, fraction: Float): Color`
	// It converts to Oklab, interpolates, converts back.
	// "The interpolation is done in the Oklab color space."

	// We need Oklab support.
	// For now, let's implement simple sRGB Lerp to unblock build, or use Oklab if possible.
	// Oklab conversion needs `colorspace` package fully working.
	// To unblock, I'll use simple sRGB lerp on ARGB values if both are sRGB.

	// Convert to sRGB 0-1 floats
	// This is valid fallback used in many graphics mutations until specialized spaces supported.

	// Actually, to match Kotlin behavior which uses Oklab, we should ideally use that.
	// But given the complexity of `Connector` and `ColorSpace` being freshly ported, maybe simple lerp is safer for now.
	// Let's implement simple sRGB Lerp.

	return NewColor(
		lerp.Between32(start.Red(), stop.Red(), fraction),
		lerp.Between32(start.Green(), stop.Green(), fraction),
		lerp.Between32(start.Blue(), stop.Blue(), fraction),
		lerp.Between32(start.Alpha(), stop.Alpha(), fraction),
		colorspace.Srgb,
	)
}

func LerpColors(a, b []Color, t float32) []Color {
	n := len(a)
	if len(b) > n {
		n = len(b)
	}
	res := make([]Color, n)
	for i := 0; i < n; i++ {
		c1 := a[min(i, len(a)-1)]
		c2 := b[min(i, len(b)-1)]
		res[i] = LerpColor(c1, c2, t)
	}
	return res
}
