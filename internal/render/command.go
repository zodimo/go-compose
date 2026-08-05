package render

import "gioui.org/op"

// DrawCommand is the result returned by runtime.Run. It wraps a recorded
// draw sequence that can be applied to an output surface.
//
// The interface is sealed: only types within the render package can
// implement it (via the unexported applyTo method). This prevents
// external packages from creating spurious DrawCommand implementations.
type DrawCommand interface {
	// applyTo applies the draw command to the active backend.
	// This is the backend-agnostic path that will be used once the
	// backend is wired in. It is unexported to seal the interface.
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

// ApplyToGio replays the recorded draw calls into gioui ops.
// It is package-level glue for runtime/ to replay the frame's draw macro.
func ApplyToGio(cmd DrawCommand, ops *op.Ops) {
	if d, ok := cmd.(*gioDrawCommand); ok {
		d.callOp.Add(ops)
	}
}

// applyTo is a no-op for the gio path.
func (d *gioDrawCommand) applyTo(b Backend) {
	// intentionally empty — current path uses ApplyToGio directly
}

// --- Backend-based DrawCommand (software / future backends) ---

// backendDrawCommand wraps a render.Backend for the software rendering path.
type backendDrawCommand struct {
	backend Backend
}

// Compile-time check that backendDrawCommand implements DrawCommand.
var _ DrawCommand = (*backendDrawCommand)(nil)

// NewDrawCommandForBackend creates a DrawCommand from any Backend.
func NewDrawCommandForBackend(b Backend) DrawCommand {
	return &backendDrawCommand{backend: b}
}

// applyTo is unused on the software path.
func (d *backendDrawCommand) applyTo(b Backend) {
	// intentionally empty
}
