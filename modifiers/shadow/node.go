package shadow

import (
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	"github.com/zodimo/go-compose/internal/clipconvert"
	"github.com/zodimo/go-compose/internal/layoutnode"
	node "github.com/zodimo/go-compose/internal/node"
	"github.com/zodimo/go-compose/internal/render"
	"github.com/zodimo/go-compose/internal/unitconvert"

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
					elevation := n.shadowData.Elevation
					if elevation <= 0 {
						return widget.Layout(gtx)
					}

					if b := gtx.DrawBackend(); b != nil {
						// Software path:
						// 1. Measure child
						rec := b.Record()
						dims := widget.Layout(gtx)
						b.Apply(rec) // replay child into canvas

						// 2. Draw shadow layers
						shadowSize := float32(gtx.ToGio().Metric.Dp(unitconvert.DpToGioUnitUnsafe(elevation)))
						col := shadowColorToRenderColor(n.shadowData.AmbientColor.TakeOrElse(graphics.ColorBlack))

						shadowLayersCount := float32(8)
						for layerIndex := shadowLayersCount; layerIndex > 0; layerIndex-- {
							sWidth := 0.75 + shadowSize*layerIndex*0.4/shadowLayersCount
							// Record a FillShape for each shadow layer
							// The exact scale/offset math mirrors the gio path
							// but we just record the shape bounds.
							shadowShape := render.Shape{
								Kind:   render.ShapeRectangle,
								Bounds: render.Rect{
									Min: render.Point{X: -sWidth, Y: -sWidth},
									Max: render.Point{X: float32(dims.Size.X) + sWidth, Y: float32(dims.Size.Y) + sWidth},
								},
							}
							b.FillShape(shadowShape, col)
						}

						// 3. Draw child on top
						childRec := b.Record()
						finalDims := widget.Layout(gtx)
						b.Apply(childRec)
						return finalDims
					}

					// Gio path: unchanged.
					macro := op.Record(gtx.ToGio().Ops)
					dims := widget.Layout(gtx)
					call := macro.Stop()

					if elevation <= 0 {
						call.Add(gtx.ToGio().Ops)
						return dims
					}

					// Draw Shadow
					shadowSize := float32(gtx.ToGio().Metric.Dp(unitconvert.DpToGioUnitUnsafe(elevation)))

					col := graphics.ColorToNRGBA(n.shadowData.AmbientColor.TakeOrElse(graphics.ColorBlack))
					if col.A == 255 {
						col.A = 30
					}

					shadowShapeBounds := f32.Point{
						X: float32(dims.Size.X),
						Y: float32(dims.Size.Y),
					}

					outline := n.shadowData.Shape.CreateOutline(dims.Size, shape.Metric{PxPerDp: gtx.ToGio().Metric.PxPerDp, PxPerSp: gtx.ToGio().Metric.PxPerSp})

					baseMacro := op.Record(gtx.ToGio().Ops)
					paint.FillShape(gtx.ToGio().Ops, col, outline.ClipOp(clipconvert.NewOps(gtx.ToGio().Ops)).ToClipOp())
					baseCall := baseMacro.Stop()

					var stack op.TransformStack
					shadowLayersCount := float32(8)

					for layerIndex := shadowLayersCount; layerIndex > 0; layerIndex-- {
						sWidth := 0.75 + shadowSize*layerIndex*0.4/shadowLayersCount
						finalSize := shadowShapeBounds.Add(f32.Point{X: sWidth, Y: sWidth})

						if shadowShapeBounds.X == 0 || shadowShapeBounds.Y == 0 {
							continue
						}

						scaleFactor := f32.Pt(finalSize.X/shadowShapeBounds.X, finalSize.Y/shadowShapeBounds.Y)
						xOffset := (shadowShapeBounds.X - finalSize.X) / 2
						yOffset := sWidth - 0.75

						scaleOrigin := f32.Point{X: scaleFactor.X / 2, Y: 0}
						sOffset := f32.Pt(xOffset, yOffset)

						stack = op.Affine(f32.AffineId().Offset(sOffset).Scale(scaleOrigin, scaleFactor)).Push(gtx.ToGio().Ops)
						baseCall.Add(gtx.ToGio().Ops)
						stack.Pop()
					}

					// Draw content on top
					call.Add(gtx.ToGio().Ops)

					return dims
				})
			})
		},
	)
	return n
}

// shadowColorToRenderColor converts a framework Color to a render.Color.
func shadowColorToRenderColor(c graphics.Color) render.Color {
	nrgba := graphics.ColorToNRGBA(c)
	// Apply the alpha adjustment matching the gio path
	if nrgba.A == 255 {
		nrgba.A = 30
	}
	return render.Color{
		R: float32(nrgba.R) / 255,
		G: float32(nrgba.G) / 255,
		B: float32(nrgba.B) / 255,
		A: float32(nrgba.A) / 255,
	}
}
