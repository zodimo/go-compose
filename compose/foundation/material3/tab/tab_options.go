package tab

import (
	"github.com/zodimo/go-compose/internal/modifier"
	"github.com/zodimo/go-compose/theme"

	"github.com/zodimo/go-maybe"
)

// TabRowOptions holds the configuration for TabRow.
type TabRowOptions struct {
	Modifier       modifier.Modifier
	ContainerColor maybe.Maybe[theme.ColorDescriptor]
	ContentColor   maybe.Maybe[theme.ColorDescriptor]
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
		ContainerColor: maybe.None[theme.ColorDescriptor](),
		ContentColor:   maybe.None[theme.ColorDescriptor](),
		Indicator:      TabRowDefaults.Indicator(),
	}
}

// TabOptions holds the configuration for Tab.
type TabOptions struct {
	Modifier               modifier.Modifier
	Selected               bool
	OnClick                func()
	Enabled                bool
	SelectedContentColor   theme.ColorDescriptor
	UnselectedContentColor theme.ColorDescriptor
}

type TabOption func(options *TabOptions)

func WithModifier(modifier modifier.Modifier) TabOption {
	return func(options *TabOptions) {
		options.Modifier = options.Modifier.Then(modifier)
	}
}

// DefaultTabOptions returns the default options for Tab.
func DefaultTabOptions() TabOptions {
	return TabOptions{
		Modifier:               modifier.EmptyModifier,
		Enabled:                true,
		SelectedContentColor:   TabDefaults.SelectedContentColor(),
		UnselectedContentColor: TabDefaults.UnselectedContentColor(),
	}
}
