package screenshot

import (
	"image"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/store"
	"github.com/zodimo/go-compose/theme"
)

// RenderResult is the observable output of running a Composable through the
// real composition + layout path, captured without a GPU.
type RenderResult struct {
	// Size is the measured size of the composed UI. It is stable across runs
	// for a given Composable and window size and is therefore the primary
	// CPU-only fingerprint of the layout.
	Size image.Point

	// Baseline is the measured baseline reported alongside Size.
	Baseline int

	// DrawOp is the recorded draw operation list produced by the layout pass.
	// On a machine with a working EGL/GPU stack it can be rasterised via
	// TakeScreenshot.
	DrawOp op.CallOp

	// Ops is the op.Ops the DrawOp was replayed into. It is retained so callers
	// can rasterise the frame without re-running the Composable.
	Ops *op.Ops
}

// Render runs ui through the same composition + layout path the desktop demos
// use, but stops short of GPU rasterisation. It performs no GPU access, so it
// runs in environments where headless.NewWindow would hang.
//
// It mirrors runtime.Run (compose the tree, lay it out through a coordinator,
// capture the draw ops) with one addition: it keeps the returned
// layout.Dimensions, which runtime.Run discards, so callers get the measured
// size as a CPU-only signal.
func Render(ui api.Composable, width, height int) RenderResult {
	var ops op.Ops

	gtx := layout.Context{
		Ops:         &ops,
		Constraints: layout.Exact(image.Pt(width, height)),
		Metric:      unit.Metric{PxPerDp: 1, PxPerSp: 1},
		Now:         time.Now(),
	}

	themeManager := theme.GetThemeManager()
	gtx = themeManager.Material3ThemeInit(gtx)

	composer := compose.NewComposer(api.ComposerWithStore(store.NewPersistentState()))

	// Mirror runtime.Run: compose to a node, wrap it in a coordinator, then lay
	// it out. The layout pass records draw operations into the coordinator's
	// call chain; Draw materialises them as an op.CallOp.
	gtx.Constraints.Min = image.Point{}
	node := ui(composer).Build()
	coordinator := layoutnode.NewNodeCoordinator(node)

	dims := coordinator.Layout(gtx)
	drawOp := coordinator.Draw(gtx)
	drawOp.Add(&ops)

	return RenderResult{
		Size:     dims.Size,
		Baseline: dims.Baseline,
		DrawOp:   drawOp,
		Ops:      &ops,
	}
}
