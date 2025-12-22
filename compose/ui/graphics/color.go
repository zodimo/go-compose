package graphics

import "github.com/zodimo/go-compose/theme"

// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-graphics/src/commonMain/kotlin/androidx/compose/ui/graphics/Color.kt;drc=4970f6e96cdb06089723da0ab8ec93ae3f067c7a;l=655
/**
 * Alternative to `() -> Color` that's useful for avoiding boxing.
 *
 * Can be used as:
 *
 * fun nonBoxedArgs(color: ColorProducer?)
 */
type ColorProducer = func() Color

type Color = theme.ColorDescriptor

type OpacityLevel = theme.OpacityLevel

var ColorUnspecified = theme.UnspecifiedColor()
