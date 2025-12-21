package graphics

import (
	"fmt"
	"image/color"

	"github.com/zodimo/go-compose/compose/ui/geometry"
	"github.com/zodimo/go-compose/compose/ui/utils/lerp"
	"github.com/zodimo/go-compose/theme"
)

// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-graphics/src/commonMain/kotlin/androidx/compose/ui/graphics/Shadow.kt;drc=4970f6e96cdb06089723da0ab8ec93ae3f067c7a;l=27

// Shadow represents a single shadow with color, offset, and blur radius.
// All fields are immutable after creation; use Copy() to create modified versions.
type Shadow struct {
	Color      Color   // Color of the shadow (typically with alpha)
	Offset     Offset  // Offset from the element
	BlurRadius float32 // Blur radius in pixels
}

// None represents no shadow. Use this constant instead of allocating a new zero Shadow.
var None = Shadow{
	Color:      theme.ColorHelper.SpecificColor(color.RGBA{R: 0, G: 0, B: 0, A: 255}), // Define Black as: Color{R:0, G:0, B:0, A:1}
	Offset:     ZeroOffset,
	BlurRadius: 0,
}

// NewShadow creates a new Shadow instance with the given properties.
func NewShadow(color Color, offset geometry.Offset, blurRadius float32) Shadow {
	return Shadow{
		Color:      color,
		Offset:     offset,
		BlurRadius: blurRadius,
	}
}

// Copy creates a new Shadow with optional field overrides.
// Pass nil for any parameter you want to keep unchanged.
func (s Shadow) Copy(color *Color, offset *geometry.Offset, blurRadius *float32) Shadow {
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

// Equal performs deep equality check on all Shadow fields.
func (s Shadow) Equal(other Shadow) bool {
	return s.Color.Compare(other.Color) &&
		s.Offset.Equal(other.Offset) &&
		s.BlurRadius == other.BlurRadius
}

// String returns a human-readable representation of the Shadow.
func (s Shadow) String() string {
	return fmt.Sprintf("Shadow(color=%v, offset=%v, blurRadius=%.2f)",
		s.Color, s.Offset, s.BlurRadius)
}

// Lerp performs linear interpolation between two shadows.
// Each field (color, offset, blurRadius) is interpolated independently.
func Lerp(start, stop Shadow, fraction float32) Shadow {
	// Color interpolation (components)
	colorLerped := theme.ColorLerp(start.Color, stop.Color, fraction)

	// Offset interpolation (each axis)
	offset := geometry.Offset{
		X: lerp.Float32(start.Offset.X, stop.Offset.X, fraction),
		Y: lerp.Float32(start.Offset.Y, stop.Offset.Y, fraction),
	}

	// Blur radius interpolation
	blurRadius := lerp.Float32(start.BlurRadius, stop.BlurRadius, fraction)

	return NewShadow(colorLerped, offset, blurRadius)
}
