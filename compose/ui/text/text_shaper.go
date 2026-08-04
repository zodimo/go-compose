package text

import "github.com/zodimo/go-compose/internal/textshaper"

// TextShaper wraps the text shaping engine. Use Shaper() to obtain the
// internal wrapper; call ToGio() on it at the engine seam.
type TextShaper struct {
	shaper *textshaper.Shaper
}

// NewTextShaper creates a TextShaper from a framework-owned textshaper.Shaper.
func NewTextShaper(shaper *textshaper.Shaper) *TextShaper {
	return &TextShaper{shaper: shaper}
}

// Shaper returns the framework-owned text shaper wrapper.
// Call ToGio() on it only at the engine seam.
func (t *TextShaper) Shaper() *textshaper.Shaper {
	return t.shaper
}
