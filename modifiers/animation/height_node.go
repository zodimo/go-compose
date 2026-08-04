package animation

import (
	"image"

	"github.com/zodimo/go-compose/internal/layoutnode"
	node "github.com/zodimo/go-compose/internal/node"
	"github.com/zodimo/go-compose/internal/render"

	"gioui.org/op"
	"gioui.org/op/clip"
	gioUnit "gioui.org/unit"
)

var _ node.ChainNode = (*AnimatedHeightNode)(nil)

func NewAnimatedHeightNode(element AnimatedHeightElement) node.ChainNode {
	n := &AnimatedHeightNode{
		element: element,
	}
	n.ChainNode = node.NewChainNode(
		node.NewNodeID(),
		node.NodeKindLayout,
		node.LayoutPhase,
		func(t node.TreeNode) {
			no := t.(layoutnode.LayoutModifierNode)
			no.AttachLayoutModifier(func(widget layoutnode.LayoutWidget) layoutnode.LayoutWidget {
				return layoutnode.NewLayoutWidget(func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
					anim := n.element.Anim

					// Calculate progress
					progress := anim.Revealed(*gtx.ToGio())
					if progress == 0 && !anim.Visible() {
						return layoutnode.LayoutDimensions{}
					}

					if b := gtx.DrawBackend(); b != nil {
						// Software path:
						// 1. Record child layout (captures into a recording)
						childRec := b.Record()
						// Apply max height constraint
						childConstraints := gtx.ToGio().Constraints
						if n.element.MaxHeight > 0 {
							childConstraints.Max.Y = gtx.ToGio().Dp(gioUnit.Dp(n.element.MaxHeight))
						}
						g := *gtx.ToGio()
						g.Constraints = childConstraints
						dims := widget.Layout(layoutnode.NewLayoutContextWithBackend(&g, gtx.DrawBackend()))
						// childRec is not yet Applied

						targetHeight := dims.Size.Y
						currentHeight := int(float32(targetHeight) * progress)

						if progress < 1.0 {
							// Push clip for animated height
							saveLevel := b.Save()
							b.ClipRect(render.Rect{
								Min: render.Point{},
								Max: render.Point{X: float32(dims.Size.X), Y: float32(currentHeight)},
							}, render.CornerRadius{})
							b.Apply(childRec)
							b.Restore(saveLevel)
						} else {
							b.Apply(childRec)
						}

						return layoutnode.LayoutDimensions{
							Size:     image.Point{X: dims.Size.X, Y: currentHeight},
							Baseline: dims.Baseline,
						}
					}

					// Gio path: unchanged.
					// Measure content first
					macro := op.Record(gtx.ToGio().Ops)
					// Apply max height constraint constraint
					childConstraints := gtx.ToGio().Constraints
					if n.element.MaxHeight > 0 {
						childConstraints.Max.Y = gtx.ToGio().Dp(gioUnit.Dp(n.element.MaxHeight))
					}
					// Pass modified constraints (copy semantics: the wrapper holds a pointer)
					g := *gtx.ToGio()
					g.Constraints = childConstraints

					dims := widget.Layout(layoutnode.NewLayoutContextWithBackend(&g, gtx.DrawBackend()))
					call := macro.Stop()

					// Apply animation to height
					targetHeight := dims.Size.Y
					currentHeight := int(float32(targetHeight) * progress)

					// Clip to current height
					if progress < 1.0 {
						defer clip.Rect{Max: image.Point{X: dims.Size.X, Y: currentHeight}}.Push(gtx.ToGio().Ops).Pop()
					}

					// Draw Child
					call.Add(gtx.ToGio().Ops)

					return layoutnode.LayoutDimensions{
						Size:     image.Point{X: dims.Size.X, Y: currentHeight},
						Baseline: dims.Baseline,
					}
				})
			})
		},
	)
	return n
}

type AnimatedHeightNode struct {
	node.ChainNode
	element AnimatedHeightElement
}
