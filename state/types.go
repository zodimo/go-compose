package state

type StateOption func(*StateOptions)

type StateOptions struct {
	Compare func(any, any) bool
}

type StateTypedOptions[T any] struct {
	Compare func(T, T) bool
}

type StateTypedOption[T any] func(*StateTypedOptions[T])

func WithTypedCompare[T any](compare func(T, T) bool) StateTypedOption[T] {
	if compare == nil {
		panic("compare cannot be nil")
	}
	return func(o *StateTypedOptions[T]) {
		o.Compare = compare
	}
}

type SupportState interface {
	Remember(key string, initial func() any, options ...StateOption) MutableValue // persistent state

	// deprecate this one to get closer to jetpack compose
	// alias for Remember
	// Deprecated: use Remember instead
	State(key string, initial func() any, options ...StateOption) MutableValue // persistent state
}

type Value interface {
	Get() any
	Subscribe(callback func()) Subscription
}

type ValueTyped[T any] interface {
	Get() T
	Subscribe(callback func()) Subscription
}

type MutableValue interface {
	Value
	Set(value any)
	CompareAndSet(expect, update any) bool
	Update(func(any) any)
	UpdateAndGet(func(any) any) any
	GetAndUpdate(func(any) any) any
}

type MutableValueTyped[T any] interface {
	ValueTyped[T]
	Set(value T)
	CompareAndSet(expect, update T) bool
	Update(func(T) T)
	UpdateAndGet(func(T) T) T
	GetAndUpdate(func(T) T) T

	Unwrap() MutableValue
}
