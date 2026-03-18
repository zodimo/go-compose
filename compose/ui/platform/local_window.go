package platform

import (
	"gioui.org/app"
	"github.com/zodimo/go-compose/compose"
)

// LocalWindow is a CompositionLocal that provides the app.Window to the composition.
// This allows components to access the primary window without creating dummy windows.
var LocalWindow = compose.StaticCompositionLocalOf[*app.Window](func() *app.Window {
	return nil
})
