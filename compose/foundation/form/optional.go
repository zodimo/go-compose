package fform

import (
	"errors"

	formengine "github.com/zodimo/go-compose/compose/foundation/form/internal"
	"github.com/zodimo/go-maybe"
)

// Optional is a value that may be absent. It is the form engine's "empty is
// representable" primitive: use Control[Optional[T]] when a typed field must
// distinguish "no value entered" from a legitimate zero value (for example an
// empty number field versus 0).
//
// Optional is an alias for github.com/zodimo/go-maybe.Maybe, so it interoperates
// with the rest of the repository (json omitzero, FromPtr, Match, Map, ...).
type Optional[T any] = maybe.Maybe[T]

// Some returns an Optional holding value.
func Some[T any](value T) Optional[T] {
	return maybe.Some(value)
}

// None returns an empty Optional.
func None[T any]() Optional[T] {
	return maybe.None[T]()
}

// OptionalOf returns Some(value) when set is true, otherwise None. It is the
// natural adapter for a (value, hasValue) pair.
func OptionalOf[T any](value T, set bool) Optional[T] {
	if set {
		return maybe.Some(value)
	}
	return maybe.None[T]()
}

// RequiredOptional returns a validator that fails when the value is empty
// (None). It is the Optional-aware counterpart of Required: Required(None[T]())
// also works because Optional[T] is comparable for comparable T, but
// RequiredOptional reads more clearly at call sites and validates the intent
// rather than comparing against a sentinel.
func RequiredOptional[T any]() ValidatorFunc[Optional[T]] {
	return func(value Optional[T]) error {
		if value.IsNone() {
			return errors.New("value is required")
		}
		return nil
	}
}

// NewOptionalControl creates a leaf control holding an Optional[T], initialized
// to None (empty) unless a validator-needing initial value is supplied through
// initial. It is a convenience over NewControl for the common "starts empty"
// case.
func NewOptionalControl[T any](store ValueStore[Optional[T]]) *Control[Optional[T]] {
	return formengine.NewControl(store, maybe.None[T]())
}

// OptionalValue unwraps an Optional control's value, returning the inner value
// and whether it is present.
func OptionalValue[T any](control *Control[Optional[T]]) (T, bool) {
	m := control.Value()
	if m.IsNone() {
		var zero T
		return zero, false
	}
	return m.UnwrapUnsafe(), true
}
