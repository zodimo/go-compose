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
