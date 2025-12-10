package size

import (
	node "go-compose-dev/internal/Node"
	"go-compose-dev/internal/layoutnode"
	"image"

	"gioui.org/layout"
)

var _ ChainNode = (*SizeNode)(nil)

// NodeKind should also implement the interface of the LayoutNode for that phase

type SizeNode struct {
	ChainNode
	size SizeData
}

func NewSizeNode(size SizeData) ChainNode {
	return SizeNode{
		ChainNode: node.NewChainNode(
			node.NewNodeID(),
			node.NodeKindLayout,
			node.LayoutPhase,
			//OnAttach
			func(n TreeNode) {
				// how should the tree now be updated when attached
				// tree nde is the layout tree

				no := n.(layoutnode.LayoutModifierNode)
				// we can now work with the layoutNode
				no.AttachLayoutModifier(func(gtx layoutnode.LayoutContext, widget layoutnode.LayoutWidget) layoutnode.LayoutWidget {
					return layoutnode.NewLayoutWidget(
						func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
							constraints := gtx.Constraints
							if size.Required {
								constraints.Min.X = size.Width
								constraints.Min.Y = size.Height
								constraints.Max.X = size.Width
								constraints.Max.Y = size.Height

							} else {
								//clamp
								constraints.Min.X = Clamp(size.Width, gtx.Constraints.Min.X, gtx.Constraints.Max.X)
								constraints.Min.Y = Clamp(size.Height, gtx.Constraints.Min.Y, gtx.Constraints.Max.Y)
								constraints.Max.X = Clamp(size.Width, gtx.Constraints.Min.X, gtx.Constraints.Max.X)
								constraints.Max.Y = Clamp(size.Height, gtx.Constraints.Min.Y, gtx.Constraints.Max.Y)
							}
							gtx.Constraints = constraints
							widget.Layout(gtx)
							return layout.Dimensions{
								Size: image.Point{
									X: size.Width,
									Y: size.Height,
								},
							}
						},
					)
				})

			},
		),
		size: size,
	}
}
