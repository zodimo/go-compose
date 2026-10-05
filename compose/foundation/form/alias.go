package fform

import (
	formengine "github.com/zodimo/go-compose/compose/foundation/form/internal"
	"github.com/zodimo/go-compose/state"
)

// This file re-exports the form engine's public types and constructors so the
// public package exposes a stable API while the engine itself lives in
// internal/ (design D1).

// FormNode is the universal interface implemented by every node in a form tree.
type FormNode = formengine.FormNode

// Status is the validation/lifecycle state of a form node.
type Status = formengine.Status

// Group is a container node with named children.
type Group = formengine.Group

// Array is a container node with indexed children.
type Array = formengine.Array

// Control is a leaf node holding a single typed value.
type Control[T any] = formengine.Control[T]

// ValueStore is the read/write seam between the form engine and typed value storage.
type ValueStore[T any] = formengine.ValueStore[T]

// ValidatorFunc validates a control value, returning a non-nil error when invalid.
type ValidatorFunc[T any] = formengine.ValidatorFunc[T]

// GroupValidatorFunc validates a group as a whole, typically by reading its
// children's values. Group validators are the home for cross-field rules whose
// subject is the group (at-least-one-of, mutually-exclusive).
type GroupValidatorFunc = formengine.GroupValidatorFunc

// GroupOption configures a Group at construction.
type GroupOption = formengine.GroupOption

// WithGroupValidator attaches a group-level validator to a Group.
func WithGroupValidator(fn GroupValidatorFunc) GroupOption {
	return formengine.WithValidator(fn)
}

// NewGroup creates a group with the given named children, wiring each child's
// parent pointer. Options attach group-level validators.
func NewGroup(children map[string]FormNode, options ...GroupOption) *Group {
	return formengine.NewGroup(children, options...)
}

// NewArray creates an array with the given indexed children.
func NewArray(children ...FormNode) *Array {
	return formengine.NewArray(children...)
}

// NewControl creates a leaf control backed by store with the given initial value
// and validators.
func NewControl[T any](store ValueStore[T], initial T, validators ...ValidatorFunc[T]) *Control[T] {
	return formengine.NewControl(store, initial, validators...)
}

// NewPlainValueStore creates a mutex-guarded in-memory value store.
func NewPlainValueStore[T any](initial T) ValueStore[T] {
	return formengine.NewPlainValueStore(initial)
}

// NewMutableValueValueStore adapts a state.MutableValueTyped[T] to a ValueStore[T].
func NewMutableValueValueStore[T any](mv state.MutableValueTyped[T]) ValueStore[T] {
	return formengine.NewMutableValueValueStore(mv)
}

// Required returns a validator that fails when value equals zero. Failures carry
// CodeRequired.
func Required[T comparable](zero T) ValidatorFunc[T] {
	return func(value T) error {
		if value == zero {
			return NewValidationError(CodeRequired, "value is required")
		}
		return nil
	}
}

// MinLength returns a validator for string values that fails below min
// characters. Failures carry CodeMinLength.
func MinLength(min int) ValidatorFunc[string] {
	return func(value string) error {
		if len(value) < min {
			return NewValidationError(CodeMinLength, "minimum length is %d", min)
		}
		return nil
	}
}

// ResolvePath resolves a dotted path such as "a.b.c" or "items[0].name" against
// the given root node.
func ResolvePath(root FormNode, path string) (FormNode, bool) {
	return formengine.ResolvePath(root, path)
}

// Walk visits root and every descendant depth-first in root-to-leaf
// (pre-order) order.
func Walk(root FormNode, visit func(node FormNode)) {
	formengine.Walk(root, visit)
}

// WalkUp visits every descendant of node depth-first in leaf-to-root
// (post-order) order.
func WalkUp(node FormNode, visit func(node FormNode)) {
	formengine.WalkUp(node, visit)
}

// Status constants are re-declared explicitly because constants cannot be
// type-aliased.
const (
	StatusValid    = formengine.StatusValid
	StatusInvalid  = formengine.StatusInvalid
	StatusPending  = formengine.StatusPending
	StatusDisabled = formengine.StatusDisabled
)
