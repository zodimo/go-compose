package button

import "git.sr.ht/~schnwalter/gio-mw/widget/button"

type ButtonOptions struct {
	Modifier Modifier
	Button   *button.Button
	Enabled  bool
}

type ButtonOption func(o *ButtonOptions)

func WithModifier(m Modifier) ButtonOption {
	return func(o *ButtonOptions) {
		o.Modifier = m
	}
}

func WithButton(button *button.Button) ButtonOption {
	return func(o *ButtonOptions) {
		o.Button = button
	}
}

func WithEnabled(enabled bool) ButtonOption {
	return func(o *ButtonOptions) {
		o.Enabled = enabled
	}
}
