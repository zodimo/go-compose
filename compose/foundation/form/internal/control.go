package formengine

import (
	"errors"
	"sync"

	"github.com/zodimo/go-compose/state"
)

// Control is a leaf node holding a single typed value.
type Control[T any] struct {
	mu                sync.RWMutex
	parent            FormNode
	valueStore        ValueStore[T]
	initialValue      T
	validators        []ValidatorFunc[T]
	errors            []error
	ownDisabled       bool
	inheritedDisabled bool
	touched           bool
	dirty             bool

	valueStream  *typedStream[T]
	statusStream *typedStream[Status]
	touchStream  *typedStream[bool]
}

var _ FormNode = (*Control[int])(nil)

// NewControl creates a leaf control backed by store, with the given initial
// value and validators. store may be a plain holder (pure Go) or a
// MutableValueTyped adapter (UI).
func NewControl[T any](store ValueStore[T], initial T, validators ...ValidatorFunc[T]) *Control[T] {
	return &Control[T]{
		valueStore:   store,
		initialValue: initial,
		validators:   validators,
		valueStream:  newTypedStream(initial),
		statusStream: newTypedStream(StatusValid),
		touchStream:  newTypedStream(false),
	}
}

// Value returns the current typed value read from the store.
func (c *Control[T]) Value() T {
	return c.valueStore.Get()
}

// Set stores v, marks the control dirty, re-runs validation on a snapshot, and
// notifies listeners before propagating up. It is a no-op when the control is
// effectively disabled.
func (c *Control[T]) Set(v T) {
	if c.effectiveDisabled() {
		return
	}

	c.mu.Lock()
	c.valueStore.Set(v)
	c.dirty = true
	c.mu.Unlock()

	c.runValidation()

	value := c.Value()
	c.notifyValue(value)
	c.notifyStatus(c.Status())
	c.notifyTouch(c.IsTouched())

	c.propagateUp()
}

// runValidation runs every validator on a snapshot of the value and stores the
// results under the node lock. Validators never run while the node lock is held
// (design D5), which keeps reentrant validators (that call back into the
// control's read methods) deadlock-free.
func (c *Control[T]) runValidation() {
	snapshot := c.Value()
	var errs []error
	for _, vf := range c.validators {
		if err := vf(snapshot); err != nil {
			errs = append(errs, err)
		}
	}
	c.mu.Lock()
	c.errors = errs
	c.mu.Unlock()
}

// Validate runs validation and reports whether the control is valid. A disabled
// control always validates true and clears its errors.
func (c *Control[T]) Validate() bool {
	if c.effectiveDisabled() {
		c.mu.Lock()
		c.errors = nil
		c.mu.Unlock()
		return true
	}
	c.runValidation()
	c.mu.RLock()
	n := len(c.errors)
	c.mu.RUnlock()
	return n == 0
}

// Errors returns the control's own path keyed to the errors.Join aggregation
// of all validator failures, or nil when valid or disabled. The joined message
// contains every validator error so the UI can surface all of them at once.
func (c *Control[T]) Errors() map[string]string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.ownDisabled || c.inheritedDisabled || len(c.errors) == 0 {
		return nil
	}
	return map[string]string{c.Path(): errors.Join(c.errors...).Error()}
}

// HasErrors reports whether the control currently has any validation errors.
func (c *Control[T]) HasErrors() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return !c.ownDisabled && !c.inheritedDisabled && len(c.errors) > 0
}

// ValidationErrors returns a copy of the control's raw validator failures, or
// nil when valid or disabled. Unlike Errors it does not join or stringify, so
// callers can inspect structured error values (for example codes via errors.As).
func (c *Control[T]) ValidationErrors() []error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.ownDisabled || c.inheritedDisabled || len(c.errors) == 0 {
		return nil
	}
	return append([]error(nil), c.errors...)
}

// Status returns DISABLED when effectively disabled, INVALID when errors exist,
// otherwise VALID.
func (c *Control[T]) Status() Status {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.ownDisabled || c.inheritedDisabled {
		return StatusDisabled
	}
	if len(c.errors) > 0 {
		return StatusInvalid
	}
	return StatusValid
}

// RawValue returns the control's value, or nil when effectively disabled.
func (c *Control[T]) RawValue() any {
	if c.effectiveDisabled() {
		return nil
	}
	return c.Value()
}

// Reset restores the initial value, clears errors, and marks the control
// pristine and untouched, then notifies and propagates.
func (c *Control[T]) Reset() {
	c.mu.Lock()
	c.valueStore.Set(c.initialValue)
	c.dirty = false
	c.touched = false
	c.errors = nil
	c.mu.Unlock()

	c.notifyValue(c.initialValue)
	c.notifyStatus(c.Status())
	c.notifyTouch(false)

	c.propagateUp()
}

// MarkAsTouched marks the control touched and notifies + propagates.
func (c *Control[T]) MarkAsTouched() {
	c.mu.Lock()
	if c.touched {
		c.mu.Unlock()
		return
	}
	c.touched = true
	c.mu.Unlock()

	c.notifyTouch(true)
	c.propagateUp()
}

// MarkAsUntouched clears the touched flag and notifies + propagates.
func (c *Control[T]) MarkAsUntouched() {
	c.mu.Lock()
	if !c.touched {
		c.mu.Unlock()
		return
	}
	c.touched = false
	c.mu.Unlock()

	c.notifyTouch(false)
	c.propagateUp()
}

// MarkAsPristine clears the dirty flag and propagates (ancestors recompute
// their derived IsDirty on read).
func (c *Control[T]) MarkAsPristine() {
	c.mu.Lock()
	if !c.dirty {
		c.mu.Unlock()
		return
	}
	c.dirty = false
	c.mu.Unlock()

	c.propagateUp()
}

// SetDisabled sets only the control's own disabled flag (design D4); inherited
// disabled state is managed separately via setInheritedDisabled.
func (c *Control[T]) SetDisabled(disabled bool) {
	c.mu.Lock()
	if c.ownDisabled == disabled {
		c.mu.Unlock()
		return
	}
	c.ownDisabled = disabled
	c.mu.Unlock()

	c.notifyStatus(c.Status())
	c.propagateUp()
}

// IsTouched reports whether the control has been touched.
func (c *Control[T]) IsTouched() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.touched
}

// IsDirty reports whether the control's value differs from its initial value.
func (c *Control[T]) IsDirty() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.dirty
}

// IsPristine is the negation of IsDirty.
func (c *Control[T]) IsPristine() bool {
	return !c.IsDirty()
}

// IsEnabled reports whether the control is neither own-disabled nor
// inherited-disabled.
func (c *Control[T]) IsEnabled() bool {
	return !c.effectiveDisabled()
}

// InheritedDisabled reports whether disabled state was inherited from an
// ancestor.
func (c *Control[T]) InheritedDisabled() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.inheritedDisabled
}

func (c *Control[T]) effectiveDisabled() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.ownDisabled || c.inheritedDisabled
}

// setInheritedDisabled wires inherited disabled state (set by an ancestor).
// It notifies the control's own status subscribers but does not propagate up:
// the ancestor that initiated the change already propagates up the chain.
func (c *Control[T]) setInheritedDisabled(inh bool) {
	c.mu.Lock()
	if c.inheritedDisabled == inh {
		c.mu.Unlock()
		return
	}
	c.inheritedDisabled = inh
	c.mu.Unlock()

	c.notifyStatus(c.Status())
}

// Parent returns the control's parent container.
func (c *Control[T]) Parent() FormNode {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.parent
}

func (c *Control[T]) setParent(p FormNode) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.parent = p
}

// Children returns nil: a control is a leaf.
func (c *Control[T]) Children() []FormNode {
	return nil
}

// Path returns the dotted path from the root.
func (c *Control[T]) Path() string {
	return pathOf(c)
}

// Get resolves a path relative to this leaf: the empty path returns the control
// itself; any non-empty path returns a miss.
func (c *Control[T]) Get(path string) (FormNode, bool) {
	if path == "" {
		return c, true
	}
	return nil, false
}

// onChildChanged is a no-op: a leaf has no children whose changes would alter
// its derived state.
func (c *Control[T]) onChildChanged() {}

// propagateUp notifies the parent that this node's derived state changed.
func (c *Control[T]) propagateUp() {
	if c.Parent() != nil {
		c.Parent().onChildChanged()
	}
}

// OnValueChange subscribes to value changes.
func (c *Control[T]) OnValueChange(fn func(T)) state.Subscription {
	return c.valueStream.subscribe(fn)
}

// OnStatusChange subscribes to status changes.
func (c *Control[T]) OnStatusChange(fn func(Status)) state.Subscription {
	return c.statusStream.subscribe(fn)
}

// OnTouchChange subscribes to touched-state changes.
func (c *Control[T]) OnTouchChange(fn func(bool)) state.Subscription {
	return c.touchStream.subscribe(fn)
}

func (c *Control[T]) notifyValue(v T)       { c.valueStream.notify(v) }
func (c *Control[T]) notifyStatus(s Status) { c.statusStream.notify(s) }
func (c *Control[T]) notifyTouch(t bool)    { c.touchStream.notify(t) }
