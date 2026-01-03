package radiobutton

type RadioButtonOptions struct {
	Modifier Modifier
	Enabled  bool
	Colors   RadioButtonColors
}

type RadioButtonOption func(*RadioButtonOptions)

func DefaultRadioButtonOptions(c Composer) RadioButtonOptions {
	return RadioButtonOptions{
		Modifier: EmptyModifier,
		Enabled:  true,
		Colors:   Defaults.Colors(c), // Use nil/defaults
	}
}

func WithModifier(m Modifier) RadioButtonOption {
	return func(o *RadioButtonOptions) {
		o.Modifier = m
	}
}

func WithEnabled(enabled bool) RadioButtonOption {
	return func(o *RadioButtonOptions) {
		o.Enabled = enabled
	}
}

func WithColors(colors RadioButtonColors) RadioButtonOption {
	return func(o *RadioButtonOptions) {
		o.Colors = colors
	}
}
