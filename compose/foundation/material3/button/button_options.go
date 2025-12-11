package button

import "git.sr.ht/~schnwalter/gio-mw/widget/button"

type ButtonOptions struct {
	Modifier Modifier
	Button   *button.Button
}

type ButtonOption func(o *ButtonOptions)

func WithModifier(modifier Modifier) ButtonOption {
	return func(o *ButtonOptions) {
		o.Modifier = o.Modifier.Then(modifier)
	}
}

func WithButton(button *button.Button) ButtonOption {
	return func(o *ButtonOptions) {
		o.Button = button
	}
}
