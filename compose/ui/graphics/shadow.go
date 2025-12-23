// Package graphics provides UI graphics primitives.
package graphics

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/ui/geometry"
	"github.com/zodimo/go-compose/pkg/floatutils"
	"github.com/zodimo/go-compose/pkg/floatutils/lerp"
)

var ShadowUnspecified = Shadow{
	Color:      ColorUnspecified,
	Offset:     geometry.OffsetUnspecified,
	BlurRadius: floatutils.Float32Unspecified,
}

// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-graphics/src/commonMain/kotlin/androidx/compose/ui/graphics/Shadow.kt;drc=4970f6e96cdb06089723da0ab8ec93ae3f067c7a;l=27

// Shadow represents a single shadow with color, offset, and blur radius.
// All fields are immutable after creation; use Copy() to create modified versions.
type Shadow struct {
	Color      Color   // Color of the shadow (typically with alpha)
	Offset     Offset  // Offset from the element
	BlurRadius float32 // Blur radius in pixels
}

// Zero constants (define these once)
var (
	// ShadowNone represents no shadow. Use this constant instead of allocating a new zero Shadow.
	ShadowNone = NewShadow(ColorBlack, ZeroOffset, 0)
)

// NewShadow creates a new Shadow instance.
func NewShadow(color Color, offset geometry.Offset, blurRadius float32) Shadow {
	return Shadow{
		Color:      color,
		Offset:     offset,
		BlurRadius: blurRadius,
	}
}

// Copy creates a new Shadow with optional field overrides.
// Pass nil for parameters you want unchanged.
func (s Shadow) Copy(color *Color, offset *Offset, blurRadius *float32) Shadow {
	result := s
	if color != nil {
		result.Color = *color
	}
	if offset != nil {
		result.Offset = *offset
	}
	if blurRadius != nil {
		result.BlurRadius = *blurRadius
	}
	return result
}

func (s Shadow) IsSpecified() bool {
	return !(s.Color == ColorUnspecified &&
		s.Offset == geometry.OffsetUnspecified &&
		s.BlurRadius == floatutils.Float32Unspecified)
}

// Equal performs deep equality check with epsilon tolerance for all fields.
func (s Shadow) Equal(other Shadow) bool {
	if !s.IsSpecified() && !other.IsSpecified() {
		return true
	}

	return s.Color == other.Color && s.Offset.Equal(other.Offset) &&
		float32Equals(s.BlurRadius, other.BlurRadius, float32EqualityThreshold)
}

// String returns a human-readable representation.
func (s Shadow) String() string {
	if !s.IsSpecified() {
		return "ShadowUnspecified"
	}
	return fmt.Sprintf("Shadow(color=%v, offset=%v, blurRadius=%.2f)",
		s.Color, s.Offset, s.BlurRadius)
}

// LerpShadow interpolates between two Shadows.
func LerpShadow(start, stop *Shadow, fraction float32) Shadow {
	if start == nil || stop == nil {
		panic("Shadow must be specified")
	}
	if fraction == 0 {
		return *start
	}
	if fraction == 1 {
		return *stop
	}
	return NewShadow(
		Lerp(start.Color, stop.Color, fraction),
		geometry.LerpOffset(start.Offset, stop.Offset, fraction),
		lerp.Between32(start.BlurRadius, stop.BlurRadius, fraction),
	)
}
