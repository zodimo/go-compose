package clickable

import (
	"fmt"
	node "go-compose-dev/internal/Node"
	"go-compose-dev/internal/layoutnode"

	"gioui.org/widget/material"
)

var _ ChainNode = (*ClickableNode)(nil)

// NodeKind should also implement the interface of the LayoutNode for that phase

func NewClickableNode(element ClickableElement) ChainNode {
	return ClickableNode{
		ChainNode: node.NewChainNode(
			node.NewNodeID(),
			node.NodeKindPointerInput,
			node.PointerInputPhase, // bit mask and && node.DrawPhase,
			//OnAttach
			func(n TreeNode) {
				// how should the tree now be updated when attached
				// tree nde is the layout tree

				if element.clickableData.Clickable == nil {
					// we need persistent storage
					lno := n.(layoutnode.LayoutNode)
					key := lno.GenerateID()

					// path := lno.GetPath()

					clickablePath := fmt.Sprintf("%d/clickable", key)
					clickableValue := lno.State(clickablePath, func() any { return &GioClickable{} })
					clickable := clickableValue.Get().(*GioClickable)
					element.clickableData.Clickable = clickable
				}

				no := n.(layoutnode.PointerInputModifierNode)
				// we can now work with the layoutNode
				no.AttachPointerInputModifier(func(widget LayoutWidget) layoutnode.LayoutWidget {

					return layoutnode.NewLayoutWidget(func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
						clickable := element.clickableData.Clickable
						onClick := element.clickableData.OnClick
						if clickable.Clicked(gtx) {
							onClick()
						}
						return layoutnode.LayoutDimensions{}
					})
				})

				dno := n.(layoutnode.DrawModifierNode)

				dno.AttachDrawModifier(func(widget LayoutWidget) layoutnode.LayoutWidget {
					return layoutnode.NewLayoutWidget(func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
						clickable := element.clickableData.Clickable
						return material.Clickable(gtx, clickable, widget.Layout)
					})
				})

			},
		),
		clickableData: element.clickableData,
	}
}

type ClickableNode struct {
	ChainNode
	clickableData ClickableData
}
