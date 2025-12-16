package navigationbar

import (
	"github.com/zodimo/go-compose/theme"

	"gioui.org/unit"
)

// NavigationBarColors represents the colors used by a NavigationBar.
type NavigationBarColors struct {
	ContainerColor theme.ColorDescriptor
	ContentColor   theme.ColorDescriptor
	IndicatorColor theme.ColorDescriptor
}

// NavigationBarDefaults holds the default values for NavigationBar.
var NavigationBarDefaults = navigationBarDefaults{}

type navigationBarDefaults struct{}

// Colors returns the default colors for a NavigationBar.
func (d navigationBarDefaults) Colors() NavigationBarColors {
	return NavigationBarColors{
		// Surface Container
		ContainerColor: theme.ColorHelper.ColorSelector().SurfaceRoles.Container,
		// On Surface Variant
		ContentColor: theme.ColorHelper.ColorSelector().SurfaceRoles.OnVariant,
		// Secondary Container
		IndicatorColor: theme.ColorHelper.ColorSelector().SecondaryRoles.Container,
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
