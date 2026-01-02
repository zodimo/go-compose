package appbar

// TopAppBarOptions configuration
// TopAppBarOptions configuration
type TopAppBarOptions struct {
	Modifier       Modifier
	NavigationIcon Composable
	Actions        []Composable
	Colors         TopAppBarColors
}

type TopAppBarOption func(*TopAppBarOptions)

func DefaultTopAppBarOptions(c Composer) TopAppBarOptions {
	return TopAppBarOptions{
		Modifier: EmptyModifier,
		Colors:   TopAppBarDefaults.Colors(c),
	}
}

func WithModifier(m Modifier) TopAppBarOption {
	return func(o *TopAppBarOptions) {
		o.Modifier = m
	}
}

func WithNavigationIcon(icon Composable) TopAppBarOption {
	return func(o *TopAppBarOptions) {
		o.NavigationIcon = icon
	}
}

func WithActions(actions ...Composable) TopAppBarOption {
	return func(o *TopAppBarOptions) {
		o.Actions = actions
	}
}

func WithColors(colors TopAppBarColors) TopAppBarOption {
	return func(o *TopAppBarOptions) {
		o.Colors = colors
	}
}
