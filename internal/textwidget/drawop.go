package textwidget

import "gioui.org/op"

// DrawOp wraps op.CallOp for paint materials, hiding it from the public API surface.
type DrawOp struct {
	O op.CallOp
}

// Add adds the recorded operation to the operation list.
func (d DrawOp) Add(ops *op.Ops) {
	d.O.Add(ops)
}

// IsEmpty reports whether the DrawOp was recorded.
func (d DrawOp) IsEmpty() bool {
	return d.O == op.CallOp{}
}

// NewDrawOp creates a DrawOp from a raw op.CallOp.
func NewDrawOp(o op.CallOp) DrawOp {
	return DrawOp{O: o}
}
