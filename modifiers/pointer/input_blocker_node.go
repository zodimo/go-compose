package pointer

import (
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/op/clip"
	"github.com/zodimo/go-compose/internal/layoutnode"
	node "github.com/zodimo/go-compose/internal/node"
)

var _ node.ChainNode = (*InputBlockerNode)(nil)

type InputBlockerNode struct {
	node.ChainNode
	tag *int
}

func NewInputBlockerNode(element InputBlockerElement) node.ChainNode {
	tag := new(int)
	return &InputBlockerNode{
		ChainNode: node.NewChainNode(
			node.NewNodeID(),
			node.NodeKindPointerInput,
			node.PointerInputPhase,
			func(n node.TreeNode) {
				no := n.(layoutnode.PointerInputModifierNode)
				no.AttachPointerInputModifier(func(widget layoutnode.LayoutWidget) layoutnode.LayoutWidget {
					return layoutnode.NewLayoutWidget(func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
						dims := widget.Layout(gtx)

						// Block all input events within the dimensions of the widget
						area := clip.Rect{Max: dims.Size}.Push(gtx.Ops)
						event.Op(gtx.Ops, tag)
						defer area.Pop()

						// Drain events to block them effectively
						for {
							_, ok := gtx.Event(pointer.Filter{
								Target: tag,
								Kinds:  pointer.Press | pointer.Release | pointer.Move | pointer.Drag | pointer.Scroll | pointer.Enter | pointer.Leave,
							})
							if !ok {
								break
							}
						}

						return dims
					})
				})
			},
		),
		tag: tag,
	}
}
