package textwidget

import "gioui.org/text"

// Params wraps gioui text.Parameters, hiding it from the public API surface.
// It holds the text shaping parameters used by Layout calls.
type Params struct {
	P text.Parameters
}
