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
					progress := n.element.Anim.Revealed(*gtx.ToGio())

					width := int(float32(n.element.MaxWidth) * progress)

					// Apply width constraint
					// We force the width to be exactly 'width'
					c := gtx.ToGio().Constraints
					c.Min.X = width
					c.Max.X = width

					// Override Gtx (copy semantics: the wrapper holds a pointer)
					g := *gtx.ToGio()
					g.Constraints = c

					dims := widget.Layout(layoutnode.NewLayoutContextWithBackend(&g, gtx.DrawBackend()))

					return dims
				})
			})
		},
	)
	return n
}
