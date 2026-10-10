package text

import (
	"math"

	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/text"
	"github.com/zodimo/go-maybe"
)

func DefaultTextOptions() TextOptions {

	return TextOptions{
		Modifier:  ui.EmptyModifier,
		TextStyle: text.TextStyleCompanion.Unspecified(),
		MaxLines:  math.MaxInt32,

		// BACKWARDS COMPATIBILITY
		LineHeightScale: maybe.None[float32](),
		Selectable:      maybe.None[bool](),
		SelectionColor:  graphics.ColorUnspecified,
	}
}
