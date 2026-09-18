package formengine

import "errors"

// ValidatorFunc validates a control value, returning a non-nil error when the
// value is invalid. Validators must not depend on the control's internal lock:
// they run on a snapshot of the value with no node lock held (design D5).
type ValidatorFunc[T any] func(value T) error

// Required returns a validator that fails when value equals zero.
func Required[T comparable](zero T) ValidatorFunc[T] {
	return func(value T) error {
		if value == zero {
			return errors.New("value is required")
		}
		return nil
	}
}

// MinLength returns a validator for string values that fails when the value is
// shorter than min characters.
func MinLength(min int) ValidatorFunc[string] {
	return func(value string) error {
		if len(value) < min {
			return errors.New("minimum length not met")
		}
		return nil
	}
}
