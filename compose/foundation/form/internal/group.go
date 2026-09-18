package formengine

import (
	"sync"

	"github.com/zodimo/go-compose/state"
)

// Group is a container node with named children.
type Group struct {
	mu                sync.RWMutex
	parent            FormNode
	children          map[string]FormNode
	ownDisabled       bool
	inheritedDisabled bool

	valueStream  *typedStream[map[string]any]
	statusStream *typedStream[Status]
	touchStream  *typedStream[bool]
}

var _ FormNode = (*Group)(nil)

// NewGroup creates a group with the given named children and wires each child's
// parent pointer to the group.
func NewGroup(children map[string]FormNode) *Group {
	g := &Group{
		children:     children,
		valueStream:  newTypedStream(map[string]any{}),
		statusStream: newTypedStream(StatusValid),
		touchStream:  newTypedStream(false),
	}
	for _, child := range children {
		child.setParent(g)
	}
	return g
}

// Status aggregates child statuses: DISABLED if the group is itself effectively
// disabled; else INVALID if any child is INVALID; else PENDING if any child is
// PENDING; else DISABLED if all children are disabled; else VALID. An empty
// group is VALID.
func (g *Group) Status() Status {
	if g.effectiveDisabled() {
		return StatusDisabled
	}
	children := g.childrenList()
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
func (g *Group) Validate() bool {
	if g.effectiveDisabled() {
		return true
	}
	valid := true
	for _, child := range g.childrenList() {
		if !child.Validate() {
			valid = false
		}
	}
	return valid
}

// Errors merges all descendant control errors keyed by their dotted paths.
func (g *Group) Errors() map[string]string {
	if g.effectiveDisabled() {
		return nil
	}
	result := make(map[string]string)
	for _, child := range g.childrenList() {
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
func (g *Group) IsTouched() bool {
	for _, child := range g.childrenList() {
		if child.IsTouched() {
			return true
		}
	}
	return false
}

// IsDirty reports whether any descendant is dirty.
func (g *Group) IsDirty() bool {
	for _, child := range g.childrenList() {
		if child.IsDirty() {
			return true
		}
	}
	return false
}

// IsPristine is the negation of IsDirty.
func (g *Group) IsPristine() bool {
	return !g.IsDirty()
}

// IsEnabled reports whether the group is neither own- nor inherited-disabled.
func (g *Group) IsEnabled() bool {
	return !g.effectiveDisabled()
}

// InheritedDisabled reports whether disabled state was inherited from an
// ancestor.
func (g *Group) InheritedDisabled() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.inheritedDisabled
}

// MarkAsTouched propagates the touched marker to all descendants.
func (g *Group) MarkAsTouched() {
	for _, child := range g.childrenList() {
		child.MarkAsTouched()
	}
}

// MarkAsUntouched propagates the untouched marker to all descendants.
func (g *Group) MarkAsUntouched() {
	for _, child := range g.childrenList() {
		child.MarkAsUntouched()
	}
}

// MarkAsPristine propagates the pristine marker to all descendants.
func (g *Group) MarkAsPristine() {
	for _, child := range g.childrenList() {
		child.MarkAsPristine()
	}
}

// SetDisabled sets only the group's own disabled flag (design D4), then
// propagates the group's effective disabled state down to children. Children
// keep their own flag, so re-enabling restores their prior state.
func (g *Group) SetDisabled(disabled bool) {
	g.mu.Lock()
	if g.ownDisabled == disabled {
		g.mu.Unlock()
		return
	}
	g.ownDisabled = disabled
	g.mu.Unlock()

	g.applyDisabledToChildren()
	g.notifyStatus(g.Status())
	g.propagateUp()
}

// Reset resets all descendants.
func (g *Group) Reset() {
	for _, child := range g.childrenList() {
		child.Reset()
	}
}

// RawValue exports the group as map[string]any keyed by child name, omitting
// effectively-disabled children.
func (g *Group) RawValue() any {
	return g.rawValueMap()
}

func (g *Group) rawValueMap() map[string]any {
	if g.effectiveDisabled() {
		return nil
	}
	children := g.childrenMap()
	m := make(map[string]any, len(children))
	for name, child := range children {
		if !child.IsEnabled() {
			continue
		}
		m[name] = child.RawValue()
	}
	return m
}

// Parent returns the group's parent container.
func (g *Group) Parent() FormNode {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.parent
}

func (g *Group) setParent(p FormNode) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.parent = p
}

// Children returns the group's child nodes (map values, unordered).
func (g *Group) Children() []FormNode {
	return g.childrenList()
}

// Path returns the dotted path from the root.
func (g *Group) Path() string {
	return pathOf(g)
}

// Get resolves path relative to this group via the path resolver.
func (g *Group) Get(path string) (FormNode, bool) {
	return ResolvePath(g, path)
}

// child returns the named child under the group lock.
func (g *Group) child(name string) (FormNode, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	n, ok := g.children[name]
	return n, ok
}

// onChildChanged recomputes the group's derived state and notifies its own
// subscribers before recursing up to the parent (design D3).
func (g *Group) onChildChanged() {
	g.notifyStatus(g.Status())
	g.notifyTouch(g.IsTouched())
	g.notifyValue(g.rawValueMap())
	g.propagateUp()
}

func (g *Group) propagateUp() {
	if g.Parent() != nil {
		g.Parent().onChildChanged()
	}
}

func (g *Group) effectiveDisabled() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.ownDisabled || g.inheritedDisabled
}

// setInheritedDisabled wires inherited disabled state from an ancestor and
// propagates the new effective state down. It notifies own status subscribers
// but does not propagate up (the initiating ancestor handles that).
func (g *Group) setInheritedDisabled(inh bool) {
	g.mu.Lock()
	if g.inheritedDisabled == inh {
		g.mu.Unlock()
		return
	}
	g.inheritedDisabled = inh
	g.mu.Unlock()

	g.applyDisabledToChildren()
	g.notifyStatus(g.Status())
}

func (g *Group) applyDisabledToChildren() {
	eff := g.effectiveDisabled()
	for _, child := range g.childrenList() {
		child.setInheritedDisabled(eff)
	}
}

func (g *Group) childrenList() []FormNode {
	g.mu.RLock()
	defer g.mu.RUnlock()
	list := make([]FormNode, 0, len(g.children))
	for _, c := range g.children {
		list = append(list, c)
	}
	return list
}

func (g *Group) childrenMap() map[string]FormNode {
	g.mu.RLock()
	defer g.mu.RUnlock()
	m := make(map[string]FormNode, len(g.children))
	for k, v := range g.children {
		m[k] = v
	}
	return m
}

// OnValueChange subscribes to value changes.
func (g *Group) OnValueChange(fn func(map[string]any)) state.Subscription {
	return g.valueStream.subscribe(fn)
}

// OnStatusChange subscribes to status changes.
func (g *Group) OnStatusChange(fn func(Status)) state.Subscription {
	return g.statusStream.subscribe(fn)
}

// OnTouchChange subscribes to touched-state changes.
func (g *Group) OnTouchChange(fn func(bool)) state.Subscription {
	return g.touchStream.subscribe(fn)
}

func (g *Group) notifyValue(v map[string]any) { g.valueStream.notify(v) }
func (g *Group) notifyStatus(s Status)        { g.statusStream.notify(s) }
func (g *Group) notifyTouch(t bool)           { g.touchStream.notify(t) }
