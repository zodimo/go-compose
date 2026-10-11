package divider

import (
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/layout"
)

func DefaultDividerOptions() DividerOptions {
	return DividerOptions{
		Modifier:  ui.EmptyModifier,
		Thickness: 1,
		Color:     graphics.ColorCompanion.Unspecified(),
		Axis:      layout.AxisUnspecified,
	}
}
