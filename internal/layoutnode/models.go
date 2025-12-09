package layoutnode

import "go-compose-dev/internal/immap"

var _ LayoutNode = (*layoutNode)(nil)

type layoutNode struct {
	id       NodeID
	key      string
	slots    immap.ImmutableMap[any]
	children []LayoutNode
	modifier Modifier
}

// Node
func (ln *layoutNode) GetID() NodeID {
	return ln.id
}

// TreeNode
func (ln *layoutNode) Children() []TreeNode {
	treeNodeChildren := []TreeNode{}
	for _, child := range ln.children {
		treeNodeChildren = append(treeNodeChildren, child)
	}
	return treeNodeChildren
}

func (ln *layoutNode) LayoutNodeChildren() []LayoutNode {
	return ln.children
}

//LayoutNode

func (ln *layoutNode) IsEmpty() bool {
	return len(ln.children) == 0
}
func (ln *layoutNode) Modifier(apply func(modifier Modifier) Modifier) {
	ln.modifier = apply(ln.modifier)
}

func (ln *layoutNode) WithChildren(children []LayoutNode) LayoutNode {
	ln.children = children
	return ln
}

func (n *layoutNode) WithSlotsAssoc(k string, v any) LayoutNode {
	n.slots = n.slots.Assoc(k, v)
	return n
}

func NewLayoutNode(id NodeID, key string, slotStore immap.ImmutableMap[any]) LayoutNode {
	return &layoutNode{
		id:       id,
		key:      key,
		children: []LayoutNode{},
		modifier: EmptyModifier,
		slots:    slotStore,
	}
}
