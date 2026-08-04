package render

import "gioui.org/op"

// DrawCommand is the result returned by runtime.Run. It wraps a recorded
// draw sequence that can be applied to an output surface.
//
// The interface is sealed: only types within the render package can
// implement it (via the unexported applyTo method). This prevents
// external packages from creating spurious DrawCommand implementations.
type DrawCommand interface {
	// Apply replays the recorded draw calls into gioui ops.
	// This is the current runtime path's compatibility method:
	// the app shell calls this to replay the frame's draw macro.
	Apply(ops *op.Ops)

	// applyTo applies the draw command to the active backend.
	// This is the backend-agnostic path that will be used once the
	// backend is wired in (task 4.5). It is unexported to seal
	// the interface.
	applyTo(b Backend)
}

// gioDrawCommand wraps a gioui op.CallOp for the current runtime path.
type gioDrawCommand struct {
	callOp op.CallOp
}

// Compile-time check that gioDrawCommand implements DrawCommand.
var _ DrawCommand = (*gioDrawCommand)(nil)

// NewDrawCommand wraps a gioui macro handle into an opaque DrawCommand.
// This is the current entry point used by runtime.Run.
func NewDrawCommand(callOp op.CallOp) DrawCommand {
	return &gioDrawCommand{callOp: callOp}
}

// Apply replays the recorded draw calls into ops.
// It is app-shell glue: the demo/CLI shell passes the gio Ops of the
// current frame.
func (d *gioDrawCommand) Apply(ops *op.Ops) {
	d.callOp.Add(ops)
}

// applyTo is a no-op for the gio path. The real backend-driven wiring
// happens in task 4.5 when runtime.Run calls the active backend directly.
func (d *gioDrawCommand) applyTo(b Backend) {
	// intentionally empty — current path uses Apply(*op.Ops) directly
}
