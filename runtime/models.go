package runtime

import (
	"image"

	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/pkg/api"

	"gioui.org/op"
)

var _ Runtime = (*runtime)(nil)

type runtime struct {
	onStartFrameFunc func()
	onEndFrameFunc   func()
	frameStarted     bool
}

func (r *runtime) Run(gtx LayoutContext, composer api.Composer, ui api.Composable) op.CallOp {

	gtx.Constraints.Min = image.Point{X: 0, Y: 0}

	composer.StartFrame()
	defer composer.EndFrame()

	node := ui(composer).Build()
	nodeCoordinator := layoutnode.NewNodeCoordinator(node)

	nodeCoordinator.Layout(gtx)
	nodeCoordinator.PointerPhase(gtx)
	return nodeCoordinator.Draw(gtx)
}
