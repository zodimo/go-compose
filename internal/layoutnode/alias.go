package layoutnode

import (
	node "github.com/zodimo/go-compose/internal/Node"
	"github.com/zodimo/go-compose/internal/identity"
	"github.com/zodimo/go-compose/internal/immap"
	"github.com/zodimo/go-compose/internal/modifier"
	"github.com/zodimo/go-compose/state"

	"gioui.org/layout"
	"gioui.org/op"
)

type TreeNode = node.TreeNode
type ChainNode = node.ChainNode

type LayoutContext = layout.Context
type LayoutDimensions = layout.Dimensions
type LayoutConstraints = layout.Constraints
type GioLayoutWidget = layout.Widget

type DrawOp = op.CallOp

type NodeID = node.NodeID

type Modifier = modifier.Modifier
type ModifierElement = modifier.ModifierElement
type InspectableModifier = modifier.InspectableModifier

type Element = modifier.Element

type ElementStore = modifier.ElementStore

var EmptyElementStore = modifier.EmptyElementStore

var EmptyModifier = modifier.EmptyModifier

type Identifier = identity.Identifier
type IdentityManager = identity.IdentityManager

var GetScopedIdentityManager = identity.GetScopedIdentityManager

type SupportState = state.SupportState
type Memo = state.Memo
type PersistentState = state.PersistentState
type MutableValue = state.MutableValue

type Slots = immap.ImmutableMap[any]
