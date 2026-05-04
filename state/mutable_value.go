package state

import (
	"reflect"

	"github.com/zodimo/go-compose/lifecycle"
	"github.com/zodimo/go-compose/state/core"
	"github.com/zodimo/go-compose/state/internal"
)

var _ lifecycle.RememberObserver = (*internal.MutableValue)(nil)

type MutableValueOptions struct {
	OnForgotten func()
}

type MutableValueOption func(*MutableValueOptions)

func MutableValueWithOnForgotten(onForgotten func()) MutableValueOption {
	if onForgotten == nil {
		panic("onForgotten cannot be nil")
	}
	return func(o *MutableValueOptions) {
		o.OnForgotten = onForgotten
	}
}

func NewMutableValue(initial any, changeNotifier func(any), compare func(any, any) bool, options ...MutableValueOption) MutableValue {
	opts := MutableValueOptions{
		OnForgotten: func() {},
	}
	for _, option := range options {
		if option != nil {
			option(&opts)
		}
	}

	if changeNotifier == nil {
		changeNotifier = func(any) {}
	}

	if compare == nil {
		compare = func(v1, v2 any) bool {
			return reflect.DeepEqual(v1, v2)
		}
	}

	return internal.NewMutableValue(
		initial,
		changeNotifier,
		compare,
		core.NewSubscriptionManager(),
		opts.OnForgotten,
	)
}
