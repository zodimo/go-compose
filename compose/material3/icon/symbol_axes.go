package icon

import (
	"github.com/go-text/typesetting/font"
	"github.com/go-text/typesetting/font/opentype"
)

// Supported ranges for the Material Symbols variable axes, in design units, as
// declared by the font's fvar table.
const (
	FillMin    float32 = 0
	FillMax    float32 = 1
	GradeMin   float32 = -50
	GradeMax   float32 = 200
	OpticalMin float32 = 20
	OpticalMax float32 = 48
	WeightMin  float32 = 100
	WeightMax  float32 = 700
)

// axisValues holds the effective, clamped variation-axis values for one render.
// It is the canonical form the cache keys on, so two calls that differ only in a
// clamped-away out-of-range value resolve to the same entry.
type axisValues struct {
	Fill        float32
	Grade       float32
	OpticalSize float32
	Weight      float32
}

// defaultAxisValues returns the font's default instance: FILL 0, GRAD 0,
// opsz 24, wght 400.
func defaultAxisValues() axisValues {
	return axisValues{
		Fill:        0,
		Grade:       0,
		OpticalSize: 24,
		Weight:      400,
	}
}

// resolveAxes clamps the caller's axis options and fills in defaults. Explicit
// values are authoritative; anything unprovided keeps the font default, except
// optical size, which adapts to the rendered size so small icons get the
// small-size design and large icons get the large-size design.
func resolveAxes(opts IconOptions, sizePx int) axisValues {
	out := defaultAxisValues()

	if opts.hasFill {
		out.Fill = clampAxis(opts.Fill, FillMin, FillMax)
	}
	if opts.hasGrade {
		out.Grade = clampAxis(opts.Grade, GradeMin, GradeMax)
	}
	if opts.hasWeight {
		out.Weight = clampAxis(opts.Weight, WeightMin, WeightMax)
	}

	if opts.hasOpticalSz {
		// An explicit optical size always wins over the size-derived one.
		out.OpticalSize = clampAxis(opts.OpticalSize, OpticalMin, OpticalMax)
	} else {
		out.OpticalSize = opticalSizeForPx(sizePx)
	}

	return out
}

// opticalSizeForPx maps a rendered icon size onto the optical size axis.
//
// The axis is expressed in the same units as a CSS pixel size, so the icon's
// pixel height is the natural input. Values outside the axis range clamp to its
// ends: a 16dp icon takes the smallest optical design, a 64dp icon the largest.
func opticalSizeForPx(sizePx int) float32 {
	return clampAxis(float32(sizePx), OpticalMin, OpticalMax)
}

// clampAxis restricts v to [lo, hi]. NaN resolves to lo so a caller cannot poison
// the cache key with an unusable value.
func clampAxis(v, lo, hi float32) float32 {
	if v != v { // NaN
		return lo
	}
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// variations converts design-space axis values into the variation settings the
// font face expects.
//
// SetVariations takes design units, not normalized coordinates, so the values
// map straight through. Only axes that differ from the font default are sent,
// which keeps the common case (a plain 24dp outlined icon) on the default
// instance.
func (a axisValues) variations() []font.Variation {
	def := defaultAxisValues()
	out := make([]font.Variation, 0, 4)

	if a.Fill != def.Fill {
		out = append(out, font.Variation{Tag: opentype.MustNewTag(axisFill), Value: a.Fill})
	}
	if a.Grade != def.Grade {
		out = append(out, font.Variation{Tag: opentype.MustNewTag(axisGrad), Value: a.Grade})
	}
	if a.OpticalSize != def.OpticalSize {
		out = append(out, font.Variation{Tag: opentype.MustNewTag(axisOpsz), Value: a.OpticalSize})
	}
	if a.Weight != def.Weight {
		out = append(out, font.Variation{Tag: opentype.MustNewTag(axisWght), Value: a.Weight})
	}

	return out
}
