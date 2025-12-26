// Package graphics provides UI graphics primitives.
package graphics

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/ui/geometry"
	"github.com/zodimo/go-compose/pkg/floatutils"
	"github.com/zodimo/go-compose/pkg/floatutils/lerp"
)

// ShadowUnspecified is the singleton sentinel for unspecified/empty Shadow.
// It is allocated in the data segment (global) and used as a pointer to avoid allocations.
var ShadowUnspecified = &Shadow{
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
func NewShadow(color Color, offset geometry.Offset, blurRadius float32) *Shadow {
	return &Shadow{
		Color:      color,
		Offset:     offset,
		BlurRadius: blurRadius,
	}
}

// Copy creates a new Shadow with optional field overrides.
func (s Shadow) Copy(options ...ShadowOption) *Shadow {
	opt := ShadowOptionsDefault
	for _, option := range options {
		option(&opt)
	}

	return &Shadow{
		Color:      opt.Color.TakeOrElse(s.Color),
		Offset:     opt.Offset.TakeOrElse(s.Offset),
		BlurRadius: floatutils.TakeOrElse(opt.BlurRadius, s.BlurRadius),
	}
}

// Helper functions

func StringShadow(s *Shadow) string {
	if !IsSpecifiedShadow(s) {
		return "EmptyShadow"
	}
	return fmt.Sprintf("Shadow(color=%v, offset=%v, blurRadius=%.2f)",
		s.Color, s.Offset, s.BlurRadius)
}

// LerpShadow interpolates between two Shadows.
func LerpShadow(start, stop *Shadow, fraction float32) *Shadow {

	if fraction == 0 {
		return start
	}
	if fraction == 1 {
		return stop
	}

	start = CoalesceShadow(start, ShadowUnspecified)
	stop = CoalesceShadow(stop, ShadowUnspecified)

	start = TakeOrElseShadow(start, ShadowUnspecified)
	stop = TakeOrElseShadow(stop, ShadowUnspecified)

	return NewShadow(
		LerpColor(start.Color, stop.Color, fraction),
		geometry.LerpOffset(start.Offset, stop.Offset, fraction),
		lerp.Between32(start.BlurRadius, stop.BlurRadius, fraction),
	)
}

func takeOrElse[T comparable](a, b, sentinel T) T {
	if a != sentinel {
		return a
	}
	return b
}

// Short for IsSpecifiedShadow
func IsSpecifiedShadow(s *Shadow) bool {
	return s != nil && s != ShadowUnspecified
}
func TakeOrElseShadow(s, def *Shadow) *Shadow {
	if !IsSpecifiedShadow(s) {
		return def
	}
	return s
}

// Identity (2 ns)
func SameShadow(a, b *Shadow) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return b == ShadowUnspecified
	}
	if b == nil {
		return a == ShadowUnspecified
	}
	return a == b
}

func EqualShadow(a, b *Shadow) bool {
	if !SameShadow(a, b) {
		return SemanticEqualShadow(a, b)
	}
	return true
}

// Semantic equality (field-by-field, 20 ns)
func SemanticEqualShadow(a, b *Shadow) bool {

	a = CoalesceShadow(a, ShadowUnspecified)
	b = CoalesceShadow(b, ShadowUnspecified)

	return a.Color == b.Color &&
		a.Offset.Equal(b.Offset) &&
		floatutils.Float32Equals(a.BlurRadius, b.BlurRadius, float32EqualityThreshold)
}

func MergeShadow(a, b *Shadow) *Shadow {
	a = CoalesceShadow(a, ShadowUnspecified)
	b = CoalesceShadow(b, ShadowUnspecified)

	if a == ShadowUnspecified {
		return b
	}
	if b == ShadowUnspecified {
		return a
	}

	// Both are custom: allocate new merged style
	return &Shadow{
		Color:      takeOrElse(a.Color, b.Color, ColorUnspecified),
		Offset:     takeOrElse(a.Offset, b.Offset, geometry.OffsetUnspecified),
		BlurRadius: takeOrElse(a.BlurRadius, b.BlurRadius, floatutils.Float32Unspecified),
	}
}

func CoalesceShadow(ptr, def *Shadow) *Shadow {
	if ptr == nil {
		return def
	}
	return ptr
}
