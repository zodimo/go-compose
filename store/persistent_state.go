package store

import (
	"context"
	"reflect"

	"github.com/zodimo/go-compose/state"
)

type PersistentStateInterface = state.PersistentState

type PersistentState struct {
	scopes      map[string]state.MutableValue
	subscribers *state.SubscriptionManager

	frameLifecycleHandler *FrameLifecycleHandler
}

func NewPersistentState(options ...PersistentStateOption) PersistentStateInterface {
	opts := PersistentStateOptions{
		StartFrameTriggerReceiver: func(trigger func()) {},
		EndFrameTriggerReceiver:   func(trigger func() []string) {},
		DebugMode:                 false,
		RootContext:               context.Background(),
		Storage:                   map[string]state.MutableValue{},
	}

	for _, option := range options {
		if option != nil {
			option(&opts)
		}
	}

	frameLifecycleHandler := NewFrameLifecycleHandler(
		WithLifecycleTriggerReceivers(
			opts.StartFrameTriggerReceiver,
			opts.EndFrameTriggerReceiver,
		),
		WithLifecycleDebugMode(opts.DebugMode),
		WithLifecycleRootContext(opts.RootContext),
	)

	return &PersistentState{
		scopes:                opts.Storage,
		subscribers:           state.NewSubscriptionManager(),
		frameLifecycleHandler: frameLifecycleHandler,
	}
}

// Deprecated: use Subscribe instead
func (ps *PersistentState) SetOnStateChange(callback func()) {
	ps.subscribers.Subscribe(callback)
}

func (ps *PersistentState) Subscribe(callback func()) state.Subscription {
	return ps.subscribers.Subscribe(callback)
}

func (ps *PersistentState) GetState(id string, initial func() any, options ...state.StateOption) state.MutableValue {

	opts := state.StateOptions{
		Compare:     reflect.DeepEqual,
		OnForgotten: func() {},
	}
	for _, option := range options {
		option(&opts)
	}

	if v, ok := ps.scopes[id]; ok {
		ps.frameLifecycleHandler.Ping(id)
		return v
	}

	initialValue := initial()

	ps.scopes[id] = state.NewMutableValue(
		initialValue,
		func(any) {
			ps.subscribers.NotifyAll()
		},
		opts.Compare,
		state.MutableValueWithOnForgotten(opts.OnForgotten),
	)

	// Mark this key as accessed for GC tracking
	ps.frameLifecycleHandler.AttachValueToLifecycle(id, ps.scopes[id])

	return ps.scopes[id]
}

func (ps *PersistentState) StartFrame() {
	ps.frameLifecycleHandler.StartFrame()
}

func (ps *PersistentState) EndFrame() {
	keysToRemove := ps.frameLifecycleHandler.EndFrame()
	for _, key := range keysToRemove {
		delete(ps.scopes, key)
	}
}
