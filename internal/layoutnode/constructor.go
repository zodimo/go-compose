package layoutnode

import "go-compose-dev/internal/immap"

func NewLayoutNode(id NodeID, key string, slotStore immap.ImmutableMap[any]) LayoutNode {
	return &layoutNode{
		id:       id,
		key:      key,
		children: []LayoutNode{},
		modifier: EmptyModifier,
		slots:    slotStore,
	}
}

func NewNodeCoordinator(node layoutNode) NodeCoordinator {
	return &nodeCoordinator{
		layoutNode: node,
	}
}
