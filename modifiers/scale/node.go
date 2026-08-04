package scale

import (
	"github.com/zodimo/go-compose/internal/layoutnode"
	node "github.com/zodimo/go-compose/internal/node"

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
					if b := gtx.DrawBackend(); b != nil {
						// Software path: push scale transform, layout child, restore.
						saveLevel := b.Save()
						b.Scale(n.data.ScaleX, n.data.ScaleY)
						dims := widget.Layout(gtx)
						b.Restore(saveLevel)
						return dims
					}
					// Gio path: unchanged.
					macro := op.Record(gtx.ToGio().Ops)
					dims := widget.Layout(gtx)
					call := macro.Stop()

					cx := float32(dims.Size.X) / 2
					cy := float32(dims.Size.Y) / 2
					center := f32.Pt(cx, cy)

					// Use captured n.data to support updates
					t := f32.AffineId().Scale(center, f32.Pt(n.data.ScaleX, n.data.ScaleY))

					stack := op.Affine(t).Push(gtx.ToGio().Ops)
					call.Add(gtx.ToGio().Ops)
					stack.Pop()

					return dims
				})
			})
		},
	)
	return n
}
