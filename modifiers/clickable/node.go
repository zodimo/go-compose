package clickable

import (
	"fmt"

	"github.com/zodimo/go-compose/internal/layoutnode"
	node "github.com/zodimo/go-compose/internal/node"
	"github.com/zodimo/go-compose/state"

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

					clickablePath := fmt.Sprintf("%d/clickable", key)
					clickableValue := state.MustRemember(lno, clickablePath, func() *GioClickable { return &GioClickable{} })
					clickable := clickableValue.Get()
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
