package appbar

import (
	"github.com/zodimo/go-compose/theme"
)

// TopAppBarColors represents the colors used by a TopAppBar in different states.
type TopAppBarColors struct {
	ContainerColor             theme.ColorDescriptor
	ScrolledContainerColor     theme.ColorDescriptor
	NavigationIconContentColor theme.ColorDescriptor
	TitleContentColor          theme.ColorDescriptor
	ActionIconContentColor     theme.ColorDescriptor
}

// TopAppBarDefaults holds the default values for TopAppBar.
var TopAppBarDefaults = topAppBarDefaults{}

type topAppBarDefaults struct{}

// Colors returns the default colors for a TopAppBar.
func (d topAppBarDefaults) Colors() TopAppBarColors {
	return TopAppBarColors{
		// Surface
		ContainerColor: theme.ColorHelper.ColorSelector().SurfaceRoles.Surface,
		// Surface Container
		ScrolledContainerColor: theme.ColorHelper.ColorSelector().SurfaceRoles.Container,
		// On Surface
		NavigationIconContentColor: theme.ColorHelper.ColorSelector().SurfaceRoles.OnSurface,
		// On Surface
		TitleContentColor: theme.ColorHelper.ColorSelector().SurfaceRoles.OnSurface,
		// On Surface Variant
		ActionIconContentColor: theme.ColorHelper.ColorSelector().SurfaceRoles.OnVariant,
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
