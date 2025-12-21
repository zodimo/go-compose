package lerp

import (
	"math"

	"github.com/zodimo/go-compose/internal/color/f32color"
)

// lerp calculates linear interpolation with color b and p.
// Prefer ColorLerpPrecise for better color interpolation.
func ColorLerp(a, b f32color.RGBA, p float32) f32color.RGBA {
	return f32color.RGBA{
		R: a.R*(1-p) + b.R*p,
		G: a.G*(1-p) + b.G*p,
		B: a.B*(1-p) + b.B*p,
		A: a.A*(1-p) + b.A*p,
	}
}

func ColorLerpPrecise(a, b f32color.RGBA, p float32) f32color.RGBA {
	invP := 1 - p
	return f32color.RGBA{
		// Interpolating squares preserves perceived brightness
		R: sqrt(a.R*a.R*invP + b.R*b.R*p),
		G: sqrt(a.G*a.G*invP + b.G*b.G*p),
		B: sqrt(a.B*a.B*invP + b.B*b.B*p),
		A: a.A*invP + b.A*p,
	}
}

// Linearly interpolate between [start] and [stop] with [fraction] fraction between them.
func FloatLerp(start, stop, fraction float32) float32 {
	return (1-fraction)*start + fraction*stop
}

// Linearly interpolate between [start] and [stop] with [fraction] fraction between them.
func IntLerp(start, stop int, fraction float32) int {
	return int(FloatLerp(float32(start), float32(stop), fraction))
}

// Linearly interpolate between [start] and [stop] with [fraction] fraction between them.
// fraction is 0-256 (where 256 is 100%)
func IntLerpFixed(start, stop, fraction int) int {
	return start + (stop-start)*fraction>>8
}

// Linearly interpolate between [start] and [stop] with [fraction] fraction between them.
func IntLerpPrecise(start, stop int, fraction float32) int {
	f64 := float64(start) + float64(stop-start)*float64(fraction)
	return int(f64 + 0.5)
}

func sqrt(v float32) float32 {
	return float32(math.Sqrt(float64(v)))
}
