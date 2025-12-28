package compose

import (
	"gioui.org/text"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/theme"
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
var LocalContentColor = CompositionLocalOf(func() graphics.Color {
	return graphics.ColorBlack
})

var LocalTextShaper = CompositionLocalOf(func() *text.Shaper {
	return theme.GetThemeManager().MaterialTheme().Shaper
})
