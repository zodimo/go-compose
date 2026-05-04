package state

import (
	"github.com/zodimo/go-compose/state/core"
)

type StateOption = core.StateOption

type StateOptions = core.StateOptions

var WithCompare = core.WithCompare

var WithOnForgotten = core.WithOnForgotten

type StateTypedOption[T any] = core.StateTypedOption[T]

type StateTypedOptions[T any] = core.StateTypedOptions[T]

func WithTypedCompare[T any](compare func(T, T) bool) StateTypedOption[T] {
	return core.WithTypedCompare(compare)
}

func WithTypedOnForgotten[T any](onForgotten func()) StateTypedOption[T] {
	return core.WithTypedOnForgotten[T](onForgotten)
}

type SupportState = core.SupportState

type Value = core.Value

type ValueTyped[T any] = core.ValueTyped[T]

type MutableValue = core.MutableValue

type MutableValueTyped[T any] = core.MutableValueTyped[T]
