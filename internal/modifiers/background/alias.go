package background

import (
	node "go-compose-dev/internal/Node"
	"go-compose-dev/internal/layoutnode"
	"go-compose-dev/internal/modifier"
	"go-compose-dev/internal/modifiers/helpers"
)

type Modifier = modifier.Modifier
type Element = modifier.Element
type InspectableModifier = modifier.InspectableModifier
type ModifierInspectorInfo = modifier.InspectorInfo

type Node = node.Node
type TreeNode = node.TreeNode
type ChainNode = node.ChainNode

var EmptyModifier = modifier.EmptyModifier

type DrawModifierNode = layoutnode.DrawModifierNode
type LayoutModifierNode = layoutnode.LayoutModifierNode
type PointerModifierNode = layoutnode.PointerModifierNode

type LayoutContext = layoutnode.LayoutContext
type LayoutWidget = layoutnode.LayoutWidget

var ToNRGBA = helpers.ToNRGBA
