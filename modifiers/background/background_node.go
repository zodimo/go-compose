package background

import (
	node "github.com/zodimo/go-compose/internal/Node"
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/theme"

	"gioui.org/layout"
	"gioui.org/op/paint"
)

var _ ChainNode = (*BackgroundNode)(nil)

// var _ DrawModifierNode = (*BackgroundNode)(nil)

// NodeKind should also implement the interface of the LayoutNode for that phase

func NewBackGroundNode(background BackgroundData) ChainNode {
	return BackgroundNode{
		ChainNode: node.NewChainNode(
			node.NewNodeID(),
			node.NodeKindDraw,
			node.DrawPhase,
			//OnAttach
			func(n TreeNode) {
				// how should the tree now be updated when attached
				// tree nde is the layout tree

				no := n.(layoutnode.DrawModifierNode)
				// we can now work with the layoutNode
				no.AttachDrawModifier(func(widget LayoutWidget) layoutnode.LayoutWidget {

					return layoutnode.NewLayoutWidget(func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
						// Resolve ColorDescriptor at layout time
						themeManager := theme.GetThemeManager()
						themeColor := themeManager.ResolveColorDescriptor(background.Color)
						nrgba := themeColor.AsNRGBA()

						return layout.Background{}.Layout(gtx,
							func(gtx layout.Context) layout.Dimensions {
								// shape
								// color
								defer background.Shape.CreateOutline(gtx.Constraints.Min, gtx.Metric).Push(gtx.Ops).Pop()

								paint.Fill(gtx.Ops, nrgba)

								return layout.Dimensions{Size: gtx.Constraints.Min}

							},
							func(gtx layout.Context) layout.Dimensions {
								return widget.Layout(gtx)
							},
						)
					})
				})

			},
		),
		background: background,
	}
}

type BackgroundNode struct {
	ChainNode
	background BackgroundData
}
