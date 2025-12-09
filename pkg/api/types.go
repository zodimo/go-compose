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

	layoutnode.TreeBuilderComposer

	state.StatefulComposer
}

// Public Modifier interface
type Modifier interface {
	// Then chains this modifier with another
	Then(other Modifier) Modifier
}

type RootNode interface{}
