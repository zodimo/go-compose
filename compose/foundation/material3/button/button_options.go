package button

type ButtonOptions struct {
	Modifier Modifier
}

type ButtonOption func(o *ButtonOptions)

func WithModifier(modifier Modifier) ButtonOption {
	return func(o *ButtonOptions) {
		o.Modifier = o.Modifier.Then(modifier)
	}
}
