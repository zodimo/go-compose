package animation

import (
	"image"

	node "github.com/zodimo/go-compose/internal/Node"
	"github.com/zodimo/go-compose/internal/layoutnode"

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
					progress := anim.Revealed(gtx)
					if progress == 0 && !anim.Visible() {
						return layoutnode.LayoutDimensions{}
					}

					// Measure content first
					macro := op.Record(gtx.Ops)
					// Apply max height constraint constraint
					childConstraints := gtx.Constraints
					if n.element.MaxHeight > 0 {
						childConstraints.Max.Y = gtx.Dp(gioUnit.Dp(n.element.MaxHeight))
					}
					// Pass modified constraints
					gtxChild := gtx
					gtxChild.Constraints = childConstraints

					dims := widget.Layout(gtxChild)
					call := macro.Stop()

					// Apply animation to height
					targetHeight := dims.Size.Y
					currentHeight := int(float32(targetHeight) * progress)

					// Clip to current height
					// Only clip if we are animating. If fully visible (progress=1.0), avoid clipping
					// to prevent cutting off shadows or hover effects that might extend outside bounds.
					if progress < 1.0 {
						defer clip.Rect{Max: image.Point{X: dims.Size.X, Y: currentHeight}}.Push(gtx.Ops).Pop()
					}

					// Draw Child
					call.Add(gtx.Ops)

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
