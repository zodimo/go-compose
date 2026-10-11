package slider

import (
	gioUnit "gioui.org/unit"
	"github.com/zodimo/go-compose/compose/ui/graphics"
)

// SliderDefaults holds default values for the Slider component.
var SliderDefaults = sliderDefaults{}

type sliderDefaults struct{}

// Colors returns the default SliderColors.
func (d sliderDefaults) Colors() SliderColors {
	return SliderColors{
		ThumbColor:            graphics.ColorCompanion.Unspecified(),
		ActiveTrackColor:      graphics.ColorCompanion.Unspecified(),
		ActiveTickColor:       graphics.ColorCompanion.Unspecified(),
		InactiveTrackColor:    graphics.ColorCompanion.Unspecified(),
		InactiveTickColor:     graphics.ColorCompanion.Unspecified(),
		DisabledThumbColor:    graphics.ColorCompanion.Unspecified(),
		DisabledActiveTrack:   graphics.ColorCompanion.Unspecified(),
		DisabledActiveTick:    graphics.ColorCompanion.Unspecified(),
		DisabledInactiveTrack: graphics.ColorCompanion.Unspecified(),
		DisabledInactiveTick:  graphics.ColorCompanion.Unspecified(),
	}
}

// Dimensions constants
// @TODO this should be compose.ui.unit
var (
	TrackHeight     = gioUnit.Dp(4)
	ThumbSize       = gioUnit.Dp(20)
	ActiveThumbSize = gioUnit.Dp(28) // M3 State Layer/Enlarged handle
	TickSize        = gioUnit.Dp(2)
	ThumbTrackGap   = gioUnit.Dp(6) // Approximate
)
