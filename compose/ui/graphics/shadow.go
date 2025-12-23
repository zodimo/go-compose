// Package graphics provides UI graphics primitives.
package graphics

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/ui/geometry"
	"github.com/zodimo/go-compose/pkg/floatutils"
	"github.com/zodimo/go-compose/pkg/floatutils/lerp"
)

// EmptyShadow is the singleton sentinel for unspecified/empty Shadow.
// It is allocated in the data segment (global) and used as a pointer to avoid allocations.
var EmptyShadow = &Shadow{
	Color:      ColorUnspecified,
	Offset:     geometry.OffsetUnspecified,
	BlurRadius: floatutils.Float32Unspecified,
}

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

// TakeOrElse returns the existing shadow if it's specified, otherwise the default.
// This is a package-level function to handle nil receivers safely.
func TakeOrElse(s, defaultShadow *Shadow) *Shadow {
	if s == nil || s == EmptyShadow {
		return defaultShadow
	}
	return s
}

// Equal performs deep equality check with epsilon tolerance for all fields.
func (s *Shadow) Equal(other *Shadow) bool {
	return ShadowEquals(s, other)
}

func (s *Shadow) IsSpecified() bool {
	return ShadowIsSpecified(s)
}

// String returns a human-readable representation.
func (s *Shadow) String() string {
	if !s.IsSpecified() {
		return "EmptyShadow"
	}
	return fmt.Sprintf("Shadow(color=%v, offset=%v, blurRadius=%.2f)",
		s.Color, s.Offset, s.BlurRadius)
}

// LerpShadow interpolates between two Shadows.
func LerpShadow(start, stop *Shadow, fraction float32) Shadow {
	start = TakeOrElse(start, EmptyShadow)
	stop = TakeOrElse(stop, EmptyShadow)

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

func ShadowIsSpecified(s *Shadow) bool {
	shallow := s != nil && s != EmptyShadow
	if shallow {
		// making sure s is not empty
		return s.Color.IsSpecified() && s.Offset.IsSpecified() && floatutils.IsSpecified(s.BlurRadius)
	}
	return shallow
}

func ShadowEquals(s1, s2 *Shadow) bool {
	if !ShadowIsSpecified(s1) && !ShadowIsSpecified(s2) {
		return true
	}

	return s1.Color == s2.Color && s1.Offset.Equal(s2.Offset) &&
		float32Equals(s1.BlurRadius, s2.BlurRadius, float32EqualityThreshold)
}
