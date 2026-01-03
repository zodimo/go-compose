package navigationbar

import (
	"github.com/zodimo/go-compose/compose/material3"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

// NavigationBarColors represents the colors used by a NavigationBar.
type NavigationBarColors struct {
	ContainerColor graphics.Color
	ContentColor   graphics.Color
	IndicatorColor graphics.Color
}

// NavigationBarDefaults holds the default values for NavigationBar.
var NavigationBarDefaults = navigationBarDefaults{}

type navigationBarDefaults struct{}

// Colors returns the default colors for a NavigationBar.
func (d navigationBarDefaults) Colors(c Composer) NavigationBarColors {
	theme := material3.Theme(c)
	return NavigationBarColors{
		// Surface Container
		ContainerColor: theme.ColorScheme().SurfaceContainer, //theme.ColorHelper.ColorSelector().SurfaceRoles.Container,
		// On Surface Variant
		ContentColor: theme.ColorScheme().SurfaceVariant.OnColor, //theme.ColorHelper.ColorSelector().SurfaceRoles.OnVariant,
		// Secondary Container
		IndicatorColor: theme.ColorScheme().SecondaryContainer.Color, //theme.ColorHelper.ColorSelector().SecondaryRoles.Container,
	}
}

// ContainerElevation returns the default elevation for a NavigationBar.
func (d navigationBarDefaults) ContainerElevation() unit.Dp {
	return unit.Dp(3) // Level 2: 3.0dp
}

// Height returns the default height for a NavigationBar.
func (d navigationBarDefaults) Height() unit.Dp {
	return unit.Dp(80)
}

func DefaultNavigationBarItemOptions() NavigationBarItemOptions {
	return NavigationBarItemOptions{
		Modifier: EmptyModifier,
	}
}
