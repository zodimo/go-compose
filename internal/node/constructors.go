package node

type NewChainNodeOptions struct {
	onDetach func()
	nextNode ChainNode
}

type NewChainNodeOption func(options *NewChainNodeOptions)

func NewChainNodeWithOnDetach(onDetach func()) NewChainNodeOption {
	return func(options *NewChainNodeOptions) {
		options.onDetach = onDetach
	}
}

func NewChainNodeWithNextNode(nextNode ChainNode) NewChainNodeOption {
	return func(options *NewChainNodeOptions) {
		options.nextNode = nextNode
	}
}

func NewChainNode(
	id NodeID,
	kind NodeKind,
	phases Phases,
	onAttach func(node TreeNode),
	options ...NewChainNodeOption,

) ChainNode {
	opts := &NewChainNodeOptions{}
	for _, option := range options {
		if option == nil {
			continue
		}
		option(opts)
	}
	return &chainNode{
		id:       id,
		kind:     kind,
		phases:   phases,
		onAttach: onAttach,
		onDetach: opts.onDetach,
		nextNode: opts.nextNode,
	}
}

func NewNodePath(ids []NodeID) NodePath {
	return NodePath{
		ids: ids,
	}
}
