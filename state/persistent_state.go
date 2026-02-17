package state

import "github.com/zodimo/go-compose/lifecycle"

type PersistentState interface {
	StateChangeNotifier
	lifecycle.FrameLifecycleAware
	GetState(key string, initial func() any, options ...StateOption) MutableValue
	// Deprecated: use Subscribe instead
	SetOnStateChange(callback func())
}
