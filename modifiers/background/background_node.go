package background

import (
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	"github.com/zodimo/go-compose/internal/clipconvert"
	"github.com/zodimo/go-compose/internal/layoutnode"
	node "github.com/zodimo/go-compose/internal/node"

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
						nrgba := graphics.ColorToNRGBA(background.Color)
							return layoutnode.FromGioDimensions(layout.Background{}.Layout(*gtx.ToGio(),
								func(gtx layout.Context) layout.Dimensions {
									// shape
									// color
							defer background.Shape.CreateOutline(gtx.Constraints.Min, shape.Metric{PxPerDp: gtx.Metric.PxPerDp, PxPerSp: gtx.Metric.PxPerSp}).Push(clipconvert.NewOps(gtx.Ops)).Pop()

									paint.Fill(gtx.Ops, nrgba)

									return layout.Dimensions{Size: gtx.Constraints.Min}

								},
								func(gtx layout.Context) layout.Dimensions {
									return layoutnode.ToGioDimensions(widget.Layout(layoutnode.NewLayoutContext(&gtx)))
							},
						))
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
