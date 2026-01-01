package padding

import (
	node "github.com/zodimo/go-compose/internal/Node"
	"github.com/zodimo/go-compose/internal/layoutnode"

	"gioui.org/layout"
	gioUnit "gioui.org/unit"
)

func ComparePadding(a, b PaddingData) bool {
	return a.Start == b.Start && a.Top == b.Top && a.End == b.End && a.Bottom == b.Bottom && a.RtlAware == b.RtlAware
}

var _ Element = (*paddingElement)(nil)

// Hold the behavior
type paddingElement struct {
	padding PaddingData
}

// Create creates a new Chain Node instance
func (pe paddingElement) Create() Node {
	//chainNode
	return PaddingNode{
		ChainNode: node.NewChainNode(
			node.NewNodeID(),
			node.NodeKindLayout,
			node.LayoutPhase,
			func(n TreeNode) {
				// how should the tree now be updated when attached
				// tree nde is the layout tree

				no := n.(layoutnode.LayoutModifierNode)
				// we can now work with the layoutNode
				no.AttachLayoutModifier(func(widget layoutnode.LayoutWidget) layoutnode.LayoutWidget {
					return layoutnode.NewLayoutWidget(func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
						// Default is LTR
						left := gioUnit.Dp(pe.padding.Start)
						right := gioUnit.Dp(pe.padding.End)

						if pe.padding.RtlAware {
							// if RTL then we should swap left and right
							if gtx.Locale.Direction == RTL {
								left = gioUnit.Dp(pe.padding.End)
								right = gioUnit.Dp(pe.padding.Start)
							}
						}

						return layout.Inset{
							Top:    gioUnit.Dp(pe.padding.Top),
							Bottom: gioUnit.Dp(pe.padding.Bottom),
							Left:   left,
							Right:  right,
						}.Layout(gtx, widget.Layout)
					})
				})

			},
		),
		padding: pe.padding,
	}

}

// Update updates an existing Chain node for efficiency
func (pe paddingElement) Update(node Node) {
	if node == nil {
		panic("node cannot be nil")
	}

	pn := node.(PaddingNode)

	if pe.padding.Start != NotSet {
		pn.padding.Start = pe.padding.Start
	}
	if pe.padding.Top != NotSet {
		pn.padding.Top = pe.padding.Top
	}
	if pe.padding.End != NotSet {
		pn.padding.End = pe.padding.End
	}
	if pe.padding.Bottom != NotSet {
		pn.padding.Bottom = pe.padding.Bottom
	}
	pn.padding.RtlAware = pe.padding.RtlAware

}

// Equals checks if this element is equivalent to another
// used during filter operations like Modifier.Any
func (pe paddingElement) Equals(other Element) bool {
	if otherElement, ok := other.(paddingElement); ok {
		return ComparePadding(pe.padding, otherElement.padding)
	}
	return false
}

func (pe paddingElement) Padding() PaddingData {
	return pe.padding
}
