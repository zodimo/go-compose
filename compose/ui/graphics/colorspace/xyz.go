package colorspace

import (
	"github.com/zodimo/go-compose/compose/ui/util"
)

type Xyz struct {
	BaseColorSpace
}

func NewXyz(name string, id int) *Xyz {
	return &Xyz{
		BaseColorSpace: NewBaseColorSpace(name, ColorModelXyz, id),
	}
}

func (x *Xyz) IsWideGamut() bool {
	return true
}

func (x *Xyz) MinValue(component int) float32 {
	return -2.0
}

func (x *Xyz) MaxValue(component int) float32 {
	return 2.0
}

func (x *Xyz) ToXyz(v []float32) []float32 {
	v[0] = util.FastCoerceIn(v[0], -2.0, 2.0)
	v[1] = util.FastCoerceIn(v[1], -2.0, 2.0)
	v[2] = util.FastCoerceIn(v[2], -2.0, 2.0)
	return v
}

func (x *Xyz) FromXyz(v []float32) []float32 {
	v[0] = util.FastCoerceIn(v[0], -2.0, 2.0)
	v[1] = util.FastCoerceIn(v[1], -2.0, 2.0)
	v[2] = util.FastCoerceIn(v[2], -2.0, 2.0)
	return v
}
