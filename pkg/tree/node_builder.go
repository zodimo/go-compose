package tree

var _ TreeBuilder = (*treeNodeBuilder)(nil)

type treeNodeBuilder struct {
	*TreeNode
}

func (lb *treeNodeBuilder) AddProperty(key string, value any) {
	lb.TreeNode.properties[key] = value
}

// func (lb *treeNodeBuilder) BuildMutableTree(c api.Composer) core.MutableValueTyped[Tree] {
// 	// build children first
// 	children := []core.MutableValueTyped[Tree]{}
// 	for _, tree := range lb.children {
// 		children = append(children, remeberTree(c, tree))
// 	}

// 	return nil
// }
// func (lb *treeNodeBuilder) BuildTree(c api.Composer) core.ValueTyped[Tree] {
// 	// build children first
// 	return nil
// }
