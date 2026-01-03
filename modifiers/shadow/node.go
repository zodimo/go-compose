package shadow

import (
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/unit"
	node "github.com/zodimo/go-compose/internal/Node"
	"github.com/zodimo/go-compose/internal/layoutnode"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/paint"
)

type ShadowNode struct {
	ChainNode
	shadowData ShadowData
}

func NewShadowNode(element ShadowElement) *ShadowNode {
	n := &ShadowNode{
		shadowData: element.shadowData,
	}
	n.ChainNode = node.NewChainNode(
		node.NewNodeID(),
		node.NodeKindDraw,
		node.DrawPhase,
		func(t TreeNode) {
			no := t.(DrawModifierNode)
			no.AttachDrawModifier(func(widget LayoutWidget) LayoutWidget {
				return layoutnode.NewLayoutWidget(func(gtx LayoutContext) LayoutDimensions {
					// Layout children first to know the size?
					// Wait, we need the size to draw the shadow.
					// But usually modifiers wrap. If we are outer, we call inner widget.Layout.
					// But we want to draw shadow BEHIND content.

					// Strategy:
					// 1. Record content layout.
					// 2. Draw shadow based on content size.
					// 3. Draw content.

					macro := op.Record(gtx.Ops)
					dims := widget.Layout(gtx)
					call := macro.Stop()

					elevation := n.shadowData.Elevation
					if elevation <= 0 {
						call.Add(gtx.Ops)
						return dims
					}

					// Draw Shadow
					// Adapted from gio-mw wdk.Elevation.Layout

					shadowSize := float32(gtx.Metric.Dp(unit.DpToGioUnit(elevation)))

					//@TODO get shadow from theme, for now default to black
					col := graphics.ColorToNRGBA(n.shadowData.AmbientColor.TakeOrElse(graphics.ColorBlack))
					// Apply some opacity if it's fully opaque?
					// gio-mw uses 0.12*255 approx 30 alpha.
					if col.A == 255 {
						col.A = 30
					}

					shadowShapeBounds := f32.Point{
						X: float32(dims.Size.X),
						Y: float32(dims.Size.Y),
					}

					// Create Outline for the shape
					// We need the outline path.
					outline := n.shadowData.Shape.CreateOutline(dims.Size, gtx.Metric)

					// Draw base layer
					baseMacro := op.Record(gtx.Ops)
					paint.FillShape(gtx.Ops, col, outline.Op(gtx.Ops))
					baseCall := baseMacro.Stop()

					var stack op.TransformStack
					shadowLayersCount := float32(8)

					// We need to offset/scale the *base layer drawing*.
					// But outline.Op(gtx.Ops) is just the path.
					// We can't reuse the path Op easily with different transforms unless we rebuild it or use transform on the fill?
					// Actually gio-mw records a macro of the FillShape and then replays it with transforms.

					// Replicate gio-mw loop
					for layerIndex := shadowLayersCount; layerIndex > 0; layerIndex-- {
						sWidth := 0.75 + shadowSize*layerIndex*0.4/shadowLayersCount
						finalSize := shadowShapeBounds.Add(f32.Point{X: sWidth, Y: sWidth})

						// Avoid division by zero
						if shadowShapeBounds.X == 0 || shadowShapeBounds.Y == 0 {
							continue
						}

						scaleFactor := f32.Pt(finalSize.X/shadowShapeBounds.X, finalSize.Y/shadowShapeBounds.Y)
						xOffset := (shadowShapeBounds.X - finalSize.X) / 2
						yOffset := sWidth - 0.75

						scaleOrigin := f32.Point{X: scaleFactor.X / 2, Y: 0}
						sOffset := f32.Pt(xOffset, yOffset)

						stack = op.Affine(f32.AffineId().Offset(sOffset).Scale(scaleOrigin, scaleFactor)).Push(gtx.Ops)
						baseCall.Add(gtx.Ops)
						stack.Pop()
					}

					// Draw content on top
					call.Add(gtx.Ops)

					return dims
				})
			})
		},
	)
	return n
}
