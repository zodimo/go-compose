// Package render is the draw-IR seam: the framework-facing surface of the
// rendering backends. In this change it is a stub — the opaque types that
// flow across the seam — with real behavior (backend registration, gio and
// software backends) arriving in Phase 4 of the abstract-render-backend
// change.
package render

import "gioui.org/op"

// CallOp is an opaque handle to a recorded sequence of draw calls. It is
// produced by a backend's Record and applied with Add; only the backend that
// created it can apply it. Phase 4 gives it real behavior.
type CallOp struct{}

// DrawCommand is an opaque frame result returned by runtime.Run. The active
// backend applies it to the output surface. Until Phase 4, it carries the
// gioui macro handle produced by the current runtime path; the demo app shell
// applies it with Apply.
type DrawCommand struct {
	callOp op.CallOp
}

// NewDrawCommand wraps a gioui macro handle into an opaque DrawCommand.
func NewDrawCommand(callOp op.CallOp) DrawCommand {
	return DrawCommand{callOp: callOp}
}

// Apply replays the recorded draw calls into ops. It is app-shell glue: the
// demo/CLI shell passes the gio Ops of the current frame.
func (d DrawCommand) Apply(ops *op.Ops) {
	d.callOp.Add(ops)
}
