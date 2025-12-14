package state

type MutableValue interface {
	Get() any
	Set(value any)
}
