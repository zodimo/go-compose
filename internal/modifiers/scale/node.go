package scale

import (
	node "go-compose-dev/internal/Node"
	"go-compose-dev/internal/layoutnode"

	"gioui.org/f32"
	"gioui.org/op"
)

type ScaleNode struct {
	node.ChainNode
	data ScaleData
}

func NewScaleNode(data ScaleData) *ScaleNode {
	n := &ScaleNode{
		data: data,
	}
	n.ChainNode = node.NewChainNode(
		node.NewNodeID(),
		node.NodeKindLayout,
		node.LayoutPhase,
		func(t node.TreeNode) {
			no := t.(layoutnode.LayoutModifierNode)
			no.AttachLayoutModifier(func(widget layoutnode.LayoutWidget) layoutnode.LayoutWidget {
				return layoutnode.NewLayoutWidget(func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
					macro := op.Record(gtx.Ops)
					dims := widget.Layout(gtx)
					call := macro.Stop()

					cx := float32(dims.Size.X) / 2
					cy := float32(dims.Size.Y) / 2
					center := f32.Pt(cx, cy)

					// Use captured n.data to support updates
					t := f32.AffineId().Scale(center, f32.Pt(n.data.ScaleX, n.data.ScaleY))

					stack := op.Affine(t).Push(gtx.Ops)
					call.Add(gtx.Ops)
					stack.Pop()

					return dims
				})
			})
		},
	)
	return n
}
