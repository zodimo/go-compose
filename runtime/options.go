package runtime

type RuntimeOptions struct {
	OnStartFrame func()
	OnEndFrame   func()
}

func WithOnStartFrame(onStartFrame func()) RuntimeOption {
	return func(options *RuntimeOptions) {
		options.OnStartFrame = onStartFrame
	}
}

func WithOnEndFrame(onEndFrame func()) RuntimeOption {
	return func(options *RuntimeOptions) {
		options.OnEndFrame = onEndFrame
	}
}

type RuntimeOption func(*RuntimeOptions)

func DefaultRuntimeOptions() RuntimeOptions {
	return RuntimeOptions{
		OnStartFrame: func() {},
		OnEndFrame:   func() {},
	}
}
