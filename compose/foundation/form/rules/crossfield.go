package rules

import (
	"cmp"
	"fmt"
	"reflect"

	fform "github.com/zodimo/go-compose/compose/foundation/form"
)

// EqualTo fails when the value does not equal the value produced by other. It
// raises CodeCompareFields. Use it for confirmation fields (password, email).
//
// other is evaluated on each validation, so it reflects the sibling's current
// value at validation time:
//
//	fform.NewControl(store, "", rules.EqualTo(func() string { return password.Value() }))
func EqualTo[T comparable](other func() T) fform.ValidatorFunc[T] {
	return func(value T) error {
		want := other()
		if value != want {
			return fail(CodeCompareFields, "must match the other field")
		}
		return nil
	}
}

// NotEqualTo fails when the value equals the value produced by other. It raises
// CodeCompareFields.
func NotEqualTo[T comparable](other func() T) fform.ValidatorFunc[T] {
	return func(value T) error {
		if value == other() {
			return fail(CodeCompareFields, "must differ from the other field")
		}
		return nil
	}
}

// GreaterThanField fails when the value is not greater than the value produced
// by other. It raises CodeCompareFields. Use it for ordered pairs (end > start).
func GreaterThanField[T cmp.Ordered](other func() T) fform.ValidatorFunc[T] {
	return func(value T) error {
		if value <= other() {
			return fail(CodeCompareFields, "must be greater than the other field")
		}
		return nil
	}
}

// LessThanField fails when the value is not less than the value produced by
// other. It raises CodeCompareFields.
func LessThanField[T cmp.Ordered](other func() T) fform.ValidatorFunc[T] {
	return func(value T) error {
		if value >= other() {
			return fail(CodeCompareFields, "must be less than the other field")
		}
		return nil
	}
}

// RequiredWhen fails when the value is zero while cond is true. It raises
// CodeRequired. It expresses conditional-required fields:
//
//	rules.RequiredWhen("", func() bool { return hasBilling.Value() })
func RequiredWhen[T comparable](zero T, cond func() bool) fform.ValidatorFunc[T] {
	return func(value T) error {
		if cond() && value == zero {
			return fail(CodeRequired, "value is required")
		}
		return nil
	}
}

// ForbiddenWhen fails when the value is non-zero while cond is true. It raises
// CodeCrossField.
func ForbiddenWhen[T comparable](zero T, cond func() bool) fform.ValidatorFunc[T] {
	return func(value T) error {
		if cond() && value != zero {
			return fail(CodeCrossField, "value must not be set")
		}
		return nil
	}
}

// AtLeastOneSet fails when every value produced by the getters is zero. It
// raises CodeCrossField. It expresses "provide at least one of" groups:
//
//	rules.AtLeastOneSet(func() string { return email.Value() }, func() string { return phone.Value() })
//
// All getters must produce the same type.
func AtLeastOneSet[T comparable](getters ...func() T) fform.ValidatorFunc[T] {
	return func(_ T) error {
		var zero T
		for _, get := range getters {
			if get() != zero {
				return nil
			}
		}
		return fail(CodeCrossField, "at least one of the related fields must be set")
	}
}

// MutuallyExclusive fails when more than one of the values produced by the
// getters is non-zero. It raises CodeCrossField.
func MutuallyExclusive[T comparable](getters ...func() T) fform.ValidatorFunc[T] {
	return func(_ T) error {
		var zero T
		set := 0
		for _, get := range getters {
			if get() != zero {
				set++
			}
		}
		if set > 1 {
			return fail(CodeCrossField, "only one of the related fields may be set")
		}
		return nil
	}
}

// DeepEqual fails when the value is not reflect.DeepEqual to want. It raises
// CodeEQ. Use it for non-comparable values (slices, maps, structs).
func DeepEqual[T any](want T) fform.ValidatorFunc[T] {
	return func(value T) error {
		if !reflect.DeepEqual(value, want) {
			return fail(CodeEQ, "must equal %v", want)
		}
		return nil
	}
}

// Custom builds a rule from an arbitrary predicate, attaching code. It is the
// escape hatch for one-off cross-field logic that the named rules do not cover.
func Custom[T any](code Code, ok func(T) bool, format string, args ...any) fform.ValidatorFunc[T] {
	message := fmt.Sprintf(format, args...)
	return func(value T) error {
		if !ok(value) {
			return fail(code, "%s", message)
		}
		return nil
	}
}
