package clickable

import (
	"gioui.org/layout"
	"gioui.org/widget"
)

// Impl wraps widget.Clickable, keeping gioui confined to internal packages.
// Public packages access click behavior through the GioClickable wrapper in
// modifiers/clickable, which converts framework-owned types at the seam.
type Impl struct {
	W *widget.Clickable
}

// New creates a new Impl wrapping a fresh widget.Clickable.
func New() *Impl {
	return &Impl{W: &widget.Clickable{}}
}

// Clicked reports whether the clickable was clicked in this frame.
func (i *Impl) Clicked(gtx layout.Context) bool {
	return i.W.Clicked(gtx)
}

// Hovered reports whether the pointer is hovering over the clickable.
func (i *Impl) Hovered() bool {
	return i.W.Hovered()
}

// Pressed reports whether the pointer is pressed down on the clickable.
func (i *Impl) Pressed() bool {
	return i.W.Pressed()
}

// Update updates the clickable state for the current frame.
func (i *Impl) Update(gtx layout.Context) {
	i.W.Update(gtx)
}

// Layout implements layout.Widget for the clickable.
func (i *Impl) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	return i.W.Layout(gtx, w)
}

// Raw returns the underlying widget.Clickable for engine-bound code.
func (i *Impl) Raw() *widget.Clickable {
	return i.W
}
