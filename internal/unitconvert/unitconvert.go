// Package unitconvert provides conversion functions between go-compose unit types
// and gioui.org unit types. These functions were moved from compose/ui/unit to
// keep gioui.org types out of the public API surface.
package unitconvert

import (
	"errors"

	gioLayout "gioui.org/layout"
	gioUnit "gioui.org/unit"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

// DpToGioUnit converts a go-compose Dp to a gio unit.Dp.
func DpToGioUnit(u unit.Dp) (gioUnit.Dp, error) {
	if u.IsUnspecified() {
		return 0, errors.New("Dp is an unspecified unit, cannot convert to GioUnit")
	}
	return gioUnit.Dp(float32(u)), nil
}

// DpToGioUnitUnsafe converts a go-compose Dp to a gio unit.Dp, panicking on error.
func DpToGioUnitUnsafe(u unit.Dp) gioUnit.Dp {
	g, err := DpToGioUnit(u)
	if err != nil {
		panic(err)
	}
	return g
}

// TextUnitToGioSp converts a go-compose TextUnit to a gio unit.Sp.
func TextUnitToGioSp(tu unit.TextUnit) (gioUnit.Sp, error) {
	if tu.IsUnspecified() {
		return 0, errors.New("TextUnit is an unspecified unit, cannot convert to Sp")
	}
	if tu.IsSp() {
		return gioUnit.Sp(tu.Value()), nil
	}
	return 0, errors.New("TextUnit is an EM unit, cannot convert to Sp")
}

// TextUnitToGioDp converts a go-compose TextUnit to a gio unit.Dp using density.
func TextUnitToGioDp(tu unit.TextUnit, density float32) (gioUnit.Dp, error) {
	if tu.IsUnspecified() {
		return 0, errors.New("TextUnit is an unspecified unit, cannot convert to Dp")
	}
	if tu.IsSp() {
		return 0, errors.New("TextUnit is an Sp unit, cannot convert to Dp")
	}
	return gioUnit.Dp(tu.Value() * density), nil
}

// TextUnitToGioSpUnsafe converts a go-compose TextUnit to a gio unit.Sp, panicking on error.
func TextUnitToGioSpUnsafe(tu unit.TextUnit) gioUnit.Sp {
	s, err := TextUnitToGioSp(tu)
	if err != nil {
		panic(err)
	}
	return s
}

// DensityFromLayoutContext creates a Density from a Gio Layout Context.
func DensityFromLayoutContext(gtx gioLayout.Context) unit.Density {
	return unit.NewDensity(gtx.Metric.PxPerDp, gtx.Metric.PxPerSp)
}
