package appbar

import (
	"github.com/zodimo/go-compose/compose/material3"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/pkg/api"
)

// TopAppBarColors represents the colors used by a TopAppBar in different states.
type TopAppBarColors struct {
	ContainerColor             graphics.Color
	ScrolledContainerColor     graphics.Color
	NavigationIconContentColor graphics.Color
	TitleContentColor          graphics.Color
	ActionIconContentColor     graphics.Color
}

// TopAppBarDefaults holds the default values for TopAppBar.
var TopAppBarDefaults = topAppBarDefaults{}

type topAppBarDefaults struct{}

// Colors returns the default colors for a TopAppBar.
func (d topAppBarDefaults) Colors(c api.Composer) TopAppBarColors {
	theme := material3.Theme(c)
	return TopAppBarColors{
		// Surface
		ContainerColor: theme.ColorScheme().Surface.Color, //theme.ColorHelper.ColorSelector().SurfaceRoles.Surface,
		// Surface Container
		ScrolledContainerColor: theme.ColorScheme().SurfaceContainer, //theme.ColorHelper.ColorSelector().SurfaceRoles.Container,
		// On Surface
		NavigationIconContentColor: theme.ColorScheme().Surface.OnColor, //theme.ColorHelper.ColorSelector().SurfaceRoles.OnSurface,
		// On Surface
		TitleContentColor: theme.ColorScheme().Surface.OnColor, //theme.ColorHelper.ColorSelector().SurfaceRoles.OnSurface,
		// On Surface Variant
		ActionIconContentColor: theme.ColorScheme().SurfaceVariant.OnColor, //theme.ColorHelper.ColorSelector().SurfaceRoles.OnVariant,
	}
}

// MediumTopAppBarColors returns the default colors for a MediumTopAppBar.
func (d topAppBarDefaults) MediumTopAppBarColors(c api.Composer) TopAppBarColors {
	return d.Colors(c)
}

// LargeTopAppBarColors returns the default colors for a LargeTopAppBar.
func (d topAppBarDefaults) LargeTopAppBarColors(c api.Composer) TopAppBarColors {
	return d.Colors(c)
}
