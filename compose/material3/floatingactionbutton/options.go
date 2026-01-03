package floatingactionbutton

import (
	"github.com/zodimo/go-compose/compose/material3"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	"github.com/zodimo/go-compose/internal/modifier"
	"github.com/zodimo/go-compose/pkg/api"

	// "github.com/zodimo/go-compose/pkg/api" // Removed unused api import

	"git.sr.ht/~schnwalter/gio-mw/token"
)

type FabSize int

const (
	// FabSizeMedium is the default 56dp FAB.
	// Available in both Original M3 and M3 Expressive.
	FabSizeMedium FabSize = iota
	// FabSizeSmall is the 40dp FAB.
	// Available in Original M3, but Deprecated in M3 Expressive (use a larger size).
	FabSizeSmall
	// FabSizeLarge is the 96dp FAB.
	// Available in both Original M3 and M3 Expressive.
	FabSizeLarge
)

type FloatingActionButtonOptions struct {
	ContainerColor graphics.Color
	ContentColor   graphics.Color
	Elevation      token.ElevationLevel
	Modifier       modifier.Modifier
	Shape          shape.Shape
	Size           FabSize
}

type FloatingActionButtonOption func(*FloatingActionButtonOptions)

// DefaultFABElevation is the default elevation for a FAB (Level 3).
var DefaultFABElevation = token.ElevationLevel3

func DefaultFloatingActionButtonOptions(c api.Composer) FloatingActionButtonOptions {
	theme := material3.Theme(c)
	return FloatingActionButtonOptions{
		ContainerColor: theme.ColorScheme().PrimaryContainer.Color,
		ContentColor:   theme.ColorScheme().PrimaryContainer.OnColor,
		Elevation:      DefaultFABElevation,
		Modifier:       modifier.EmptyModifier,
		Shape:          shape.ShapeCircle,
		Size:           FabSizeMedium,
	}
}

func WithContainerColor(col graphics.Color) FloatingActionButtonOption {
	return func(o *FloatingActionButtonOptions) {
		o.ContainerColor = col
	}
}

func WithContentColor(col graphics.Color) FloatingActionButtonOption {
	return func(o *FloatingActionButtonOptions) {
		o.ContentColor = col
	}
}

// WithElevation sets the base elevation.
func WithElevation(elevation token.ElevationLevel) FloatingActionButtonOption {
	return func(o *FloatingActionButtonOptions) {
		o.Elevation = elevation
	}
}

func WithModifier(m modifier.Modifier) FloatingActionButtonOption {
	return func(o *FloatingActionButtonOptions) {
		o.Modifier = m
	}
}

func WithShape(shape shape.Shape) FloatingActionButtonOption {
	return func(o *FloatingActionButtonOptions) {
		o.Shape = shape
	}
}

func WithSize(size FabSize) FloatingActionButtonOption {
	return func(o *FloatingActionButtonOptions) {
		o.Size = size
	}
}
