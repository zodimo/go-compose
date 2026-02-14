package state

func WithCompare(compare func(any, any) bool) StateOption {
	if compare == nil {
		panic("WithCompare: compare cannot be nil")
	}
	return func(opts *StateOptions) {
		opts.Compare = compare
	}
}

type PersistentState interface {
	StateChangeNotifier
	GetState(key string, initial func() any, options ...StateOption) MutableValue
	// Deprecated: use Subscribe instead
	SetOnStateChange(callback func())
}
