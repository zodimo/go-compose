package border

import (
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	"github.com/zodimo/go-compose/compose/ui/unit"
	node "github.com/zodimo/go-compose/internal/Node"
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/internal/modifier"
	"github.com/zodimo/go-compose/modifiers/helpers"
)

type Modifier = modifier.Modifier
type Element = modifier.Element
type InspectableModifier = modifier.InspectableModifier
type ModifierInspectorInfo = modifier.InspectorInfo

type Node = node.Node
type TreeNode = node.TreeNode
type ChainNode = node.ChainNode

type DrawModifierNode = layoutnode.DrawModifierNode
type LayoutWidget = layoutnode.LayoutWidget
type LayoutContext = layoutnode.LayoutContext
type LayoutDimensions = layoutnode.LayoutDimensions

var ToNRGBA = helpers.ToNRGBA

type Shape = shape.Shape
type Dp = unit.Dp
