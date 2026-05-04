package tree

var _ TreeBuilder = (*treeLeafBuilder)(nil)

type treeLeafBuilder struct {
	*TreeLeaf
}

func (lb *treeLeafBuilder) AddProperty(key string, value any) {
	lb.TreeLeaf.properties[key] = value
}

// func (lb *treeLeafBuilder) BuildMutableTree(c api.Composer) core.MutableValueTyped[Tree] {
// 	return remeberTree(c, &lb.TreeLeaf)
// }
// func (lb *treeLeafBuilder) BuildTree(c api.Composer) core.ValueTyped[Tree] {
// 	return lb.BuildMutableTree(c).(core.ValueTyped[Tree])
// }
