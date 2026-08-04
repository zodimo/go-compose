package textwidget

import (
	"time"

	"gioui.org/gesture"
	"gioui.org/io/input"
	"gioui.org/io/pointer"
	"gioui.org/op"
	"gioui.org/unit"
)

// Clicker wraps gioui gesture.Click, hiding it from the public API surface.
type Clicker struct {
	inner gesture.Click
}

// Add adds the click recognizer to the operation list.
func (c *Clicker) Add(ops *op.Ops) {
	c.inner.Add(ops)
}

// Update processes pointer events and returns the click event, if any.
func (c *Clicker) Update(q input.Source) (gesture.ClickEvent, bool) {
	return c.inner.Update(q)
}

// Dragger wraps gioui gesture.Drag, hiding it from the public API surface.
type Dragger struct {
	inner gesture.Drag
}

// Add adds the drag recognizer to the operation list.
func (d *Dragger) Add(ops *op.Ops) {
	d.inner.Add(ops)
}

// Update processes pointer events and returns the drag event, if any.
func (d *Dragger) Update(
	cfg unit.Metric,
	q input.Source,
	axis gesture.Axis,
) (pointer.Event, bool) {
	return d.inner.Update(cfg, q, axis)
}

// Scroller wraps gioui gesture.Scroll, hiding it from the public API surface.
type Scroller struct {
	inner gesture.Scroll
}

// Add adds the scroll recognizer to the operation list.
func (s *Scroller) Add(ops *op.Ops) {
	s.inner.Add(ops)
}

// Update processes pointer events and returns the scroll distance along axis.
func (s *Scroller) Update(
	cfg unit.Metric,
	q input.Source,
	t time.Time,
	axis gesture.Axis,
	scrollx, scrolly pointer.ScrollRange,
) int {
	return s.inner.Update(cfg, q, t, axis, scrollx, scrolly)
}

// State returns the current scroll state.
func (s *Scroller) State() gesture.ScrollState {
	return s.inner.State()
}

// Stop stops any active scroll gesture.
func (s *Scroller) Stop() {
	s.inner.Stop()
}
