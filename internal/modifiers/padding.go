package modifiers

import (
	node "go-compose-dev/internal/Node"
	"go-compose-dev/internal/layoutnode"

	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/unit"
)

const PaddingNotSet = -1

const RTL = system.RTL

type Padding struct {
	Start    int
	Top      int
	End      int
	Bottom   int
	RtlAware bool // future proofing for RTL support
}

func DefaultPadding() Padding {
	return Padding{
		Start:    PaddingNotSet,
		Top:      PaddingNotSet,
		End:      PaddingNotSet,
		Bottom:   PaddingNotSet,
		RtlAware: false,
	}
}

func ComparePadding(a, b Padding) bool {
	return a.Start == b.Start && a.Top == b.Top && a.End == b.End && a.Bottom == b.Bottom && a.RtlAware == b.RtlAware
}

var _ ChainNode = (*PaddingNode)(nil)

// NodeKind should also implement the interface of the LayoutNode for that phase

type PaddingNode struct {
	ChainNode
	padding Padding
}

var _ Element = (*paddingElement)(nil)

// Hold the behavior
type paddingElement struct {
	padding Padding
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
				no.AttachLayoutModifier(func(gtx layoutnode.LayoutContext, widget layoutnode.LayoutWidget) layoutnode.LayoutDimensions {

					// Default is LTR
					left := unit.Dp(pe.padding.Start)
					right := unit.Dp(pe.padding.End)

					if pe.padding.RtlAware {
						// if RTL then we should swap left and right
						if gtx.Locale.Direction == RTL {
							left = unit.Dp(pe.padding.End)
							right = unit.Dp(pe.padding.Start)
						}
					}

					return layout.Inset{
						Top:    unit.Dp(pe.padding.Top),
						Bottom: unit.Dp(pe.padding.Bottom),
						Left:   left,
						Right:  right,
					}.Layout(gtx, widget)
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

	if pe.padding.Start != PaddingNotSet {
		pn.padding.Start = pe.padding.Start
	}
	if pe.padding.Top != PaddingNotSet {
		pn.padding.Top = pe.padding.Top
	}
	if pe.padding.End != PaddingNotSet {
		pn.padding.End = pe.padding.End
	}
	if pe.padding.Bottom != PaddingNotSet {
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
