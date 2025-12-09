package layoutnode

import (
	node "go-compose-dev/internal/Node"
	"go-compose-dev/internal/modifier"

	"gioui.org/layout"
	"gioui.org/op"
)

type TreeNode = node.TreeNode
type ChainNode = node.ChainNode

type LayoutContext = layout.Context
type LayoutDimensions = layout.Dimensions
type LayoutWidget = layout.Widget
type DrawOp = op.CallOp

type NodeID = node.NodeID

type Modifier = modifier.Modifier

var EmptyModifier = modifier.EmptyModifier
