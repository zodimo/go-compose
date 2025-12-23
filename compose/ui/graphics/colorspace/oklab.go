package colorspace

import (
	"github.com/zodimo/go-compose/compose/ui/util"
)

type Oklab struct {
	BaseColorSpace
}

var (
	// Matrices for M1 and M2 in Oklab spec
	M1 = []float32{
		0.8189330101, 0.3618667424, -0.1288597137,
		0.0329845436, 0.9293118715, 0.0361456387,
		0.0482003018, 0.2643662691, 0.6338517070,
	}
	M2 = []float32{
		0.2104542553, 1.9779984951, 0.0259040371,
		0.7936177850, -2.4285922050, 0.7827717662,
		-0.0040720468, 0.4505937, -0.8086757660,
	}
	InverseM1 = util.Inverse3x3(M1)
	InverseM2 = util.Inverse3x3(M2)
)

func NewOklab(name string, id int) *Oklab {
	return &Oklab{
		BaseColorSpace: NewBaseColorSpace(name, ColorModelLab, id),
	}
}

func (o *Oklab) IsWideGamut() bool {
	return true
}

func (o *Oklab) MinValue(component int) float32 {
	if component == 0 {
		return 0.0
	}
	return -0.5 // Roughly -0.5 to 0.5 for a and b
}

func (o *Oklab) MaxValue(component int) float32 {
	if component == 0 {
		return 1.0
	}
	return 0.5
}

func (o *Oklab) ToXyz(v []float32) []float32 {
	v[0] = util.FastCoerceIn(v[0], 0, 1)
	v[1] = util.FastCoerceIn(v[1], -0.5, 0.5)
	v[2] = util.FastCoerceIn(v[2], -0.5, 0.5)

	util.Mul3x3Float3(InverseM2, v)
	v[0] = v[0] * v[0] * v[0]
	v[1] = v[1] * v[1] * v[1]
	v[2] = v[2] * v[2] * v[2]
	util.Mul3x3Float3(InverseM1, v)

	return v
}

func (o *Oklab) FromXyz(v []float32) []float32 {
	util.Mul3x3Float3(M1, v)

	v[0] = util.FastCbrt(v[0])
	v[1] = util.FastCbrt(v[1])
	v[2] = util.FastCbrt(v[2])

	util.Mul3x3Float3(M2, v)
	return v
}

func (o *Oklab) ToXy(v0, v1, v2 float32) int64 {
	v := []float32{v0, v1, v2}
	xyz := o.ToXyz(v)
	return util.PackFloats(xyz[0], xyz[1])
}

func (o *Oklab) ToZ(v0, v1, v2 float32) float32 {
	v := []float32{v0, v1, v2}
	xyz := o.ToXyz(v)
	return xyz[2]
}

func (o *Oklab) XyzaToColor(x, y, z, a float32, colorSpace ColorSpace) (uint64, error) {
	return 0, nil
}
