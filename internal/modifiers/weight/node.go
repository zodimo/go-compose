package weight

import (
	node "go-compose-dev/internal/Node"
	"go-compose-dev/internal/layoutnode"
)

const WeightElementKey = "weight"

var _ ChainNode = (*WeightNode)(nil)

// NodeKind should also implement the interface of the LayoutNode for that phase

type WeightNode struct {
	ChainNode
	weight WeightData
}

func NewWeightNode(element WeightElement) ChainNode {
	return WeightNode{
		ChainNode: node.NewChainNode(
			node.NewNodeID(),
			node.NodeKindLayout,
			node.LayoutPhase,
			//OnAttach
			func(n TreeNode) {
				// how should the tree now be updated when attached
				// tree nde is the layout tree

				no := n.(layoutnode.ParentDataModifierNode)
				// we can now work with the layoutNode

				no.AttachParentDataModifier(func(store layoutnode.ElementStore) layoutnode.ElementStore {
					return store.SetElement(WeightElementKey, element)
				})

			},
		),
		weight: element.weight,
	}
}
