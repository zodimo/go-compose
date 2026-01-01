package unit

import (
	gioUnit "gioui.org/unit"
)

func DpToGioUnit(u Dp) gioUnit.Dp {
	return gioUnit.Dp(float32(u))
}

func TextUnitToGioSp(tu TextUnit) gioUnit.Sp {
	if tu.IsUnspecified() {
		return gioUnit.Sp(0)
	}
	if tu.IsEm() {
		panic("TextUnit is an EM unit, cannot convert to Sp")
	}
	return gioUnit.Sp(tu.Value())
}
