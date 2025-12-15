package api

import (
	node "github.com/zodimo/go-compose/internal/Node"
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/internal/modifier"
	idApi "github.com/zodimo/go-compose/pkg/compose-identifier/api"
	"github.com/zodimo/go-compose/state"
)

// compose-identifier.api.Identifier
type Identifier = idApi.Identifier

// type Composition = func(Composable) Composable
type Composable func(Composer) Composer

type MutableValue = state.MutableValue

type NodePath = node.NodePath

// Public API of the composer
type Composer interface {
	// --
	GetID() Identifier
	GetPath() NodePath

	modifier.ModifierAwareComposer

	// -- id management
	GenerateID() Identifier

	EmitSlot(k string, v any) Composer

	TreeBuilderComposer
	GioLayoutNodeAwareComposer

	state.SupportState

	WithComposable(composable Composable) Composer

	If(condition bool, ifTrue Composable, ifFalse Composable) Composable
	When(condition bool, ifTrue Composable) Composable
	Else(condition bool, ifFalse Composable) Composable

	Sequence(contents ...Composable) Composable

	// Control Flow
	Key(key any, content Composable) Composable
	Range(count int, fn func(int) Composable) Composable
}

// Public Modifier interface
type Modifier interface {
	// Then chains this modifier with another
	Then(other Modifier) Modifier
}

type LayoutNode = layoutnode.LayoutNode

type TreeBuilderComposer interface {
	StartBlock(id string) Composer
	EndBlock() Composer
	Build() LayoutNode
}

type GioLayoutNodeAwareComposer interface {
	SetWidgetConstructor(constructor layoutnode.LayoutNodeWidgetConstructor)
}
