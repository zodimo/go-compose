package slider

import (
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/internal/layoutnode"
)

// SliderColors represents the colors used by a Slider in different states.
type SliderColors struct {
	ThumbColor            graphics.Color
	ActiveTrackColor      graphics.Color
	ActiveTickColor       graphics.Color
	InactiveTrackColor    graphics.Color
	InactiveTickColor     graphics.Color
	DisabledThumbColor    graphics.Color
	DisabledActiveTrack   graphics.Color
	DisabledActiveTick    graphics.Color
	DisabledInactiveTrack graphics.Color
	DisabledInactiveTick  graphics.Color
}

// ThumbColor returns the color of the thumb based on the enabled state.
func (c SliderColors) Thumb(enabled bool) graphics.Color {
	if enabled {
		return c.ThumbColor
	}
	return c.DisabledThumbColor
}

// TrackColor returns the color of the track based on the enabled and active state.
func (c SliderColors) Track(enabled, active bool) graphics.Color {
	if enabled {
		if active {
			return c.ActiveTrackColor
		}
		return c.InactiveTrackColor
	}
	if active {
		return c.DisabledActiveTrack
	}
	return c.DisabledInactiveTrack
}

// TickColor returns the color of the ticks based on the enabled and active state.
func (c SliderColors) Tick(enabled, active bool) graphics.Color {
	if enabled {
		if active {
			return c.ActiveTickColor
		}
		return c.InactiveTickColor
	}
	if active {
		return c.DisabledActiveTick
	}
	return c.DisabledInactiveTick
}

func (c SliderColors) Layout(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
	return layoutnode.LayoutDimensions{}
}
