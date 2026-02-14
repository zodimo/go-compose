package alpha

import (
	"gioui.org/op/paint"
	"github.com/zodimo/go-compose/internal/layoutnode"
	node "github.com/zodimo/go-compose/internal/node"
)

type AlphaState struct {
	Alpha float32
}

type AlphaNode struct {
	ChainNode
	state *AlphaState
}

var _ ChainNode = (*AlphaNode)(nil)

func NewAlphaNode(element AlphaElement) ChainNode {
	state := &AlphaState{Alpha: element.Alpha}

	return &AlphaNode{
		ChainNode: node.NewChainNode(
			node.NewNodeID(),
			node.NodeKindDraw,
			node.DrawPhase,
			//OnAttach
			func(n TreeNode) {
				no := n.(layoutnode.DrawModifierNode)

				no.AttachDrawModifier(func(widget LayoutWidget) layoutnode.LayoutWidget {
					return layoutnode.NewLayoutWidget(func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
						macro := paint.PushOpacity(gtx.Ops, state.Alpha)
						defer macro.Pop()
						return widget.Layout(gtx)
					})
				})
			},
		),
		state: state,
	}
}
