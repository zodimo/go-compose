package modifiers

import (
	"go-compose-dev/compose/ui/graphics"
	node "go-compose-dev/internal/Node"
	"go-compose-dev/internal/layoutnode"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type Background struct {
	color color.Color
	shape graphics.Shape
}

func CompareBackground(a, b Background) bool {
	return a.color == b.color && a.shape == b.shape
}

var _ ChainNode = (*BackgroundNode)(nil)

// NodeKind should also implement the interface of the LayoutNode for that phase

type BackgroundNode struct {
	ChainNode
	background Background
}

var _ Element = (*backgroundElement)(nil)

// Hold the behavior
type backgroundElement struct {
	background Background
}

// Create creates a new Chain Node instance
func (be backgroundElement) Create() Node {
	//chainNode
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
				no.AttachDrawModifier(func(gtx layoutnode.LayoutContext, widget layoutnode.LayoutWidget) {
					layout.Background{}.Layout(gtx,
						func(gtx layout.Context) layout.Dimensions {
							// shape
							// color
							defer clip.Rect{Max: gtx.Constraints.Min}.Push(gtx.Ops).Pop()

							paint.Fill(gtx.Ops, ToNRGBA(be.background.color))

							return layout.Dimensions{Size: gtx.Constraints.Min}

						},
						func(gtx layout.Context) layout.Dimensions {
							return widget(gtx)
						},
					)
				})

			},
		),
		background: be.background,
	}
}

// Update updates an existing Chain node for efficiency
func (be backgroundElement) Update(node Node) {
	if node == nil {
		panic("node cannot be nil")
	}

	bn := node.(BackgroundNode)
	bn.background = be.background

}

// Equals checks if this element is equivalent to another
// used during filter operations like Modifier.Any
func (be backgroundElement) Equals(other Element) bool {
	if otherElement, ok := other.(backgroundElement); ok {
		return CompareBackground(be.background, otherElement.background)
	}
	return false
}
