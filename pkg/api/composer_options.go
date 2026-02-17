package api

import (
	"time"

	"github.com/zodimo/go-compose/state"
	"github.com/zodimo/go-compose/store"
)

type ComposerOptions struct {
	TimeNowFunc func() time.Time
	Store       state.PersistentState

	StartFrameTriggerReceiver func(trigger func())
	EndFrameTriggerReceiver   func(trigger func())
}

type ComposerOption func(*ComposerOptions)

func ComposerWithTimeNow(now time.Time) ComposerOption {
	return func(options *ComposerOptions) {
		options.TimeNowFunc = func() time.Time { return now }
	}
}

func ComposerWithStore(store state.PersistentState) ComposerOption {
	return func(opts *ComposerOptions) {
		opts.Store = store
	}
}

func ComposerWithLifecycleTriggerReceivers(
	startFrameTriggerReceiver func(trigger func()),
	endFrameTriggerReceiver func(trigger func()),
) ComposerOption {
	return func(opts *ComposerOptions) {
		opts.StartFrameTriggerReceiver = startFrameTriggerReceiver
		opts.EndFrameTriggerReceiver = endFrameTriggerReceiver
	}
}

func DefaultComposerOptions() ComposerOptions {
	return ComposerOptions{
		TimeNowFunc:               time.Now,
		Store:                     store.NewPersistentState(),
		StartFrameTriggerReceiver: func(trigger func()) {},
		EndFrameTriggerReceiver:   func(trigger func()) {},
	}
}
