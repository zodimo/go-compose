package formengine

// VisitFunc is called for each node during traversal.
type VisitFunc func(node FormNode)

// Walk visits node and all descendants depth-first in root-to-leaf (pre-order)
// order. It is a companion to the parent-pointer propagation used by onChildChanged
// and is available for external traversal and collection.
func Walk(root FormNode, visit VisitFunc) {
	if root == nil || visit == nil {
		return
	}
	visit(root)
	for _, child := range root.Children() {
		Walk(child, visit)
	}
}

// WalkUp visits all descendants depth-first in leaf-to-root (post-order) order.
func WalkUp(node FormNode, visit VisitFunc) {
	if node == nil || visit == nil {
		return
	}
	for _, child := range node.Children() {
		WalkUp(child, visit)
	}
	visit(node)
}
