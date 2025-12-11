package clip

import (
	"go-compose-dev/compose/ui/graphics/shape"
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
type PointerModifierNode = layoutnode.PointerInputModifierNode

type LayoutContext = layoutnode.LayoutContext
type LayoutWidget = layoutnode.LayoutWidget

var ToNRGBA = helpers.ToNRGBA

type Shape = shape.Shape

var ShapeCircle = shape.ShapeCircle
var ShapeRectangle = shape.ShapeRectangle
