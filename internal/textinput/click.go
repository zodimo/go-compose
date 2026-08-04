package textinput

import (
	"gioui.org/gesture"
	"gioui.org/io/key"
	"gioui.org/op"

	"github.com/zodimo/go-compose/internal/layoutnode"
)

// Click wraps gioui.org/gesture.Click to hide gioui types from public API.
type Click struct {
	C gesture.Click
}

// Add registers the click area with the given ops.
func (ck *Click) Add(ops *op.Ops) {
	ck.C.Add(ops)
}

// ProcessEvents drains click events and requests focus if a press is detected.
func (ck *Click) ProcessEvents(gtx layoutnode.LayoutContext, focusTag interface{}) {
	for {
		ev, ok := ck.C.Update(gtx.ToGio().Source)
		if !ok {
			break
		}
		switch ev.Kind {
		case gesture.KindPress:
			gtx.ToGio().Execute(key.FocusCmd{Tag: focusTag})
		}
	}
}

// Hovered reports whether the pointer is hovering over the click area.
func (ck *Click) Hovered() bool {
	return ck.C.Hovered()
}
