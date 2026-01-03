package tab

import (
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/internal/modifier"
)

// TabRowOptions holds the configuration for TabRow.
type TabRowOptions struct {
	Modifier       modifier.Modifier
	ContainerColor graphics.Color
	ContentColor   graphics.Color
	Indicator      Composable
}

type TabRowOption func(options *TabRowOptions)

func WithTabRowModifier(modifier modifier.Modifier) TabRowOption {
	return func(options *TabRowOptions) {
		options.Modifier = options.Modifier.Then(modifier)
	}
}

// DefaultTabRowOptions returns the default options for TabRow.
func DefaultTabRowOptions() TabRowOptions {
	return TabRowOptions{
		Modifier:       modifier.EmptyModifier,
		ContainerColor: graphics.ColorUnspecified,
		ContentColor:   graphics.ColorUnspecified,
		Indicator:      TabRowDefaults.Indicator(),
	}
}

// TabOptions holds the configuration for Tab.
type TabOptions struct {
	Modifier               modifier.Modifier
	Selected               bool
	OnClick                func()
	Enabled                bool
	SelectedContentColor   graphics.Color
	UnselectedContentColor graphics.Color
}

type TabOption func(options *TabOptions)

func WithModifier(m modifier.Modifier) TabOption {
	return func(options *TabOptions) {
		options.Modifier = m
	}
}

// DefaultTabOptions returns the default options for Tab.
func DefaultTabOptions(c Composer) TabOptions {
	return TabOptions{
		Modifier:               modifier.EmptyModifier,
		Enabled:                true,
		SelectedContentColor:   TabDefaults.SelectedContentColor(c),
		UnselectedContentColor: TabDefaults.UnselectedContentColor(c),
	}
}
