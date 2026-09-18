package formengine

// FormNode is the universal interface implemented by every node in a form tree:
// leaf Control[T] nodes, Group containers (named children), and Array containers
// (indexed children).
//
// The trailing three methods (setParent, onChildChanged, setInheritedDisabled)
// are unexported. In Go, unexported interface methods restrict implementations of
// this interface to this package, which is intentional: the tree engine is
// internal. setParent and onChildChanged are specified by design D6;
// setInheritedDisabled is an additional internal hook required by design D4's
// inherited-disabled propagation (D6's interface omits it despite D4 requiring
// downward propagation).
type FormNode interface {
	// Status & validation
	Status() Status
	Validate() bool
	Errors() map[string]string // keyed by path for groups; controls join all errors via errors.Join

	// Lifecycle
	IsTouched() bool
	IsDirty() bool
	IsPristine() bool
	IsEnabled() bool

	// Actions
	MarkAsTouched()
	MarkAsUntouched()
	MarkAsPristine()
	SetDisabled(disabled bool)
	Reset()

	// Tree
	Parent() FormNode
	Children() []FormNode
	Path() string // dotted path from root
	Get(path string) (FormNode, bool)

	// Export
	RawValue() any

	// Internal hooks (unexported on impls; promoted via interface)
	setParent(p FormNode)
	onChildChanged() // propagation (D3)
	setInheritedDisabled(bool)
}
