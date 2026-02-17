package runtime

type RuntimeOptions struct {
}

type RuntimeOption func(*RuntimeOptions)

func DefaultRuntimeOptions() RuntimeOptions {
	return RuntimeOptions{}
}
