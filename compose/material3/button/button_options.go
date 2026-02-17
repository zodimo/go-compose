package button

import (
	"git.sr.ht/~schnwalter/gio-mw/widget/button"
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/compose/ui/graphics"
)

type ButtonOptions struct {
	Modifier     ui.Modifier
	Button       *button.Button
	Enabled      bool
	ContentColor graphics.Color
}

type ButtonOption func(o *ButtonOptions)

func WithModifier(m ui.Modifier) ButtonOption {
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

func WithContentColor(color graphics.Color) ButtonOption {
	return func(o *ButtonOptions) {
		o.ContentColor = color
	}
}
