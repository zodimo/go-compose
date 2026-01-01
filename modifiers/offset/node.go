package offset

import (
	"image"

	"github.com/zodimo/go-compose/compose/ui/unit"
	node "github.com/zodimo/go-compose/internal/Node"
	"github.com/zodimo/go-compose/internal/layoutnode"

	"gioui.org/op"
)

type OffsetNode struct {
	node.ChainNode
	data OffsetData
}

func NewOffsetNode(data OffsetData) *OffsetNode {
	n := &OffsetNode{
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
					// Convert dp to pixels
					offsetX := gtx.Dp(unit.DpToGioUnit(n.data.X))
					offsetY := gtx.Dp(unit.DpToGioUnit(n.data.Y))

					// Apply translation offset using op.Offset
					stack := op.Offset(image.Point{X: offsetX, Y: offsetY}).Push(gtx.Ops)

					// Layout the content
					dims := widget.Layout(gtx)

					stack.Pop()

					// Return original dimensions (offset doesn't change the element's size)
					return dims
				})
			})
		},
	)
	return n
}
