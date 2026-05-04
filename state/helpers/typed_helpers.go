package helpers

import (
	"fmt"
	"reflect"

	"github.com/zodimo/go-compose/state/core"
	"github.com/zodimo/go-compose/state/internal"
)

func Remember[T any](c core.SupportState, key string, initial func() T, options ...core.StateTypedOption[T]) (core.MutableValueTyped[T], error) {
	opts := core.StateTypedOptions[T]{
		Compare: func(t1, t2 T) bool {
			return reflect.DeepEqual(t1, t2)
		},
		OnForgotten: func() {},
	}
	for _, option := range options {
		option(&opts)
	}

	isNillable := isNillableType[T]()

	mv := c.State(
		key,
		func() any { return initial() },
		core.WithCompare(func(a, b any) bool {
			return opts.Compare(typeAssertUnsafe[T](a, isNillable), typeAssertUnsafe[T](b, isNillable))
		}),
		core.WithOnForgotten(opts.OnForgotten),
	)
	anyMv, ok := mv.(*internal.MutableValue)
	if !ok {
		return nil, fmt.Errorf("mutable value is not of type %T", internal.MutableValue{})
	}
	return internal.MutableValueToTyped[T](anyMv)
}

func RememberUnsafe[T any](c core.SupportState, key string, initial func() T, options ...core.StateTypedOption[T]) core.MutableValueTyped[T] {
	mv, err := Remember[T](c, key, initial, options...)
	if err != nil {
		panic(err)
	}
	return mv
}

func MustRemember[T any](c core.SupportState, key string, initial func() T, options ...core.StateTypedOption[T]) core.MutableValueTyped[T] {
	return RememberUnsafe[T](c, key, initial, options...)
}

// Deprecated: use Remember instead
func State[T any](c core.SupportState, key string, initial func() T, options ...core.StateTypedOption[T]) (core.MutableValueTyped[T], error) {
	opts := core.StateTypedOptions[T]{
		Compare: func(t1, t2 T) bool {
			return reflect.DeepEqual(t1, t2)
		},
		OnForgotten: func() {},
	}
	for _, option := range options {
		option(&opts)
	}

	isNillable := isNillableType[T]()

	mv := c.State(
		key,
		func() any { return initial() },
		core.WithCompare(func(a, b any) bool {
			return opts.Compare(typeAssertUnsafe[T](a, isNillable), typeAssertUnsafe[T](b, isNillable))
		}),
		core.WithOnForgotten(opts.OnForgotten),
	)
	anyMv, ok := mv.(*internal.MutableValue)
	if !ok {
		return nil, fmt.Errorf("mutable value is not of type %T", internal.MutableValue{})
	}
	return internal.MutableValueToTyped[T](anyMv)
}

// Deprecated: use RememberUnsafe instead
func StateUnsafe[T any](c core.SupportState, key string, initial func() T, options ...core.StateTypedOption[T]) core.MutableValueTyped[T] {
	mv, err := State[T](c, key, initial, options...)
	if err != nil {
		panic(err)
	}
	return mv
}

// Deprecated: use MustState instead
func MustState[T any](c core.SupportState, key string, initial func() T, options ...core.StateTypedOption[T]) core.MutableValueTyped[T] {
	return StateUnsafe[T](c, key, initial, options...)
}

func typeAssertUnsafe[T any](a any, nillable bool) T {
	val, err := typeAssert[T](a, nillable)
	if err != nil {
		panic(err)
	}
	return val
}

func typeAssert[T any](a any, nillable bool) (out T, err error) {
	if a == nil && nillable {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			if recoveredError, ok := r.(error); ok {
				err = recoveredError
			} else {
				err = fmt.Errorf("failed to assert type %T to %T: %v", a, out, r)
			}
		}
	}()

	out = a.(T)
	return
}

func isNillableType[T any]() bool {
	var zero T
	t := reflect.TypeOf(&zero).Elem()
	if t != nil {
		switch t.Kind() {
		case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Func, reflect.Chan, reflect.Interface:
			return true
		default:
			return false
		}
	}
	// if we cannot determine type then assume not nillable
	return false
}
