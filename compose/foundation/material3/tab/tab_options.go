package tab

import (
	"go-compose-dev/internal/modifier"
	"image/color"

	"github.com/zodimo/go-maybe"
)

// TabRowOptions holds the configuration for TabRow.
type TabRowOptions struct {
	Modifier       modifier.Modifier
	ContainerColor maybe.Maybe[color.Color]
	ContentColor   maybe.Maybe[color.Color]
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
		ContainerColor: maybe.None[color.Color](),
		ContentColor:   maybe.None[color.Color](),
		Indicator:      TabRowDefaults.Indicator(),
	}
}

// TabOptions holds the configuration for Tab.
type TabOptions struct {
	Modifier               modifier.Modifier
	Selected               bool
	OnClick                func()
	Enabled                bool
	SelectedContentColor   color.NRGBA
	UnselectedContentColor color.NRGBA
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
