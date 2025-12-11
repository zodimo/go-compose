package card

type CardOptions struct {
	Modifier Modifier
}

type CardOption func(o *CardOptions)

func WithModifier(modifier Modifier) CardOption {
	return func(o *CardOptions) {
		o.Modifier = o.Modifier.Then(modifier)
	}
}
