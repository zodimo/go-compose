// Package textshaper provides a framework-owned wrapper around gioui.org/text.Shaper
// to keep gioui types out of the public API surface.
package textshaper

import "gioui.org/text"

// Shaper wraps gioui.org/text.Shaper for the public API surface.
// Use ToGio() to obtain the underlying shaper at the engine seam.
type Shaper struct {
	s *text.Shaper
}

// NewShaper wraps a gioui text.Shaper into a framework-owned Shaper.
func NewShaper(s *text.Shaper) *Shaper {
	return &Shaper{s: s}
}

// ToGio returns the underlying gioui text.Shaper.
// This is a seam-only accessor for engine-bound code.
func (s *Shaper) ToGio() *text.Shaper {
	return s.s
}
