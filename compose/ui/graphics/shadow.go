// Package graphics provides UI graphics primitives.
package graphics

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/ui/geometry"
	"github.com/zodimo/go-compose/pkg/floatutils"
	"github.com/zodimo/go-compose/pkg/floatutils/lerp"
	"github.com/zodimo/go-sentinel-helper/sentinel/companions"
)

// shadowUnspecified is the singleton sentinel for unspecified/empty Shadow.
// It is allocated in the data segment (global) and used as a pointer to avoid
// allocations. It is deliberately package-private: the only public way to obtain
// it is ShadowCompanion.Unspecified().
var shadowUnspecified = &Shadow{
	Color:      colorUnspecified,
	Offset:     geometry.OffsetCompanion.Unspecified(),
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

// Helper functions

// Copy creates a new Shadow with optional field overrides.
func CopyShadow(s *Shadow, options ...ShadowOption) *Shadow {
	opt := *shadowUnspecified // dereference unspecified

	for _, option := range options {
		option(&opt)
	}

	return &Shadow{
		Color:      takeOrElseColor(opt.Color, s.Color),
		Offset:     geometry.OffsetCompanion.TakeOrElse(opt.Offset, s.Offset),
		BlurRadius: floatutils.TakeOrElse(opt.BlurRadius, s.BlurRadius),
	}
}

func StringShadow(s *Shadow) string {
	if !isSpecifiedShadow(s) {
		return "Shadow{Unspecified}"
	}
	return fmt.Sprintf("Shadow{color=%v, offset=%v, blurRadius=%.2f}",
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

	start = coalesceShadow(start, shadowUnspecified)
	stop = coalesceShadow(stop, shadowUnspecified)

	start = takeOrElseShadow(start, shadowUnspecified)
	stop = takeOrElseShadow(stop, shadowUnspecified)

	return NewShadow(
		LerpColor(start.Color, stop.Color, fraction),
		geometry.LerpOffset(start.Offset, stop.Offset, fraction),
		lerp.Between32(start.BlurRadius, stop.BlurRadius, fraction),
	)
}

// === Sentinel operations ===

func isSpecifiedShadow(s *Shadow) bool {
	return s != nil && s != shadowUnspecified
}

func isUnspecifiedShadow(s *Shadow) bool {
	return !isSpecifiedShadow(s)
}

func takeOrElseShadow(s, def *Shadow) *Shadow {
	if !isSpecifiedShadow(s) {
		return def
	}
	return s
}

// mergeShadow layers b onto a: a value from b wins only where b is specified.
func mergeShadow(a, b *Shadow) *Shadow {
	a = coalesceShadow(a, shadowUnspecified)
	b = coalesceShadow(b, shadowUnspecified)

	if a == shadowUnspecified {
		return b
	}
	if b == shadowUnspecified {
		return a
	}

	// Both are custom: allocate new merged style. b's value wins where specified.
	return &Shadow{
		Color:      takeOrElseColor(b.Color, a.Color),
		Offset:     geometry.OffsetCompanion.TakeOrElse(b.Offset, a.Offset),
		BlurRadius: floatutils.TakeOrElse(b.BlurRadius, a.BlurRadius),
	}
}

func stringShadow(s *Shadow) string {
	return StringShadow(s)
}

// sameShadow is identity (2 ns).
func sameShadow(a, b *Shadow) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return b == shadowUnspecified
	}
	if b == nil {
		return a == shadowUnspecified
	}
	return a == b
}

// semanticEqualShadow is field-by-field equality (20 ns). Two unspecified
// shadows are equal even though their NaN payloads would not compare equal.
func semanticEqualShadow(a, b *Shadow) bool {

	a = coalesceShadow(a, shadowUnspecified)
	b = coalesceShadow(b, shadowUnspecified)

	if a == shadowUnspecified && b == shadowUnspecified {
		return true
	}

	return a.Color == b.Color &&
		geometry.OffsetCompanion.SemanticEqual(a.Offset, b.Offset) &&
		floatutils.Float32Equals(a.BlurRadius, b.BlurRadius, float32EqualityThreshold)
}

func equalShadow(a, b *Shadow) bool {
	if !sameShadow(a, b) {
		return semanticEqualShadow(a, b)
	}
	return true
}

func coalesceShadow(ptr, def *Shadow) *Shadow {
	if ptr == nil {
		return def
	}
	return ptr
}

// copyShadow returns a copy of s with the given options applied on top; the
// source is never mutated.
func copyShadow(s *Shadow, options ...func(*Shadow)) *Shadow {
	opt := *shadowUnspecified
	for _, option := range options {
		option(&opt)
	}
	return &Shadow{
		Color:      takeOrElseColor(opt.Color, s.Color),
		Offset:     geometry.OffsetCompanion.TakeOrElse(opt.Offset, s.Offset),
		BlurRadius: floatutils.TakeOrElse(opt.BlurRadius, s.BlurRadius),
	}
}

// ShadowCompanion is the single exported sentinel-family symbol for Shadow.
var ShadowCompanion = companions.NewStructCompanion(
	func() *Shadow { return shadowUnspecified },
	stringShadow,
	isSpecifiedShadow,
	isUnspecifiedShadow,
	takeOrElseShadow,
	mergeShadow,
	sameShadow,
	semanticEqualShadow,
	equalShadow,
	coalesceShadow,
	copyShadow,
)

// Short for isSpecifiedShadow
func IsSpecifiedShadow(s *Shadow) bool {
	return isSpecifiedShadow(s)
}
func TakeOrElseShadow(s, def *Shadow) *Shadow {
	return takeOrElseShadow(s, def)
}

// Identity (2 ns)
func SameShadow(a, b *Shadow) bool {
	return sameShadow(a, b)
}

// Semantic equality (field-by-field, 20 ns)
func SemanticEqualShadow(a, b *Shadow) bool {
	return semanticEqualShadow(a, b)
}

func EqualShadow(a, b *Shadow) bool {
	return equalShadow(a, b)
}

func MergeShadow(a, b *Shadow) *Shadow {
	return mergeShadow(a, b)
}

func CoalesceShadow(ptr, def *Shadow) *Shadow {
	return coalesceShadow(ptr, def)
}
