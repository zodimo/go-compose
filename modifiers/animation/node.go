package animation

import (
	"github.com/zodimo/go-compose/internal/layoutnode"
	node "github.com/zodimo/go-compose/internal/node"
)

type AnimatedWidthNode struct {
	node.ChainNode
	element AnimatedWidthElement
}

func NewAnimatedWidthNode(element AnimatedWidthElement) *AnimatedWidthNode {
	n := &AnimatedWidthNode{
		element: element,
	}
	n.ChainNode = node.NewChainNode(
		node.NewNodeID(),
		node.NodeKindLayout,
		node.LayoutPhase,
		func(t node.TreeNode) {
			no := t.(layoutnode.LayoutModifierNode)
			no.AttachLayoutModifier(func(widget layoutnode.LayoutWidget) layoutnode.LayoutWidget {
				return layoutnode.NewLayoutWidget(func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
					// Logic
					progress := n.element.Anim.Revealed(gtx)

					width := int(float32(n.element.MaxWidth) * progress)

					// Apply width constraint
					// We force the width to be exactly 'width'
					c := gtx.Constraints
					c.Min.X = width
					c.Max.X = width

					// Override Gtx
					childGtx := gtx
					childGtx.Constraints = c

					dims := widget.Layout(childGtx)

					return dims
				})
			})
		},
	)
	return n
}
