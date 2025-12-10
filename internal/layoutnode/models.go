package layoutnode

import (
	"go-compose-dev/internal/modifier"

	"gioui.org/op"
)

var _ NodeCoordinator = (*nodeCoordinator)(nil)

type nodeCoordinator struct {
	LayoutNode
	layoutCallChain  LayoutWidget
	drawCallChain    LayoutWidget
	pointerCallChain LayoutWidget
	elementStore     ElementStore
}

func (nc *nodeCoordinator) Expand() {
	modifierChain := nc.LayoutNode.UnwrapModifier().AsChain()
	*nc = *modifier.Fold(modifierChain, nc, func(nc *nodeCoordinator, mod Modifier) *nodeCoordinator {
		if inspectable, ok := mod.(InspectableModifier); ok {
			mod = inspectable.Unwrap()
		}
		modifierElement := mod.(ModifierElement)

		modifierNode := modifierElement.Create()
		modifierChainNode := modifierNode.(ChainNode)
		modifierChainNode.Attach(nc)

		return nc
	})

}

func (nc *nodeCoordinator) AttachLayoutModifier(attach func(gtx LayoutContext, widget LayoutWidget) LayoutWidget) {

	nc.layoutCallChain = nc.layoutCallChain.Map(func(in LayoutWidget) LayoutWidget {
		return NewLayoutWidget(func(gtx LayoutContext) LayoutDimensions {
			return attach(gtx, in).Layout(gtx)
		})
	})
}
func (nc *nodeCoordinator) AttachDrawModifier(attach func(gtx LayoutContext, widget LayoutWidget) LayoutWidget) {
	nc.drawCallChain = nc.drawCallChain.Map(func(in LayoutWidget) LayoutWidget {
		return NewLayoutWidget(func(gtx LayoutContext) LayoutDimensions {
			return attach(gtx, in).Layout(gtx)
		})
	})
}
func (nc *nodeCoordinator) AttachPointerModifier(attach func(gtx LayoutContext, widget LayoutWidget) LayoutWidget) {
	nc.pointerCallChain = nc.pointerCallChain.Map(func(in LayoutWidget) LayoutWidget {
		return NewLayoutWidget(func(gtx LayoutContext) LayoutDimensions {
			return attach(gtx, in).Layout(gtx)
		})
	})
}

func (nc *nodeCoordinator) LayoutPhase(gtx LayoutContext) {
	defer op.Record(gtx.Ops).Stop()
	nc.LayoutSelf(gtx)
}

func (nc *nodeCoordinator) PointerPhase(gtx LayoutContext) {
	defer op.Record(gtx.Ops).Stop()
	nc.pointerCallChain.Layout(gtx)
}

func (nc *nodeCoordinator) DrawPhase(gtx LayoutContext) {
	nc.drawCallChain.Layout(gtx)
}

func (nc *nodeCoordinator) Elements() ElementStore {
	return nc.elementStore
}

func (nc *nodeCoordinator) LayoutSelf(gtx LayoutContext) LayoutDimensions {
	return nc.layoutCallChain.Layout(gtx)

}

type LayoutContextReceiver = func(gtx LayoutContext)
