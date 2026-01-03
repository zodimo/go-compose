package radiobutton

import (
	"github.com/zodimo/go-compose/compose/material3"
	"github.com/zodimo/go-compose/compose/ui/graphics"

	"git.sr.ht/~schnwalter/gio-mw/token"
)

// RadioButtonColors represents the colors used by a RadioButton in different states.
type RadioButtonColors struct {
	SelectedColor   graphics.Color
	UnselectedColor graphics.Color
	DisabledColor   graphics.Color
}

// Defaults contains the default values for RadioButton.
var Defaults = radioButtonDefaults{}

type radioButtonDefaults struct{}

// Colors returns the default colors for a RadioButton.
func (d radioButtonDefaults) Colors(c Composer) RadioButtonColors {
	theme := material3.Theme(c)
	return RadioButtonColors{
		SelectedColor:   theme.ColorScheme().Primary.Color,                                                      //theme.ColorHelper.ColorSelector().PrimaryRoles.Primary,
		UnselectedColor: theme.ColorScheme().SurfaceVariant.OnColor,                                             //theme.ColorHelper.ColorSelector().SurfaceRoles.OnVariant,
		DisabledColor:   graphics.SetOpacity(theme.ColorScheme().Surface.OnColor, float32(token.OpacityLevel5)), //theme.ColorHelper.ColorSelector().SurfaceRoles.OnSurface.SetOpacity(token.OpacityLevel5),
	}
}

// Color returns the resolved color based on state (used internally for drawing).
func (c RadioButtonColors) Color(enabled, selected bool) graphics.Color {
	if !enabled {
		return c.DisabledColor
	}
	if selected {
		return c.SelectedColor
	}
	return c.UnselectedColor
}
