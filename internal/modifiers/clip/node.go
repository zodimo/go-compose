package clip

import (
	"fmt"
	node "go-compose-dev/internal/Node"
	"go-compose-dev/internal/layoutnode"

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
						stack := ClipShape(element.clipData.Shape, gtx, dimensions)
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

func ClipShape(shape Shape, gtx layout.Context, dimensions layoutnode.LayoutDimensions) clip.Stack {
	switch shape {
	case ShapeCircle:
		return clip.Ellipse{Max: dimensions.Size}.Push(gtx.Ops)

	case ShapeRectangle:
		return clip.Rect{Max: dimensions.Size}.Push(gtx.Ops)
	default:
		panic(fmt.Sprintf("clip unsupported shape: %s", shape))
	}
}
