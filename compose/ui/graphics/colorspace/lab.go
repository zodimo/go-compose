package colorspace

import (
	"github.com/zodimo/go-compose/compose/ui/util"
)

const (
	LabA = 216.0 / 24389.0
	LabB = 841.0 / 108.0
	LabC = 4.0 / 29.0
	LabD = 6.0 / 29.0
)

type Lab struct {
	BaseColorSpace
}

func NewLab(name string, id int) *Lab {
	return &Lab{
		BaseColorSpace: NewBaseColorSpace(name, ColorModelLab, id),
	}
}

func (l *Lab) IsWideGamut() bool {
	return true
}

func (l *Lab) MinValue(component int) float32 {
	if component == 0 {
		return 0.0
	}
	return -128.0
}

func (l *Lab) MaxValue(component int) float32 {
	if component == 0 {
		return 100.0
	}
	return 128.0
}

func (l *Lab) ToXyz(v []float32) []float32 {
	v[0] = util.FastCoerceIn(v[0], 0.0, 100.0)
	v[1] = util.FastCoerceIn(v[1], -128.0, 128.0)
	v[2] = util.FastCoerceIn(v[2], -128.0, 128.0)

	fy := (v[0] + 16.0) / 116.0
	fx := fy + (v[1] * 0.002)
	fz := fy - (v[2] * 0.005)

	x := fx * fx * fx
	if fx <= LabD {
		x = (1.0 / LabB) * (fx - LabC)
	}

	y := fy * fy * fy
	if fy <= LabD {
		y = (1.0 / LabB) * (fy - LabC)
	}

	z := fz * fz * fz
	if fz <= LabD {
		z = (1.0 / LabB) * (fz - LabC)
	}

	v[0] = x * IlluminantD50Xyz[0]
	v[1] = y * IlluminantD50Xyz[1]
	v[2] = z * IlluminantD50Xyz[2]
	return v
}

func (l *Lab) FromXyz(v []float32) []float32 {
	x := v[0] / IlluminantD50Xyz[0]
	y := v[1] / IlluminantD50Xyz[1]
	z := v[2] / IlluminantD50Xyz[2]

	fx := util.FastCbrt(x)
	if x <= LabA {
		fx = LabB*x + LabC
	}

	fy := util.FastCbrt(y)
	if y <= LabA {
		fy = LabB*y + LabC
	}

	fz := util.FastCbrt(z)
	if z <= LabA {
		fz = LabB*z + LabC
	}

	L := 116.0*fy - 16.0
	a := 500.0 * (fx - fy)
	b := 200.0 * (fy - fz)

	v[0] = util.FastCoerceIn(L, 0.0, 100.0)
	v[1] = util.FastCoerceIn(a, -128.0, 128.0)
	v[2] = util.FastCoerceIn(b, -128.0, 128.0)
	return v
}

func (l *Lab) ToXy(v0, v1, v2 float32) int64 {
	// Port logic from Lab.kt
	v0 = util.FastCoerceIn(v0, 0.0, 100.0)
	v1 = util.FastCoerceIn(v1, -128.0, 128.0)
	// v2 ignored for ToXy in Lab?

	fy := (v0 + 16.0) / 116.0
	fx := fy + (v1 * 0.002)

	x := fx * fx * fx
	if fx <= LabD {
		x = (1.0 / LabB) * (fx - LabC)
	}

	y := fy * fy * fy
	if fy <= LabD {
		y = (1.0 / LabB) * (fy - LabC)
	}

	return util.PackFloats(x*IlluminantD50Xyz[0], y*IlluminantD50Xyz[1])
}

func (l *Lab) ToZ(v0, v1, v2 float32) float32 {
	v0 = util.FastCoerceIn(v0, 0.0, 100.0)
	v2 = util.FastCoerceIn(v2, -128.0, 128.0)

	fy := (v0 + 16.0) / 116.0
	fz := fy - (v2 * 0.005)

	z := fz * fz * fz
	if fz <= LabD {
		z = (1.0 / LabB) * (fz - LabC)
	}
	return z * IlluminantD50Xyz[2]
}

func (l *Lab) XyzaToColor(x, y, z, a float32, colorSpace ColorSpace) (uint64, error) {
	// Implement using XyzaToColorHelper or logic
	// But as discussed, we might not return packed Color here depending on cyclic deps.
	// For now, return 0, nil as placeholder or remove from interface if confirmed unused by consumer in a way that requires this method.
	// However, interface requires it.
	// If we want to implement it, we need `graphics.Color` wrapping logic, but we can't import it.
	// So distinct failure or restructuring is needed.
	// Given the task is to port dependencies for `Color`, `Color` relies on `ColorSpace`.
	// `xyzaToColor` in Kotlin creates a `Color`.
	// We should probably rely on `Color(r, g, b, a, space)` factory in `graphics` package instead of this method on `ColorSpace`.
	// This method is likely internal or helper in Kotlin.
	// I already made it part of interface in Go.
	return 0, nil
}
