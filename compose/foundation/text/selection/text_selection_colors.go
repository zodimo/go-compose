package selection

import (
	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/ui/graphics"
)

// TextSelectionColors represents the colors used for text selection.
type TextSelectionColors struct {
	// BackgroundColor is the color used for the selection background highlight.
	BackgroundColor graphics.Color
}

// NewTextSelectionColors creates new TextSelectionColors with the given colors.
func NewTextSelectionColors(backgroundColor graphics.Color) TextSelectionColors {
	return TextSelectionColors{
		BackgroundColor: backgroundColor,
	}
}

// DefaultTextSelectionColors returns the default text selection colors.
// Uses a light blue for selection background similar to Android defaults.
func DefaultTextSelectionColors() TextSelectionColors {
	// Default selection colors - light blue background, blue handle
	// 0x6633B5E5 is ARGB: alpha=0x66, R=0x33, G=0xB5, B=0xE5
	return TextSelectionColors{
		BackgroundColor: graphics.NewColorSrgb(0x33, 0xB5, 0xE5, 0x66), // Light blue with alpha
	}
}

// Copy returns a copy of the TextSelectionColors.
func (c TextSelectionColors) Copy() TextSelectionColors {
	return TextSelectionColors{
		BackgroundColor: c.BackgroundColor,
	}
}

// LocalTextSelectionColors is a CompositionLocal for TextSelectionColors.
// Provides selection colors to text composables within the composition.
var LocalTextSelectionColors = compose.CompositionLocalOf(DefaultTextSelectionColors)
