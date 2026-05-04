package internal

import "reflect"

// ReferentialEqualityPolicy compares values using Go's == operator.
// Only works with comparable types.
type ReferentialEqualityPolicy[T comparable] struct{}

func (ReferentialEqualityPolicy[T]) Equivalent(a, b T) bool {
	return a == b
}

func (ReferentialEqualityPolicy[T]) Merge(previous, current, applied T) (T, bool) {
	var zero T
	return zero, false
}

// StructuralEqualityPolicy compares values using reflect.DeepEqual.
type StructuralEqualityPolicy[T any] struct{}

func (StructuralEqualityPolicy[T]) Equivalent(a, b T) bool {
	return reflect.DeepEqual(a, b)
}

func (StructuralEqualityPolicy[T]) Merge(previous, current, applied T) (T, bool) {
	var zero T
	return zero, false
}

// NeverEqualPolicy treats all values as different.
type NeverEqualPolicy[T any] struct{}

func (NeverEqualPolicy[T]) Equivalent(a, b T) bool {
	return false
}

func (NeverEqualPolicy[T]) Merge(previous, current, applied T) (T, bool) {
	var zero T
	return zero, false
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

func NewFunctionalPolicy[T any](
	equivalent func(a, b T) bool,
	merge func(previous, current, applied T) (T, bool),
) *FunctionalPolicy[T] {
	return &FunctionalPolicy[T]{
		equivalent: equivalent,
		merge:      merge,
	}
}
