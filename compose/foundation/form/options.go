package fform

import (
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/internal/modifier"
)

// FormOptions configures the Form composable.
type FormOptions struct {
	Modifier ui.Modifier
}

// FormOption mutates FormOptions.
type FormOption func(o *FormOptions)

// DefaultFormOptions returns the default form options.
func DefaultFormOptions() FormOptions {
	return FormOptions{
		Modifier: modifier.EmptyModifier,
	}
}

// WithModifier appends a modifier to the form's LazyColumn.
func WithModifier(m ui.Modifier) FormOption {
	return func(o *FormOptions) {
		o.Modifier = o.Modifier.Then(m)
	}
}
