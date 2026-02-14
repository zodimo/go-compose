package clip

import (
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	"github.com/zodimo/go-compose/internal/layoutnode"
	node "github.com/zodimo/go-compose/internal/node"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
)

type ClipNode struct {
	ChainNode
	clipData ClipData
}

var _ ChainNode = (*ClipNode)(nil)

func NewClipNode(element ClipElement) ChainNode {
	return ClipNode{
		ChainNode: node.NewChainNode(
			node.NewNodeID(),
			node.NodeKindDraw,
			node.DrawPhase,
			//OnAttach
			func(n TreeNode) {

				no := n.(layoutnode.DrawModifierNode)
				// we can now work with the layoutNode
				no.AttachDrawModifier(func(widget LayoutWidget) layoutnode.LayoutWidget {
					return layoutnode.NewLayoutWidget(func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
						//clip to the shape
						macro := op.Record(gtx.Ops)
						dimensions := widget.Layout(gtx)
						callOp := macro.Stop()
						// Clip Shape here
						clipDimensions := dimensions
						if element.clipData.ClipToBounds {
							clipDimensions = layoutnode.LayoutDimensions{
								Size: gtx.Constraints.Max,
							}
						}

						stack := ClipShape(element.clipData.Shape, gtx, clipDimensions)

						callOp.Add(gtx.Ops)
						stack.Pop()

						return dimensions
					})
				})

			},
		),
		clipData: element.ClipData(),
	}
}

func ClipShape(shape shape.Shape, gtx layout.Context, dimensions layoutnode.LayoutDimensions) clip.Stack {
	return shape.CreateOutline(dimensions.Size, gtx.Metric).Push(gtx.Ops)
}
