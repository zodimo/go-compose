package rules

import (
	"cmp"
	"slices"

	fform "github.com/zodimo/go-compose/compose/foundation/form"
)

// OneOf fails when the value is not in the allowed set. It raises CodeOneOf.
//
// A nil/empty set never fails — use it only when a set is genuinely configured,
// otherwise an unconfigured rule silently passes.
func OneOf[T comparable](allowed ...T) fform.ValidatorFunc[T] {
	return func(value T) error {
		if len(allowed) == 0 {
			return nil
		}
		if !slices.Contains(allowed, value) {
			return fail(CodeOneOf, "must be one of %v", allowed)
		}
		return nil
	}
}

// NoneOf fails when the value is in the forbidden set. It raises CodeNoneOf.
func NoneOf[T comparable](forbidden ...T) fform.ValidatorFunc[T] {
	return func(value T) error {
		if slices.Contains(forbidden, value) {
			return fail(CodeNoneOf, "must not be one of %v", forbidden)
		}
		return nil
	}
}

// EQ fails when the value does not equal want. It raises CodeEQ.
func EQ[T comparable](want T) fform.ValidatorFunc[T] {
	return func(value T) error {
		if value != want {
			return fail(CodeEQ, "must equal %v", want)
		}
		return nil
	}
}

// NEQ fails when the value equals forbidden. It raises CodeNEQ.
func NEQ[T comparable](forbidden T) fform.ValidatorFunc[T] {
	return func(value T) error {
		if value == forbidden {
			return fail(CodeNEQ, "must not equal %v", forbidden)
		}
		return nil
	}
}

// GT fails when the value is not greater than bound. It raises CodeGT.
func GT[T cmp.Ordered](bound T) fform.ValidatorFunc[T] {
	return func(value T) error {
		if value <= bound {
			return fail(CodeGT, "must be greater than %v", bound)
		}
		return nil
	}
}

// GTE fails when the value is less than bound. It raises CodeGTE.
func GTE[T cmp.Ordered](bound T) fform.ValidatorFunc[T] {
	return func(value T) error {
		if value < bound {
			return fail(CodeGTE, "must be greater than or equal to %v", bound)
		}
		return nil
	}
}

// LT fails when the value is not less than bound. It raises CodeLT.
func LT[T cmp.Ordered](bound T) fform.ValidatorFunc[T] {
	return func(value T) error {
		if value >= bound {
			return fail(CodeLT, "must be less than %v", bound)
		}
		return nil
	}
}

// LTE fails when the value is greater than bound. It raises CodeLTE.
func LTE[T cmp.Ordered](bound T) fform.ValidatorFunc[T] {
	return func(value T) error {
		if value > bound {
			return fail(CodeLTE, "must be less than or equal to %v", bound)
		}
		return nil
	}
}

// InRange fails when the value is outside [min, max] inclusive. It raises
// CodeGTE (below min) or CodeLTE (above max).
func InRange[T cmp.Ordered](min, max T) fform.ValidatorFunc[T] {
	return func(value T) error {
		if value < min {
			return fail(CodeGTE, "must be between %v and %v", min, max)
		}
		if value > max {
			return fail(CodeLTE, "must be between %v and %v", min, max)
		}
		return nil
	}
}
