package clickable

import (
	"fmt"

	node "github.com/zodimo/go-compose/internal/Node"
	"github.com/zodimo/go-compose/internal/layoutnode"

	"gioui.org/layout"
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

					// fmt.Printf("Clickable Key : %s", key)
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

						return layout.Background{}.Layout(gtx,
							func(gtx layout.Context) layout.Dimensions {
								backgroundWidget := func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
									return layout.Dimensions{Size: gtx.Constraints.Min}
								}
								return material.Clickable(gtx, clickable, backgroundWidget)
							},
							widget.Layout,
						)
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
