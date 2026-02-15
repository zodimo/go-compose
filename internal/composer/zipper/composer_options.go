package zipper

import "time"

type ComposerOptions struct {
	TimeNowFunc  func() time.Time
	OnStartFrame func()
	OnEndFrame   func()
}

func WithTimeNow(now time.Time) ComposerOption {
	return func(options *ComposerOptions) {
		options.TimeNowFunc = func() time.Time { return now }
	}
}

func WithOnStartFrame(onStartFrame func()) ComposerOption {
	return func(options *ComposerOptions) {
		options.OnStartFrame = onStartFrame
	}
}

func WithOnEndFrame(onEndFrame func()) ComposerOption {
	return func(options *ComposerOptions) {
		options.OnEndFrame = onEndFrame
	}
}

type ComposerOption func(*ComposerOptions)

func DefaultComposerOptions() ComposerOptions {
	return ComposerOptions{
		TimeNowFunc:  time.Now,
		OnStartFrame: func() {},
		OnEndFrame:   func() {},
	}
}
