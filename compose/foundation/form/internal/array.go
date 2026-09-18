package formengine

import (
	"sync"

	"github.com/zodimo/go-compose/state"
)

// Array is a container node with indexed children.
//
// Deviation from the design sketch's "Array.Get(index)": the FormNode interface
// (design D6) reserves the method name Get for path resolution, Get(path string)
// (FormNode, bool). Go forbids method overloading, so indexed access is exposed
// as At(index) (FormNode, bool) instead of Get(index).
type Array struct {
	mu                sync.RWMutex
	parent            FormNode
	children          []FormNode
	ownDisabled       bool
	inheritedDisabled bool

	valueStream  *typedStream[[]any]
	statusStream *typedStream[Status]
	touchStream  *typedStream[bool]
}

var _ FormNode = (*Array)(nil)

// NewArray creates an array with the given children and wires their parent
// pointers to the array.
func NewArray(children ...FormNode) *Array {
	a := &Array{
		children:     children,
		valueStream:  newTypedStream([]any{}),
		statusStream: newTypedStream(StatusValid),
		touchStream:  newTypedStream(false),
	}
	for _, child := range children {
		child.setParent(a)
	}
	return a
}

// Add appends node, fires notification and propagates up.
func (a *Array) Add(node FormNode) {
	node.setParent(a)
	a.mu.Lock()
	a.children = append(a.children, node)
	a.mu.Unlock()

	a.notifyValue(a.rawValueSlice())
	a.notifyStatus(a.Status())
	a.notifyTouch(a.IsTouched())
	a.propagateUp()
}

// Insert inserts node at index (clamped to [0, len]), shifting later children.
func (a *Array) Insert(index int, node FormNode) {
	if index < 0 {
		index = 0
	}
	node.setParent(a)
	a.mu.Lock()
	if index > len(a.children) {
		index = len(a.children)
	}
	a.children = append(a.children, nil)
	copy(a.children[index+1:], a.children[index:])
	a.children[index] = node
	a.mu.Unlock()

	a.notifyValue(a.rawValueSlice())
	a.notifyStatus(a.Status())
	a.notifyTouch(a.IsTouched())
	a.propagateUp()
}

// Remove removes the child at index (no-op if out of range), shifting later
// children left, and fires notification up the chain.
func (a *Array) Remove(index int) {
	a.mu.Lock()
	if index < 0 || index >= len(a.children) {
		a.mu.Unlock()
		return
	}
	removed := a.children[index]
	a.children = append(a.children[:index], a.children[index+1:]...)
	a.mu.Unlock()

	removed.setParent(nil)

	a.notifyValue(a.rawValueSlice())
	a.notifyStatus(a.Status())
	a.notifyTouch(a.IsTouched())
	a.propagateUp()
}

// Length returns the number of children.
func (a *Array) Length() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.children)
}

// At returns the child at index, or (nil, false) when out of range.
func (a *Array) At(index int) (FormNode, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if index < 0 || index >= len(a.children) {
		return nil, false
	}
	return a.children[index], true
}

// Status aggregates child statuses with the same precedence as Group.
func (a *Array) Status() Status {
	if a.effectiveDisabled() {
		return StatusDisabled
	}
	children := a.childrenList()
	if len(children) == 0 {
		return StatusValid
	}
	allDisabled := true
	anyPending := false
	for _, child := range children {
		cs := child.Status()
		if cs == StatusInvalid {
			return StatusInvalid
		}
		if cs == StatusPending {
			anyPending = true
		}
		if cs != StatusDisabled {
			allDisabled = false
		}
	}
	if anyPending {
		return StatusPending
	}
	if allDisabled {
		return StatusDisabled
	}
	return StatusValid
}

// Validate validates all children and reports whether the subtree is valid.
func (a *Array) Validate() bool {
	if a.effectiveDisabled() {
		return true
	}
	valid := true
	for _, child := range a.childrenList() {
		if !child.Validate() {
			valid = false
		}
	}
	return valid
}

// Errors merges all descendant control errors keyed by dotted path.
func (a *Array) Errors() map[string]string {
	if a.effectiveDisabled() {
		return nil
	}
	result := make(map[string]string)
	for _, child := range a.childrenList() {
		for k, v := range child.Errors() {
			result[k] = v
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

// IsTouched reports whether any descendant is touched.
func (a *Array) IsTouched() bool {
	for _, child := range a.childrenList() {
		if child.IsTouched() {
			return true
		}
	}
	return false
}

// IsDirty reports whether any descendant is dirty.
func (a *Array) IsDirty() bool {
	for _, child := range a.childrenList() {
		if child.IsDirty() {
			return true
		}
	}
	return false
}

// IsPristine is the negation of IsDirty.
func (a *Array) IsPristine() bool {
	return !a.IsDirty()
}

// IsEnabled reports whether the array is neither own- nor inherited-disabled.
func (a *Array) IsEnabled() bool {
	return !a.effectiveDisabled()
}

// InheritedDisabled reports whether disabled state was inherited from an
// ancestor.
func (a *Array) InheritedDisabled() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.inheritedDisabled
}

// MarkAsTouched propagates the touched marker to all descendants.
func (a *Array) MarkAsTouched() {
	for _, child := range a.childrenList() {
		child.MarkAsTouched()
	}
}

// MarkAsUntouched propagates the untouched marker to all descendants.
func (a *Array) MarkAsUntouched() {
	for _, child := range a.childrenList() {
		child.MarkAsUntouched()
	}
}

// MarkAsPristine propagates the pristine marker to all descendants.
func (a *Array) MarkAsPristine() {
	for _, child := range a.childrenList() {
		child.MarkAsPristine()
	}
}

// SetDisabled sets only the array's own disabled flag and propagates its
// effective disabled state down to children.
func (a *Array) SetDisabled(disabled bool) {
	a.mu.Lock()
	if a.ownDisabled == disabled {
		a.mu.Unlock()
		return
	}
	a.ownDisabled = disabled
	a.mu.Unlock()

	a.applyDisabledToChildren()
	a.notifyStatus(a.Status())
	a.propagateUp()
}

// Reset resets all descendants.
func (a *Array) Reset() {
	for _, child := range a.childrenList() {
		child.Reset()
	}
}

// RawValue exports the array as []any in child order, omitting
// effectively-disabled children.
func (a *Array) RawValue() any {
	return a.rawValueSlice()
}

func (a *Array) rawValueSlice() []any {
	if a.effectiveDisabled() {
		return nil
	}
	children := a.childrenList()
	out := make([]any, 0, len(children))
	for _, child := range children {
		if !child.IsEnabled() {
			continue
		}
		out = append(out, child.RawValue())
	}
	return out
}

// Parent returns the array's parent container.
func (a *Array) Parent() FormNode {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.parent
}

func (a *Array) setParent(p FormNode) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.parent = p
}

// Children returns the array's child nodes in order.
func (a *Array) Children() []FormNode {
	return a.childrenList()
}

// Path returns the dotted path from the root.
func (a *Array) Path() string {
	return pathOf(a)
}

// Get resolves path relative to this array via the path resolver.
func (a *Array) Get(path string) (FormNode, bool) {
	return ResolvePath(a, path)
}

// onChildChanged recomputes the array's derived state and notifies its own
// subscribers before recursing up to the parent.
func (a *Array) onChildChanged() {
	a.notifyStatus(a.Status())
	a.notifyTouch(a.IsTouched())
	a.notifyValue(a.rawValueSlice())
	a.propagateUp()
}

func (a *Array) propagateUp() {
	if a.Parent() != nil {
		a.Parent().onChildChanged()
	}
}

func (a *Array) effectiveDisabled() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.ownDisabled || a.inheritedDisabled
}

func (a *Array) setInheritedDisabled(inh bool) {
	a.mu.Lock()
	if a.inheritedDisabled == inh {
		a.mu.Unlock()
		return
	}
	a.inheritedDisabled = inh
	a.mu.Unlock()

	a.applyDisabledToChildren()
	a.notifyStatus(a.Status())
}

func (a *Array) applyDisabledToChildren() {
	eff := a.effectiveDisabled()
	for _, child := range a.childrenList() {
		child.setInheritedDisabled(eff)
	}
}

func (a *Array) childrenList() []FormNode {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return append([]FormNode(nil), a.children...)
}

// OnValueChange subscribes to value changes.
func (a *Array) OnValueChange(fn func([]any)) state.Subscription {
	return a.valueStream.subscribe(fn)
}

// OnStatusChange subscribes to status changes.
func (a *Array) OnStatusChange(fn func(Status)) state.Subscription {
	return a.statusStream.subscribe(fn)
}

// OnTouchChange subscribes to touched-state changes.
func (a *Array) OnTouchChange(fn func(bool)) state.Subscription {
	return a.touchStream.subscribe(fn)
}

func (a *Array) notifyValue(v []any)   { a.valueStream.notify(v) }
func (a *Array) notifyStatus(s Status) { a.statusStream.notify(s) }
func (a *Array) notifyTouch(t bool)    { a.touchStream.notify(t) }
