package platform

import (
	"github.com/zodimo/go-compose/compose"
)

// Window is the framework-owned window type. It wraps the platform window
// in an opaque struct so that gioui.org/app.Window does not appear in
// the public API surface.
type Window struct {
	// win holds the platform window pointer. Unexported to keep
	// gioui.org types out of the public API.
	win any
}

// NewWindow creates a framework Window wrapping a platform window.
// The platformWindow parameter is the gioui *app.Window (passed as any
// to avoid gioui imports in this package's signature).
func NewWindow(platformWindow any) *Window {
	return &Window{win: platformWindow}
}

// PlatformWindow returns the underlying platform window.
// This is seam-only: for engine-bound code.
func (w *Window) PlatformWindow() any {
	if w == nil {
		return nil
	}
	return w.win
}

// LocalWindow is a CompositionLocal that provides the Window to the composition.
// This allows components to access the primary window without creating dummy windows.
var LocalWindow = compose.StaticCompositionLocalOf[*Window](func() *Window {
	return nil
})
