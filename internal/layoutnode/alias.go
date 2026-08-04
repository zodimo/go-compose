package layoutnode

import (
	"github.com/zodimo/go-compose/internal/identity"
	"github.com/zodimo/go-compose/internal/modifier"
	node "github.com/zodimo/go-compose/internal/node"
	"github.com/zodimo/go-compose/state"
)

type TreeNode = node.TreeNode
type ChainNode = node.ChainNode

type NodeID = node.NodeID

type ModifierElement = modifier.ModifierElement
type InspectableModifier = modifier.InspectableModifier

type Element = modifier.Element

type ElementStore = modifier.ElementStore

var EmptyElementStore = modifier.EmptyElementStore

type Identifier = identity.Identifier
type IdentityManager = identity.IdentityManager

var GetScopedIdentityManager = identity.GetScopedIdentityManager

type Memo = state.Memo
type PersistentState = state.PersistentState
type MutableValue = state.MutableValue
type StateOption = state.StateOption
type StateOptions = state.StateOptions
