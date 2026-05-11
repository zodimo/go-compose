package tree

import (
	"slices"
	"strings"
)

type TreeNodeScope interface {
	AddChild(builder func(scope TreeScope))
	AddProperty(key string, value any)
}

var _ TreeNodeScope = (*treeNodeScopeImpl)(nil)

type treeNodeScopeImpl struct {
	parent      *TreeNode
	treeBuilder TreeBuilder
}

func (t *treeNodeScopeImpl) AddChild(builder func(scope TreeScope)) {
	treeScope := &treeScopeImpl{
		treeBuilder: t.treeBuilder,
	}
	builder(treeScope)
}
func (t *treeNodeScopeImpl) AddProperty(key string, value any) {
	t.treeBuilder.AddProperty(key, value)
}

var _ Tree = (*TreeNode)(nil)

type TreeNode struct {
	id         TreeId
	parent     Tree
	children   []Tree
	properties map[string]any
}

// Tree
func (n *TreeNode) ID() TreeId {
	return n.id
}
func (n *TreeNode) Label() string {
	return n.id.Label()
}

func (n *TreeNode) String() string {
	return "Node"
}

func (n *TreeNode) FindById(id TreeId) (Tree, bool) {
	if n.id.Equals(id) {
		return n, true
	}
	if strings.HasPrefix(id.String(), n.id.String()) {
		for _, child := range n.children {
			found, ok := child.FindById(id)
			if ok {
				return found, true
			}
		}
	}
	return nil, false
}

func (n *TreeNode) FoldIn(initial interface{}, operation func(carry interface{}, item Tree) interface{}) interface{} {
	carry := initial

	// Recurse into children (Folding Inwward)
	// Iterate backwards to simulate foldr without mutating the original slice
	for i := len(n.children) - 1; i >= 0; i-- {
		carry = n.children[i].FoldIn(carry, operation)
	}

	// Process current node
	carry = operation(carry, n)
	return carry

}

func (n *TreeNode) FoldOut(initial interface{}, operation func(element Tree, carry interface{}) interface{}) interface{} {
	// Process current node
	carry := operation(n, initial)

	// Recurse into children (Folding Outward)
	for _, child := range n.children {
		carry = child.FoldOut(carry, operation)
	}
	return carry
}

func (n *TreeNode) isTree() {}

// TreeNode
func (n *TreeNode) Parent() *TreeNode {
	if n.parent == nil {
		return nil
	}
	return n.parent.(*TreeNode)
}

func (n *TreeNode) HasParent() bool {
	return n.Parent() != nil
}

func (n *TreeNode) Children() []Tree {
	return slices.Clone(n.children)
}

func (n *TreeNode) Properties() map[string]any {
	//clone map
	properties := make(map[string]any, len(n.properties))
	for k, v := range n.properties {
		properties[k] = v
	}
	return properties
}
