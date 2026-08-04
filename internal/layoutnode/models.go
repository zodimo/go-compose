package layoutnode

import (
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/internal/modifier"

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
	*nc = *modifier.FoldIn(modifierChain, nc, func(nc *nodeCoordinator, mod ui.Modifier) *nodeCoordinator {

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
	if b := gtx.DrawBackend(); b != nil {
		// Software path: discard pointer ops. Pointer modifiers (pointer,
		// clickable) emit only gtx.ToGio().Ops event registrations, not
		// Backend draw calls, so the recording is effectively empty.
		callOp := b.Record()
		nc.pointerCallChain.Layout(gtx)
		b.Apply(callOp)
		return
	}
	// Gio path: unchanged.
	defer op.Record(gtx.ToGio().Ops).Stop()
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
	if b := gtx.DrawBackend(); b != nil {
		// Software path: all draw ops were captured during Layout.
		// Draw is a no-op; the caller does not use the return value.
		return DrawOp{}
	}
	// Gio path: unchanged.
	macro := op.Record(gtx.ToGio().Ops)
	nc.layoutCallChain.Layout(gtx)
	return macro.Stop()
}

func (n *nodeCoordinator) GetWidget() GioLayoutWidget {
	maybeLayoutResult := n.GetLayoutResult()
	if maybeLayoutResult.IsSome() {
		return func(gtx LayoutContext) LayoutDimensions {
			layoutResult := maybeLayoutResult.UnwrapUnsafe()
			if b := gtx.DrawBackend(); b != nil {
				// Software path: DrawOp is zero; all draw ops were
				// captured during the Layout that produced this result.
				return layoutResult.Dimensions
			}
			// Gio path: unchanged.
			layoutResult.DrawOp.Add(gtx.ToGio().Ops)
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
