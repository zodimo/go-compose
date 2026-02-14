package background

import (
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/internal/modifier"
	node "github.com/zodimo/go-compose/internal/node"
	"github.com/zodimo/go-compose/modifiers/helpers"
)

type Element = modifier.Element
type InspectableModifier = modifier.InspectableModifier
type ModifierInspectorInfo = modifier.InspectorInfo

type Node = node.Node
type TreeNode = node.TreeNode
type ChainNode = node.ChainNode

type DrawModifierNode = layoutnode.DrawModifierNode
type LayoutModifierNode = layoutnode.LayoutModifierNode
type PointerModifierNode = layoutnode.PointerInputModifierNode

type LayoutContext = layoutnode.LayoutContext
type LayoutWidget = layoutnode.LayoutWidget

var ToNRGBA = helpers.ToNRGBA

type Shape = shape.Shape

var ShapeCircle = shape.ShapeCircle
var ShapeRectangle = shape.ShapeRectangle
