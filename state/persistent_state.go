package state

type PersistentState interface {
	StateChangeNotifier
	GetState(key string, initial func() any, options ...StateOption) MutableValue
	// Deprecated: use Subscribe instead
	SetOnStateChange(callback func())
}
