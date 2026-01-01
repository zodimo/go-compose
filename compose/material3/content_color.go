package material3

import (
	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/ui/graphics"
)

// LocalContentColor is a CompositionLocal containing the preferred content color for a given
// position in the hierarchy. This typically represents the "on" color for a color in ColorScheme.
// For example, if the background color is ColorScheme.surface, this color is typically set to
// ColorScheme.onSurface.
//
// This color should be used for any typography / iconography, to ensure that the color of these
// adjusts when the background color changes. For example, on a dark background, text should be
// light, and on a light background, text should be dark.
//
// Defaults to Color.Black if no color has been explicitly set.
var LocalContentColor = compose.CompositionLocalOf(func() graphics.Color {
	return graphics.ColorBlack
})
