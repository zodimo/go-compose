package navigationbar

import (
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

// NavigationBarOptions configuration
type NavigationBarOptions struct {
	Modifier       Modifier
	ContainerColor graphics.Color
	ContentColor   graphics.Color
	TonalElevation unit.Dp
	Height         unit.Dp
}

// NavigationBarOption is a function that configures a NavigationBar.
type NavigationBarOption func(*NavigationBarOptions)

// DefaultNavigationBarOptions returns the default options.
func DefaultNavigationBarOptions(c Composer) NavigationBarOptions {
	return NavigationBarOptions{
		Modifier:       EmptyModifier,
		ContainerColor: NavigationBarDefaults.Colors(c).ContainerColor,
		ContentColor:   NavigationBarDefaults.Colors(c).ContentColor,
		TonalElevation: NavigationBarDefaults.ContainerElevation(),
		Height:         NavigationBarDefaults.Height(),
	}
}

func WithModifier(m Modifier) NavigationBarOption {
	return func(o *NavigationBarOptions) {
		o.Modifier = m
	}
}

func WithContainerColor(col graphics.Color) NavigationBarOption {
	return func(o *NavigationBarOptions) {
		o.ContainerColor = col
	}
}

func WithContentColor(col graphics.Color) NavigationBarOption {
	return func(o *NavigationBarOptions) {
		o.ContentColor = col
	}
}

func WithTonalElevation(e unit.Dp) NavigationBarOption {
	return func(o *NavigationBarOptions) {
		o.TonalElevation = e
	}
}

func WithHeight(h unit.Dp) NavigationBarOption {
	return func(o *NavigationBarOptions) {
		o.Height = h
	}
}
