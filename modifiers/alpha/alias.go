package alpha

import (
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/internal/modifier"
	node "github.com/zodimo/go-compose/internal/node"
)

type Element = modifier.Element
type InspectableModifier = modifier.InspectableModifier
type ModifierInspectorInfo = modifier.InspectorInfo

type Node = node.Node
type TreeNode = node.TreeNode
type ChainNode = node.ChainNode

type DrawModifierNode = layoutnode.DrawModifierNode
type LayoutContext = layoutnode.LayoutContext
type LayoutWidget = layoutnode.LayoutWidget
type LayoutDimensions = layoutnode.LayoutDimensions
