// Package clipconvert provides framework-owned wrapper types for gioui.org clip
// operations. These types hide gioui types from the public API surface while
// allowing internal code to interoperate with the rendering engine.
package clipconvert

import (
	"gioui.org/op"
	"gioui.org/op/clip"
)

// Ops wraps *op.Ops for the public API surface.
type Ops struct {
	ops *op.Ops
}

// NewOps wraps a gioui op.Ops into a framework-owned Ops.
func NewOps(ops *op.Ops) *Ops {
	return &Ops{ops: ops}
}

// ToGio returns the underlying gioui op.Ops. Seam-only: for engine-bound code.
func (o *Ops) ToGio() *op.Ops {
	return o.ops
}

// Stack wraps clip.Stack for the public API surface.
type Stack struct {
	stack clip.Stack
}

// NewStack wraps a gioui clip.Stack into a framework-owned Stack.
func NewStack(s clip.Stack) Stack {
	return Stack{stack: s}
}

// Pop removes the clip applied by Push.
func (s Stack) Pop() {
	s.stack.Pop()
}

// ToGio returns the underlying gioui clip.Stack. Seam-only: for engine-bound code.
func (s Stack) ToGio() clip.Stack {
	return s.stack
}

// ClipOp wraps clip.Op for the public API surface.
// Use ToClipOp() to obtain the underlying gioui clip.Op for engine-bound code.
type ClipOp struct {
	op clip.Op
}

// NewClipOp wraps a gioui clip.Op into a framework-owned ClipOp.
func NewClipOp(o clip.Op) ClipOp {
	return ClipOp{op: o}
}

// ToClipOp returns the underlying gioui clip.Op. Seam-only: for engine-bound code.
func (o ClipOp) ToClipOp() clip.Op {
	return o.op
}

// PathSpec wraps clip.PathSpec for the public API surface.
type PathSpec struct {
	spec clip.PathSpec
}

// NewPathSpec wraps a gioui clip.PathSpec into a framework-owned PathSpec.
func NewPathSpec(p clip.PathSpec) PathSpec {
	return PathSpec{spec: p}
}

// ToClipPathSpec returns the underlying gioui clip.PathSpec. Seam-only: for engine-bound code.
func (p PathSpec) ToClipPathSpec() clip.PathSpec {
	return p.spec
}
