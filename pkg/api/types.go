package api

import (
	"time"

	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/internal/modifier"
	node "github.com/zodimo/go-compose/internal/node"
	"github.com/zodimo/go-compose/lifecycle"
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
	// frame time
	TimeNow() time.Time

	// --
	GetID() Identifier
	GetPath() NodePath

	modifier.ModifierAwareComposer

	// -- Composition Locals
	StartProviders(values []ProvidedValue) Composer
	EndProviders() Composer
	Consume(key interface{}) interface{}

	// -- id management
	GenerateID() Identifier

	EmitSlot(k string, v any) Composer

	TreeBuilderComposer
	GioLayoutNodeAwareComposer

	state.SupportState

	WithComposable(composable Composable) Composer

	If(condition bool, ifTrue Composable, ifFalse Composable) Composable
	IfLazy(condition bool, ifTrue func() Composable, ifFalse func() Composable) Composable

	When(condition bool, ifTrue Composable) Composable
	WhenLazy(condition bool, ifTrue func() Composable) Composable

	Else(condition bool, ifFalse Composable) Composable
	ElseLazy(condition bool, ifFalse func() Composable) Composable

	Sequence(contents ...Composable) Composable

	// Control Flow
	Key(key any, content Composable) Composable
	Range(count int, fn func(int) Composable) Composable

	// frame lifecycle
	lifecycle.FrameLifecycleAware
}

type LayoutNode = layoutnode.LayoutNode

type TreeBuilderComposer interface {
	StartBlock(key string) Composer
	EndBlock() Composer
	Build() LayoutNode
}

// UI Specific awareness
type GioLayoutNodeAwareComposer interface {
	//Should this not just emit the constructor to a known key in the slottble?
	SetWidgetConstructor(constructor layoutnode.LayoutNodeWidgetConstructor)
}

// Deprecated: use ui.EmptyModifier
var EmptyModifier = ui.EmptyModifier
