package runtime

func NewRuntime(options ...RuntimeOption) Runtime {
	opts := DefaultRuntimeOptions()
	for _, option := range options {
		if option != nil {
			option(&opts)
		}
	}
	return &runtime{}
}
