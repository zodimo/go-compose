package widgetstate

import (
	"gioui.org/layout"
	"gioui.org/widget"
)

// List wraps gioui.org/widget.List to hide gioui types from public API.
type List struct {
	L widget.List
}

// NewVerticalList creates a new vertical list.
func NewVerticalList() *List {
	return &List{
		L: widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
	}
}
