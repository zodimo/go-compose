package icon

import (
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

type IconOptions struct {
	Modifier ui.Modifier
	Color    graphics.Color
	Size     unit.Dp       // Unified size for both IconBytes and SymbolName
	FontSize unit.TextUnit // Deprecated: Use Size instead
}

type IconOption func(*IconOptions)

func DefaultIconOptions() IconOptions {
	return IconOptions{
		Modifier: ui.EmptyModifier,
		Color:    graphics.ColorCompanion.Unspecified(),
		Size:     unit.DpCompanion.Unspecified(),
		FontSize: unit.TextUnitCompanion.Unspecified(),
	}
}

func WithModifier(m ui.Modifier) IconOption {
	return func(o *IconOptions) {
		o.Modifier = m
	}
}

func WithColor(col graphics.Color) IconOption {
	return func(o *IconOptions) {
		o.Color = col
	}
}

// WithSize sets the icon size in Dp.
// For IconBytes, this constrains the icon dimensions.
// For SymbolName, this sets the font size (Dp is treated as Sp).
func WithSize(size unit.Dp) IconOption {
	return func(o *IconOptions) {
		o.Size = size
	}
}

// Deprecated: Use WithSize instead.
// WithSymbolSize sets the font size for symbol icons.
func WithSymbolSize(size unit.TextUnit) IconOption {
	return func(o *IconOptions) {
		o.FontSize = size
	}
}
