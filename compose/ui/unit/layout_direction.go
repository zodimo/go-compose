package unit

import (
	textunit "github.com/zodimo/go-compose/compose/ui/text/unit"
)

// LayoutDirection represents the horizontal layout direction of content.
//
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-unit/src/commonMain/kotlin/androidx/compose/ui/unit/LayoutDirection.kt
type LayoutDirection = textunit.LayoutDirection

const (
	// LayoutDirectionLtr represents left-to-right layout direction.
	LayoutDirectionLtr = textunit.LayoutDirectionLtr

	// LayoutDirectionRtl represents right-to-left layout direction.
	LayoutDirectionRtl = textunit.LayoutDirectionRtl
)
