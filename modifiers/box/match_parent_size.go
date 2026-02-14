package box

import (
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/internal/modifier"
	node "github.com/zodimo/go-compose/internal/node"
)

const MatchParentSizeKey = "box_match_parent_size"

// --- Element ---

type matchParentSizeElement struct {
}

var _ modifier.Element = (*matchParentSizeElement)(nil)

func (e *matchParentSizeElement) Create() node.Node {
	return NewMatchParentSizeNode()
}

func (e *matchParentSizeElement) Update(n node.Node) {
	// Nothing to update, it's a flag
}

func (e *matchParentSizeElement) Equals(other modifier.Element) bool {
	_, ok := other.(*matchParentSizeElement)
	return ok
}

// --- Node ---

type matchParentSizeNode struct {
	node.ChainNode
}

func NewMatchParentSizeNode() node.ChainNode {
	// matchParentSizeElement is stateless, so we can use a shared instance or new one.
	element := &matchParentSizeElement{}
	return &matchParentSizeNode{
		ChainNode: node.NewChainNode(
			node.NewNodeID(),
			node.NodeKindLayout,
			node.LayoutPhase,
			func(n node.TreeNode) {
				no := n.(layoutnode.ParentDataModifierNode)
				no.AttachParentDataModifier(func(store layoutnode.ElementStore) layoutnode.ElementStore {
					return store.SetElement(MatchParentSizeKey, element)
				})
			},
		),
	}
}

// --- Public API ---

func MatchParentSize() ui.Modifier {
	return modifier.NewInspectableModifier(
		modifier.NewModifier(&matchParentSizeElement{}),
		modifier.NewInspectorInfo("MatchParentSize", nil),
	)
}
