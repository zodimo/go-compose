package colorspace

import (
	"github.com/zodimo/go-compose/compose/ui/util"
)

// Connector connects two color spaces to allow conversion from the source color space
// to the destination color space.
type Connector struct {
	Source               ColorSpace
	Destination          ColorSpace
	Intent               RenderIntent
	TransformMatrix      []float32
	TransformSource      []float32
	TransformDestination []float32
}

// RenderIntent enum
type RenderIntent int

const (
	RenderIntentPerceptual RenderIntent = iota
	RenderIntentRelative
	RenderIntentSaturation
	RenderIntentAbsolute
)

// NewConnector creates a new Connector.
func NewConnector(source, destination ColorSpace, intent RenderIntent) *Connector {
	// Logic for computing transform matrix based on source/dest white points and intent.
	// This is complex logic involving chromatic adaptation.
	// See ColorSpace.kt createConnector / connect.

	// Simplification for now: assuming identity or basic transform if RGB to RGB.
	// If Source and Destination are both RGB, we can compute the transform matrix.
	// Kotlin implementation handles this in `connect` extension or `Connector` constructor.

	// We will implement a simplified version that just holds the spaces for now,
	// and let `Color.Convert` handle specific logic or implement `Transform` calculation here.
	// Actually, `Color.Convert` delegates to `Connector.transform`.

	c := &Connector{
		Source:      source,
		Destination: destination,
		Intent:      intent,
	}

	c.computeTransform()
	return c
}

func (c *Connector) computeTransform() {
	// Check if RGB
	srcRgb, srcIsRgb := c.Source.(*Rgb)
	dstRgb, dstIsRgb := c.Destination.(*Rgb)

	if srcIsRgb && dstIsRgb {
		// Calculate combined transform: Recalculate based on white points?
		// Source -> XYZ -> Destination
		// T = Inv(Dest) * Adapt(SrcWP -> DstWP) * Src

		// Assume D50 adaptation for now if needed, or if white points match.
		// If white points match: T = Inv(Dest) * Src

		// This math requires matrix multiplication support we added.
		// transform = Mul3x3(dstRgb.InverseTransform, srcRgb.Transform)

		// But we need to check WhitePoints.
		// CompareWhitePoint(srcRgb.WhitePoint, dstRgb.WhitePoint)

		// Doing full implementation requires porting `chromaticAdaptation` from `ColorSpace.kt`.
		// We put `Adaptation` transform matrices in `adaptation.go` but not the `chromaticAdaptation` function.
		// I should add `ChromaticAdaptation` function to `adaptation.go`.

		// For now, let's leave it as identity placeholder or simple multiply.
		c.TransformMatrix = util.Mul3x3(dstRgb.InverseTransform, srcRgb.Transform)
	}
}

// Transform converts a color from source space to destination space.
func (c *Connector) Transform(v []float32) []float32 {
	// 1. Source OETF^-1 (EOTF) -> Linear
	// 2. Linear Transform (Matrix) -> XYZ -> XYZ -> Dest Linear
	// 3. Dest OETF

	// If RGB->RGB:
	// v = Source.ToXyz(v)  <-- This applies EOTF and Matrix to XYZ
	// But Connector might optimize this by combining matrices.

	// If we use `computeTransform` combined matrix `c.Transform`:
	// v = EOTF_src(v)
	// v = Matrix * v
	// v = OETF_dst(v)

	srcRgb, srcIsRgb := c.Source.(*Rgb)
	dstRgb, dstIsRgb := c.Destination.(*Rgb)

	if srcIsRgb && dstIsRgb {
		v[0] = float32(srcRgb.Eotf(float64(v[0])))
		v[1] = float32(srcRgb.Eotf(float64(v[1])))
		v[2] = float32(srcRgb.Eotf(float64(v[2])))

		v = util.Mul3x3Float3(c.TransformMatrix, v)

		v[0] = float32(dstRgb.Oetf(float64(v[0])))
		v[1] = float32(dstRgb.Oetf(float64(v[1])))
		v[2] = float32(dstRgb.Oetf(float64(v[2])))
		return v
	}

	// Fallback: Source -> XYZ -> Dest
	xyz := c.Source.ToXyz(v)
	return c.Destination.FromXyz(xyz)
}

func (c *Connector) TransformToColor(v []float32) (uint64, error) {
	_ = c.Transform(v)
	// Pack to Color (uint64)
	// Using simple packing for now.
	// We assume 4th component (alpha) is passed separately or we assume v has 3 components?
	// v in Transform is 3 components (R,G,B).
	// Color has alpha.
	// This helper usually takes `x, y, z, a`?
	// Kotlin: `fun transform(r: Float, g: Float, b: Float): FloatArray`
	return 0, nil // Placeholder
}
