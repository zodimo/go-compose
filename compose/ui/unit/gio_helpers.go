package unit

import (
	"errors"

	gioLayout "gioui.org/layout"
	gioUnit "gioui.org/unit"
)

func DpToGioUnit(u Dp) (gioUnit.Dp, error) {
	if u.IsUnspecified() {
		return 0, errors.New("Dp is an unspecified unit, cannot convert to GioUnit")
	}
	return gioUnit.Dp(float32(u)), nil
}

func DpToGioUnitUnsafe(u Dp) gioUnit.Dp {
	unit, err := DpToGioUnit(u)
	if err != nil {
		panic(err)
	}
	return unit
}

func TextUnitToGioSp(tu TextUnit) (gioUnit.Sp, error) {
	if tu.IsUnspecified() {
		return 0, errors.New("TextUnit is an unspecified unit, cannot convert to Sp")
	}
	if tu.IsSp() {
		return gioUnit.Sp(tu.Value()), nil
	}

	return 0, errors.New("TextUnit is an EM unit, cannot convert to Sp")
}

func TextUnitToGioDp(tu TextUnit, density float32) (gioUnit.Dp, error) {
	if tu.IsUnspecified() {
		return 0, errors.New("TextUnit is an unspecified unit, cannot convert to Dp")
	}
	if tu.IsSp() {
		return 0, errors.New("TextUnit is an Sp unit, cannot convert to Dp")
	}
	return gioUnit.Dp(tu.Value() * density), nil
}

func TextUnitToGioSpUnsafe(tu TextUnit) gioUnit.Sp {
	unit, err := TextUnitToGioSp(tu)
	if err != nil {
		panic(err)
	}
	return unit
}

// DensityFromLayoutContext creates a Density from a Gio Layout Context.
// gtx: The Gio Layout Context.
func DensityFromLayoutContext(gtx gioLayout.Context) Density {
	return NewDensity(gtx.Metric.PxPerDp, gtx.Metric.PxPerSp)
}
