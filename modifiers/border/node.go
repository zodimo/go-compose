package border

import (
	"image"

	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	"github.com/zodimo/go-compose/internal/clipconvert"
	"github.com/zodimo/go-compose/internal/layoutnode"
	node "github.com/zodimo/go-compose/internal/node"
	"github.com/zodimo/go-compose/internal/render"
	"github.com/zodimo/go-compose/internal/unitconvert"

	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type BorderNode struct {
	ChainNode
	borderData BorderData
}

func NewBorderNode(element BorderElement) *BorderNode {
	n := &BorderNode{
		borderData: element.borderData,
	}
	n.ChainNode = node.NewChainNode(
		node.NewNodeID(),
		node.NodeKindDraw,
		node.DrawPhase,
		func(t TreeNode) {
			no := t.(DrawModifierNode)
			no.AttachDrawModifier(func(widget LayoutWidget) LayoutWidget {
				return layoutnode.NewLayoutWidget(func(gtx LayoutContext) LayoutDimensions {
					// Layout content first
					dims := widget.Layout(gtx)

					width := n.borderData.Width
					if width <= 0 {
						return dims
					}

					if !shape.IsSpecifiedShape(n.borderData.Shape) {
						panic("BorderNode: Shape is not specified")
					}

					if b := gtx.DrawBackend(); b != nil {
						// Software path: emit a stroke FillShape through the Backend.
						strokeWidth := float32(gtx.ToGio().Metric.Dp(unitconvert.DpToGioUnitUnsafe(width)))
						rShape := borderStrokeToRenderShape(n.borderData.Shape, dims.Size, strokeWidth, gtx)
						rColor := borderColorToRenderColor(n.borderData.Color)
						b.FillShape(rShape, rColor)
						return dims
					}

					// Gio path: unchanged.
					outline := n.borderData.Shape.CreateOutline(dims.Size, shape.Metric{PxPerDp: gtx.ToGio().Metric.PxPerDp, PxPerSp: gtx.ToGio().Metric.PxPerSp})
					macro := op.Record(gtx.ToGio().Ops)

					strokeWidth := float32(gtx.ToGio().Metric.Dp(unitconvert.DpToGioUnitUnsafe(width)))

					pathSpec := outline.Path(clipconvert.NewOps(gtx.ToGio().Ops)).ToClipPathSpec()

					// Create stroke op
					strokeOp := clip.Stroke{
						Path:  pathSpec,
						Width: strokeWidth,
					}.Op()

					nrgba := graphics.ColorToNRGBA(n.borderData.Color)

					// Paint the stroke
					paint.FillShape(gtx.ToGio().Ops, nrgba, strokeOp)

					call := macro.Stop()
					call.Add(gtx.ToGio().Ops)

					return dims
				})
			})
		},
	)
	return n
}

// borderStrokeToRenderShape converts a border's shape + stroke width into a
// render.Shape that represents the stroke. For golden tests we record the
// shape bounds and stroke width; the software backend does not rasterize
// the stroke path, only records the DrawCall.
func borderStrokeToRenderShape(s shape.Shape, size image.Point, strokeWidth float32, gtx LayoutContext) render.Shape {
	metric := shape.Metric{
		PxPerDp: gtx.ToGio().Metric.PxPerDp,
		PxPerSp: gtx.ToGio().Metric.PxPerSp,
	}
	outline := s.CreateOutline(size, metric)
	// For golden tests we represent the stroke as a rectangle with the
	// outline bounds. The stroke width is encoded as the Radius.TL field.
	_ = outline
	return render.Shape{
		Kind:   render.ShapeRectangle,
		Bounds: render.Rect{Min: render.Point{}, Max: render.Point{X: float32(size.X), Y: float32(size.Y)}},
		Radius: render.CornerRadius{TL: strokeWidth}, // stroke width sentinel
	}
}

// borderColorToRenderColor converts a framework Color to a render.Color.
func borderColorToRenderColor(c graphics.Color) render.Color {
	nrgba := graphics.ColorToNRGBA(c)
	return render.Color{
		R: float32(nrgba.R) / 255,
		G: float32(nrgba.G) / 255,
		B: float32(nrgba.B) / 255,
		A: float32(nrgba.A) / 255,
	}
}
