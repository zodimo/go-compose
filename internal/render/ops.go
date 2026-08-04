package render

import "sync/atomic"

var frameCounter uint64

// Ops is an opaque per-frame draw list handle.
// Created by Backend.BeginFrame, it tracks the frame token for
// staleness detection (EC1: drawing with a stale token panics).
type Ops struct {
	token uint64
}

// newOps creates a new Ops with a unique frame token.
func newOps() *Ops {
	return &Ops{token: atomic.AddUint64(&frameCounter, 1)}
}

// FrameToken returns this frame's unique token.
// Returns 0 for a nil Ops.
func (o *Ops) FrameToken() uint64 {
	if o == nil {
		return 0
	}
	return o.token
}

// NewFrameOps creates a new Ops with a unique frame token.
// Backend implementations use this in BeginFrame().
func NewFrameOps() *Ops {
	return newOps()
}
