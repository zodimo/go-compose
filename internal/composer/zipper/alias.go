package zipper

import (
	node "github.com/zodimo/go-compose/internal/Node"
	"github.com/zodimo/go-compose/internal/identity"
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/internal/modifier"
	"github.com/zodimo/go-compose/internal/sequence"
	"github.com/zodimo/go-compose/pkg/api"
	idApi "github.com/zodimo/go-compose/pkg/compose-identifier/api"
	"github.com/zodimo/go-compose/state"
)

type LayoutNode = layoutnode.LayoutNode
type Memo = state.MemoTyped[any]
type ElementMemo = state.MemoTyped[Element]

type PersistentState = state.PersistentState

var EmptyMemo = state.EmptyMemo[any]()
var EmptyElementMemo = state.EmptyMemo[Element]()

type IdentityManager = identity.IdentityManager

var GetScopedIdentityManager = identity.GetScopedIdentityManager

// compose-identifier.api.Identifier
type Identifier = idApi.Identifier // Public API of the composer
type NodePath = node.NodePath
type Modifier = modifier.Modifier
type Element = modifier.Element
type MutableValue = state.MutableValue
type RootNode = node.TreeNode

type Composable = api.Composable
type Composer = api.Composer

type ProvidedValue = api.ProvidedValue

var Sequence = sequence.Sequence
