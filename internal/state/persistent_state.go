package state

type PersistentState interface {
	GetState(key string, initial func() any) MutableValue
	SetOnStateChange(callback func())
}
