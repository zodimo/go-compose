package divider

import (
	"github.com/zodimo/go-compose/compose/ui/graphics"
)

func DefaultDividerOptions() DividerOptions {
	return DividerOptions{
		Modifier:  EmptyModifier,
		Thickness: 1,
		Color:     graphics.ColorUnspecified,
	}
}
