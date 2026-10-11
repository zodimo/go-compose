package foundation

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

// BorderStroke specifies the stroke to draw a border with.
//
// See: https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/foundation/foundation/src/commonMain/kotlin/androidx/compose/foundation/BorderStroke.kt
type BorderStroke struct {
	// Width of the border. Use unit.DpHairline for a one-pixel border.
	Width unit.Dp
	// Brush to paint the border with.
	Brush graphics.Brush
}

// 1. BorderStrokeUnspecified is the sentinel singleton for an unspecified BorderStroke.
var BorderStrokeUnspecified = &BorderStroke{
	Width: unit.DpCompanion.Unspecified(),
	Brush: graphics.BrushCompanion.Unspecified(),
}

// 2. IsSpecifiedBorderStroke returns true if the BorderStroke is specified (not nil and not the Unspecified singleton).
func IsSpecifiedBorderStroke(bs *BorderStroke) bool {
	return bs != nil && bs != BorderStrokeUnspecified
}

// 3. TakeOrElseBorderStroke returns bs if specified, otherwise returns defaultBS.
func TakeOrElseBorderStroke(bs, defaultBS *BorderStroke) *BorderStroke {
	if bs == nil || bs == BorderStrokeUnspecified {
		return defaultBS
	}
	return bs
}

// 4. MergeBorderStroke merges two BorderStrokes, preferring b's specified values over a's.
func MergeBorderStroke(a, b *BorderStroke) *BorderStroke {
	a = CoalesceBorderStroke(a, BorderStrokeUnspecified)
	b = CoalesceBorderStroke(b, BorderStrokeUnspecified)

	if a == BorderStrokeUnspecified {
		return b
	}
	if b == BorderStrokeUnspecified {
		return a
	}

	// Both are custom: allocate new merged struct
	return &BorderStroke{
		Width: unit.DpCompanion.TakeOrElse(b.Width, a.Width),
		Brush: graphics.TakeOrElseBrush(b.Brush, a.Brush),
	}
}

// 5. StringBorderStroke returns a string representation of BorderStroke.
func StringBorderStroke(bs *BorderStroke) string {
	if !IsSpecifiedBorderStroke(bs) {
		return "BorderStroke{Unspecified}"
	}

	return fmt.Sprintf(
		"BorderStroke{Width: %s, Brush: %v}",
		bs.Width,
		bs.Brush,
	)
}

// 6. CoalesceBorderStroke returns ptr if not nil, otherwise returns def.
func CoalesceBorderStroke(ptr, def *BorderStroke) *BorderStroke {
	if ptr == nil {
		return def
	}
	return ptr
}

// 7. SameBorderStroke returns true if a and b are the same pointer or both unspecified.
func SameBorderStroke(a, b *BorderStroke) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return b == BorderStrokeUnspecified
	}
	if b == nil {
		return a == BorderStrokeUnspecified
	}
	return a == b
}

// 8. SemanticEqualBorderStroke checks field-by-field equality.
func SemanticEqualBorderStroke(a, b *BorderStroke) bool {
	a = CoalesceBorderStroke(a, BorderStrokeUnspecified)
	b = CoalesceBorderStroke(b, BorderStrokeUnspecified)

	return a.Width == b.Width && graphics.EqualBrush(a.Brush, b.Brush)
}

// 9. EqualBorderStroke returns true if a and b are semantically equal.
func EqualBorderStroke(a, b *BorderStroke) bool {
	if SameBorderStroke(a, b) {
		return true
	}
	return SemanticEqualBorderStroke(a, b)
}

// BorderStrokeOption is a functional option for CopyBorderStroke.
type BorderStrokeOption func(*BorderStroke)

// WithBorderStrokeWidth sets the width option.
func WithBorderStrokeWidth(width unit.Dp) BorderStrokeOption {
	return func(o *BorderStroke) {
		o.Width = width
	}
}

// WithBorderStrokeBrush sets the brush option.
func WithBorderStrokeBrush(brush graphics.Brush) BorderStrokeOption {
	return func(o *BorderStroke) {
		o.Brush = brush
	}
}

// 10. CopyBorderStroke creates a copy with optional modifications.
func CopyBorderStroke(bs *BorderStroke, options ...BorderStrokeOption) *BorderStroke {
	opt := *BorderStrokeUnspecified

	for _, option := range options {
		option(&opt)
	}

	bs = CoalesceBorderStroke(bs, BorderStrokeUnspecified)

	return &BorderStroke{
		Width: unit.DpCompanion.TakeOrElse(opt.Width, bs.Width),
		Brush: graphics.TakeOrElseBrush(opt.Brush, bs.Brush),
	}
}

// Convenience constructors

// NewBorderStroke creates a new BorderStroke with the given width and brush.
func NewBorderStroke(width unit.Dp, brush graphics.Brush) *BorderStroke {
	return &BorderStroke{
		Width: width,
		Brush: brush,
	}
}

// NewBorderStrokeWithColor creates a new BorderStroke with the given width and solid color.
// This is a convenience function equivalent to NewBorderStroke(width, &graphics.SolidColor{Value: color}).
func NewBorderStrokeWithColor(width unit.Dp, color graphics.Color) *BorderStroke {
	return &BorderStroke{
		Width: width,
		Brush: &graphics.SolidColor{Value: color},
	}
}
