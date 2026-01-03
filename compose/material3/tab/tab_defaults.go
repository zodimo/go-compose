package tab

import (
	"github.com/zodimo/go-compose/compose/material3"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

// TabRowDefaults holds default values for the TabRow and Tab components.
var TabRowDefaults = tabRowDefaults{}
var TabDefaults = tabDefaults{}

type tabRowDefaults struct{}
type tabDefaults struct{}

func (tabRowDefaults) ContainerColor(c Composer) graphics.Color {
	return material3.Theme(c).ColorScheme().Surface.Color
}

func (tabRowDefaults) ContentColor(c Composer) graphics.Color {
	return material3.Theme(c).ColorScheme().Surface.OnColor
}

func (tabRowDefaults) IndicatorColor(c Composer) graphics.Color {
	return material3.Theme(c).ColorScheme().Primary.Color
}

func (tabRowDefaults) IndicatorHeight() unit.Dp {
	return unit.Dp(3)
}

// Indicator returns a default indicator composable.
// In the future this should be a proper composable function.
func (tabRowDefaults) Indicator() Composable {
	return nil // Default handled in Tab or TabRow if nil
}

func (tabDefaults) SelectedContentColor(c Composer) graphics.Color {
	return material3.Theme(c).ColorScheme().Primary.Color
}

func (tabDefaults) UnselectedContentColor(c Composer) graphics.Color {
	return material3.Theme(c).ColorScheme().Surface.OnColor
}
