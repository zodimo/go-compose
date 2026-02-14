package node

// // Build the path towards the root
// func GetNodePath(node Node) NodePath {
// 	nodeIDs := nodePathRecursive(node, []NodeID{})
// 	return NodePath{
// 		ids: nodeIDs,
// 	}

// }

// func nodePathRecursive(node Node, visitedNodes []NodeID) []NodeID {
// 	//walk the tree
// 	path := []NodeID{}
// 	switch n := node.(type) {
// 	case TreeNode:
// 		path = append(visitedNodes, node.GetID())
// 		for _, child := range n.Children() {
// 			path = append(path, nodePathRecursive(child, path)...)
// 		}

// 	case ChainNode:
// 		path = append(path, n.GetID())
// 		path = append(path, nodePathRecursive(n.Next(), path)...)
// 	}
// 	return append(visitedNodes, path...)
// }
