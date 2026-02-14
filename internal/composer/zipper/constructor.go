package zipper

import "time"

type ComposerOptions struct {
	TimeNowFunc func() time.Time
}

func WithTimeNow(now time.Time) ComposerOption {
	return func(options *ComposerOptions) {
		options.TimeNowFunc = func() time.Time { return now }
	}
}

type ComposerOption func(*ComposerOptions)

func DefaultComposerOptions() ComposerOptions {
	return ComposerOptions{
		TimeNowFunc: time.Now,
	}
}

func NewComposer(state PersistentState, options ...ComposerOption) Composer {
	opts := DefaultComposerOptions()
	for _, option := range options {
		if option != nil {
			option(&opts)
		}
	}

	idManager := GetScopedIdentityManager("composer")
	idManager.ResetKeyCounter()

	return &composer{
		timeNowFunc:    opts.TimeNowFunc,
		focus:          nil,
		path:           []pathItem{},
		memo:           EmptyMemo,
		state:          state,
		idManager:      idManager,
		locals:         make(map[interface{}]interface{}),
		providersStack: []map[interface{}]interface{}{},
	}
}
