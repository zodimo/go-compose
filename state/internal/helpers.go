package internal

import (
	"fmt"
	"reflect"
)

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
