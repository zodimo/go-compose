package modifier

// Element represents a modifier that can be applied to create a Node
type Element interface {
	// Create creates a new Chain Node instance
	Create() Node

	// Update updates an existing Chain node for efficiency
	Update(node Node)

	// Equals checks if this element is equivalent to another
	// used during filter operations like Modifier.Any
	Equals(other Element) bool
}

type Chain interface {
	// Fold folds the modifier chain, alias for FoldIn
	Fold(initial interface{}, operation func(interface{}, Modifier) interface{}) interface{}

	FoldIn(initial interface{}, operation func(interface{}, Modifier) interface{}) interface{}

	FoldOut(initial interface{}, operation func(interface{}, Modifier) interface{}) interface{}

	// Any returns true if any element matches the predicate
	Any(predicate func(Modifier) bool) bool
}

type ModifierChain interface {
	Modifier
	Chain
}

type ModifierElementChain interface {
	Modifier
	Element
	Chain
}

// The base type of the NodeModifiers
// Modifier represents a chain of modifier elements
type Modifier interface {
	// Then chains this modifier with another
	Then(other Modifier) Modifier
	AsChain() ModifierChain
}

type InspectableModifier interface {
	Modifier
	InspectorInfo() *InspectorInfo
	Unwrap() Modifier
}

type ModifierElement interface {
	Modifier
	Element
}
