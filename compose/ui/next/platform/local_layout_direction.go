package platform

import (
	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

// LocalLayoutDirection is a CompositionLocal that provides the layout direction to the composition.
// This allows components to adapt their layout for left-to-right (LTR) or right-to-left (RTL) languages.
var LocalLayoutDirection = compose.StaticCompositionLocalOf[unit.LayoutDirection](func() unit.LayoutDirection {
	return unit.LayoutDirectionLtr
})
