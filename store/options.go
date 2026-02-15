package store

import "context"

type PersistentStateOptions struct {
	StartFrameTriggerReceiver func(trigger func())
	EndFrameTriggerReceiver   func(trigger func() []string)
	RootContext               context.Context
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

func WithRootContext(rootContext context.Context) PersistentStateOption {
	return func(opts *PersistentStateOptions) {
		opts.RootContext = rootContext
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
	RootContext               context.Context
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

func WithLifecycleRootContext(rootContext context.Context) FrameLifecycleHandlerOption {
	return func(opts *FrameLifecycleHandlerOptions) {
		opts.RootContext = rootContext
	}
}

func WithLifecycleDebugMode(debugMode bool) FrameLifecycleHandlerOption {
	return func(opts *FrameLifecycleHandlerOptions) {
		opts.DebugMode = debugMode
	}
}
