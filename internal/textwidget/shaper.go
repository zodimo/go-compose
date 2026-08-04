package textwidget

import "gioui.org/text"

// Shaper wraps the gioui text shaper, hiding it from the public API surface.
type Shaper struct {
	S *text.Shaper
}

// NewShaper creates a Shaper wrapping the given gio text shaper.
func NewShaper(s *text.Shaper) *Shaper {
	if s == nil {
		return nil
	}
	return &Shaper{S: s}
}

// ToGio returns the underlying gio text shaper. For engine-bound callers only.
func (s *Shaper) ToGio() *text.Shaper {
	if s == nil {
		return nil
	}
	return s.S
}
