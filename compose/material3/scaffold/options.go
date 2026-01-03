package scaffold

import (
	"github.com/zodimo/go-compose/compose/material3"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/internal/modifier"
)

type ScaffoldOptions struct {
	Modifier                     Modifier
	TopBar                       Composable
	BottomBar                    Composable
	SnackbarHost                 Composable
	FloatingActionButton         Composable
	FloatingActionButtonPosition FabPosition
	ContainerColor               graphics.Color
	ContentColor                 graphics.Color
}

type ScaffoldOption func(*ScaffoldOptions)

// DefaultScaffoldOptions returns the default options for a Scaffold.
func DefaultScaffoldOptions(c Composer) ScaffoldOptions {
	theme := material3.Theme(c)
	return ScaffoldOptions{
		Modifier:                     modifier.EmptyModifier,
		TopBar:                       nil,
		BottomBar:                    nil,
		SnackbarHost:                 nil,
		FloatingActionButton:         nil,
		FloatingActionButtonPosition: FabPositionEnd,
		ContainerColor:               theme.ColorScheme().Surface.Color,   //theme.ColorHelper.ColorSelector().SurfaceRoles.Surface,
		ContentColor:                 theme.ColorScheme().Surface.OnColor, //theme.ColorHelper.ColorSelector().SurfaceRoles.OnSurface,
	}
}

// WithModifier sets the modifier for the Scaffold.
func WithModifier(m Modifier) ScaffoldOption {
	return func(o *ScaffoldOptions) {
		if o.Modifier == nil {
			o.Modifier = m
		} else {
			o.Modifier = m
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
func WithContainerColor(col graphics.Color) ScaffoldOption {
	return func(o *ScaffoldOptions) { o.ContainerColor = col }
}

// WithContentColor sets the content color of the scaffold.
func WithContentColor(col graphics.Color) ScaffoldOption {
	return func(o *ScaffoldOptions) { o.ContentColor = col }
}
