package animation

import (
	"image"

	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/internal/layoutnode"
	node "github.com/zodimo/go-compose/internal/node"
	"github.com/zodimo/go-compose/internal/render"

	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type AnimatedBackgroundNode struct {
	node.ChainNode
	element AnimatedBackgroundElement
}

func NewAnimatedBackgroundNode(element AnimatedBackgroundElement) *AnimatedBackgroundNode {
	n := &AnimatedBackgroundNode{
		element: element,
	}
	n.ChainNode = node.NewChainNode(
		node.NewNodeID(),
		node.NodeKindLayout,
		node.DrawPhase|node.LayoutPhase,
		func(t node.TreeNode) {
			no := t.(layoutnode.LayoutModifierNode)
			no.AttachLayoutModifier(func(widget layoutnode.LayoutWidget) layoutnode.LayoutWidget {
				return layoutnode.NewLayoutWidget(func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
					// 1. Logic
					progress := n.element.Anim.Revealed(*gtx.ToGio())

					// 2. Layout Child
					dims := widget.Layout(gtx)

					// 3. Paint Background
					nrgba := graphics.ColorToNRGBA(n.element.Color)
					nrgba.A = uint8(float32(nrgba.A) * progress)

					if b := gtx.DrawBackend(); b != nil {
						// Software path: emit FillShape through Backend.
						rShape := render.Shape{
							Kind:   render.ShapeRectangle,
							Bounds: render.Rect{
								Min: render.Point{},
								Max: render.Point{X: float32(dims.Size.X), Y: float32(dims.Size.Y)},
							},
						}
						rColor := render.Color{
							R: float32(nrgba.R) / 255,
							G: float32(nrgba.G) / 255,
							B: float32(nrgba.B) / 255,
							A: float32(nrgba.A) / 255,
						}
						b.FillShape(rShape, rColor)
						return dims
					}

					// Gio path: unchanged.
					rect := image.Rectangle{Max: dims.Size}
					paint.FillShape(gtx.ToGio().Ops, nrgba, clip.Rect(rect).Op())

					return dims
				})
			})
		},
	)
	return n
}
