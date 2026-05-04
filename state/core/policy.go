package core

import "reflect"

// MutationPolicy controls change detection and conflict resolution.
// This matches Kotlin's SnapshotMutationPolicy interface from androidx.compose.runtime.
type MutationPolicy[T any] interface {
	// Equivalent returns true if a and b should be considered the same value.
	// If Equivalent returns true, setting a value will not trigger notifications.
	Equivalent(a, b T) bool

	// Merge resolves conflicts between concurrent modifications.
	// Returns (merged, true) if the conflict was resolved, or (zero, false) if unresolvable.
	// Most policies return (zero, false) to indicate no merge support.
	Merge(previous, current, applied T) (T, bool)
}

// --- Built-in Policies ---

// referentialEqualityPolicy compares values using Go's == operator.
// Only works with comparable types.
type referentialEqualityPolicy[T comparable] struct{}

func (referentialEqualityPolicy[T]) Equivalent(a, b T) bool {
	return a == b
}

func (referentialEqualityPolicy[T]) Merge(previous, current, applied T) (T, bool) {
	var zero T
	return zero, false
}

// ReferentialEqualityPolicy returns a policy that compares values using Go's == operator.
// This is suitable for primitive types and types where identity comparison is sufficient.
func ReferentialEqualityPolicy[T comparable]() MutationPolicy[T] {
	return referentialEqualityPolicy[T]{}
}

// structuralEqualityPolicy compares values using reflect.DeepEqual.
type structuralEqualityPolicy[T any] struct{}

func (structuralEqualityPolicy[T]) Equivalent(a, b T) bool {
	return reflect.DeepEqual(a, b)
}

func (structuralEqualityPolicy[T]) Merge(previous, current, applied T) (T, bool) {
	var zero T
	return zero, false
}

// StructuralEqualityPolicy returns a policy that compares values using reflect.DeepEqual.
// This is the default policy and is suitable for complex structs and slices.
func StructuralEqualityPolicy[T any]() MutationPolicy[T] {
	return structuralEqualityPolicy[T]{}
}

// neverEqualPolicy treats all values as different.
type neverEqualPolicy[T any] struct{}

func (neverEqualPolicy[T]) Equivalent(a, b T) bool {
	return false
}

func (neverEqualPolicy[T]) Merge(previous, current, applied T) (T, bool) {
	var zero T
	return zero, false
}

// NeverEqualPolicy returns a policy that always treats values as different.
// Setting any value will always trigger notifications, even if the value is the same.
func NeverEqualPolicy[T any]() MutationPolicy[T] {
	return neverEqualPolicy[T]{}
}

// --- Functional Policy for Custom Comparisons ---

// FunctionalPolicy creates a MutationPolicy from simple functions.
type FunctionalPolicy[T any] struct {
	equivalent func(a, b T) bool
	merge      func(previous, current, applied T) (T, bool)
}

func (p FunctionalPolicy[T]) Equivalent(a, b T) bool {
	return p.equivalent(a, b)
}

func (p FunctionalPolicy[T]) Merge(previous, current, applied T) (T, bool) {
	if p.merge == nil {
		var zero T
		return zero, false
	}
	return p.merge(previous, current, applied)
}

// NewMutationPolicy creates a MutationPolicy from an equivalence function.
// The merge function is optional (pass nil to disable merging).
func NewMutationPolicy[T any](equivalent func(a, b T) bool, merge func(previous, current, applied T) (T, bool)) MutationPolicy[T] {
	return FunctionalPolicy[T]{
		equivalent: equivalent,
		merge:      merge,
	}
}
