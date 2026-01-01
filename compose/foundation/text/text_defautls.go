package text

import (
	"math"

	gioFont "gioui.org/font"
	"gioui.org/unit"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/text"
	"github.com/zodimo/go-maybe"
)

func DefaultTextOptions() TextOptions {

	return TextOptions{
		Modifier:  EmptyModifier,
		TextStyle: text.TextStyleUnspecified,
		MaxLines:  math.MaxInt32,

		// BACKWARDS COMPATIBILITY
		Font:            maybe.None[gioFont.Font](),
		WrapPolicy:      maybe.None[WrapPolicy](),
		LineHeight:      maybe.None[unit.Sp](),
		LineHeightScale: maybe.None[float32](),
		Selectable:      maybe.None[bool](),
		SelectionColor:  graphics.ColorUnspecified,
		TextSize:        maybe.None[unit.Sp](),
		Strikethrough:   maybe.None[bool](),
	}
}
