package tab

import (
	"image/color"

	"gioui.org/unit"
)

// TabRowDefaults holds default values for the TabRow and Tab components.
var TabRowDefaults = tabRowDefaults{}
var TabDefaults = tabDefaults{}

type tabRowDefaults struct{}
type tabDefaults struct{}

// ContainerColor is the default container color for a TabRow.
// M3: Surface
var (
	// TODO: Replace with theme lookups
	defaultContainerColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255} // Surface
	defaultContentColor   = color.NRGBA{R: 29, G: 27, B: 32, A: 255}    // OnSurface
	defaultIndicatorColor = color.NRGBA{R: 103, G: 80, B: 164, A: 255}  // Primary
)

func (tabRowDefaults) ContainerColor() color.NRGBA {
	return defaultContainerColor
}

func (tabRowDefaults) ContentColor() color.NRGBA {
	return defaultContentColor
}

func (tabRowDefaults) IndicatorColor() color.NRGBA {
	return defaultIndicatorColor
}

func (tabRowDefaults) IndicatorHeight() unit.Dp {
	return unit.Dp(3)
}

// Indicator returns a default indicator composable.
// In the future this should be a proper composable function.
func (tabRowDefaults) Indicator() Composable {
	return nil // Default handled in Tab or TabRow if nil
}

func (tabDefaults) SelectedContentColor() color.NRGBA {
	return defaultIndicatorColor // Primary
}

func (tabDefaults) UnselectedContentColor() color.NRGBA {
	return color.NRGBA{R: 73, G: 69, B: 79, A: 255} // OnSurfaceVariant
}
