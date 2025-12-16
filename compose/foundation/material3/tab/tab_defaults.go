package tab

import (
	"gioui.org/unit"
	"github.com/zodimo/go-compose/theme"
)

// TabRowDefaults holds default values for the TabRow and Tab components.
var TabRowDefaults = tabRowDefaults{}
var TabDefaults = tabDefaults{}

type tabRowDefaults struct{}
type tabDefaults struct{}

func (tabRowDefaults) ContainerColor() theme.ColorDescriptor {
	return theme.ColorHelper.ColorSelector().SurfaceRoles.Surface
}

func (tabRowDefaults) ContentColor() theme.ColorDescriptor {
	return theme.ColorHelper.ColorSelector().SurfaceRoles.OnSurface
}

func (tabRowDefaults) IndicatorColor() theme.ColorDescriptor {
	return theme.ColorHelper.ColorSelector().PrimaryRoles.Primary
}

func (tabRowDefaults) IndicatorHeight() unit.Dp {
	return unit.Dp(3)
}

// Indicator returns a default indicator composable.
// In the future this should be a proper composable function.
func (tabRowDefaults) Indicator() Composable {
	return nil // Default handled in Tab or TabRow if nil
}

func (tabDefaults) SelectedContentColor() theme.ColorDescriptor {
	return theme.ColorHelper.ColorSelector().PrimaryRoles.Primary
}

func (tabDefaults) UnselectedContentColor() theme.ColorDescriptor {
	return theme.ColorHelper.ColorSelector().SurfaceRoles.OnVariant
}
