package slider

import (
	gioUnit "gioui.org/unit"
	"github.com/zodimo/go-compose/theme"
)

// SliderDefaults holds default values for the Slider component.
var SliderDefaults = sliderDefaults{}

type sliderDefaults struct{}

// Colors returns the default SliderColors.
func (d sliderDefaults) Colors() SliderColors {
	return SliderColors{
		ThumbColor:            theme.ColorHelper.UnspecifiedColor(),
		ActiveTrackColor:      theme.ColorHelper.UnspecifiedColor(),
		ActiveTickColor:       theme.ColorHelper.UnspecifiedColor(),
		InactiveTrackColor:    theme.ColorHelper.UnspecifiedColor(),
		InactiveTickColor:     theme.ColorHelper.UnspecifiedColor(),
		DisabledThumbColor:    theme.ColorHelper.UnspecifiedColor(),
		DisabledActiveTrack:   theme.ColorHelper.UnspecifiedColor(),
		DisabledActiveTick:    theme.ColorHelper.UnspecifiedColor(),
		DisabledInactiveTrack: theme.ColorHelper.UnspecifiedColor(),
		DisabledInactiveTick:  theme.ColorHelper.UnspecifiedColor(),
	}
}

// Dimensions constants
var (
	TrackHeight     = gioUnit.Dp(4)
	ThumbSize       = gioUnit.Dp(20)
	ActiveThumbSize = gioUnit.Dp(28) // M3 State Layer/Enlarged handle
	TickSize        = gioUnit.Dp(2)
	ThumbTrackGap   = gioUnit.Dp(6) // Approximate
)
