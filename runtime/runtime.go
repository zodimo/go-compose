package runtime

import (
	"github.com/zodimo/go-compose/internal/render"
	"github.com/zodimo/go-compose/pkg/api"
)

// Runtime drives the composition and rendering of a composable tree for a
// single frame.
//
// The runtime interface is engine-agnostic: Run takes the engine frame
// context as an opaque value (the app shell passes the gio *layout.Context)
// and returns an opaque draw command. The seam conversion between the engine
// context and the framework-owned LayoutContext happens inside the runtime.
// Phase 4 replaces the opaque ctx with backend-driven frame construction.
type Runtime interface {
	Run(ctx any, composer api.Composer, ui api.Composable) render.DrawCommand
}
