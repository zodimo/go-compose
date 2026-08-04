// Package explorerwrap provides a framework-owned wrapper around gioui.org/x/explorer.Explorer
// to keep gioui types out of the public API surface.
package explorerwrap

import "gioui.org/x/explorer"

// Explorer wraps gioui.org/x/explorer.Explorer for the public API surface.
// Use ToGio() to obtain the underlying explorer at the engine seam.
type Explorer struct {
	e *explorer.Explorer
}

// NewExplorer wraps a gioui explorer.Explorer into a framework-owned Explorer.
func NewExplorer(e *explorer.Explorer) *Explorer {
	return &Explorer{e: e}
}

// ToGio returns the underlying gioui explorer.Explorer.
// This is a seam-only accessor for engine-bound code.
func (e *Explorer) ToGio() *explorer.Explorer {
	return e.e
}
