package rules

import (
	"slices"

	fform "github.com/zodimo/go-compose/compose/foundation/form"
)

// Unique fails when the slice contains duplicate values. It raises CodeUnique.
// Nil elements of a comparable type are compared normally; use UniqueBy for
// non-comparable elements or a custom key.
func Unique[T comparable]() fform.ValidatorFunc[[]T] {
	return func(value []T) error {
		seen := make(map[T]struct{}, len(value))
		for _, item := range value {
			if _, ok := seen[item]; ok {
				return fail(CodeUnique, "must not contain duplicate values")
			}
			seen[item] = struct{}{}
		}
		return nil
	}
}

// UniqueBy fails when two elements map to the same key. It raises CodeUnique.
// It works for any element type, using key to project a comparable identity.
func UniqueBy[T any, K comparable](key func(T) K) fform.ValidatorFunc[[]T] {
	return func(value []T) error {
		seen := make(map[K]struct{}, len(value))
		for _, item := range value {
			k := key(item)
			if _, ok := seen[k]; ok {
				return fail(CodeUnique, "must not contain duplicate values")
			}
			seen[k] = struct{}{}
		}
		return nil
	}
}

// SliceContainsAll fails when the slice does not contain every wanted element.
// It raises CodeSliceContains.
func SliceContainsAll[T comparable](wanted ...T) fform.ValidatorFunc[[]T] {
	return func(value []T) error {
		for _, w := range wanted {
			if !slices.Contains(value, w) {
				return fail(CodeSliceContains, "must contain %v", w)
			}
		}
		return nil
	}
}

// SliceContainsNone fails when the slice contains any forbidden element. It
// raises CodeSliceContains.
func SliceContainsNone[T comparable](forbidden ...T) fform.ValidatorFunc[[]T] {
	return func(value []T) error {
		for _, f := range forbidden {
			if slices.Contains(value, f) {
				return fail(CodeSliceContains, "must not contain %v", f)
			}
		}
		return nil
	}
}

// MapKeyPresent fails when the map does not contain every key. It raises
// CodeMapKeyPresent. Empty maps fail for any non-empty key set.
func MapKeyPresent[K comparable, V any](keys ...K) fform.ValidatorFunc[map[K]V] {
	return func(value map[K]V) error {
		for _, k := range keys {
			if _, ok := value[k]; !ok {
				return fail(CodeMapKeyPresent, "must contain key %v", k)
			}
		}
		return nil
	}
}

// MapValueRequired fails when any of the given keys maps to its zero value. It
// raises CodeRequired.
func MapValueRequired[K comparable, V comparable](zero V, keys ...K) fform.ValidatorFunc[map[K]V] {
	return func(value map[K]V) error {
		for _, k := range keys {
			if v, ok := value[k]; !ok || v == zero {
				return fail(CodeRequired, "value for key %v is required", k)
			}
		}
		return nil
	}
}
