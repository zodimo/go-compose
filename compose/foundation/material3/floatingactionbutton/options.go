package floatingactionbutton

import (
	"image/color"

	"go-compose-dev/compose/ui/graphics/shape"
	"go-compose-dev/internal/modifier"

	// "go-compose-dev/pkg/api" // Removed unused api import

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
	ContainerColor color.Color
	ContentColor   color.Color
	Elevation      token.ElevationLevel
	Modifier       modifier.Modifier
	Shape          shape.Shape
	Size           FabSize
}

type FloatingActionButtonOption func(*FloatingActionButtonOptions)

// DefaultFABElevation is the default elevation for a FAB (Level 3).
var DefaultFABElevation = token.ElevationLevel3

func DefaultFloatingActionButtonOptions() FloatingActionButtonOptions {
	return FloatingActionButtonOptions{
		// Defaults will be resolved in the component if zero values are detected,
		// or we can set them here if we have access to theme.
		Elevation: DefaultFABElevation,
		Modifier:  modifier.EmptyModifier,
		Shape:     shape.ShapeCircle,
		Size:      FabSizeMedium,
	}
}

func WithContainerColor(color color.Color) FloatingActionButtonOption {
	return func(o *FloatingActionButtonOptions) {
		o.ContainerColor = color
	}
}

func WithContentColor(color color.Color) FloatingActionButtonOption {
	return func(o *FloatingActionButtonOptions) {
		o.ContentColor = color
	}
}

// WithElevation sets the base elevation.
func WithElevation(elevation token.ElevationLevel) FloatingActionButtonOption {
	return func(o *FloatingActionButtonOptions) {
		o.Elevation = elevation
	}
}

func WithModifier(modifier modifier.Modifier) FloatingActionButtonOption {
	return func(o *FloatingActionButtonOptions) {
		o.Modifier = modifier
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
