package api

import (
	node "go-compose-dev/internal/Node"
	"go-compose-dev/internal/layoutnode"
	"go-compose-dev/internal/modifier"
	"go-compose-dev/internal/state"
	idApi "go-compose-dev/pkg/compose-identifier/api"
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

	TreeBuilderComposer
	GioLayoutNodeAwareComposer

	state.StatefulComposer

	WithComposable(composable Composable) Composer
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
