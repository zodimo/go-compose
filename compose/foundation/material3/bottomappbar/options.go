package bottomappbar

import (
	"gioui.org/unit"
	"github.com/zodimo/go-compose/theme"
)

// PaddingValues describes the padding to be applied to the content.
type PaddingValues struct {
	Start, Top, End, Bottom unit.Dp
}

// BottomAppBarOptions configuration
type BottomAppBarOptions struct {
	Modifier             Modifier
	ContainerColor       theme.ColorDescriptor
	ContentColor         theme.ColorDescriptor
	TonalElevation       unit.Dp
	ContentPadding       PaddingValues
	FloatingActionButton Composable
}

// BottomAppBarOption is a function that configures a BottomAppBar.
type BottomAppBarOption func(*BottomAppBarOptions)

// DefaultBottomAppBarOptions returns the default options.
func DefaultBottomAppBarOptions() BottomAppBarOptions {
	s, t, e, b := BottomAppBarDefaults.ContentPadding()
	return BottomAppBarOptions{
		Modifier:       EmptyModifier,
		ContainerColor: BottomAppBarDefaults.Colors().ContainerColor,
		ContentColor:   BottomAppBarDefaults.Colors().ContentColor,
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

func WithContainerColor(colorDesc theme.ColorDescriptor) BottomAppBarOption {
	return func(o *BottomAppBarOptions) {
		o.ContainerColor = colorDesc
	}
}

func WithContentColor(colorDesc theme.ColorDescriptor) BottomAppBarOption {
	return func(o *BottomAppBarOptions) {
		o.ContentColor = colorDesc
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
