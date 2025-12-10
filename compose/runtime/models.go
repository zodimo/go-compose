package runtime

import (
	"go-compose-dev/internal/layoutnode"
	"image"

	"gioui.org/op"
)

var _ Runtime = (*runtime)(nil)

type runtime struct {
}

func (r *runtime) Run(gtx LayoutContext, node LayoutNode) op.CallOp {

	gtx.Constraints.Min = image.Point{X: 0, Y: 0}
	nodeCoordinator := layoutnode.NewNodeCoordinator(node)

	nodeCoordinator.Layout(gtx)
	nodeCoordinator.PointerPhase(gtx)
	return nodeCoordinator.Draw(gtx)
}
