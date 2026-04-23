package runtime

import (
	"image"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/pkg/api"

	"gioui.org/op"
)

var _ Runtime = (*runtime)(nil)

type runtime struct {
}

func (r *runtime) Run(gtx LayoutContext, composer api.Composer, ui api.Composable) op.CallOp {

	gtx.Constraints.Min = image.Point{X: 0, Y: 0}

	density := unit.DensityFromLayoutContext(gtx)

	composer.StartFrame()
	defer composer.EndFrame()

	finalComposable := compose.CompositionLocalProvider(
		[]api.ProvidedValue{compose.LocalDensity.Provides(density)},
		ui,
	)

	node := finalComposable(composer).Build()
	nodeCoordinator := layoutnode.NewNodeCoordinator(node)

	// Stop double rendering
	macro := op.Record(gtx.Ops)
	nodeCoordinator.Layout(gtx)
	nodeCoordinator.PointerPhase(gtx)
	_ = macro.Stop()

	return nodeCoordinator.Draw(gtx)
}
