package snackbar

import (
	"context"

	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
)

// SnackbarOptions holds configuration for a snackbar.
type SnackbarOptions struct {
	Modifier ui.Modifier

	Context           context.Context
	Duration          SnackbarDuration
	ActionLabel       string
	WithDismissAction bool
	OnResult          func(SnackbarResult)

	Shape                     shape.Shape
	ContainerColor            graphics.Color
	ContentColor              graphics.Color
	ActionColor               graphics.Color
	ActionContentColor        graphics.Color
	DismissActionContentColor graphics.Color
}

// SnackbarOption is a functional option for configuring a snackbar.
type SnackbarOption func(*SnackbarOptions)

func WithModifier(m ui.Modifier) SnackbarOption {
	return func(o *SnackbarOptions) {
		o.Modifier = m
	}
}
func WithContext(ctx context.Context) SnackbarOption {
	return func(o *SnackbarOptions) {
		o.Context = ctx
	}
}

// WithDuration sets the display duration of the snackbar.
func WithDuration(duration SnackbarDuration) SnackbarOption {
	return func(o *SnackbarOptions) {
		o.Duration = duration
	}
}

// WithActionLabel sets an action button label on the snackbar.
func WithActionLabel(label string) SnackbarOption {
	return func(o *SnackbarOptions) {
		o.ActionLabel = label
	}
}

// WithDismissAction enables a dismiss action on the snackbar.
func WithDismissAction() SnackbarOption {
	return func(o *SnackbarOptions) {
		o.WithDismissAction = true
	}
}

// WithOnResult sets a callback invoked when the snackbar is dismissed or its action is performed.
func WithOnResult(callback func(SnackbarResult)) SnackbarOption {
	return func(o *SnackbarOptions) {
		o.OnResult = callback
	}
}

// DefaultOptions returns the default snackbar options.
// Matches Kotlin: Short duration when no action, Indefinite when action present.
func DefaultOptions() SnackbarOptions {
	return SnackbarOptions{
		Modifier:                  ui.EmptyModifier,
		Context:                   context.Background(),
		Duration:                  SnackbarDurationShort,
		Shape:                     shape.ShapeUnspecified,
		ContainerColor:            graphics.ColorUnspecified,
		ContentColor:              graphics.ColorUnspecified,
		ActionColor:               graphics.ColorUnspecified,
		ActionContentColor:        graphics.ColorUnspecified,
		DismissActionContentColor: graphics.ColorUnspecified,
	}
}

// resolveDefaults applies Kotlin's default duration logic:
// if an action label is present and duration is Short, upgrade to Indefinite.
func (o *SnackbarOptions) resolveDefaults() {
	if o.ActionLabel != "" && o.Duration == SnackbarDurationShort {
		o.Duration = SnackbarDurationIndefinite
	}
}
