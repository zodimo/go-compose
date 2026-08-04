package layoutnode

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
)

// This file converts the transparent alias wall (internal/layoutnode/alias.go)
// into framework-owned defined types. These types flow through every
// exported signature of the public packages; making them defined types (not
// aliases) removes gioui from the public API surface while keeping engine
// access confined to the seam.
//
// The engine context is held as an unexported pointer; consumers above the
// seam interact only with framework-owned values and reach the engine through
// ToGio(). The design's EC1 frame-token guard and the render.Ops-based
// representation (design D2) arrive with the Phase-4 backend; until then
// ToGio() is a plain accessor.

// LayoutContext is the framework-owned layout context passed to components
// and modifiers. It wraps the engine context; ToGio is the only sanctioned
// crossing from framework types to engine types.
type LayoutContext struct {
	gtx *layout.Context
}

// NewLayoutContext wraps a gio layout.Context into a framework-owned
// LayoutContext. It must be called during the frame that owns gtx.
func NewLayoutContext(gtx *layout.Context) LayoutContext {
	return LayoutContext{gtx: gtx}
}

// ToGio returns the underlying engine context. It is seam-only: nothing above
// the seam should call it except engine-bound components (classified Door 2).
func (c LayoutContext) ToGio() *layout.Context {
	return c.gtx
}

// LayoutDimensions is the framework-owned result of a layout pass.
type LayoutDimensions struct {
	Size     image.Point
	Baseline int
}

// ToGioDimensions converts framework dimensions to engine dimensions.
func ToGioDimensions(d LayoutDimensions) layout.Dimensions {
	return layout.Dimensions{Size: d.Size, Baseline: d.Baseline}
}

// FromGioDimensions converts engine dimensions to framework dimensions.
func FromGioDimensions(d layout.Dimensions) LayoutDimensions {
	return LayoutDimensions{Size: d.Size, Baseline: d.Baseline}
}

// LayoutConstraints is the framework-owned measurement constraint pair.
type LayoutConstraints struct {
	Min image.Point
	Max image.Point
}

// ToGioConstraints converts framework constraints to engine constraints.
func ToGioConstraints(c LayoutConstraints) layout.Constraints {
	return layout.Constraints{Min: c.Min, Max: c.Max}
}

// FromGioConstraints converts engine constraints to framework constraints.
func FromGioConstraints(c layout.Constraints) LayoutConstraints {
	return LayoutConstraints{Min: c.Min, Max: c.Max}
}

// GioLayoutWidget is the framework-owned widget constructor signature. The
// name is historical (the project's original name for the gioui layout.Widget
// alias); it is a defined type over framework types.
type GioLayoutWidget func(gtx LayoutContext) LayoutDimensions

// DrawOp is the recorded engine draw result of a frame. It remains an alias
// to the engine macro handle; the seam is the only place that sees it.
type DrawOp = op.CallOp
