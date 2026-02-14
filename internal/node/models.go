package node

import "strings"

type NodePath struct {
	ids []NodeID
}

func (np NodePath) Unwrap() []NodeID {
	return np.ids
}

func (np NodePath) String() string {
	var stringIDs = []string{}
	for _, NodeID := range np.ids {
		stringIDs = append(stringIDs, NodeID.String())
	}
	return strings.Join(stringIDs, "/")
}

var _ ChainNode = (*chainNode)(nil)

type chainNode struct {
	id       NodeID
	kind     NodeKind
	phases   Phases
	onAttach func(node TreeNode)
	onDetach func()
	nextNode ChainNode
}

// Node Interface
func (cn *chainNode) GetID() NodeID {
	return cn.id
}

// Kind returns a unique identifier for this node type
// should this be here or should be type assert ?
func (cn *chainNode) Kind() NodeKind {
	return cn.kind
}

// Bitset of phases a node cares about.
func (cn *chainNode) Phases() Phases {
	return cn.phases
}

// Attach is called when the node is attached to a node tree
func (cn *chainNode) Attach(node TreeNode) {
	cn.onAttach(node)
}

// Detach is called when the node is detached from a node tree
func (cn *chainNode) Detach() {
	cn.onDetach()
}

// Next returns the next node in the chain
func (cn *chainNode) Next() ChainNode {
	return cn.nextNode
}

// SetNext sets the next node in the chain
func (cn *chainNode) SetNext(node ChainNode) {
	cn.nextNode = node
}
