package state

import (
	"fmt"
)

func Remember[T any](c SupportState, key string, calc func() T) (T, error) {
	anyCalc := func() any { return calc() }
	anyValue := c.Remember(key, anyCalc)

	tValue, ok := anyValue.(T)
	if !ok {
		var zero T
		return zero, fmt.Errorf("value is not of type %T", zero)
	}
	return tValue, nil
}

func RememberUnsafe[T any](c SupportState, key string, calc func() T) T {
	anyCalc := func() any { return calc() }
	anyValue := c.Remember(key, anyCalc)

	tValue, ok := anyValue.(T)
	if !ok {
		var zero T
		panic(fmt.Errorf("value is not of type %T", zero))
	}
	return tValue
}

func MustRemember[T any](c SupportState, key string, calc func() T) T {
	return RememberUnsafe[T](c, key, calc)
}
