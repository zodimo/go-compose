package scaffold

import (
	"go-compose-dev/internal/modifier"
	"image/color"
)

type ScaffoldOptions struct {
	Modifier                     Modifier
	TopBar                       Composable
	BottomBar                    Composable
	SnackbarHost                 Composable
	FloatingActionButton         Composable
	FloatingActionButtonPosition FabPosition
	ContainerColor               color.Color
	ContentColor                 color.Color
}

type ScaffoldOption func(*ScaffoldOptions)

// DefaultScaffoldOptions returns the default options for a Scaffold.
func DefaultScaffoldOptions() ScaffoldOptions {
	return ScaffoldOptions{
		Modifier:                     modifier.EmptyModifier,
		TopBar:                       nil,
		BottomBar:                    nil,
		SnackbarHost:                 nil,
		FloatingActionButton:         nil,
		FloatingActionButtonPosition: FabPositionEnd,
		ContainerColor:               color.Transparent, // Will be resolved by Surface if nil/transparent
		ContentColor:                 color.NRGBA{},     // Will be resolved by Surface
	}
}

// WithModifier sets the modifier for the Scaffold.
func WithModifier(m Modifier) ScaffoldOption {
	return func(o *ScaffoldOptions) {
		if o.Modifier == nil {
			o.Modifier = m
		} else {
			o.Modifier = o.Modifier.Then(m)
		}
	}
}

// WithTopBar sets the top app bar.
func WithTopBar(c Composable) ScaffoldOption {
	return func(o *ScaffoldOptions) { o.TopBar = c }
}

// WithBottomBar sets the bottom app bar.
func WithBottomBar(c Composable) ScaffoldOption {
	return func(o *ScaffoldOptions) { o.BottomBar = c }
}

// WithSnackbarHost sets the snackbar host.
func WithSnackbarHost(c Composable) ScaffoldOption {
	return func(o *ScaffoldOptions) { o.SnackbarHost = c }
}

// WithFloatingActionButton sets the FAB.
func WithFloatingActionButton(c Composable) ScaffoldOption {
	return func(o *ScaffoldOptions) { o.FloatingActionButton = c }
}

// WithFloatingActionButtonPosition sets the position of the FAB.
func WithFloatingActionButtonPosition(pos FabPosition) ScaffoldOption {
	return func(o *ScaffoldOptions) { o.FloatingActionButtonPosition = pos }
}

// WithContainerColor sets the background color of the scaffold.
func WithContainerColor(c color.Color) ScaffoldOption {
	return func(o *ScaffoldOptions) { o.ContainerColor = c }
}

// WithContentColor sets the content color of the scaffold.
func WithContentColor(c color.Color) ScaffoldOption {
	return func(o *ScaffoldOptions) { o.ContentColor = c }
}
