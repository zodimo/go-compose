package slider

import (
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

// SliderDefaults holds default values for the Slider component.
var SliderDefaults = sliderDefaults{}

type sliderDefaults struct{}

// Colors returns the default SliderColors.
func (d sliderDefaults) Colors() SliderColors {
	return SliderColors{
		ThumbColor:            graphics.ColorUnspecified,
		ActiveTrackColor:      graphics.ColorUnspecified,
		ActiveTickColor:       graphics.ColorUnspecified,
		InactiveTrackColor:    graphics.ColorUnspecified,
		InactiveTickColor:     graphics.ColorUnspecified,
		DisabledThumbColor:    graphics.ColorUnspecified,
		DisabledActiveTrack:   graphics.ColorUnspecified,
		DisabledActiveTick:    graphics.ColorUnspecified,
		DisabledInactiveTrack: graphics.ColorUnspecified,
		DisabledInactiveTick:  graphics.ColorUnspecified,
	}
}

// Dimensions constants
var (
	TrackHeight     = unit.Dp(4)
	ThumbSize       = unit.Dp(20)
	ActiveThumbSize = unit.Dp(28) // M3 State Layer/Enlarged handle
	TickSize        = unit.Dp(2)
	ThumbTrackGap   = unit.Dp(6) // Approximate
)
