package formengine

import "errors"

// Test-only validators. The engine ships no rules of its own (the public fform
// package and its rules subpackage do), so the engine tests define the small
// validators they need here.

// Required fails when value equals zero.
func Required[T comparable](zero T) ValidatorFunc[T] {
	return func(value T) error {
		if value == zero {
			return errors.New("value is required")
		}
		return nil
	}
}

// MinLength fails when the value is shorter than min characters.
func MinLength(min int) ValidatorFunc[string] {
	return func(value string) error {
		if len(value) < min {
			return errors.New("minimum length not met")
		}
		return nil
	}
}
