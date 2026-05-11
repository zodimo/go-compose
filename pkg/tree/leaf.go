package tree

var _ Tree = (*TreeLeaf)(nil)

type TreeLeaf struct {
	id         TreeId
	parent     Tree
	properties map[string]any
}

// Tree
func (n *TreeLeaf) ID() TreeId {
	return n.id
}

func (n *TreeLeaf) Label() string {
	return n.id.Label()
}

func (l *TreeLeaf) String() string {
	return "Leaf"
}

func (n *TreeLeaf) FindById(id TreeId) (Tree, bool) {
	if n.id.Equals(id) {
		return n, true
	}
	return nil, false
}

func (n *TreeLeaf) FoldIn(initial interface{}, operation func(carry interface{}, item Tree) interface{}) interface{} {
	return operation(initial, n)
}
func (n *TreeLeaf) FoldOut(initial interface{}, operation func(element Tree, carry interface{}) interface{}) interface{} {
	return operation(n, initial)
}

func (l *TreeLeaf) isTree() {}

// Leaf
func (n *TreeLeaf) Parent() *TreeNode {
	if n.parent == nil {
		return nil
	}
	return n.parent.(*TreeNode)
}

func (n *TreeLeaf) HasParent() bool {
	return n.Parent() != nil
}

func (n *TreeLeaf) Properties() map[string]any {
	//clone map
	properties := make(map[string]any, len(n.properties))
	for k, v := range n.properties {
		properties[k] = v
	}
	return properties
}

var _ TreeLeafScope = (*treeLeafScopeImpl)(nil)

type treeLeafScopeImpl struct {
	parent      *TreeNode
	treeBuilder TreeBuilder
}

func (t *treeLeafScopeImpl) AddProperty(key string, value any) {
	t.treeBuilder.AddProperty(key, value)
}
