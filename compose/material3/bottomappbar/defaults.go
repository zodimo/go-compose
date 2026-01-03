package bottomappbar

import (
	"github.com/zodimo/go-compose/compose/material3"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

// BottomAppBarColors represents the colors used by a BottomAppBar.
type BottomAppBarColors struct {
	ContainerColor graphics.Color
	ContentColor   graphics.Color
}

// BottomAppBarDefaults holds the default values for BottomAppBar.
var BottomAppBarDefaults = bottomAppBarDefaults{}

type bottomAppBarDefaults struct{}

// Colors returns the default colors for a BottomAppBar.
func (d bottomAppBarDefaults) Colors(c Composer) BottomAppBarColors {
	theme := material3.Theme(c)
	return BottomAppBarColors{
		ContainerColor: theme.ColorScheme().SurfaceContainer,
		ContentColor:   theme.ColorScheme().SurfaceVariant.OnColor,
	}
}

// ContainerElevation returns the default elevation for a BottomAppBar.
func (d bottomAppBarDefaults) ContainerElevation() unit.Dp {
	return unit.Dp(3) // Level 2: 3.0dp
}

// ContentPadding returns the default content padding for a BottomAppBar.
// Standard is Horizontal 16dp, Vertical 12dp (for tall bar 80dp) or 16dp (for short 64dp).
// Let's use internal values in implementation, but expose a default here if we define PaddingValues.
// For now, I'll return specific values.
func (d bottomAppBarDefaults) ContentPadding() (start, top, end, bottom unit.Dp) {
	return unit.Dp(16), unit.Dp(12), unit.Dp(16), unit.Dp(12)
}
