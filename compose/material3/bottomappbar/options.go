package bottomappbar

import (
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

// PaddingValues describes the padding to be applied to the content.
type PaddingValues struct {
	Start, Top, End, Bottom unit.Dp
}

// BottomAppBarOptions configuration
type BottomAppBarOptions struct {
	Modifier             Modifier
	ContainerColor       graphics.Color
	ContentColor         graphics.Color
	TonalElevation       unit.Dp
	ContentPadding       PaddingValues
	FloatingActionButton Composable
}

// BottomAppBarOption is a function that configures a BottomAppBar.
type BottomAppBarOption func(*BottomAppBarOptions)

// DefaultBottomAppBarOptions returns the default options.
func DefaultBottomAppBarOptions(c Composer) BottomAppBarOptions {

	s, t, e, b := BottomAppBarDefaults.ContentPadding()
	return BottomAppBarOptions{
		Modifier:       EmptyModifier,
		ContainerColor: BottomAppBarDefaults.Colors(c).ContainerColor,
		ContentColor:   BottomAppBarDefaults.Colors(c).ContentColor,
		TonalElevation: BottomAppBarDefaults.ContainerElevation(),
		ContentPadding: PaddingValues{
			Start:  s,
			Top:    t,
			End:    e,
			Bottom: b,
		},
	}
}

func WithModifier(m Modifier) BottomAppBarOption {
	return func(o *BottomAppBarOptions) {
		o.Modifier = m
	}
}

func WithContainerColor(col graphics.Color) BottomAppBarOption {
	return func(o *BottomAppBarOptions) {
		o.ContainerColor = col
	}
}

func WithContentColor(col graphics.Color) BottomAppBarOption {
	return func(o *BottomAppBarOptions) {
		o.ContentColor = col
	}
}

func WithTonalElevation(e unit.Dp) BottomAppBarOption {
	return func(o *BottomAppBarOptions) {
		o.TonalElevation = e
	}
}

func WithContentPadding(p PaddingValues) BottomAppBarOption {
	return func(o *BottomAppBarOptions) {
		o.ContentPadding = p
	}
}

func WithFloatingActionButton(fab Composable) BottomAppBarOption {
	return func(o *BottomAppBarOptions) {
		o.FloatingActionButton = fab
	}
}
