package layoutnode

import (
	"go-compose-dev/internal/modifier"

	"gioui.org/op"
)

var _ NodeCoordinator = (*nodeCoordinator)(nil)

type nodeCoordinator struct {
	LayoutNode
	layoutCallChain  LayoutWidget
	pointerCallChain LayoutWidget
	elementStore     ElementStore
	wrappedChildren  []TreeNode
	expanded         bool
}

func (nc *nodeCoordinator) WrapChildren() {
	children := nc.LayoutNode.Children()
	wrappedChildren := []TreeNode{}

	for _, child := range children {
		wrappedChild := NewNodeCoordinator(child.(LayoutNode))
		wrappedChildren = append(wrappedChildren, wrappedChild)
	}
	nc.wrappedChildren = wrappedChildren
}

func (nc *nodeCoordinator) Expand() {

	modifierChain := nc.LayoutNode.UnwrapModifier().AsChain()
	*nc = *modifier.FoldIn(modifierChain, nc, func(nc *nodeCoordinator, mod Modifier) *nodeCoordinator {

		if inspectable, ok := mod.(InspectableModifier); ok {
			mod = inspectable.Unwrap()
		}

		modifierElement, ok := mod.(ModifierElement)
		if !ok {
			// probably EmptyModifier
			return nc
		}

		modifierNode := modifierElement.Create()
		modifierChainNode := modifierNode.(ChainNode)

		modifierChainNode.Attach(nc)

		return nc
	})

	for _, child := range nc.wrappedChildren {
		nodeCoordinatorChild := child.(NodeCoordinator)
		nodeCoordinatorChild.Expand()
	}
	nc.expanded = true

}

func (nc *nodeCoordinator) Children() []TreeNode {
	return nc.wrappedChildren
}

func (nc *nodeCoordinator) AttachLayoutModifier(attach func(widget LayoutWidget) LayoutWidget) {
	nc.layoutCallChain = nc.layoutCallChain.Map(func(in LayoutWidget) LayoutWidget {
		return NewLayoutWidget(func(gtx LayoutContext) LayoutDimensions {
			return attach(in).Layout(gtx)
		})
	})
}
func (nc *nodeCoordinator) AttachDrawModifier(attach func(widget LayoutWidget) LayoutWidget) {
	nc.layoutCallChain = nc.layoutCallChain.Map(func(in LayoutWidget) LayoutWidget {
		return NewLayoutWidget(func(gtx LayoutContext) LayoutDimensions {
			return attach(in).Layout(gtx)
		})
	})

}
func (nc *nodeCoordinator) AttachPointerInputModifier(attach func(widget LayoutWidget) LayoutWidget) {
	nc.layoutCallChain = nc.layoutCallChain.Map(func(in LayoutWidget) LayoutWidget {
		return NewLayoutWidget(func(gtx LayoutContext) LayoutDimensions {
			return attach(in).Layout(gtx)
		})
	})
}
func (nc *nodeCoordinator) AttachParentDataModifier(attach func(elements ElementStore) ElementStore) {
	nc.elementStore = attach(nc.elementStore)
}
func (nc *nodeCoordinator) PointerPhase(gtx LayoutContext) {
	defer op.Record(gtx.Ops).Stop()
	nc.pointerCallChain.Layout(gtx)
}

func (nc *nodeCoordinator) Elements() ElementStore {
	return nc.elementStore
}

func (nc *nodeCoordinator) Layout(gtx LayoutContext) LayoutDimensions {

	if !nc.expanded {
		nc.Expand()
	}

	return nc.layoutCallChain.Layout(gtx)
}

func (nc *nodeCoordinator) Draw(gtx LayoutContext) DrawOp {
	macro := op.Record(gtx.Ops)
	nc.layoutCallChain.Layout(gtx)
	return macro.Stop()
}

func (n *nodeCoordinator) GetWidget() GioLayoutWidget {
	maybeLayoutResult := n.GetLayoutResult()
	if maybeLayoutResult.IsSome() {
		return func(gtx LayoutContext) LayoutDimensions {
			layoutResult := maybeLayoutResult.UnwrapUnsafe()
			layoutResult.DrawOp.Add(gtx.Ops)
			return layoutResult.Dimensions
		}
	}
	return n.GetWidgetConstructor().Make(n)
}

type LayoutContextReceiver = func(gtx LayoutContext)

type LayoutResult struct {
	Dimensions LayoutDimensions
	DrawOp     op.CallOp
}
