package api

import (
	node "github.com/zodimo/go-compose/internal/Node"
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/internal/modifier"
	idApi "github.com/zodimo/go-compose/pkg/compose-identifier/api"
	"github.com/zodimo/go-compose/state"
)

// Identifier is a unique identifier for a composable node in the composition tree.
// It is used to track nodes across recompositions.
type Identifier = idApi.Identifier

// Composable is the fundamental building block of the UI.
// It is a function that takes a Composer and returns a Composer, representing a transformation
// or emission of UI elements into the composition tree.
//
// Unlike Jetpack Compose where @Composable is an annotation, here it is an explicit function signature.
//
// Example:
//
//	func MyComponent(c api.Composer) api.Composer {
//	    return c.Sequence(
//	        m3Text.Text("Hello"),
//	        m3Button.Filled(func() { ... }, "Click Me"),
//	    )
//	}
type Composable func(Composer) Composer

// MutableValue is a value holder that can be observed for changes.
// When the value changes, it triggers a recomposition of the scopes that read it.
// This is an alias to state.MutableValue.
type MutableValue = state.MutableValue

// NodePath represents the path to a node in the composition tree.
type NodePath = node.NodePath

// Composer is the interface that orchestrates the composition process.
// It manages the tree of composables, handles state, and builds the final layout tree.
//
// Users interact with the Composer primarily to:
//   - Define the structure of the UI using methods like Sequence, If, When.
//   - Manage state using Remember and State (via state.SupportState).
//   - Apply modifiers to components.
//   - Control flow using Key and Range.
type Composer interface {
	// GetID returns the unique identifier of the current composable node.
	GetID() Identifier

	// GetPath returns the path of the current node in the composition tree.
	GetPath() NodePath

	modifier.ModifierAwareComposer

	// GenerateID generates a new unique identifier for a child node.
	GenerateID() Identifier

	// EmitSlot emits a value into a named slot of the current node.
	// This is used for advanced component composition where data needs to be passed
	// to the underlying layout node.
	EmitSlot(k string, v any) Composer

	TreeBuilderComposer
	GioLayoutNodeAwareComposer

	state.SupportState

	// WithComposable executes the given Composable with this Composer.
	// It is equivalent to calling composable(c).
	WithComposable(composable Composable) Composer

	// If conditionally executes one of two Composables based on the boolean condition.
	// If condition is true, ifTrue is executed; otherwise, ifFalse is executed.
	If(condition bool, ifTrue Composable, ifFalse Composable) Composable

	// When conditionally executes a Composable if the boolean condition is true.
	// If condition is false, it behaves like an empty composable.
	When(condition bool, ifTrue Composable) Composable

	// Else conditionally executes a Composable if the boolean condition is false.
	// It acts as the inverse of When.
	Else(condition bool, ifFalse Composable) Composable

	// Sequence executes a list of Composables in order.
	// This is the primary way to group multiple components together.
	Sequence(contents ...Composable) Composable

	// Key identifies a block of execution with a specific key.
	// This is useful for maintaining state when the order of items in a list changes.
	Key(key any, content Composable) Composable

	// Range loops count times and executes the function fn for each index.
	// It is used for rendering lists or repeating elements.
	Range(count int, fn func(int) Composable) Composable
}

// Modifier is an interface for objects that can modify the behavior or appearance of a UI element.
// Modifiers are chained together to apply multiple effects.
type Modifier interface {
	// Then chains this modifier with another modifier.
	// It returns a new Modifier that represents the combination of the two.
	Then(other Modifier) Modifier
}

// LayoutNode represents a node in the layout tree produced by the composition.
// It contains the information needed by the runtime to measure, layout, and draw the UI.
type LayoutNode = layoutnode.LayoutNode

// TreeBuilderComposer provides methods for building the composition tree structure.
// These methods are typically used internally by framework components but are exposed
// for advanced custom component creation.
type TreeBuilderComposer interface {
	// StartBlock starts a new group or node in the composition tree with the given key.
	StartBlock(key string) Composer

	// EndBlock ends the current group or node.
	EndBlock() Composer

	// Build finalizes the composition and returns the root of the generated layout tree.
	Build() LayoutNode
}

// GioLayoutNodeAwareComposer allows setting the widget constructor for the current layout node.
// This bridges the composition world with the underlying Gio widgets.
type GioLayoutNodeAwareComposer interface {
	// SetWidgetConstructor sets the function responsible for creating the Gio widget
	// associated with the current layout node.
	SetWidgetConstructor(constructor layoutnode.LayoutNodeWidgetConstructor)
}
