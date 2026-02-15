package store

type PersistentStateOptions struct {
	StartFrameTriggerReceiver func(trigger func())
	EndFrameTriggerReceiver   func(trigger func() []string)
	DebugMode                 bool
}

func WithFrameLifecycleTriggerReceivers(
	startFrameTriggerReceiver func(trigger func()),
	endFrameTriggerReceiver func(trigger func() []string),
) PersistentStateOption {
	return func(opts *PersistentStateOptions) {
		opts.StartFrameTriggerReceiver = startFrameTriggerReceiver
		opts.EndFrameTriggerReceiver = endFrameTriggerReceiver
	}
}

func WithDebugMode() PersistentStateOption {
	return func(opts *PersistentStateOptions) {
		opts.DebugMode = true
	}
}

type PersistentStateOption func(*PersistentStateOptions)

type FrameLifecycleHandlerOptions struct {
	StartFrameTriggerReceiver func(trigger func())
	EndFrameTriggerReceiver   func(trigger func() []string)
	DebugMode                 bool
}

type FrameLifecycleHandlerOption func(*FrameLifecycleHandlerOptions)

func WithLifecycleTriggerReceivers(
	startFrameTriggerReceiver func(trigger func()),
	endFrameTriggerReceiver func(trigger func() []string),
) FrameLifecycleHandlerOption {
	return func(opts *FrameLifecycleHandlerOptions) {
		opts.StartFrameTriggerReceiver = startFrameTriggerReceiver
		opts.EndFrameTriggerReceiver = endFrameTriggerReceiver
	}
}

func WithLifecycleDebugMode(debugMode bool) FrameLifecycleHandlerOption {
	return func(opts *FrameLifecycleHandlerOptions) {
		opts.DebugMode = debugMode
	}
}
