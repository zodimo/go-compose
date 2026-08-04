package runtime

import (
	"image"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/internal/unitconvert"
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/internal/render"
	"github.com/zodimo/go-compose/pkg/api"

	"gioui.org/layout"
	"gioui.org/op"
)

var _ Runtime = (*runtime)(nil)

type runtime struct{}

// Run drives one frame. ctx is the engine layout context for this frame: the
// app shell passes the gio *layout.Context (as any, to keep the public
// interface engine-agnostic). The runtime wraps it in the framework-owned
// LayoutContext and returns an opaque draw command the shell applies.
func (r *runtime) Run(ctx any, composer api.Composer, ui api.Composable) render.DrawCommand {
	gctx, ok := ctx.(layout.Context)
	if !ok {
		panic("runtime.Run: ctx must be a gioui layout.Context")
	}

	backend := render.ActiveBackend()
	gtx := layoutnode.NewLayoutContextWithBackend(&gctx, backend)

	density := unitconvert.DensityFromLayoutContext(*gtx.ToGio())

	composer.StartFrame()
	defer composer.EndFrame()

	finalComposable := compose.CompositionLocalProvider(
		[]api.ProvidedValue{compose.LocalDensity.Provides(density)},
		ui,
	)

	node := finalComposable(composer).Build()
	nodeCoordinator := layoutnode.NewNodeCoordinator(node)

	// Branch on the active backend.
	// Gio path: preserve the existing behaviour exactly — outer discard
	// macro around Layout+PointerPhase, then coordinator.Draw into a
	// CallOp-based DrawCommand.
	if render.ActiveBackendName() == "gio" {
		gtx.ToGio().Constraints.Min = image.Point{X: 0, Y: 0}
		macro := op.Record(gtx.ToGio().Ops)
		nodeCoordinator.Layout(gtx)
		nodeCoordinator.PointerPhase(gtx)
		_ = macro.Stop()
		return render.NewDrawCommand(nodeCoordinator.Draw(gtx))
	}

	// Software / other backend path:
	// BeginFrame on the backend, single Layout pass (no outer discard
	// macro, no coordinator.Draw re-run), return a backend-based
	// DrawCommand. The harness reads Canvas() from the backend directly.
	w := float32(gtx.ToGio().Constraints.Max.X)
	h := float32(gtx.ToGio().Constraints.Max.Y)
	pxPerDp := gtx.ToGio().Metric.PxPerDp
	backend.BeginFrame(w, h, pxPerDp)
	defer backend.EndFrame()

	nodeCoordinator.Layout(gtx)
	nodeCoordinator.PointerPhase(gtx)

	return render.NewDrawCommandForBackend(backend)
}
