package state

import (
	"github.com/zodimo/go-compose/state/core"
	"github.com/zodimo/go-compose/state/internal"
)

// MutationPolicy controls change detection and conflict resolution.
// This matches Kotlin's SnapshotMutationPolicy interface from androidx.compose.runtime.
type MutationPolicy[T any] = core.MutationPolicy[T]

// --- Built-in Policies ---

// ReferentialEqualityPolicy returns a policy that compares values using Go's == operator.
// This is suitable for primitive types and types where identity comparison is sufficient.
func ReferentialEqualityPolicy[T comparable]() MutationPolicy[T] {
	return internal.ReferentialEqualityPolicy[T]{}
}

// StructuralEqualityPolicy returns a policy that compares values using reflect.DeepEqual.
// This is the default policy and is suitable for complex structs and slices.
func StructuralEqualityPolicy[T any]() MutationPolicy[T] {
	return internal.StructuralEqualityPolicy[T]{}
}

// NeverEqualPolicy returns a policy that always treats values as different.
// Setting any value will always trigger notifications, even if the value is the same.
func NeverEqualPolicy[T any]() MutationPolicy[T] {
	return internal.NeverEqualPolicy[T]{}
}

// NewMutationPolicy creates a MutationPolicy from an equivalence function.
// The merge function is optional (pass nil to disable merging).
func NewMutationPolicy[T any](equivalent func(a, b T) bool, merge func(previous, current, applied T) (T, bool)) MutationPolicy[T] {
	return internal.NewFunctionalPolicy[T](equivalent, merge)
}
