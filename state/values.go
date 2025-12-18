package state

type Value interface {
	Get() any
}

type TypedValue[T any] interface {
	Get() T
}

type MutableValue interface {
	Value
	Set(value any)
}

type TypedMutableValue[T any] interface {
	TypedValue[T]
	Set(value T)

	Unwrap() MutableValue
}
