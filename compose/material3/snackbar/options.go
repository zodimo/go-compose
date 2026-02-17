package snackbar

// SnackbarOptions holds configuration for a snackbar.
type SnackbarOptions struct {
	Duration          SnackbarDuration
	ActionLabel       string
	WithDismissAction bool
}

// SnackbarOption is a functional option for configuring a snackbar.
type SnackbarOption func(*SnackbarOptions)

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

// DefaultOptions returns the default snackbar options.
// Matches Kotlin: Short duration when no action, Indefinite when action present.
func DefaultOptions() SnackbarOptions {
	return SnackbarOptions{
		Duration: SnackbarDurationShort,
	}
}

// resolveDefaults applies Kotlin's default duration logic:
// if an action label is present and duration is Short, upgrade to Indefinite.
func (o *SnackbarOptions) resolveDefaults() {
	if o.ActionLabel != "" && o.Duration == SnackbarDurationShort {
		o.Duration = SnackbarDurationIndefinite
	}
}
