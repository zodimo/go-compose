package platform

import (
	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

// LocalDensity is a CompositionLocal that provides the Density to the composition.
var LocalDensity = compose.StaticCompositionLocalOf[unit.Density](func() unit.Density {
	panic("CompositionLocal LocalDensity not present")
})
