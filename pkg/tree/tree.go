package tree

import (
	"fmt"
)

// Tree is a sum type of Node and Leaf

type Tree interface {
	fmt.Stringer
	ID() TreeId
	Label() string
	isTree()
}

type TreeBuilder interface {
	Tree
	AddProperty(key string, value any)
	// BuildMutableTree(c api.Composer) core.MutableValueTyped[Tree]
	// BuildTree(c api.Composer) core.ValueTyped[Tree]
}

func UnwrapTreeBuilder(treeBuilder TreeBuilder) Tree {
	switch tree := treeBuilder.(type) {
	case *treeNodeBuilder:
		return tree.TreeNode
	case *treeLeafBuilder:
		return tree.TreeLeaf
	default:
		panic(fmt.Sprintf("Unknown tree type: %T", tree))
	}
}

type TreeLeafScope interface {
	AddProperty(key string, value any)
}

type TreeScope interface {
	AddNode(label string, builder func(scope TreeNodeScope))
	AddLeaf(label string, builder func(scope TreeLeafScope))
}

func NewTree(rootLabel string, builder func(scope TreeScope)) Tree {

	treeBuilder := &treeNodeBuilder{
		TreeNode: &TreeNode{
			id: NewTreeNodeId(rootLabel),
		},
	}
	treeScope := &treeScopeImpl{
		treeBuilder: treeBuilder,
	}
	builder(treeScope)
	return treeBuilder.TreeNode
}

var _ TreeScope = (*treeScopeImpl)(nil)

type treeScopeImpl struct {
	treeBuilder TreeBuilder
}

func (s *treeScopeImpl) AddNode(label string, builder func(scope TreeNodeScope)) {
	parent := UnwrapTreeBuilder(s.treeBuilder).(*TreeNode)
	child := TreeNode{
		id:         NewChildTreeId(label, s.treeBuilder.ID()),
		properties: map[string]any{},
	}
	node := treeNodeBuilder{
		TreeNode: &child,
	}
	parent.children = append(parent.children, &child)

	builder(&treeNodeScopeImpl{
		parent:      s.treeBuilder.(*treeNodeBuilder).TreeNode,
		treeBuilder: &node,
	})
}

func (s *treeScopeImpl) AddLeaf(label string, builder func(scope TreeLeafScope)) {

	parent := UnwrapTreeBuilder(s.treeBuilder).(*TreeNode)
	leaf := TreeLeaf{
		id:         NewChildTreeId(label, s.treeBuilder.ID()),
		properties: map[string]any{},
	}
	parent.children = append(parent.children, &leaf)

	treeBuilder := &treeLeafBuilder{
		&leaf,
	}

	builder(&treeLeafScopeImpl{
		treeBuilder: treeBuilder,
		parent:      parent,
	})
}
