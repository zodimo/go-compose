package appbar

import (
	"image/color"
)

// TopAppBarColors represents the colors used by a TopAppBar in different states.
type TopAppBarColors struct {
	ContainerColor             color.Color
	ScrolledContainerColor     color.Color
	NavigationIconContentColor color.Color
	TitleContentColor          color.Color
	ActionIconContentColor     color.Color
}

// TopAppBarDefaults holds the default values for TopAppBar.
var TopAppBarDefaults = topAppBarDefaults{}

type topAppBarDefaults struct{}

// CenterAlignedTopAppBarColors returns the default colors for a CenterAlignedTopAppBar.
func (d topAppBarDefaults) Colors() TopAppBarColors {
	return TopAppBarColors{
		// Surface
		ContainerColor: color.NRGBA{R: 251, G: 252, B: 254, A: 255},
		// Surface Container
		ScrolledContainerColor: color.NRGBA{R: 241, G: 244, B: 249, A: 255},
		// On Surface
		NavigationIconContentColor: color.NRGBA{R: 25, G: 28, B: 32, A: 255},
		// On Surface
		TitleContentColor: color.NRGBA{R: 25, G: 28, B: 32, A: 255},
		// On Surface Variant
		ActionIconContentColor: color.NRGBA{R: 67, G: 71, B: 78, A: 255},
	}
}

// MediumTopAppBarColors returns the default colors for a MediumTopAppBar.
func (d topAppBarDefaults) MediumTopAppBarColors() TopAppBarColors {
	return d.Colors()
}

// LargeTopAppBarColors returns the default colors for a LargeTopAppBar.
func (d topAppBarDefaults) LargeTopAppBarColors() TopAppBarColors {
	return d.Colors()
}
