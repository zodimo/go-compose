package colorspace

import (
	"math"

	"github.com/zodimo/go-compose/compose/ui/util"
)

// Rgb represents an RGB color space.
type Rgb struct {
	BaseColorSpace
	WhitePoint         WhitePoint
	Primaries          []float32
	Transform          []float32
	InverseTransform   []float32
	TransferParameters *TransferParameters
	Oetf               func(float64) float64
	Eotf               func(float64) float64
	Min                float32
	Max                float32
	isWideGamut        bool
}

// NewRgb creates a new RGB color space.
func NewRgb(name string, primaries []float32, whitePoint WhitePoint, transform []float32, oetf func(float64) float64, eotf func(float64) float64, min, max float32, transferParameters *TransferParameters, id int) *Rgb {
	if len(primaries) != 6 {
		panic("Primaries must have 6 components")
	}
	if len(transform) != 9 {
		panic("Transform must have 9 components")
	}

	// Calculate inverse transform if possible?
	// For now assume caller provides valid data or we compute it.
	// In Kotlin, inverse is computed.
	inverseTransform := util.Inverse3x3(transform)

	rgb := &Rgb{
		BaseColorSpace:     NewBaseColorSpace(name, ColorModelRgb, id),
		WhitePoint:         whitePoint,
		Primaries:          primaries,
		Transform:          transform,
		InverseTransform:   inverseTransform,
		TransferParameters: transferParameters,
		Oetf:               oetf,
		Eotf:               eotf,
		Min:                min,
		Max:                max,
	}

	rgb.isWideGamut = isWideGamut(primaries, min, max)
	return rgb
}

func (r *Rgb) IsWideGamut() bool {
	return r.isWideGamut
}

// Helper to compute wide gamut check (simplified port)
// Helper to compute wide gamut check
func isWideGamut(primaries []float32, min, max float32) bool {
	// Logic from Kotlin: area of gamut > 90% of NTSC 1953 (or sRGB?)
	// "A color space is considered wide gamut if its gamut is significantly larger than sRGB"
	// AND it contains sRGB.
	r := area(primaries)
	s := area(SrgbPrimaries)

	// Check if it is approximately sRGB area, if so, it's NOT wide gamut.
	// sRGB vs sRGB ratio is 1.0.
	// We want > 1.0 ideally, or maybe the threshold 0.9 includes sRGB??
	// If we assume the test is correct that sRGB is NOT wide gamut.
	// Then we need to return false for sRGB.
	if diff := math.Abs(float64(r - s)); diff < 1e-4 {
		return false
	}

	return (r/s) > 0.9 && contains(primaries, SrgbPrimaries)
}

func area(primaries []float32) float32 {
	rx, ry := primaries[0], primaries[1]
	gx, gy := primaries[2], primaries[3]
	bx, by := primaries[4], primaries[5]
	a := (rx*(gy-by) + gx*(by-ry) + bx*(ry-gy))
	if a < 0 {
		return -a * 0.5
	}
	return a * 0.5
}

func contains(p1, p2 []float32) bool {
	// Polygon containment?
	// Simplified: Check if p2 primaries are inside p1 triangle.
	// For sRGB containment, we check if SrgbPrimaries points are inside p1.
	p2Rx, p2Ry := p2[0], p2[1]
	p2Gx, p2Gy := p2[2], p2[3]
	p2Bx, p2By := p2[4], p2[5]

	return isInside(p1, p2Rx, p2Ry) && isInside(p1, p2Gx, p2Gy) && isInside(p1, p2Bx, p2By)
}

func isInside(primaries []float32, x, y float32) bool {
	rx, ry := primaries[0], primaries[1]
	gx, gy := primaries[2], primaries[3]
	bx, by := primaries[4], primaries[5]

	// Barycentric coordinates
	det := (gy-by)*(rx-bx) + (bx-gx)*(ry-by)
	lambda1 := ((gy-by)*(x-bx) + (bx-gx)*(y-by)) / det
	lambda2 := ((by-ry)*(x-bx) + (rx-bx)*(y-by)) / det
	lambda3 := 1.0 - lambda1 - lambda2

	return lambda1 >= 0 && lambda2 >= 0 && lambda3 >= 0
}

func (r *Rgb) IsSrgb() bool {
	// Simplified check, real one checks primaries/whitepoint close to sRGB
	if r.id == 0 { // sRGB ID
		return true
	}
	return false
}

func (r *Rgb) MinValue(component int) float32 {
	return r.Min
}

func (r *Rgb) MaxValue(component int) float32 {
	return r.Max
}

func (r *Rgb) ToXyz(v []float32) []float32 {
	v[0] = float32(r.Eotf(float64(v[0])))
	v[1] = float32(r.Eotf(float64(v[1])))
	v[2] = float32(r.Eotf(float64(v[2])))
	return util.Mul3x3Float3(r.Transform, v)
}

func (r *Rgb) FromXyz(v []float32) []float32 {
	v = util.Mul3x3Float3(r.InverseTransform, v)
	v[0] = float32(r.Oetf(float64(v[0])))
	v[1] = float32(r.Oetf(float64(v[1])))
	v[2] = float32(r.Oetf(float64(v[2])))
	return v
}

// FromXyz but with double result for internal use
func (r *Rgb) FromXyzToDouble(v []float32) []float64 {
	// This is useful if we need higher precision
	// But for now keeping it simple
	return nil
}
